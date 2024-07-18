package block

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"syscall"
	"time"

	"github.com/oklog/ulid"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/apiutil"
	kconfig "sigs.k8s.io/controller-runtime/pkg/client/config"

	"github.com/WhizardTelemetry/whizard/pkg/api/monitoring/v1alpha1"
	"github.com/WhizardTelemetry/whizard/pkg/util"
)

const (
	binPath = "/bin/thanos"

	host    = "http://0.0.0.0:10902"
	listURL = "/api/v1/blocks"
	markURL = "/api/v1/blocks/mark"
)

type BlockManager struct {
	ctx context.Context

	client.Client
	*runtime.Scheme
	cache.Cache

	tenantLabelName   string
	defaultTenantId   string
	storageConfig     string
	storageConfigFile string
	gcInterval        time.Duration
	gcCleanupTimeout  time.Duration
}

func NewBlockManager(ctx context.Context,
	tenantLabelName,
	defaultTenantId,
	storageConfig,
	storageConfigFile string,
	interval, cleanupTimeout time.Duration) *BlockManager {
	cfg, err := kconfig.GetConfig()
	if err != nil {
		klog.Errorf("Failed to get kubeconfig, %s ", err)
		os.Exit(1)
	}

	scheme := runtime.NewScheme()
	_ = v1alpha1.AddToScheme(scheme)

	informerCache, err := cache.New(cfg, cache.Options{
		Scheme: scheme,
	})
	if err != nil {
		klog.Errorf("Failed to create informer cache, %s ", err)
		os.Exit(1)
	}

	mapper, err := func(c *rest.Config) (meta.RESTMapper, error) {
		httpClient, err := rest.HTTPClientFor(c)
		if err != nil {
			return nil, err
		}
		return apiutil.NewDynamicRESTMapper(c, httpClient)
	}(cfg)
	if err != nil {
		klog.Errorf("Failed to create rest mapper, %s ", err)
		os.Exit(1)
	}

	c, err := client.New(cfg, client.Options{Scheme: scheme, Mapper: mapper})
	if err != nil {
		klog.Errorf("Failed to create kubernetes client, %s ", err)
		os.Exit(1)
	}
	/*
		delegatingClient, err := client.NewDelegatingClient(
			client.NewDelegatingClientInput{
				CacheReader: informerCache,
				Client:      c,
			})
		if err != nil {
			klog.Errorf("Failed to create delegating client, %s ", err)
			os.Exit(1)
		}
	*/
	return &BlockManager{
		ctx:               ctx,
		Client:            c,
		Scheme:            scheme,
		Cache:             informerCache,
		storageConfig:     storageConfig,
		storageConfigFile: storageConfigFile,
		gcInterval:        interval,
		gcCleanupTimeout:  cleanupTimeout,
		tenantLabelName:   tenantLabelName,
		defaultTenantId:   defaultTenantId,
	}
}

func (b *BlockManager) Run() error {

	go func() {
		_ = b.Start(b.ctx)
	}()

	if ok := b.WaitForCacheSync(b.ctx); !ok {
		return fmt.Errorf("sync cache failed")
	}

	return b.gc()
}

func (b *BlockManager) gc() error {
	for {
		timer := time.NewTimer(b.gcInterval)
		select {
		case <-b.ctx.Done():
			return nil
		case <-timer.C:
			timer.Stop()
			break
		}

		blocks, err := b.listBlocks()
		if err != nil {
			klog.Errorf("list block failed, %s", err)
			continue
		}

		tenants, err := b.listTenants()
		if err != nil {
			klog.Errorf("list tenant failed, %s", err)
			continue
		}

		for _, block := range blocks {
			tenant := block.Thanos.Labels[b.tenantLabelName]
			if tenant == "" {
				continue
			}

			if !util.Contains(tenants, tenant) {
				if err := b.markBlockForDeletion(block); err != nil {
					klog.Errorf("mark block %s for deleting failed, %s", block.ULID, err)
				}
			}
		}

		b.cleanupBlocks()
	}
}

type ListResponse struct {
	Status string     `json:"status"`
	Data   BlocksInfo `json:"data"`
}

type BlocksInfo struct {
	Label       string    `json:"label"`
	Blocks      []Meta    `json:"blocks"`
	RefreshedAt time.Time `json:"refreshedAt"`
}

type Meta struct {
	ULID   ulid.ULID `json:"ulid"`
	Thanos Thanos    `json:"thanos"`
}

type Thanos struct {
	Labels map[string]string `json:"labels"`
}

func (b *BlockManager) listBlocks() ([]Meta, error) {

	resp, err := http.Get(host + listURL)
	if err != nil {
		return nil, err
	}

	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	defer func() {
		_, _ = io.Copy(ioutil.Discard, resp.Body)
		_ = resp.Body.Close()
	}()

	listRespones := &ListResponse{}
	if err := json.Unmarshal(bs, listRespones); err != nil {
		return nil, err
	}

	return listRespones.Data.Blocks, nil
}

type response struct {
	Status    string      `json:"status"`
	Data      interface{} `json:"data,omitempty"`
	ErrorType string      `json:"errorType,omitempty"`
	Error     string      `json:"error,omitempty"`
	Warnings  []string    `json:"warnings,omitempty"`
}

func (b *BlockManager) markBlockForDeletion(m Meta) error {

	postMessageURL, err := url.Parse(host + markURL)
	if err != nil {
		return err
	}
	postMessageURL.Query()
	values := postMessageURL.Query()
	values.Set("id", m.ULID.String())
	values.Set("action", "DELETION")

	postMessageURL.RawQuery = values.Encode()

	request, err := http.NewRequest(http.MethodPost, postMessageURL.String(), nil)
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json")

	resp, err := (&http.Client{}).Do(request.WithContext(b.ctx))
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		bs, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		defer func() {
			_, _ = io.Copy(ioutil.Discard, resp.Body)
			_ = resp.Body.Close()
		}()

		res := &response{}
		if err := json.Unmarshal(bs, res); err != nil {
			return err
		}

		klog.Errorf("mark block %s failed, %s", m.ULID, res.Error)
	}

	return nil
}

func (b *BlockManager) cleanupBlocks() {

	args := []string{
		"tools",
		"bucket",
		"cleanup",
		"--delete-delay=0",
	}

	if b.storageConfig != "" {
		args = append(args, "--objstore.config="+b.storageConfig)
	} else {
		args = append(args, "--objstore.config-file="+b.storageConfigFile)
	}

	cmd := exec.Command(binPath, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	stopCh := make(chan struct{}, 1)
	go func() {
		if err := cmd.Run(); err != nil {
			klog.Errorf("run block cleanup failed, %s", err)
		}
		stopCh <- struct{}{}
	}()

	timer := time.NewTimer(b.gcCleanupTimeout)
	select {
	case <-timer.C:
		klog.Errorf("block cleanup timeout")
		timer.Stop()
		if cmd != nil {
			if err := cmd.Process.Signal(syscall.SIGTERM); err != nil {
				klog.Errorf("terminate block cleanup faile, %s", err)
			}
		}
		return
	case <-stopCh:
		return
	}
}

func (b *BlockManager) listTenants() ([]string, error) {
	tenantList := &v1alpha1.TenantList{}

	if err := b.Client.List(b.ctx, tenantList); err != nil {
		return nil, err
	}

	var tenants []string
	for _, item := range tenantList.Items {
		if item.DeletionTimestamp != nil && !item.DeletionTimestamp.IsZero() {
			continue
		}

		tenants = append(tenants, item.Name)
	}
	tenants = append(tenants, b.defaultTenantId)

	return tenants, nil
}

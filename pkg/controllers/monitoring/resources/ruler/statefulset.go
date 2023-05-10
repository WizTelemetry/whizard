package ruler

import (
	"fmt"
	"net/url"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	monitoringv1alpha1 "github.com/kubesphere/whizard/pkg/api/monitoring/v1alpha1"
	"github.com/kubesphere/whizard/pkg/constants"
	"github.com/kubesphere/whizard/pkg/controllers/monitoring/resources"
	"github.com/kubesphere/whizard/pkg/controllers/monitoring/resources/query"
	"github.com/kubesphere/whizard/pkg/controllers/monitoring/resources/queryfrontend"
	"github.com/kubesphere/whizard/pkg/controllers/monitoring/resources/router"
	"github.com/kubesphere/whizard/pkg/util"
	promoperator "github.com/prometheus-operator/prometheus-operator/pkg/operator"
	promcommonconfig "github.com/prometheus/common/config"
	promconfig "github.com/prometheus/prometheus/config"
	"github.com/prometheus/prometheus/model/relabel"
	yamlv3 "gopkg.in/yaml.v3"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/klog/v2"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	// repeatableArgs is the args that can be set repeatedly.
	// An error will occur if a non-repeatable arg is set repeatedly.
	repeatableArgs = []string{
		"--query",
		"--query.sd-files",
		"--rule-file",
		"--alertmanagers.url",
		"alert.label-drop",
	}
	// unsupportedArgs is the args that are not allowed to be set by the user.
	unsupportedArgs = []string{
		"--receive.local-endpoint",
		"--http-address",
		"--grpc-address",
	}
)

func (r *Ruler) statefulSets() (retResources []resources.Resource) {
	// for target statefulsets
	var targetNames = make(map[string]struct{}, *r.ruler.Spec.Shards)
	for i := 0; i < int(*r.ruler.Spec.Shards); i++ {
		shardSn := i
		targetNames[r.name(strconv.Itoa(shardSn))] = struct{}{}
		retResources = append(retResources, func() (runtime.Object, resources.Operation, error) {
			return r.statefulSet(shardSn)
		})
	}

	var stsList appsv1.StatefulSetList
	err := r.Client.List(r.Context, &stsList, client.InNamespace(r.ruler.Namespace))
	if err != nil {
		return errResourcesFunc(err)
	}
	// check statefulsets to be deleted.
	// the statefulsets owned by the ruler have a same name prefix and a shard sequence number suffix
	var namePrefix = r.name() + "-"
	for i := range stsList.Items {
		sts := stsList.Items[i]
		if !strings.HasPrefix(sts.Name, namePrefix) {
			continue
		}
		sn := strings.TrimPrefix(sts.Name, namePrefix)
		if sequenceNumberRegexp.MatchString(sn) {
			if _, ok := targetNames[sts.Name]; !ok {
				retResources = append(retResources, func() (runtime.Object, resources.Operation, error) {
					return &sts, resources.OperationDelete, nil
				})
			}
		}
	}
	return
}

func (r *Ruler) statefulSet(shardSn int) (runtime.Object, resources.Operation, error) {
	var sts = &appsv1.StatefulSet{ObjectMeta: r.meta(r.name(strconv.Itoa(shardSn)))}

	ls := r.labels()
	ls[constants.LabelNameRulerShardSn] = strconv.Itoa(shardSn)

	sts.Spec = appsv1.StatefulSetSpec{
		Replicas: r.ruler.Spec.Replicas,
		Selector: &metav1.LabelSelector{
			MatchLabels: ls,
		},
		ServiceName: r.name(strconv.Itoa(shardSn), constants.ServiceNameSuffix),
		Template: corev1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Labels: ls,
			},
			Spec: corev1.PodSpec{
				NodeSelector: r.ruler.Spec.NodeSelector,
				Tolerations:  r.ruler.Spec.Tolerations,
				Affinity:     r.ruler.Spec.Affinity,
			},
		},
	}

	var container = corev1.Container{
		Name:      "rule",
		Image:     r.ruler.Spec.Image,
		Args:      []string{"rule"},
		Resources: r.ruler.Spec.Resources,
		Ports: []corev1.ContainerPort{
			{
				Protocol:      corev1.ProtocolTCP,
				Name:          constants.GRPCPortName,
				ContainerPort: constants.GRPCPort,
			},
			{
				Protocol:      corev1.ProtocolTCP,
				Name:          constants.HTTPPortName,
				ContainerPort: constants.HTTPPort,
			},
		},
		LivenessProbe:  r.DefaultLivenessProbe(),
		ReadinessProbe: r.DefaultReadinessProbe(),
		Env: []corev1.EnvVar{
			{
				Name: "POD_NAME",
				ValueFrom: &corev1.EnvVarSource{
					FieldRef: &corev1.ObjectFieldSelector{
						FieldPath: "metadata.name",
					},
				},
			},
		},
	}

	var watchedDirectories []string
	var configReloaderVolumeMounts []corev1.VolumeMount

	for cmName := range r.shardsRuleConfigMapNames[shardSn] {
		vol := corev1.Volume{
			Name: "configmap-" + cmName,
			VolumeSource: corev1.VolumeSource{
				ConfigMap: &corev1.ConfigMapVolumeSource{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: cmName,
					},
				},
			},
		}
		sts.Spec.Template.Spec.Volumes = append(sts.Spec.Template.Spec.Volumes, vol)
		volMount := corev1.VolumeMount{
			Name:      vol.Name,
			MountPath: filepath.Join(rulesDir, cmName),
		}
		container.VolumeMounts = append(container.VolumeMounts, volMount)
		configReloaderVolumeMounts = append(configReloaderVolumeMounts, volMount)
		watchedDirectories = append(watchedDirectories, volMount.MountPath)
	}

	var tsdbVolume = &corev1.Volume{
		Name: "tsdb",
		VolumeSource: corev1.VolumeSource{
			EmptyDir: &corev1.EmptyDirVolumeSource{},
		},
	}
	if v := r.ruler.Spec.DataVolume; v != nil {
		if pvc := v.PersistentVolumeClaim; pvc != nil {
			if pvc.Name == "" {
				pvc.Name = sts.Name + "-tsdb"
			}
			if pvc.Spec.AccessModes == nil {
				pvc.Spec.AccessModes = []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce}
			}
			sts.Spec.VolumeClaimTemplates = append(sts.Spec.VolumeClaimTemplates, *pvc)
			container.VolumeMounts = append(container.VolumeMounts, corev1.VolumeMount{
				Name:      pvc.Name,
				MountPath: storageDir,
			})
			tsdbVolume = nil
		} else if v.EmptyDir != nil {
			tsdbVolume.EmptyDir = v.EmptyDir
		}
	}
	if tsdbVolume != nil {
		sts.Spec.Template.Spec.Volumes = append(sts.Spec.Template.Spec.Volumes, *tsdbVolume)
		container.VolumeMounts = append(container.VolumeMounts, corev1.VolumeMount{
			Name:      tsdbVolume.Name,
			MountPath: storageDir,
		})
	}

	if r.ruler.Spec.AlertmanagersConfig != nil {
		container.Args = append(container.Args, "--alertmanagers.config=$(ALERTMANAGERS_CONFIG)")
		container.Env = append(container.Env, corev1.EnvVar{
			Name: "ALERTMANAGERS_CONFIG",
			ValueFrom: &corev1.EnvVarSource{
				SecretKeyRef: r.ruler.Spec.AlertmanagersConfig,
			},
		})
	} else if len(r.ruler.Spec.AlertmanagersURL) > 0 {
		for _, url := range r.ruler.Spec.AlertmanagersURL {
			container.Args = append(container.Args, fmt.Sprintf("--alertmanagers.url=%s", url))
		}
	}

	if r.ruler.Spec.LogLevel != "" {
		container.Args = append(container.Args, "--log.level="+r.ruler.Spec.LogLevel)
	}
	if r.ruler.Spec.LogFormat != "" {
		container.Args = append(container.Args, "--log.format="+r.ruler.Spec.LogFormat)
	}
	container.Args = append(container.Args, fmt.Sprintf(`--label=%s="$(POD_NAME)"`, constants.RulerReplicaLabelName))
	for k, v := range r.ruler.Spec.Labels {
		container.Args = append(container.Args, fmt.Sprintf(`--label=%s="%s"`, k, v))
	}
	container.Args = append(container.Args, fmt.Sprintf("--data-dir=%s", storageDir))
	container.Args = append(container.Args, fmt.Sprintf("--rule-file=%s/*/*.yaml", rulesDir))
	container.Args = append(container.Args, fmt.Sprintf("--alert.label-drop=%s", constants.RulerReplicaLabelName))
	for _, lb := range r.ruler.Spec.AlertDropLabels {
		container.Args = append(container.Args, fmt.Sprintf("--alert.label-drop=%s", lb))
	}
	if r.ruler.Spec.EvaluationInterval != "" {
		container.Args = append(container.Args, fmt.Sprintf("--eval-interval=%s", r.ruler.Spec.EvaluationInterval))
	}

	namespacedName := util.ServiceNamespacedName(r.ruler)

	if namespacedName != nil {
		var service monitoringv1alpha1.Service
		if err := r.Client.Get(r.Context, *namespacedName, &service); err != nil {
			return nil, resources.OperationCreateOrUpdate, err
		}

		writeAddr, err := r.remoteWriteAddress()
		if err != nil {
			return nil, "", err
		}
		queryAddr, err := r.queryAddress()
		if err != nil {
			return nil, "", err
		}
		var rwCfg = &promconfig.RemoteWriteConfig{}
		*rwCfg = promconfig.DefaultRemoteWriteConfig
		var cfgs struct {
			RemoteWriteConfigs []*promconfig.RemoteWriteConfig `yaml:"remote_write,omitempty"`
		}

		// proxy config
		// if the tenant exists, append QueryProxy
		// otherwise, append WriteProxy
		if r.ruler.Spec.Tenant == "" {
			container.Args = append(container.Args, "--query="+queryAddr)
			writeUrl, err := url.Parse("http://127.0.0.1:8081/push")
			if err != nil {
				return nil, resources.OperationCreateOrUpdate, err
			}
			rwCfg.URL = &promcommonconfig.URL{URL: writeUrl}

			cfgs.RemoteWriteConfigs = append(cfgs.RemoteWriteConfigs, rwCfg)
			content, err := yamlv3.Marshal(&cfgs)
			if err != nil {
				return nil, resources.OperationCreateOrUpdate, err
			}
			container.Args = append(container.Args, "--remote-write.config="+string(content))

			writeProxyContainer, err := r.addWriteProxyContainer(&service.Spec, writeAddr)
			if err != nil {
				return nil, "", err
			}
			sts.Spec.Template.Spec.Containers = append(sts.Spec.Template.Spec.Containers, *writeProxyContainer)
		} else {
			container.Args = append(container.Args,
				fmt.Sprintf("--query=http://127.0.0.1:9080/%s", r.ruler.Spec.Tenant))

			writeUrl, err := url.Parse(writeAddr + "/api/v1/receive")
			if err != nil {
				return nil, resources.OperationCreateOrUpdate, err
			}
			if rwCfg.Headers == nil {
				rwCfg.Headers = make(map[string]string)
			}
			rwCfg.Headers[service.Spec.TenantHeader] = r.ruler.Spec.Tenant
			rwCfg.URL = &promcommonconfig.URL{URL: writeUrl}

			rwCfg.WriteRelabelConfigs = append(rwCfg.WriteRelabelConfigs, &relabel.Config{
				Regex:  relabel.MustNewRegexp(service.Spec.TenantLabelName),
				Action: relabel.LabelDrop,
			})
			cfgs.RemoteWriteConfigs = append(cfgs.RemoteWriteConfigs, rwCfg)
			content, err := yamlv3.Marshal(&cfgs)
			if err != nil {
				return nil, resources.OperationCreateOrUpdate, err
			}
			container.Args = append(container.Args, "--remote-write.config="+string(content))

			queryProxyContainer, _ := r.addQueryProxyContainer(&service.Spec, queryAddr)
			sts.Spec.Template.Spec.Containers = append(sts.Spec.Template.Spec.Containers, *queryProxyContainer)
		}
	}

	for _, flag := range r.ruler.Spec.Flags {
		arg := util.GetArgName(flag)
		if util.Contains(unsupportedArgs, arg) {
			klog.V(3).Infof("ignore the unsupported flag %s", arg)
			continue
		}

		if util.Contains(repeatableArgs, arg) {
			container.Args = append(container.Args, flag)
			continue
		}

		replaced := util.ReplaceInSlice(container.Args, func(v interface{}) bool {
			return util.GetArgName(v.(string)) == util.GetArgName(flag)
		}, flag)
		if !replaced {
			container.Args = append(container.Args, flag)
		}
	}

	sort.Strings(container.Args[1:])

	var reloaderConfig = promoperator.ContainerConfig{
		Image:         r.Options.PrometheusConfigReloader.Image,
		CPURequest:    r.Options.PrometheusConfigReloader.Resources.Requests.Cpu().String(),
		MemoryRequest: r.Options.PrometheusConfigReloader.Resources.Requests.Memory().String(),
		CPULimit:      r.Options.PrometheusConfigReloader.Resources.Limits.Cpu().String(),
		MemoryLimit:   r.Options.PrometheusConfigReloader.Resources.Limits.Memory().String(),
	}

	var reloadContainer = promoperator.CreateConfigReloader(
		"config-reloader",
		promoperator.ReloaderConfig(reloaderConfig),
		promoperator.ReloaderURL(url.URL{
			Scheme: "http",
			Host:   fmt.Sprintf("localhost:%d", constants.HTTPPort),
			Path:   path.Clean("/-/reload"),
		}),
		promoperator.ListenLocal(true),
		promoperator.LocalHost("localhost"),
		promoperator.LogFormat(r.ruler.Spec.LogFormat),
		promoperator.LogLevel(r.ruler.Spec.LogLevel),
		promoperator.WatchedDirectories(watchedDirectories),
		promoperator.VolumeMounts(configReloaderVolumeMounts),
		promoperator.Shard(-1),
	)

	sts.Spec.Template.Spec.Containers = append(sts.Spec.Template.Spec.Containers, container, reloadContainer)

	if r.ruler.Spec.ImagePullSecrets != nil && len(r.ruler.Spec.ImagePullSecrets) > 0 {
		sts.Spec.Template.Spec.ImagePullSecrets = r.ruler.Spec.ImagePullSecrets
	}

	return sts, resources.OperationCreateOrUpdate, ctrl.SetControllerReference(r.ruler, sts, r.Scheme)
}

func (r *Ruler) remoteWriteAddress() (string, error) {
	routerList := &monitoringv1alpha1.RouterList{}
	if err := r.Client.List(r.Context, routerList, client.MatchingLabels(util.ManagedLabelBySameService(r.ruler))); err != nil {
		return "", err
	}

	if len(routerList.Items) > 0 {
		if len(routerList.Items) > 1 {
			return "", fmt.Errorf("more than one router defined for service %s/%s", r.Service.Name, r.Service.Namespace)
		}

		o := routerList.Items[0]
		r, err := router.New(r.BaseReconciler, &o, nil)
		if err != nil {
			return "", err
		}

		return r.RemoteWriteAddr(), nil
	}

	return "", fmt.Errorf("no router defined for service %s/%s", r.Service.Name, r.Service.Namespace)
}

// queryAddress returns the address from which the ruler should query.
// The ruler should query from the QueryFrontend with a higher performance the Query
// because the feature https://thanos.io/v0.31/proposals-accepted/202205-vertical-query-sharding/
func (r *Ruler) queryAddress() (string, error) {
	queryFrontendList := &monitoringv1alpha1.QueryFrontendList{}
	if err := r.Client.List(r.Context, queryFrontendList, client.MatchingLabels(util.ManagedLabelBySameService(r.ruler))); err != nil {
		return "", err
	}

	if len(queryFrontendList.Items) > 0 {
		if len(queryFrontendList.Items) > 1 {
			return "", fmt.Errorf("more than one query frontend defined for service %s/%s", r.Service.Name, r.Service.Namespace)
		}

		q := queryFrontendList.Items[0]
		r, err := queryfrontend.New(r.BaseReconciler, &q)
		if err != nil {
			return "", err
		}

		return r.HttpAddr(), nil
	}

	queryList := &monitoringv1alpha1.QueryList{}
	if err := r.Client.List(r.Context, queryList, client.MatchingLabels(util.ManagedLabelBySameService(r.ruler))); err != nil {
		return "", err
	}

	if len(queryList.Items) > 0 {
		if len(queryList.Items) > 1 {
			return "", fmt.Errorf("more than one query defined for service %s/%s", r.Service.Name, r.Service.Namespace)
		}

		o := queryList.Items[0]
		r, err := query.New(r.BaseReconciler, &o, nil)
		if err != nil {
			return "", err
		}

		return r.HttpAddr(), nil
	}

	return "", fmt.Errorf("no query frontend or query exist for service %s/%s", r.Service.Name, r.Service.Namespace)
}

func (r *Ruler) addQueryProxyContainer(serviceSpec *monitoringv1alpha1.ServiceSpec, queryAddr string) (*corev1.Container, error) {

	var queryProxyContainer *corev1.Container

	queryProxyContainer = &corev1.Container{
		Name:  "query-proxy",
		Image: r.Options.RulerQueryProxy.Image,
		Args: []string{
			"--http-address=127.0.0.1:9080",
		},
		Resources: r.Options.RulerQueryProxy.Resources,
	}
	queryProxyContainer.Args = append(queryProxyContainer.Args, "--tenant.label-name="+serviceSpec.TenantLabelName)
	queryProxyContainer.Args = append(queryProxyContainer.Args, "--query.address="+queryAddr)
	return queryProxyContainer, nil
}

// cortex-tenant config
// https://github.com/blind-oracle/cortex-tenant/blob/main/config.go#L13
type config struct {
	Listen      string
	ListenPprof string `yaml:"listen_pprof"`

	Target string

	LogLevel        string `yaml:"log_level"`
	Timeout         time.Duration
	TimeoutShutdown time.Duration `yaml:"timeout_shutdown"`
	Concurrency     int
	Metadata        bool

	Auth struct {
		Egress struct {
			Username string
			Password string
		}
	}

	Tenant struct {
		Label       string
		LabelRemove bool `yaml:"label_remove"`
		Header      string
		Default     string
		AcceptAll   bool `yaml:"accept_all"`
	}
}

func (r *Ruler) addWriteProxyContainer(serviceSpec *monitoringv1alpha1.ServiceSpec, writeAddr string) (*corev1.Container, error) {
	var writeProxyContainer *corev1.Container
	cfg := &config{
		Listen:          "127.0.0.1:8081",
		LogLevel:        "warn",
		Timeout:         time.Second * 10,
		TimeoutShutdown: time.Second * 10,
		Concurrency:     1000,
		Metadata:        false,
	}

	writeUrl, err := url.Parse(writeAddr + "/api/v1/receive")
	if err != nil {
		return writeProxyContainer, err
	}
	cfg.Target = writeUrl.String()

	cfg.Tenant.Label = serviceSpec.TenantLabelName
	cfg.Tenant.LabelRemove = true
	cfg.Tenant.Header = serviceSpec.TenantHeader
	cfg.Tenant.Default = serviceSpec.DefaultTenantId

	cfgContent, err := yamlv3.Marshal(cfg)
	if err != nil {
		return writeProxyContainer, err
	}

	writeProxyContainer = &corev1.Container{
		Name:  "write-proxy",
		Image: r.Options.RulerWriteProxy.Image,
		Args: []string{
			"--config-content=" + string(cfgContent),
		},
		Resources: r.Options.RulerWriteProxy.Resources,
	}
	return writeProxyContainer, nil
}

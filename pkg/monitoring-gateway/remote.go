package monitoringgateway

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/pkg/errors"
	config_util "github.com/prometheus/common/config"
	"gopkg.in/yaml.v2"
)

type remoteWriteClient struct {
	Client *http.Client
	url    *config_util.URL

	timeout time.Duration
}

// LoadExternalRemoteWriteConfig loads remotewrites config, and prefers file to content
func LoadExternalRemoteWriteConfig(file, content string) ([]ExternalRemoteWriteConfig, error) {
	var buff []byte
	if file != "" {
		c, err := os.ReadFile(file)
		if err != nil {
			return nil, err
		}
		buff = c
	} else {
		buff = []byte(content)
	}
	if len(buff) == 0 {
		return nil, nil
	}
	var rws []ExternalRemoteWriteConfig
	if err := yaml.UnmarshalStrict(buff, &rws); err != nil {
		return nil, errors.Wrap(err, "")
	}
	return rws, nil
}

func NewExternalRemoteWriteClients(rwsCfg []ExternalRemoteWriteConfig) ([]*remoteWriteClient, error) {

	var clients []*remoteWriteClient
	for _, rwCfg := range rwsCfg {
		writeClient, err := newExternalRemoteWriteClient(&rwCfg)
		if err != nil {
			return nil, err
		}
		if writeClient != nil {
			clients = append(clients, writeClient)
		}
	}
	return clients, nil
}

func newExternalRemoteWriteClient(conf *ExternalRemoteWriteConfig) (*remoteWriteClient, error) {
	cfg := config_util.HTTPClientConfig{
		TLSConfig:   conf.TLSConfig,
		BearerToken: config_util.Secret(conf.BearerToken),
	}
	if conf.BasicAuth != nil {
		cfg.BasicAuth = &config_util.BasicAuth{
			Username:     conf.BasicAuth.Username,
			Password:     config_util.Secret(conf.BasicAuth.Password),
			PasswordFile: conf.BasicAuth.PasswordFile,
		}
	}
	httpClient, err := config_util.NewClientFromConfig(cfg, "remote_storage_write_client")
	if err != nil {
		return nil, err
	}
	t := httpClient.Transport

	if len(conf.Headers) > 0 {
		t = newInjectHeadersRoundTripper(conf.Headers, t)
	}
	httpClient.Transport = t
	timeout := time.Second * 30
	if conf.RemoteTimeout > 0 {
		timeout = time.Duration(conf.RemoteTimeout)
	}
	return &remoteWriteClient{
		Client:  httpClient,
		url:     conf.URL,
		timeout: timeout,
	}, nil
}

func (c *remoteWriteClient) Send(ctx context.Context, body []byte, header http.Header) result {
	httpReq, err := http.NewRequest("POST", c.url.String(), bytes.NewReader(body))
	if err != nil {
		return result{code: http.StatusBadRequest, err: err}
	}
	for k, vs := range header {
		for _, v := range vs {
			httpReq.Header.Add(k, v)
		}
	}

	httpReq.Header.Add("Content-Encoding", "snappy")
	httpReq.Header.Set("Content-Type", "application/x-protobuf")
	httpReq.Header.Set("X-Prometheus-Remote-Write-Version", "0.1.0")
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	httpResp, err := c.Client.Do(httpReq.WithContext(ctx))
	if err != nil {
		return result{code: http.StatusBadGateway, err: err}
	}
	defer func() {
		io.Copy(io.Discard, httpResp.Body)
		httpResp.Body.Close()
	}()

	if httpResp.StatusCode/100 != 2 {
		scanner := bufio.NewScanner(io.LimitReader(httpResp.Body, 1024))
		line := ""
		if scanner.Scan() {
			line = scanner.Text()
		}
		err = fmt.Errorf("server returned HTTP status %s: %s", httpResp.Status, line)
	}
	return result{code: httpResp.StatusCode, err: err}
}

func (c remoteWriteClient) Endpoint() string {
	return c.url.String()
}

type result struct {
	code int
	err  error
}

func newInjectHeadersRoundTripper(h map[string]string, underlyingRT http.RoundTripper) *injectHeadersRoundTripper {
	return &injectHeadersRoundTripper{headers: h, RoundTripper: underlyingRT}
}

type injectHeadersRoundTripper struct {
	headers map[string]string
	http.RoundTripper
}

func (t *injectHeadersRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	for key, value := range t.headers {
		req.Header.Set(key, value)
	}
	return t.RoundTripper.RoundTrip(req)
}

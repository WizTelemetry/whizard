package monitoringgateway

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	config_util "github.com/prometheus/common/config"
)

type remoteWriteClient struct {
	Client *http.Client
	url    *config_util.URL

	timeout time.Duration
}

func newRemoteWriteClient(conf *RemoteWriteConfig) (*remoteWriteClient, error) {
	httpClientConfig := config_util.HTTPClientConfig{
		TLSConfig: conf.TLSConfig,
	}
	httpClient, err := config_util.NewClientFromConfig(httpClientConfig, "remote_storage_write_client")
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

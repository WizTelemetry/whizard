package monitoringgateway

import (
	"crypto/tls"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"sync"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/pkg/errors"
	"github.com/prometheus-community/prom-label-proxy/injectproxy"
	"github.com/prometheus/common/route"
	"github.com/prometheus/prometheus/model/labels"
)

const (
	apiPrefix = "/:tenant_id/api/v1"

	epQuery       = apiPrefix + "/query"
	epQueryRange  = apiPrefix + "/query_range"
	epSeries      = apiPrefix + "/series"
	epLabels      = apiPrefix + "/labels"
	epLabelValues = apiPrefix + "/label/*path"
	epReceive     = apiPrefix + "/receive"
	epRules       = apiPrefix + "/rules"
	epAlerts      = apiPrefix + "/alerts"
)

type Options struct {
	ListenAddress string
	TLSConfig     *tls.Config

	TenantHeader    string
	TenantLabelName string

	RemoteWriteHandler http.Handler
	QueryProxy         *httputil.ReverseProxy
	RulesQueryProxy    *httputil.ReverseProxy

	CertAuthenticator *CertAuthenticator
}

type Handler struct {
	logger  log.Logger
	options *Options
	router  *route.Router

	remoteWriteHander http.Handler
	queryProxy        *httputil.ReverseProxy
	rulesQueryProxy   *httputil.ReverseProxy
}

func NewHandler(logger log.Logger, o *Options) *Handler {
	if logger == nil {
		logger = log.NewNopLogger()
	}

	h := &Handler{
		logger:            logger,
		options:           o,
		router:            route.New(),
		remoteWriteHander: o.RemoteWriteHandler,
		queryProxy:        o.QueryProxy,
		rulesQueryProxy:   o.RulesQueryProxy,
	}

	h.router.Get(epQuery, h.wrap(h.query))
	h.router.Post(epQuery, h.wrap(h.query))
	h.router.Get(epQueryRange, h.wrap(h.query))
	h.router.Post(epQueryRange, h.wrap(h.query))
	h.router.Get(epSeries, h.wrap(h.matcher(matchersParam)))
	h.router.Get(epLabels, h.wrap(h.matcher(matchersParam)))
	h.router.Get(epLabelValues, h.wrap(h.matcher(matchersParam)))
	h.router.Get(epRules, h.wrap(h.matcher(matchersParam)))
	// do provide /api/v1/alerts because thanos does not support alerts filtering as of v0.28.0
	// please filtering alerts by /api/v1/rules
	// h.router.Get(epAlerts, h.wrap(h.matcher(matchersParam)))

	h.router.Post(epReceive, h.wrap(h.remoteWrite))

	return h
}

func (h *Handler) wrap(f http.HandlerFunc) http.HandlerFunc {
	if h.options.CertAuthenticator != nil {
		f = withAuthorization(f, h.options.CertAuthenticator)
	}
	return withRequestInfo(f)
}

func (h *Handler) query(w http.ResponseWriter, req *http.Request) {
	if h.queryProxy == nil {
		http.Error(w, "The query target is not configured for the server", http.StatusNotAcceptable)
		return
	}

	var (
		query    = req.URL.Query()
		postForm url.Values
	)

	if req.Method == http.MethodPost {
		if err := req.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		postForm = req.PostForm
	}

	if query.Get(queryParam) == "" && postForm.Get(queryParam) == "" {
		return
	}

	ctx := req.Context()
	requestInfo, found := requestInfoFrom(ctx)

	if !found || requestInfo.TenantId == "" {
		http.NotFound(w, req)
		return
	}

	enforcer := injectproxy.NewEnforcer(true, &labels.Matcher{
		Type:  labels.MatchEqual,
		Name:  h.options.TenantLabelName,
		Value: requestInfo.TenantId,
	})

	q, found, err := enforceQueryValues(enforcer, query)
	if err != nil {
		switch err.(type) {
		case injectproxy.IllegalLabelMatcherError:
			http.Error(w, err.Error(), http.StatusBadRequest)
		case queryParseError:
			http.Error(w, err.Error(), http.StatusBadRequest)
		case enforceLabelError:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	if found {
		req.URL.RawQuery = q
	}

	if postForm != nil {
		q, found, err := enforceQueryValues(enforcer, postForm)
		if err != nil {
			switch err.(type) {
			case injectproxy.IllegalLabelMatcherError:
				http.Error(w, err.Error(), http.StatusBadRequest)
			case queryParseError:
				http.Error(w, err.Error(), http.StatusBadRequest)
			case enforceLabelError:
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
		if found {
			_ = req.Body.Close()
			req.Body = ioutil.NopCloser(strings.NewReader(q))
			req.ContentLength = int64(len(q))
		}
	}

	h.queryProxy.ServeHTTP(w, req)
}

func (h *Handler) matcher(matchersParam string) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if h.queryProxy == nil {
			http.Error(w, "The query target is not configured for the server", http.StatusNotAcceptable)
			return
		}

		ctx := req.Context()
		requestInfo, found := requestInfoFrom(ctx)

		if !found || requestInfo.TenantId == "" {
			http.NotFound(w, req)
			return
		}

		if requestInfo.TenantId == "" {
			h.queryProxy.ServeHTTP(w, req)
			return
		}

		matcher := &labels.Matcher{
			Type:  labels.MatchEqual,
			Name:  h.options.TenantLabelName,
			Value: requestInfo.TenantId,
		}
		q := req.URL.Query()

		if err := injectMatcher(q, matcher, matchersParam); err != nil {
			return
		}
		req.URL.RawQuery = q.Encode()
		if req.Method == http.MethodPost {
			if err := req.ParseForm(); err != nil {
				return
			}
			q = req.PostForm
			if err := injectMatcher(q, matcher, matchersParam); err != nil {
				return
			}
			_ = req.Body.Close()
			req.Body = ioutil.NopCloser(strings.NewReader(q.Encode()))
			req.ContentLength = int64(len(q))
		}

		if (strings.HasSuffix(req.URL.Path, "/rules") || strings.HasSuffix(req.URL.Path, "/alerts")) &&
			h.rulesQueryProxy != nil {
			h.rulesQueryProxy.ServeHTTP(w, req)
			return
		}
		h.queryProxy.ServeHTTP(w, req)
	}
}

func (h *Handler) remoteWrite(w http.ResponseWriter, req *http.Request) {
	if h.remoteWriteHander == nil {
		http.Error(w, "There is no remote write targets configured for the server", http.StatusNotAcceptable)
		return
	}
	ctx := req.Context()
	requestInfo, found := requestInfoFrom(ctx)

	if !found || requestInfo.TenantId == "" {
		http.NotFound(w, req)
		return
	}

	req.Header.Set(h.options.TenantHeader, requestInfo.TenantId)

	h.remoteWriteHander.ServeHTTP(w, req)
}

func (h *Handler) Run() error {
	srv := &http.Server{
		Handler:   h.router,
		Addr:      h.options.ListenAddress,
		TLSConfig: h.options.TLSConfig,
	}

	if h.options.TLSConfig != nil {
		level.Info(h.logger).Log("msg", "Serving HTTPS", "address", h.options.ListenAddress)
		return srv.ListenAndServeTLS("", "")
	}

	level.Info(h.logger).Log("msg", "Serving plain HTTP", "address", h.options.ListenAddress)
	return srv.ListenAndServe()
}

func NewSingleHostReverseProxy(target *url.URL, tlsConfig *tls.Config) *httputil.ReverseProxy {
	proxy := httputil.NewSingleHostReverseProxy(target)
	oldDirector := proxy.Director
	if target.Scheme == "https" {
		proxy.Transport = &http.Transport{
			TLSClientConfig: tlsConfig,
		}
	}
	proxy.Director = func(req *http.Request) {
		req.URL.Scheme = target.Scheme
		req.Host = target.Host
		// remove the prefix /:tenant_id from path
		if index := indexByteNth(req.URL.Path, '/', 2); index > 0 {
			req.URL.Path = req.URL.Path[index:]
		}

		oldDirector(req)
	}
	return proxy
}

type remoteWriteHandler struct {
	writeClients []*remoteWriteClient
	tenantHeader string
}

func (h remoteWriteHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if len(h.writeClients) == 0 {
		return
	}

	ctx := req.Context()

	// Forward the request to multiple targets in parallel.
	// If either forwarding fails, the errors are responded. This may result in repeated sending same data to one target.
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var results = make([]result, len(h.writeClients))
	var i = 0
	var wg sync.WaitGroup

	var tenantHeader = make(http.Header)
	if tenantId := req.Header.Get(h.tenantHeader); tenantId != "" {
		tenantHeader.Set(h.tenantHeader, tenantId)
	}
	for _, writeClient := range h.writeClients {
		wg.Add(1)
		ep := writeClient.Endpoint()

		go func(idx int, writeClient *remoteWriteClient) {
			defer wg.Done()
			result := writeClient.Send(ctx, body, tenantHeader)
			if result.err != nil {
				result.err = errors.Wrapf(result.err, "forwarding request to endpoint %v", ep)
			}
			results[idx] = result
		}(i, writeClient)
		i++
	}
	wg.Wait()

	var code int
	for _, result := range results {
		if result.code > code {
			code = result.code
			err = result.err
		}
	}
	if code <= 0 {
		code = http.StatusNoContent
	}
	if err != nil {
		http.Error(w, err.Error(), code)
	} else {
		w.WriteHeader(code)
	}
}

func NewRemoteWriteHandler(rwsCfg []RemoteWriteConfig, tenantHeader string) (http.Handler, error) {

	if len(rwsCfg) > 0 {
		var handler = remoteWriteHandler{tenantHeader: tenantHeader}
		for _, rwCfg := range rwsCfg {
			writeClient, err := newRemoteWriteClient(&rwCfg)
			if err != nil {
				return nil, err
			}
			if writeClient != nil {
				handler.writeClients = append(handler.writeClients, writeClient)
			}
		}
		return &handler, nil
	}

	return nil, nil
}

// indexByteNth returns the index of the nth instance of c in s, or -1 if the nth c is not present in s.
func indexByteNth(s string, c byte, nth int) int {
	num := 0
	for i, c := range s {
		if c == '/' {
			num++
			if num == nth {
				return i
			}
		}
	}
	return -1
}

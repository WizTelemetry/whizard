package monitoringgateway

import (
	"crypto/tls"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
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

	RemoteWriteProxy *httputil.ReverseProxy
	QueryProxy       *httputil.ReverseProxy

	CertAuthenticator *CertAuthenticator
}

type Handler struct {
	logger  log.Logger
	options *Options
	router  *route.Router

	remoteWriteProxy *httputil.ReverseProxy
	queryProxy       *httputil.ReverseProxy
}

func NewHandler(logger log.Logger, o *Options) *Handler {
	if logger == nil {
		logger = log.NewNopLogger()
	}

	h := &Handler{
		logger:           logger,
		options:          o,
		router:           route.New(),
		remoteWriteProxy: o.RemoteWriteProxy,
		queryProxy:       o.QueryProxy,
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
		h.queryProxy.ServeHTTP(w, req)
	}
}

func (h *Handler) remoteWrite(w http.ResponseWriter, req *http.Request) {
	if h.remoteWriteProxy == nil {
		http.Error(w, "The remote write target is not configured for the server", http.StatusNotAcceptable)
		return
	}

	ctx := req.Context()
	requestInfo, found := requestInfoFrom(ctx)

	if !found || requestInfo.TenantId == "" {
		http.NotFound(w, req)
		return
	}

	req.Header.Set(h.options.TenantHeader, requestInfo.TenantId)

	h.remoteWriteProxy.ServeHTTP(w, req)
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

func NewSingleHostReverseProxy(target *url.URL) *httputil.ReverseProxy {
	proxy := httputil.NewSingleHostReverseProxy(target)
	oldDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		// remove the prefix /:tenant_id from path
		if index := indexByteNth(req.URL.Path, '/', 2); index > 0 {
			req.URL.Path = req.URL.Path[index:]
		}

		oldDirector(req)
	}
	return proxy
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

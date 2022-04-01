package gateway

import (
	"crypto/tls"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"path"
	"regexp"
	"strings"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus-community/prom-label-proxy/injectproxy"
	"github.com/prometheus/common/route"
	"github.com/prometheus/prometheus/model/labels"
)

var apiPrefix = "/namespaces/:namespace/agents/:name"
var apiPrefixRegexp = regexp.MustCompile("/namespaces/[^/]+/agents/[^/]+")

type Options struct {
	ListenAddress string
	TLSConfig     *tls.Config

	TenantHeader    string
	TenantLabelName string

	RemoteWriteProxy *httputil.ReverseProxy
	QueryProxy       *httputil.ReverseProxy

	GetAgentFunc      GetAgentFunc
	CertAuthenticator *CertAuthenticator
}

type Handler struct {
	logger  log.Logger
	options *Options
	router  *route.Router

	remoteWriteProxy *httputil.ReverseProxy
	queryProxy       *httputil.ReverseProxy

	getAgentFunc GetAgentFunc
}

func NewHandler(logger log.Logger, o *Options) *Handler {
	if logger == nil {
		logger = log.NewNopLogger()
	}

	h := &Handler{
		logger:           logger,
		options:          o,
		router:           route.New().WithPrefix(apiPrefix),
		remoteWriteProxy: o.RemoteWriteProxy,
		queryProxy:       o.QueryProxy,
		getAgentFunc:     o.GetAgentFunc,
	}

	h.router.Get("/query", h.wrap(h.query))
	h.router.Post("/query", h.wrap(h.query))
	h.router.Get("/query_range", h.wrap(h.query))
	h.router.Post("/query_range", h.wrap(h.query))
	h.router.Post("/write", h.wrap(h.remoteWrite))
	h.router.Get("/series", h.wrap(h.matcher(matchersParam)))
	h.router.Get("/labels", h.wrap(h.matcher(matchersParam)))
	h.router.Get("/label/*path", h.wrap(h.matcher(matchersParam)))
	h.router.Get("/targets/metadata", h.wrap(h.matcher(targetMatchersParam)))

	return h
}

func (h *Handler) wrap(f http.HandlerFunc) http.HandlerFunc {
	if h.options.CertAuthenticator != nil {
		f = withAuthorization(f, h.options.CertAuthenticator)
	}
	return withRequestInfo(f, h.getAgentFunc)
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

	if !found || requestInfo.Agent == nil {
		http.NotFound(w, req)
		return
	}

	if requestInfo.Agent.Spec.Tenant == "" {
		h.queryProxy.ServeHTTP(w, req)
		return
	}

	enforcer := injectproxy.NewEnforcer(true, &labels.Matcher{
		Type:  labels.MatchEqual,
		Name:  h.options.TenantLabelName,
		Value: requestInfo.Agent.Spec.Tenant,
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

		if !found || requestInfo.Agent == nil {
			http.NotFound(w, req)
			return
		}

		if requestInfo.Agent.Spec.Tenant == "" {
			h.queryProxy.ServeHTTP(w, req)
			return
		}

		matcher := &labels.Matcher{
			Type:  labels.MatchEqual,
			Name:  h.options.TenantLabelName,
			Value: requestInfo.Agent.Spec.Tenant,
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

	if !found || requestInfo.Agent == nil {
		http.NotFound(w, req)
		return
	}

	tenant := requestInfo.Agent.Spec.Tenant

	if tenant == "" {
		h.remoteWriteProxy.ServeHTTP(w, req)
		return
	}

	req.Header.Set(h.options.TenantHeader, tenant)

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

func NewDirector(target *url.URL) func(*http.Request) {
	return func(req *http.Request) {
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.URL.Path = path.Join(target.Path, apiPrefixRegexp.ReplaceAllString(req.URL.Path, "/api/v1"))
		req.Host = ""
	}
}

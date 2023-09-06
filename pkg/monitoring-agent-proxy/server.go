package monitoringagentproxy

import (
	"crypto/tls"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/common/route"
)

const (
	query       = "/api/v1/query"
	queryRange  = "/api/v1/query_range"
	series      = "/api/v1/series"
	labels      = "/api/v1/labels"
	labelValues = "/api/v1/label/*path"
	receive     = "/api/v1/receive"
	rules       = "/api/v1/rules"
	alerts      = "/api/v1/alerts"
)

type Options struct {
	ListenAddress               string
	TLSConfig                   *tls.Config
	GatewayProxyEndpoint        *url.URL
	GatewayProxyClientTLSConfig *tls.Config
	MaxIdleConnsPerHost         int
	MaxConnsPerHost             int

	Tenant       string
	GatewayProxy *httputil.ReverseProxy
}

type Server struct {
	logger  log.Logger
	router  *route.Router
	options *Options

	gatewayProxy *httputil.ReverseProxy
}

func NewServer(logger log.Logger, opt *Options) *Server {

	if logger == nil {
		logger = log.NewNopLogger()
	}
	s := &Server{
		options:      opt,
		router:       route.New(),
		logger:       logger,
		gatewayProxy: opt.GatewayProxy,
	}

	s.router.Get(query, s.wrap())
	s.router.Post(query, s.wrap())
	s.router.Get(queryRange, s.wrap())
	s.router.Post(queryRange, s.wrap())
	s.router.Get(series, s.wrap())
	s.router.Get(labels, s.wrap())
	s.router.Get(labelValues, s.wrap())
	s.router.Get(rules, s.wrap())
	// do provide /api/v1/alerts because thanos does not support alerts filtering as of v0.28.0
	// please filtering alerts by /api/v1/rules
	// s.router.Get(alerts, s.wrap(alerts))

	s.router.Post(receive, s.wrap())

	return s
}

func (s *Server) Router() *route.Router {
	return s.router
}

func (s *Server) wrap() http.HandlerFunc {

	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

		if s.options.Tenant != "" {
			// add the prefix /:tenant_id from path
			req.URL.Path = "/" + s.options.Tenant + req.URL.Path
		}

		s.gatewayProxy.ServeHTTP(w, req)
	})
}

func NewSingleHostReverseProxy(target *url.URL, rt http.RoundTripper) *httputil.ReverseProxy {

	proxy := httputil.NewSingleHostReverseProxy(target)

	proxy.Transport = rt

	return proxy
}

func (s *Server) Run() error {

	srv := &http.Server{
		Handler:   s.router,
		Addr:      s.options.ListenAddress,
		TLSConfig: s.options.TLSConfig,
	}

	if s.options.TLSConfig != nil {
		level.Info(s.logger).Log("msg", "Serving HTTPS", "address", s.options.ListenAddress)
		return srv.ListenAndServeTLS("", "")
	}
	level.Info(s.logger).Log("msg", "Serving plain HTTP", "address", s.options.ListenAddress)
	return srv.ListenAndServe()
}

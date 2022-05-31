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
	targetsMeta = "/api/v1/targets/meta"
	receive     = "/api/v1/receive"
)

type Options struct {
	ListenAddress        string
	TLSConfig            *tls.Config
	Tenant               string
	GatewayProxyEndpoint *url.URL
}

type Server struct {
	logger  log.Logger
	router  *route.Router
	options *Options
}

func NewServer(logger log.Logger, opt *Options) *Server {

	if logger == nil {
		logger = log.NewNopLogger()
	}
	s := &Server{
		options: opt,
		router:  route.New(),
		logger:  logger,
	}

	s.router.Get(query, s.wrap(query))
	s.router.Post(query, s.wrap(query))
	s.router.Get(queryRange, s.wrap(queryRange))
	s.router.Post(queryRange, s.wrap(queryRange))
	s.router.Get(series, s.wrap(series))
	s.router.Get(labels, s.wrap(labels))
	s.router.Get(labelValues, s.wrap(labelValues))
	s.router.Get(targetsMeta, s.wrap(targetsMeta))

	return s
}

func (s *Server) wrap(path string) http.HandlerFunc {

	proxy := httputil.NewSingleHostReverseProxy(s.options.GatewayProxyEndpoint)
	oldDirector := proxy.Director
	if s.options.GatewayProxyEndpoint.Scheme == "https" {
		proxy.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}

	proxy.Director = func(req *http.Request) {
		req.URL.Scheme = s.options.GatewayProxyEndpoint.Scheme
		req.Host = s.options.GatewayProxyEndpoint.Host
		if s.options.Tenant != "" {
			// add the prefix /:tenant_id from path
			req.URL.Path = "/" + s.options.Tenant + req.URL.Path
		}
		oldDirector(req)
	}

	return proxy.ServeHTTP
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

package main

import (
	"net/http"
	"net/url"
	"time"

	extflag "github.com/efficientgo/tools/extkingpin"
	"github.com/go-kit/log"
	"github.com/oklog/run"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	config_util "github.com/prometheus/common/config"
	"github.com/prometheus/common/model"
	"github.com/thanos-io/thanos/pkg/clientconfig"
	"github.com/thanos-io/thanos/pkg/component"
	"github.com/thanos-io/thanos/pkg/extkingpin"
	"github.com/thanos-io/thanos/pkg/extprom"
	"github.com/thanos-io/thanos/pkg/prober"
	httpserver "github.com/thanos-io/thanos/pkg/server/http"
	"gopkg.in/yaml.v2"

	monitoringagentproxy "github.com/WhizardTelemetry/whizard/pkg/monitoring-agent-proxy"
	monitoringgateway "github.com/WhizardTelemetry/whizard/pkg/monitoring-gateway"
)

type agentProxyConfig struct {
	httpBindAddr    *string
	httpGracePeriod *model.Duration
	httpTLSConfig   *string

	serverTlsKey      string
	serverTlsCert     string
	serverTlsClientCa string

	gatewayConfigYaml []byte
	gatewayConfig     gatewayCfg

	tenant string // Tenant is the tenant name to be used for all requests.
}

type gatewayCfg struct {
	clientConfigPath extflag.PathOrContent

	address string

	clientTlsKey       string
	clientTlsCert      string
	serverTlsClientCa  string
	serverName         string
	insecureSkipVerify bool

	maxIdleConnsPerHost int
	maxConnsPerHost     int
}

func registerAgentProxy(app *extkingpin.App) {
	cmd := app.Command(AgentProxy.String(), "Proxy and forward query and remote write API requests to thanos.")

	conf := &agentProxyConfig{}
	conf.registerFlag(cmd)

	cmd.Setup(func(g *run.Group, logger log.Logger, reg *prometheus.Registry, tracer opentracing.Tracer, _ <-chan struct{}, debugLogging bool) error {
		var err error
		conf.gatewayConfigYaml, err = conf.gatewayConfig.clientConfigPath.Content()
		if err != nil {
			return err
		}
		if len(conf.gatewayConfig.address) == 0 {
			return errors.New("no --gateway.address parameter was given")
		}
		if conf.tenant == "" {
			return errors.Wrap(err, "no --tenant parameter was given")
		}

		return runAgentProxy(
			g,
			logger,
			reg,
			conf,
			AgentProxy,
		)
	})
}

func runAgentProxy(
	g *run.Group,
	logger log.Logger,
	reg *prometheus.Registry,
	conf *agentProxyConfig,
	comp component.Component,
) error {

	gatewayClientCfg := clientconfig.NewDefaultHTTPClientConfig()
	if len(conf.gatewayConfigYaml) > 0 {
		if err := yaml.UnmarshalStrict(conf.gatewayConfigYaml, &gatewayClientCfg); err != nil {
			return errors.Wrap(err, "parsing gateway config YAML file failed")
		}
	} else {
		gatewayClientCfg.TLSConfig.CertFile = conf.gatewayConfig.clientTlsCert
		gatewayClientCfg.TLSConfig.KeyFile = conf.gatewayConfig.clientTlsKey
		gatewayClientCfg.TLSConfig.CAFile = conf.gatewayConfig.serverTlsClientCa
		gatewayClientCfg.TLSConfig.ServerName = conf.gatewayConfig.serverName
		gatewayClientCfg.TLSConfig.InsecureSkipVerify = conf.gatewayConfig.insecureSkipVerify

		if conf.gatewayConfig.maxIdleConnsPerHost != gatewayClientCfg.TransportConfig.MaxIdleConnsPerHost {
			gatewayClientCfg.TransportConfig.MaxIdleConnsPerHost = conf.gatewayConfig.maxIdleConnsPerHost
		}
		if conf.gatewayConfig.maxConnsPerHost != gatewayClientCfg.TransportConfig.MaxConnsPerHost {
			gatewayClientCfg.TransportConfig.MaxConnsPerHost = conf.gatewayConfig.maxConnsPerHost
		}
	}

	roundTripper, err := newRoundTripperFromConfig(&gatewayClientCfg, "agent-proxy")
	if err != nil {
		return err
	}

	rawUrl, err := url.Parse(conf.gatewayConfig.address)
	if err != nil {
		return errors.Wrap(err, "setup query downstream service")
	}

	gatewayProxy := monitoringagentproxy.NewSingleHostReverseProxy(rawUrl, roundTripper)

	options := &monitoringagentproxy.Options{
		GatewayProxyEndpoint: rawUrl,
		GatewayProxy:         gatewayProxy,
		Tenant:               conf.tenant,
	}

	httpProbe := prober.NewHTTP()
	statusProber := prober.Combine(
		httpProbe,
		prober.NewInstrumentation(comp, logger, extprom.WrapRegistererWithPrefix("whizard_", reg)),
	)

	srv := httpserver.New(logger, reg, comp, httpProbe,
		httpserver.WithListen(*conf.httpBindAddr),
		httpserver.WithGracePeriod(time.Duration(*conf.httpGracePeriod)),
		httpserver.WithTLSConfig(*conf.httpTLSConfig),
	)

	webhandler := monitoringagentproxy.NewServer(logger, options)
	srv.Handle("/", webhandler.Router())

	g.Add(func() error {
		statusProber.Healthy()

		return srv.ListenAndServe()
	}, func(err error) {
		statusProber.NotReady(err)
		defer statusProber.NotHealthy(err)

		srv.Shutdown(err)
	})
	return nil
}

func (c *agentProxyConfig) registerFlag(cmd extkingpin.FlagClause) {
	c.httpBindAddr, c.httpGracePeriod, c.httpTLSConfig = monitoringgateway.RegisterHTTPFlags(cmd)

	c.gatewayConfig.clientConfigPath = *extflag.RegisterPathOrContent(cmd, "gateway.config", "YAML file that contains downstream tripper configuration.", extflag.WithEnvSubstitution())
	cmd.Flag("gateway.address", "Address to connect whizard monitor-gateway").Default("").StringVar(&c.gatewayConfig.address)
	cmd.Flag("gateway.client-tls-key", "TLS key for gateway client authentication (if the scheme is https).").Default("").StringVar(&c.gatewayConfig.clientTlsKey)
	cmd.Flag("gateway.client-tls-cert", "TLS cert for gateway client authentication (if the scheme is https)(Deprecated, please use gateway.config[/config-file] instead).").Default("").StringVar(&c.gatewayConfig.clientTlsCert)
	cmd.Flag("gateway.server-tls-client-ca", "TLS CA cert for gateway client authentication (if the scheme is https)(Deprecated, please use gateway.config[/config-file] instead).").Default("").StringVar(&c.gatewayConfig.serverTlsClientCa)
	cmd.Flag("gateway.server-name", "Server name used to verify the hostname returned by TLS handshake (if the scheme is https)(Deprecated, please use gateway.config[/config-file] instead).").Default("").StringVar(&c.gatewayConfig.serverName)
	cmd.Flag("gateway.insecure-skip-verify", "Skip verification of gateway TLS certificates (not recommended!)(Deprecated, please use gateway.config[/config-file] instead).").Default("false").BoolVar(&c.gatewayConfig.insecureSkipVerify)

	cmd.Flag("maxIdleConnsPerHost", "Max idle connections per Host(Deprecated, please use gateway.config[/config-file] instead).").Default("100").IntVar(&c.gatewayConfig.maxIdleConnsPerHost)
	cmd.Flag("maxConnsPerHost", "Max connections per Host(Deprecated, please use gateway.config[/config-file] instead).").Default("0").IntVar(&c.gatewayConfig.maxConnsPerHost)

	cmd.Flag("server-tls-key", "TLS Key for HTTP server, leave blank to disable TLS(Deprecated, please use http.config instead).").Default("").StringVar(&c.serverTlsKey)
	cmd.Flag("server-tls-cert", "TLS Certificate for HTTP server, leave blank to disable TLS(Deprecated, please use http.config instead).").Default("").StringVar(&c.serverTlsCert)
	cmd.Flag("server-tls-client-ca", "TLS CA to verify clients against. If no client CA is specified, there is no client verification on server side. (tls.NoClientCert)(Deprecated, please use http.config instead).").Default("").StringVar(&c.serverTlsClientCa)

	cmd.Flag("tenant", "Tenant is the tenant name to be used for all requests.").Default("").StringVar(&c.tenant)
}

func newRoundTripperFromConfig(cfg *clientconfig.HTTPClientConfig, name string) (http.RoundTripper, error) {
	httpClientConfig := config_util.HTTPClientConfig{
		BearerToken:     config_util.Secret(cfg.BearerToken),
		BearerTokenFile: cfg.BearerTokenFile,
		TLSConfig: config_util.TLSConfig{
			CAFile:             cfg.TLSConfig.CAFile,
			CertFile:           cfg.TLSConfig.CertFile,
			KeyFile:            cfg.TLSConfig.KeyFile,
			ServerName:         cfg.TLSConfig.ServerName,
			InsecureSkipVerify: cfg.TLSConfig.InsecureSkipVerify,
		},
	}
	if cfg.ProxyURL != "" {
		var proxy config_util.URL
		err := yaml.Unmarshal([]byte(cfg.ProxyURL), &proxy)
		if err != nil {
			return nil, err
		}
		httpClientConfig.ProxyURL = proxy
	}
	if !cfg.BasicAuth.IsZero() {
		httpClientConfig.BasicAuth = &config_util.BasicAuth{
			Username:     cfg.BasicAuth.Username,
			Password:     config_util.Secret(cfg.BasicAuth.Password),
			PasswordFile: cfg.BasicAuth.PasswordFile,
		}
	}

	if cfg.BearerToken != "" {
		httpClientConfig.BearerToken = config_util.Secret(cfg.BearerToken)
	}

	if cfg.BearerTokenFile != "" {
		httpClientConfig.BearerTokenFile = cfg.BearerTokenFile
	}

	if err := httpClientConfig.Validate(); err != nil {
		return nil, err
	}

	rt, err := clientconfig.NewRoundTripperFromConfig(
		httpClientConfig,
		cfg.TransportConfig,
		name,
	)

	return rt, err
}

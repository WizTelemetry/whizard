package main

import (
	"context"
	"net/url"
	"time"

	extflag "github.com/efficientgo/tools/extkingpin"
	"github.com/go-kit/log"
	"github.com/oklog/run"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/model"
	"github.com/thanos-io/thanos/pkg/component"
	"github.com/thanos-io/thanos/pkg/extkingpin"
	"github.com/thanos-io/thanos/pkg/extprom"
	"github.com/thanos-io/thanos/pkg/prober"
	httpserver "github.com/thanos-io/thanos/pkg/server/http"

	monitoringgateway "github.com/WhizardTelemetry/whizard/pkg/monitoring-gateway"
)

type gatewayConfig struct {
	httpBindAddr    *string
	httpGracePeriod *model.Duration
	httpTLSConfig   *string

	debugEnabledUI bool

	tenantsFilePath    string
	tenantsFileContent string
	refreshInterval    *model.Duration

	tenantHeader    string
	tenantLabelName string

	ExternalRemoteWrites struct {
		ConfigPathOrContent extflag.PathOrContent
	}

	queryConfig       *monitoringgateway.QueryConfig
	rulesQueryConfig  *monitoringgateway.RulesQueryConfig
	remoteWriteConfig *monitoringgateway.RemoteWriteConfig
}

func registerGateway(app *extkingpin.App) {
	cmd := app.Command(Gateway.String(), "Proxy and forward query and remote write API requests to thanos.")

	conf := &gatewayConfig{
		queryConfig:       &monitoringgateway.QueryConfig{},
		rulesQueryConfig:  &monitoringgateway.RulesQueryConfig{},
		remoteWriteConfig: &monitoringgateway.RemoteWriteConfig{},
	}
	conf.registerFlag(cmd)

	cmd.Setup(func(g *run.Group, logger log.Logger, reg *prometheus.Registry, tracer opentracing.Tracer, _ <-chan struct{}, debugLogging bool) error {

		return runGateway(
			g,
			logger,
			reg,
			conf,
			Gateway,
		)
	})
}

func runGateway(
	g *run.Group,
	logger log.Logger,
	reg *prometheus.Registry,
	conf *gatewayConfig,
	comp component.Component,
) error {

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

	options := &monitoringgateway.Options{
		TenantHeader:    conf.tenantHeader,
		TenantLabelName: conf.tenantLabelName,
		EnabledQueryUI:  conf.debugEnabledUI,
	}

	if conf.queryConfig.DownstreamURL != "" {
		downstreamURL, err := url.Parse(conf.queryConfig.DownstreamURL)
		if err != nil {
			return errors.Wrap(err, "setup query downstream service")
		}
		downstreamTripperConfContentYaml, err := conf.queryConfig.DownstreamTripperConfig.TripperPathOrContent.Content()
		if err != nil {
			return err
		}
		downstreamTripper, err := monitoringgateway.ParseTransportConfiguration(downstreamTripperConfContentYaml)
		if err != nil {
			return err
		}
		options.QueryProxy = monitoringgateway.NewSingleHostReverseProxy(downstreamURL, downstreamTripper)
	}

	if conf.rulesQueryConfig.DownstreamURL != "" {
		downstreamURL, err := url.Parse(conf.rulesQueryConfig.DownstreamURL)
		if err != nil {
			return errors.Wrap(err, "setup rules query downstream service")
		}
		downstreamTripperConfContentYaml, err := conf.rulesQueryConfig.DownstreamTripperConfig.TripperPathOrContent.Content()
		if err != nil {
			return err
		}
		downstreamTripper, err := monitoringgateway.ParseTransportConfiguration(downstreamTripperConfContentYaml)
		if err != nil {
			return err
		}
		options.RulesQueryProxy = monitoringgateway.NewSingleHostReverseProxy(downstreamURL, downstreamTripper)
	}

	if conf.remoteWriteConfig.DownstreamURL != "" {
		downstreamURL, err := url.Parse(conf.remoteWriteConfig.DownstreamURL)
		if err != nil {
			return errors.Wrap(err, "setup remote write downstream service")
		}
		downstreamTripperConfContentYaml, err := conf.remoteWriteConfig.DownstreamTripperConfig.TripperPathOrContent.Content()
		if err != nil {
			return err
		}
		downstreamTripper, err := monitoringgateway.ParseTransportConfiguration(downstreamTripperConfContentYaml)
		if err != nil {
			return err
		}
		options.RemoteWriteProxy = monitoringgateway.NewSingleHostReverseProxy(downstreamURL, downstreamTripper)
	}

	content, err := conf.ExternalRemoteWrites.ConfigPathOrContent.Content()
	if err != nil {
		return err
	}
	rwsCfg, err := monitoringgateway.LoadExternalRemoteWriteConfig("", string(content))
	if err != nil {
		return err
	}
	clients, err := monitoringgateway.NewExternalRemoteWriteClients(rwsCfg)
	if err != nil {
		return err
	}

	options.ExternalRWClients = clients

	if conf.tenantsFileContent != "" || conf.tenantsFilePath != "" {
		options.EnabledTenantsAdmission = true
	}

	webhandler := monitoringgateway.NewHandler(logger, reg, options)

	srv.Handle("/", webhandler.Router())

	//
	g.Add(func() error {
		statusProber.Healthy()

		return srv.ListenAndServe()
	}, func(err error) {
		statusProber.NotReady(err)
		defer statusProber.NotHealthy(err)

		srv.Shutdown(err)
	})

	updates := make(chan monitoringgateway.AdmissionControlConfig, 1)

	// The config file path is given initializing config watcher.
	if conf.tenantsFilePath != "" {
		cw, err := monitoringgateway.NewConfigWatcher(log.With(logger, "component", "config-watcher"), reg, conf.tenantsFilePath, *conf.refreshInterval)
		if err != nil {
			return errors.Wrap(err, "failed to initialize config watcher")
		}

		// Check the configuration on before running the watcher.
		if err := cw.ValidateConfig(); err != nil {
			cw.Stop()
			close(updates)
			return errors.Wrap(err, "failed to validate configuration file")
		}

		ctx, cancel := context.WithCancel(context.Background())
		g.Add(func() error {
			return monitoringgateway.ConfigFromWatcher(ctx, updates, cw)
		}, func(error) {
			cancel()
		})
	} else {
		var (
			cf  monitoringgateway.AdmissionControlConfig
			err error
		)
		// The config file content given initialize configuration from content.
		if len(conf.tenantsFileContent) > 0 {
			cf, err = monitoringgateway.ParseConfig([]byte(conf.tenantsFileContent))
			if err != nil {
				close(updates)
				return errors.Wrap(err, "failed to validate configuration content")
			}
		}

		cancel := make(chan struct{})
		g.Add(func() error {
			defer close(updates)
			updates <- cf
			<-cancel
			return nil
		}, func(error) {
			close(cancel)
		})
	}

	cancel := make(chan struct{})
	g.Add(func() error {

		for {
			select {
			case c, ok := <-updates:
				if !ok {
					return nil
				}

				if err := webhandler.SetAdmissionControlHandler(c); err != nil {
					return errors.Wrap(err, "failed to set tenants admission config in gateway")
				}

				// If not, just signal we are ready (this is important during first config load)
				statusProber.Ready()

			case <-cancel:
				return nil
			}
		}
	}, func(err error) {
		close(cancel)
	},
	)

	return nil
}

func (gc *gatewayConfig) registerFlag(cmd extkingpin.FlagClause) {
	gc.httpBindAddr, gc.httpGracePeriod, gc.httpTLSConfig = monitoringgateway.RegisterHTTPFlags(cmd)

	cmd.Flag("debug.enable-ui", "If true, Gateway will proxy and expose Thanos Query UI for debugging.").Default("false").BoolVar(&gc.debugEnabledUI)

	cmd.Flag("tenant.header", "HTTP header to determine tenant for write requests.").Default("WHIZARD-TENANT").StringVar(&gc.tenantHeader)
	cmd.Flag("tenant.label-name", "Label name through which the tenant will be announced.").Default("tenant_id").StringVar(&gc.tenantLabelName)
	cmd.Flag("tenant.admission-control-config-file", "Path to file that contains the configuration. A watcher is initialized to watch changes and update the dynamically.").PlaceHolder("<path>").StringVar(&gc.tenantsFilePath)
	cmd.Flag("tenant.admission-control-config", "Alternative to 'tenant.admission-control-config-file' flag (lower priority). Content of file that contains the configuration.").PlaceHolder("<content>").StringVar(&gc.tenantsFileContent)
	gc.refreshInterval = extkingpin.ModelDuration(cmd.Flag("tenant.admission-control-config-file-refresh-interval", "Refresh interval to re-read the configuration file. (used as a fallback)").Default("1m"))

	gc.ExternalRemoteWrites.ConfigPathOrContent = *extflag.RegisterPathOrContent(cmd, "external-remote-writes.config", "Path to YAML config for the external remote-write configurations, that specify servers where received remote-write requests should be forwarded to.", extflag.WithEnvSubstitution())

	gc.queryConfig.RegisterFlag(cmd)
	gc.rulesQueryConfig.RegisterFlag(cmd)
	gc.remoteWriteConfig.RegisterFlag(cmd)
}

var (
	Gateway    = customcomponent{name: "gateway"}
	AgentProxy = customcomponent{name: "agent-proxy"}
)

type customcomponent struct {
	name string
}

func (c customcomponent) String() string { return c.name }

package main

import (
	"crypto/tls"
	"io/ioutil"
	"net/url"
	"path"
	"strings"

	"github.com/alecthomas/kong"
	"github.com/pkg/errors"
	"github.com/prometheus/common/config"
	"github.com/thanos-io/thanos/pkg/logging"
	thanos_tls "github.com/thanos-io/thanos/pkg/tls"
	"gopkg.in/yaml.v2"

	monitoringgateway "github.com/kubesphere/whizard/pkg/monitoring-gateway"
)

var cli struct {
	Log struct {
		Level  string `enum:"debug,info,warn,error" default:"info" help:"Log filtering level. Possible options: ${enum}."`
		Format string `enum:"json,logfmt" default:"logfmt" help:"Log format to use. Possible options: ${enum}."`
	} `embed:"" prefix:"log."`

	HttpAddress       string `default:"0.0.0.0:9090" help:"Listen host:port for HTTP endpoints."`
	ServerTlsKey      string `default:"" help:"TLS Key for HTTP server, leave blank to disable TLS."`
	ServerTlsCert     string `default:"" help:"TLS Certificate for HTTP server, leave blank to disable TLS."`
	ServerTlsClientCa string `default:"" help:"TLS CA to verify clients against. If no client CA is specified, there is no client verification on server side. (tls.NoClientCert)"`

	RemoteWrites struct {
		ConfigFile string `default:"" help:"Path to YAML config for the remote-write configurations, that specify servers where received remote-write requests should be forwarded to."`
		Config     string `default:"" help:"Alternative to 'remote-writes.config-file' flag (mutually exclusive)" `
	} `embed:"" prefix:"remote-writes."`

	RemoteWrite struct {
		Address    string `default:"" help:"Address to send remote write requests. (Deprecated, please use remote-writes.config[/config-file] instead)"`
		ConfigFile string `default:"" help:"Downstream receive service configuration file. (Deprecated, please use remote-writes.config[/config-file] instead)"`
		Config     string `default:"" help:"Downstream receive service configuration content. (Deprecated, please use remote-writes.config[/config-file] instead)"`
	} `embed:"" prefix:"remote-write."`
	Query struct {
		Address    string `default:"" help:"Address to send query requests."`
		ConfigFile string `default:"" help:"Downstream query/query-frontend service configuration file."`
		Config     string `default:"" help:"Downstream query/query-frontend service configuration content." `
	} `embed:"" prefix:"query."`
	QueryRules struct {
		Address    string `default:"" help:"Address to send rules query requests."`
		ConfigFile string `default:"" help:"Downstream query/query-frontend service configuration file."`
		Config     string `default:"" help:"Downstream query/query-frontend service configuration content."`
	} `embed:"" prefix:"query-rules."`
	Tenant struct {
		Header    string `default:"WHIZARD-TENANT" help:"Http header to determine tenant for requests"`
		LabelName string `default:"tenant_id" help:"Label name through which the tenant will be announced"`
	} `embed:"" prefix:"tenant."`
}

type Config struct {
	TLSConfig *config.TLSConfig `yaml:"tls_config,omitempty" json:"tls_config,omitempty"`
}

func main() {

	ctx := kong.Parse(&cli)

	logger := logging.NewLogger(cli.Log.Level, cli.Log.Format, "")

	var options = monitoringgateway.Options{
		ListenAddress:   cli.HttpAddress,
		TenantHeader:    cli.Tenant.Header,
		TenantLabelName: cli.Tenant.LabelName,
	}
	var err error

	options.TLSConfig, err = thanos_tls.NewServerConfig(logger, cli.ServerTlsCert, cli.ServerTlsKey, cli.ServerTlsClientCa)
	ctx.FatalIfErrorf(err)
	if options.TLSConfig != nil && options.TLSConfig.ClientAuth == tls.RequireAndVerifyClientCert {
		options.CertAuthenticator = monitoringgateway.NewCertAuthenticator()
	}

	rwsCfg, err := monitoringgateway.LoadRemoteWritesConfig(cli.RemoteWrites.ConfigFile, cli.RemoteWrites.Config)
	if err != nil {
		ctx.FatalIfErrorf(err)
	}
	if cli.RemoteWrite.Address != "" {
		rwUrl, err := url.Parse(cli.RemoteWrite.Address)
		ctx.FatalIfErrorf(err)
		cfg, err := parseConfig(cli.RemoteWrite.ConfigFile, cli.RemoteWrite.Config)
		if err != nil {
			ctx.FatalIfErrorf(err)
		}
		if !strings.HasSuffix(strings.TrimSuffix(rwUrl.Path, "/"), "/api/v1/receive") { // to make it compactible with previous config
			rwUrl.Path = path.Join(rwUrl.Path, "/api/v1/receive")
		}
		rwCfg := monitoringgateway.RemoteWriteConfig{URL: &config.URL{URL: rwUrl}}
		if cfg != nil && cfg.TLSConfig != nil {
			rwCfg.TLSConfig = *cfg.TLSConfig

		}
		rwsCfg = append(rwsCfg, rwCfg)
	}
	options.RemoteWriteHandler, err = monitoringgateway.NewRemoteWriteHandler(rwsCfg, cli.Tenant.Header)
	if err != nil {
		ctx.FatalIfErrorf(err)
	}

	if cli.Query.Address != "" {
		qUrl, err := url.Parse(cli.Query.Address)
		ctx.FatalIfErrorf(err)
		cfg, err := parseConfig(cli.Query.ConfigFile, cli.Query.Config)
		if err != nil {
			ctx.FatalIfErrorf(err)
		}
		if cfg != nil && cfg.TLSConfig != nil {
			tlsConfig, err := config.NewTLSConfig(cfg.TLSConfig)
			if err != nil {
				ctx.FatalIfErrorf(err)
			}
			options.QueryProxy = monitoringgateway.NewSingleHostReverseProxy(qUrl, tlsConfig)
		} else {
			options.QueryProxy = monitoringgateway.NewSingleHostReverseProxy(qUrl, nil)
		}

	}

	if cli.QueryRules.Address != "" {
		qUrl, err := url.Parse(cli.QueryRules.Address)
		ctx.FatalIfErrorf(err)
		cfg, err := parseConfig(cli.QueryRules.ConfigFile, cli.QueryRules.Config)
		if err != nil {
			ctx.FatalIfErrorf(err)
		}
		if cfg != nil && cfg.TLSConfig != nil {
			tlsConfig, err := config.NewTLSConfig(cfg.TLSConfig)
			if err != nil {
				ctx.FatalIfErrorf(err)
			}
			options.QueryRulesProxy = monitoringgateway.NewSingleHostReverseProxy(qUrl, tlsConfig)
		} else {
			options.QueryRulesProxy = monitoringgateway.NewSingleHostReverseProxy(qUrl, nil)
		}

	}

	handler := monitoringgateway.NewHandler(logger, &options)

	err = handler.Run()
	ctx.FatalIfErrorf(err)
}

func parseConfig(file string, content string) (*Config, error) {
	var buff []byte
	if len(file) > 0 {
		c, err := ioutil.ReadFile(file)
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
	cfg := &Config{}
	if err := yaml.UnmarshalStrict(buff, cfg); err != nil {
		return nil, errors.Wrap(err, "")
	}
	return cfg, nil
}

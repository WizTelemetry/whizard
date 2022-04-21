package main

import (
	"crypto/tls"
	"net/url"

	"github.com/alecthomas/kong"
	"github.com/thanos-io/thanos/pkg/logging"
	thanos_tls "github.com/thanos-io/thanos/pkg/tls"

	monitoringgateway "github.com/kubesphere/paodin-monitoring/pkg/monitoring-gateway"
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

	RemoteWrite struct {
		Address string `default:"" help:"Address to send remote write requests."`
	} `embed:"" prefix:"remote-write."`
	Query struct {
		Address string `default:"" help:"Address to send query requests."`
	} `embed:"" prefix:"query."`
	Tenant struct {
		Header    string `default:"THANOS-TENANT" help:"Http header to determine tenant for requests"`
		LabelName string `default:"tenant_id" help:"Label name through which the tenant will be announced"`
	} `embed:"" prefix:"tenant."`
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

	if cli.RemoteWrite.Address != "" {
		rwUrl, err := url.Parse(cli.RemoteWrite.Address)
		ctx.FatalIfErrorf(err)
		options.RemoteWriteProxy = monitoringgateway.NewSingleHostReverseProxy(rwUrl)
	}
	if cli.Query.Address != "" {
		qUrl, err := url.Parse(cli.Query.Address)
		ctx.FatalIfErrorf(err)
		options.QueryProxy = monitoringgateway.NewSingleHostReverseProxy(qUrl)
	}

	handler := monitoringgateway.NewHandler(logger, &options)

	err = handler.Run()
	ctx.FatalIfErrorf(err)
}

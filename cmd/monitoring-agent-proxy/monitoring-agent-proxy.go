package main

import (
	"net/url"

	"github.com/alecthomas/kong"
	"github.com/thanos-io/thanos/pkg/logging"
	thanos_tls "github.com/thanos-io/thanos/pkg/tls"

	monitoringagentproxy "github.com/kubesphere/whizard/pkg/monitoring-agent-proxy"
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

	MonitorGateway struct {
		Address            string `default:"" help:"Address to connect whizard monitor-gateway"`
		ClientTlsKey       string `default:"" help:"TLS Key for HTTP client, leave blank to skip verify."`
		ClientTlsCert      string `default:"" help:"TLS Certificate for HTTP client, leave blank to skip verify."`
		ServerTlsClientCa  string `default:"" help:"TLS CA to verify clients against. If no client CA is specified, there is no client verification on server side. (tls.NoClientCert)"`
		ServerName         string `default:"" help:"TLS ServerName used to verify the hostname"`
		InsecureSkipVerify bool   `default:"true" help:"Disable certificate validation."`
	} `embed:"" prefix:"gateway."`

	Tenant              string
	MaxIdleConnsPerHost int `default:"100" name:"maxIdleConnsPerHost" help:"Max idle connections per Host"`
	MaxConnsPerHost     int `default:"0" name:"maxConnsPerHost" help:"Max connections per Host"`
}

func main() {

	ctx := kong.Parse(&cli)
	logger := logging.NewLogger(cli.Log.Level, cli.Log.Format, "")

	rawUrl, err := url.Parse(cli.MonitorGateway.Address)
	ctx.FatalIfErrorf(err)

	options := &monitoringagentproxy.Options{
		Tenant:               cli.Tenant,
		ListenAddress:        cli.HttpAddress,
		GatewayProxyEndpoint: rawUrl,
		MaxIdleConnsPerHost:  cli.MaxIdleConnsPerHost,
		MaxConnsPerHost:      cli.MaxConnsPerHost,
	}
	options.GatewayProxyClientTLSConfig, err = thanos_tls.NewClientConfig(logger, cli.MonitorGateway.ClientTlsCert, cli.MonitorGateway.ClientTlsKey, cli.MonitorGateway.ServerTlsClientCa, cli.MonitorGateway.ServerName, cli.MonitorGateway.InsecureSkipVerify)
	ctx.FatalIfErrorf(err)

	options.TLSConfig, err = thanos_tls.NewServerConfig(logger, cli.ServerTlsCert, cli.ServerTlsKey, cli.ServerTlsClientCa)
	ctx.FatalIfErrorf(err)

	server := monitoringagentproxy.NewServer(logger, options)
	err = server.Run()
	ctx.FatalIfErrorf(err)
}

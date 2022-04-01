package main

import (
	"context"
	"crypto/tls"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/alecthomas/kong"
	"github.com/thanos-io/thanos/pkg/logging"
	thanos_tls "github.com/thanos-io/thanos/pkg/tls"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/clientcmd"

	monitoringv1alpha1 "github.com/kubesphere/paodin-monitoring/pkg/api/monitoring/v1alpha1"
	monitoringclient "github.com/kubesphere/paodin-monitoring/pkg/client/clientset/versioned"
	monitoringinformers "github.com/kubesphere/paodin-monitoring/pkg/client/informers/externalversions"
	"github.com/kubesphere/paodin-monitoring/pkg/gateway"
)

var cli struct {
	Log struct {
		Level  string `enum:"debug,info,warn,error" default:"info" help:"Log filtering level. Possible options: ${enum}."`
		Format string `enum:"json,logfmt" default:"logfmt" help:"Log format to use. Possible options: ${enum}."`
	} `embed:"" prefix:"log."`

	HttpAddress       string `default:"0.0.0.0:9080" help:"Listen host:port for HTTP endpoints."`
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

	KubeConfig string `name:"kubeconfig" default:"" help:"Path to the kubeconfig file to use"`

	Agent struct {
		Namespace string `default:"" help:"Agent namespace"`
		Selector  string ` default:"" help:"Selector (agent label query) to filter on, supports '=', '==', and '!='.(e.g. -l key1=value1,key2=value2)"`
	} `embed:"" prefix:"agent."`
}

func main() {

	ctx := kong.Parse(&cli)

	logger := logging.NewLogger(cli.Log.Level, cli.Log.Format, "")

	var options = gateway.Options{
		ListenAddress:   cli.HttpAddress,
		TenantHeader:    cli.Tenant.Header,
		TenantLabelName: cli.Tenant.LabelName,
	}
	var err error

	options.TLSConfig, err = thanos_tls.NewServerConfig(logger, cli.ServerTlsCert, cli.ServerTlsKey, cli.ServerTlsClientCa)
	ctx.FatalIfErrorf(err)
	if options.TLSConfig != nil && options.TLSConfig.ClientAuth == tls.RequireAndVerifyClientCert {
		options.CertAuthenticator = gateway.NewCertAuthenticator()
	}

	if cli.RemoteWrite.Address != "" {
		rwUrl, err := url.Parse(cli.RemoteWrite.Address)
		ctx.FatalIfErrorf(err)
		options.RemoteWriteProxy = &httputil.ReverseProxy{Director: gateway.NewDirector(rwUrl)}
	}
	if cli.Query.Address != "" {
		qUrl, err := url.Parse(cli.Query.Address)
		ctx.FatalIfErrorf(err)
		options.QueryProxy = &httputil.ReverseProxy{Director: gateway.NewDirector(qUrl)}
	}

	_, err = labels.Parse(cli.Agent.Selector)
	ctx.FatalIfErrorf(err)

	k8scfg, err := clientcmd.BuildConfigFromFlags("", cli.KubeConfig)
	ctx.FatalIfErrorf(err)
	mclient, err := monitoringclient.NewForConfig(k8scfg)
	ctx.FatalIfErrorf(err)
	factory := monitoringinformers.NewSharedInformerFactoryWithOptions(mclient, time.Minute*5,
		monitoringinformers.WithNamespace(cli.Agent.Namespace),
		monitoringinformers.WithTweakListOptions(func(lo *metav1.ListOptions) {
			lo.LabelSelector = cli.Agent.Selector
		}))
	_, err = factory.ForResource(monitoringv1alpha1.SchemeGroupVersion.WithResource("agents"))
	ctx.FatalIfErrorf(err)
	stop := context.TODO().Done()
	factory.Start(stop)
	factory.WaitForCacheSync(stop)
	agentLister := factory.Monitoring().V1alpha1().Agents().Lister()
	options.GetAgentFunc = func(name *types.NamespacedName) (*monitoringv1alpha1.Agent, error) {
		return agentLister.Agents(name.Namespace).Get(name.Name)
	}

	handler := gateway.NewHandler(logger, &options)

	err = handler.Run()
	ctx.FatalIfErrorf(err)
}

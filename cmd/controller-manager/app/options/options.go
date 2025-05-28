package options

import (
	"flag"
	"strings"

	cliflag "k8s.io/component-base/cli/flag"
	"k8s.io/klog/v2"
)

type ControllerManagerOptions struct {
	MetricsAddr     string
	SecureMetrics   bool
	MetricsCertPath string
	MetricsCertName string
	MetricsCertKey  string
	EnableHTTP2     bool

	EnableLeaderElection bool
	ProbeAddr            string
}

func NewControllerManagerOptions() *ControllerManagerOptions {
	return &ControllerManagerOptions{
		EnableLeaderElection: false,
		MetricsAddr:          "0",
		SecureMetrics:        true,
		MetricsCertName:      "tls.crt",
		MetricsCertKey:       "tls.key",
		MetricsCertPath:      "",
		EnableHTTP2:          false,

		ProbeAddr: ":8081",
	}
}

func (s *ControllerManagerOptions) Flags() cliflag.NamedFlagSets {
	fss := cliflag.NamedFlagSets{}

	kfs := fss.FlagSet("klog")
	local := flag.NewFlagSet("klog", flag.ExitOnError)
	klog.InitFlags(local)
	local.VisitAll(func(fl *flag.Flag) {
		fl.Name = strings.Replace(fl.Name, "_", "-", -1)
		kfs.AddGoFlag(fl)
	})

	mfs := fss.FlagSet("metrics")
	mfs.StringVar(&s.MetricsAddr, "metrics-bind-address", s.MetricsAddr, "The address the metrics endpoint binds to. "+
		"Use :8443 for HTTPS or :8080 for HTTP, or leave as 0 to disable the metrics service.")
	mfs.BoolVar(&s.SecureMetrics, "metrics-secure", s.SecureMetrics,
		"If set, the metrics endpoint is served securely via HTTPS. Use --metrics-secure=false to use HTTP instead.")
	mfs.StringVar(&s.MetricsCertPath, "metrics-cert-path", s.MetricsCertPath,
		"The directory that contains the metrics server certificate.")
	mfs.StringVar(&s.MetricsCertName, "metrics-cert-name", s.MetricsCertName, "The name of the metrics server certificate file.")
	mfs.StringVar(&s.MetricsCertKey, "metrics-cert-key", s.MetricsCertKey, "The name of the metrics server key file.")
	mfs.BoolVar(&s.EnableHTTP2, "enable-http2", s.EnableHTTP2,
		"If set, HTTP/2 will be enabled for the metrics and webhook servers")

	ofs := fss.FlagSet("other")
	ofs.StringVar(&s.ProbeAddr, "health-probe-bind-address", s.ProbeAddr, "The address the probe endpoint binds to.")
	ofs.BoolVar(&s.EnableLeaderElection, "leader-elect", s.EnableLeaderElection,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")

	return fss
}

func (s *ControllerManagerOptions) Validate() []error {
	var errs []error

	return errs
}

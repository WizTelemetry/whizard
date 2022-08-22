package options

import (
	"flag"
	"strings"
	"time"

	"github.com/spf13/pflag"
	"k8s.io/client-go/tools/leaderelection"
	cliflag "k8s.io/component-base/cli/flag"
	"k8s.io/klog/v2"

	"github.com/kubesphere/whizard/pkg/client/k8s"
	"github.com/kubesphere/whizard/pkg/controllers/config"
	monitoring "github.com/kubesphere/whizard/pkg/controllers/monitoring/options"
)

type ControllerManagerOptions struct {
	KubernetesOptions *k8s.KubernetesOptions
	MonitoringOptions *monitoring.Options

	LeaderElect    bool
	LeaderElection *leaderelection.LeaderElectionConfig
	WebhookCertDir string

	MetricsBindAddress     string
	HealthProbeBindAddress string
}

func NewControllerManagerOptions() *ControllerManagerOptions {
	return &ControllerManagerOptions{
		KubernetesOptions: k8s.NewKubernetesOptions(),
		MonitoringOptions: monitoring.NewOptions(),

		LeaderElection: &leaderelection.LeaderElectionConfig{
			LeaseDuration: 30 * time.Second,
			RenewDeadline: 15 * time.Second,
			RetryPeriod:   5 * time.Second,
		},
		LeaderElect:    false,
		WebhookCertDir: "",

		MetricsBindAddress:     ":8080",
		HealthProbeBindAddress: ":8081",
	}
}

func (s *ControllerManagerOptions) Flags() cliflag.NamedFlagSets {
	fss := cliflag.NamedFlagSets{}
	s.KubernetesOptions.AddFlags(fss.FlagSet("kubernetes"), s.KubernetesOptions)
	s.MonitoringOptions.AddFlags(fss.FlagSet("monitoring"), s.MonitoringOptions)

	fs := fss.FlagSet("leaderelection")
	s.bindLeaderElectionFlags(s.LeaderElection, fs)

	fs.BoolVar(&s.LeaderElect, "leader-elect", s.LeaderElect, ""+
		"Whether to enable leader election. This field should be enabled when controller manager"+
		"deployed with multiple replicas.")

	fs.StringVar(&s.WebhookCertDir, "webhook-cert-dir", s.WebhookCertDir, ""+
		"Certificate directory used to setup webhooks, need tls.crt and tls.key placed inside."+
		"if not set, webhook server would look up the server key and certificate in"+
		"{TempDir}/k8s-webhook-server/serving-certs")

	kfs := fss.FlagSet("klog")
	local := flag.NewFlagSet("klog", flag.ExitOnError)
	klog.InitFlags(local)
	local.VisitAll(func(fl *flag.Flag) {
		fl.Name = strings.Replace(fl.Name, "_", "-", -1)
		kfs.AddGoFlag(fl)
	})

	ofs := fss.FlagSet("other")
	ofs.StringVar(&s.MetricsBindAddress, "metrics-bind-address", ":8080", "The address the metric endpoint binds to.")
	ofs.StringVar(&s.HealthProbeBindAddress, "health-probe-bind-address", ":8081", "The address the probe endpoint binds to.")

	return fss
}

func (s *ControllerManagerOptions) Validate() []error {
	var errs []error
	errs = append(errs, s.KubernetesOptions.Validate()...)
	errs = append(errs, s.MonitoringOptions.Validate()...)
	return errs
}

func (s *ControllerManagerOptions) bindLeaderElectionFlags(l *leaderelection.LeaderElectionConfig, fs *pflag.FlagSet) {
	fs.DurationVar(&l.LeaseDuration, "leader-elect-lease-duration", l.LeaseDuration, ""+
		"The duration that non-leader candidates will wait after observing a leadership "+
		"renewal until attempting to acquire leadership of a led but unrenewed leader "+
		"slot. This is effectively the maximum duration that a leader can be stopped "+
		"before it is replaced by another candidate. This is only applicable if leader "+
		"election is enabled.")
	fs.DurationVar(&l.RenewDeadline, "leader-elect-renew-deadline", l.RenewDeadline, ""+
		"The interval between attempts by the acting master to renew a leadership slot "+
		"before it stops leading. This must be less than or equal to the lease duration. "+
		"This is only applicable if leader election is enabled.")
	fs.DurationVar(&l.RetryPeriod, "leader-elect-retry-period", l.RetryPeriod, ""+
		"The duration the clients should wait between attempting acquisition and renewal "+
		"of a leadership. This is only applicable if leader election is enabled.")
}

// MergeConfig merge new config without validation
// When misconfigured, the app should just crash directly
func (s *ControllerManagerOptions) MergeConfig(cfg *config.Config) {
	if cfg.KubernetesOptions != nil {
		cfg.KubernetesOptions.ApplyTo(s.KubernetesOptions)
	}

	if cfg.MonitoringOptions != nil {
		cfg.MonitoringOptions.ApplyTo(s.MonitoringOptions)
	}
}

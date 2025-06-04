package app

import (
	"context"
	"crypto/tls"
	"fmt"
	"path/filepath"

	promv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	"github.com/prometheus/common/version"
	"github.com/spf13/cobra"
	apiextensions "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/errors"
	cliflag "k8s.io/component-base/cli/flag"
	"k8s.io/klog/v2"
	"k8s.io/klog/v2/textlogger"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/certwatcher"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"
	"sigs.k8s.io/controller-runtime/pkg/metrics/filters"
	metricsserver "sigs.k8s.io/controller-runtime/pkg/metrics/server"

	"github.com/WhizardTelemetry/whizard/cmd/controller-manager/app/options"
	"github.com/WhizardTelemetry/whizard/pkg/apis"
)

func NewControllerManagerCommand() *cobra.Command {
	// Here will create a default whizard controller manager options
	s := options.NewControllerManagerOptions()

	// Initialize command to run our controllers later
	cmd := &cobra.Command{
		Use:   "controller-manager",
		Short: `Whizard controller manager`,
		Run: func(cmd *cobra.Command, args []string) {
			if errs := s.Validate(); len(errs) != 0 {
				klog.Fatal(errors.NewAggregate(errs))
			}
			if err := Run(s, signals.SetupSignalHandler()); err != nil {
				klog.Fatalf("Failed to run whizard controller manager: %v", err)
			}
		},
		SilenceUsage: true,
	}

	fs := cmd.Flags()
	// Add pre-defined flags into command
	namedFlagSets := s.Flags()

	for _, f := range namedFlagSets.FlagSets {
		fs.AddFlagSet(f)
	}

	usageFmt := "Usage:\n  %s\n"
	cmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		_, _ = fmt.Fprintf(cmd.OutOrStdout(), "%s\n\n"+usageFmt, cmd.Long, cmd.UseLine())
		cliflag.PrintSections(cmd.OutOrStdout(), namedFlagSets, 0)
	})

	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Print the version of Whizard controller",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Println(version.Print("whizard controller"))
		},
	}

	cmd.AddCommand(versionCmd)

	return cmd
}

func Run(s *options.ControllerManagerOptions, ctx context.Context) error {

	disableHTTP2 := func(c *tls.Config) {
		klog.Info("disabling http/2")
		c.NextProtos = []string{"http/1.1"}
	}
	var metricsCertWatcher *certwatcher.CertWatcher
	var tlsOpts []func(*tls.Config)
	if !s.EnableHTTP2 {
		tlsOpts = append(tlsOpts, disableHTTP2)
	}

	// Metrics endpoint is enabled in 'config/default/kustomization.yaml'. The Metrics options configure the server.
	// More info:
	// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.21.0/pkg/metrics/server
	// - https://book.kubebuilder.io/reference/metrics.html
	metricsServerOptions := metricsserver.Options{
		BindAddress:   s.MetricsAddr,
		SecureServing: s.SecureMetrics,
		TLSOpts:       tlsOpts,
	}

	if s.SecureMetrics {
		// FilterProvider is used to protect the metrics endpoint with authn/authz.
		// These configurations ensure that only authorized users and service accounts
		// can access the metrics endpoint. The RBAC are configured in 'config/rbac/kustomization.yaml'. More info:
		// https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.21.0/pkg/metrics/filters#WithAuthenticationAndAuthorization
		metricsServerOptions.FilterProvider = filters.WithAuthenticationAndAuthorization
	}

	// If the certificate is not specified, controller-runtime will automatically
	// generate self-signed certificates for the metrics server. While convenient for development and testing,
	// this setup is not recommended for production.
	//
	// TODO(user): If you enable certManager, uncomment the following lines:
	// - [METRICS-WITH-CERTS] at config/default/kustomization.yaml to generate and use certificates
	// managed by cert-manager for the metrics server.
	// - [PROMETHEUS-WITH-CERTS] at config/prometheus/kustomization.yaml for TLS certification.
	if len(s.MetricsCertPath) > 0 {
		klog.Info("Initializing metrics certificate watcher using provided certificates",
			"metrics-cert-path", s.MetricsCertPath, "metrics-cert-name", s.MetricsCertName, "metrics-cert-key", s.MetricsCertKey)

		var err error
		metricsCertWatcher, err = certwatcher.New(
			filepath.Join(s.MetricsCertPath, s.MetricsCertName),
			filepath.Join(s.MetricsCertPath, s.MetricsCertKey),
		)
		if err != nil {
			return fmt.Errorf("unable to initialize metrics certificate watcher: %w", err)
		}

		metricsServerOptions.TLSOpts = append(metricsServerOptions.TLSOpts, func(config *tls.Config) {
			config.GetCertificate = metricsCertWatcher.GetCertificate
		})
	}

	mgrOptions := ctrl.Options{

		Metrics: metricsServerOptions,

		HealthProbeBindAddress: s.ProbeAddr,
		LeaderElection:         s.EnableLeaderElection,
		LeaderElectionID:       "whizard-controller-manager-leader-election",
		// LeaderElectionReleaseOnCancel defines if the leader should step down voluntarily
		// when the Manager ends. This requires the binary to immediately end when the
		// Manager is stopped, otherwise, this setting is unsafe. Setting this significantly
		// speeds up voluntary leader transitions as the new leader don't have to wait
		// LeaseDuration time first.
		//
		// In the default scaffold provided, the program ends immediately after
		// the manager stops, so would be fine to enable this option. However,
		// if you are doing or is intended to do any operation such as perform cleanups
		// after the manager stops then its usage might be unsafe.
		// LeaderElectionReleaseOnCancel: true,
	}

	klog.Info("setting up manager")
	ctrl.SetLogger(textlogger.NewLogger(textlogger.NewConfig()))
	// Use 8443 instead of 443 cause we need root permission to bind port 443
	// Init controller manager
	mgr, err := manager.New(ctrl.GetConfigOrDie(), mgrOptions)
	if err != nil {
		return fmt.Errorf("unable to set up overall controller manager: %w", err)
	}
	_ = apis.AddToScheme(mgr.GetScheme())
	_ = apiextensions.AddToScheme(mgr.GetScheme())
	_ = promv1.AddToScheme(mgr.GetScheme())

	// register common meta types into schemas.
	metav1.AddToGroupVersion(mgr.GetScheme(), metav1.SchemeGroupVersion)

	if metricsCertWatcher != nil {
		klog.Info("Adding metrics certificate watcher to manager")
		if err := mgr.Add(metricsCertWatcher); err != nil {
			return fmt.Errorf("unable to add metrics certificate watcher to manager: %w", err)
		}
	}

	if err = addControllers(mgr, ctx); err != nil {
		return fmt.Errorf("unable to register controllers to the manager: %v", err)
	}

	if err := mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		return fmt.Errorf("unable to set up health check: %w", err)
	}
	if err := mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		return fmt.Errorf("unable to set up ready check: %w", err)
	}

	klog.Info("Starting the controllers.")
	if err = mgr.Start(ctx); err != nil {
		return fmt.Errorf("unable to run the manager: %w", err)
	}

	return nil
}

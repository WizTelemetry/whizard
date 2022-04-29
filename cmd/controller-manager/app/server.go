package app

import (
	"context"
	"fmt"
	"os"

	promv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	"github.com/spf13/cobra"
	apiextensions "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	cliflag "k8s.io/component-base/cli/flag"
	"k8s.io/klog/v2"
	"k8s.io/klog/v2/klogr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"

	"github.com/kubesphere/paodin/cmd/controller-manager/app/options"
	"github.com/kubesphere/paodin/pkg/apis"
	"github.com/kubesphere/paodin/pkg/client/k8s"
	"github.com/kubesphere/paodin/pkg/controllers/config"
	"github.com/kubesphere/paodin/pkg/informers"
)

func NewControllerManagerCommand() *cobra.Command {
	// Here will create a default paodin controller manager options
	s := options.NewPaodinControllerManagerOptions()
	conf, err := config.TryLoadFromDisk()
	if err == nil {
		// make sure LeaderElection is not nil
		// override paodin controller manager options
		s.KubernetesOptions = conf.KubernetesOptions
		s.MonitoringOptions = conf.MonitoringOptions
	} else {
		klog.Fatal("Failed to load configuration from disk", err)
	}

	// Initialize command to run our controllers later
	cmd := &cobra.Command{
		Use:   "controller-manager",
		Short: `Paodin controller manager`,
		Run: func(cmd *cobra.Command, args []string) {
			if errs := s.Validate(); len(errs) != 0 {
				klog.Error(utilerrors.NewAggregate(errs))
				os.Exit(1)
			}
			if err = Run(s, config.WatchConfigChange(), signals.SetupSignalHandler()); err != nil {
				klog.Error(err)
				os.Exit(1)
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
		Short: "Print the version of Paodin controller",
		Run: func(cmd *cobra.Command, args []string) {
			// cmd.Println(version.Get())
		},
	}

	cmd.AddCommand(versionCmd)

	return cmd
}

func Run(s *options.PaodinControllerManagerOptions, configCh <-chan config.Config, ctx context.Context) error {
	ictx, cancelFunc := context.WithCancel(context.TODO())
	errCh := make(chan error)
	defer close(errCh)
	go func() {
		if err := run(s, ictx); err != nil {
			errCh <- err
		}
	}()

	// The ctx (signals.SetupSignalHandler()) is to control the entire program life cycle,
	// The ictx(internal context)  is created here to control the life cycle of the controller-manager(all controllers, sharedInformer, webhook etc.)
	// when config changed, stop server and renew context, start new server
	for {
		select {
		case <-ctx.Done():
			cancelFunc()
			return nil
		case cfg := <-configCh:
			cancelFunc()
			s.MergeConfig(&cfg)
			ictx, cancelFunc = context.WithCancel(context.TODO())
			go func() {
				if err := run(s, ictx); err != nil {
					errCh <- err
				}
			}()
		case err := <-errCh:
			cancelFunc()
			return err
		}
	}
}

func run(s *options.PaodinControllerManagerOptions, ctx context.Context) error {
	// Init k8s client
	kubernetesClient, err := k8s.NewKubernetesClient(s.KubernetesOptions)
	if err != nil {
		klog.Errorf("Failed to create kubernetes clientset %v", err)
		return err
	}

	// Init informers
	informerFactory := informers.NewInformerFactories(
		kubernetesClient.Kubernetes(),
		kubernetesClient.ApiExtensions())

	mgrOptions := manager.Options{
		CertDir: s.WebhookCertDir,
		Port:    8443,

		MetricsBindAddress:     s.MetricsBindAddress,
		HealthProbeBindAddress: s.HealthProbeBindAddress,
	}

	if s.LeaderElect {
		mgrOptions.LeaderElection = s.LeaderElect
		mgrOptions.LeaderElectionID = "paodin-controller-manager-leader-election"
		mgrOptions.LeaseDuration = &s.LeaderElection.LeaseDuration
		mgrOptions.RetryPeriod = &s.LeaderElection.RetryPeriod
		mgrOptions.RenewDeadline = &s.LeaderElection.RenewDeadline
	}

	klog.V(0).Info("setting up manager")
	ctrl.SetLogger(klogr.New())
	// Use 8443 instead of 443 cause we need root permission to bind port 443
	// Init controller manager
	mgr, err := manager.New(kubernetesClient.Config(), mgrOptions)
	if err != nil {
		klog.Fatalf("unable to set up overall controller manager: %v", err)
	}
	apis.AddToScheme(mgr.GetScheme())
	_ = apiextensions.AddToScheme(mgr.GetScheme())

	promv1.AddToScheme(mgr.GetScheme())

	// register common meta types into schemas.
	metav1.AddToGroupVersion(mgr.GetScheme(), metav1.SchemeGroupVersion)

	if err = addControllers(mgr,
		kubernetesClient,
		informerFactory,
		s,
		ctx); err != nil {
		return fmt.Errorf("unable to register controllers to the manager: %v", err)
	}

	if err := mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		klog.Error(err, "unable to set up health check")
		os.Exit(1)
	}
	if err := mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		klog.Error(err, "unable to set up ready check")
		os.Exit(1)
	}

	// Start cache data after all informer is registered
	klog.V(0).Info("Starting cache resource from apiserver...")
	informerFactory.Start(ctx.Done())

	klog.V(0).Info("Starting the controllers.")
	if err = mgr.Start(ctx); err != nil {
		klog.Fatalf("unable to run the manager: %v", err)
	}

	return nil
}

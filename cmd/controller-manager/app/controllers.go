package app

import (
	"context"

	"github.com/kubesphere/whizard/cmd/controller-manager/app/options"
	"github.com/kubesphere/whizard/pkg/client/k8s"
	"github.com/kubesphere/whizard/pkg/controllers/monitoring"
	"github.com/kubesphere/whizard/pkg/informers"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

func addControllers(mgr manager.Manager, client k8s.Client, informerFactory informers.InformerFactory,
	cmOptions *options.ControllerManagerOptions, ctx context.Context) error {

	if err := (&monitoring.GatewayReconciler{
		Client:  mgr.GetClient(),
		Scheme:  mgr.GetScheme(),
		Context: ctx,
		Options: cmOptions.MonitoringOptions,
	}).SetupWithManager(mgr); err != nil {
		klog.Errorf("Unable to create Gateway controller: %v", err)
		return err
	}

	if err := (&monitoring.QueryFrontendReconciler{
		Client:  mgr.GetClient(),
		Scheme:  mgr.GetScheme(),
		Context: ctx,
		Options: cmOptions.MonitoringOptions,
	}).SetupWithManager(mgr); err != nil {
		klog.Errorf("Unable to create Query Frontend controller: %v", err)
		return err
	}

	if err := (&monitoring.QueryReconciler{
		Client:  mgr.GetClient(),
		Scheme:  mgr.GetScheme(),
		Context: ctx,
		Options: cmOptions.MonitoringOptions,
	}).SetupWithManager(mgr); err != nil {
		klog.Errorf("Unable to create Query controller: %v", err)
		return err
	}

	if err := (&monitoring.RouterReconciler{
		Client:  mgr.GetClient(),
		Scheme:  mgr.GetScheme(),
		Context: ctx,
		Options: cmOptions.MonitoringOptions,
	}).SetupWithManager(mgr); err != nil {
		klog.Errorf("Unable to create Router controller: %v", err)
		return err
	}

	if err := (&monitoring.StoreReconciler{
		DefaulterValidator: monitoring.CreateStoreDefaulterValidator(*cmOptions.MonitoringOptions),
		Client:             mgr.GetClient(),
		Scheme:             mgr.GetScheme(),
		Context:            ctx,
		Options:            cmOptions.MonitoringOptions.Store,
	}).SetupWithManager(mgr); err != nil {
		klog.Errorf("Unable to create Store controller: %v", err)
		return err
	}

	if err := (&monitoring.CompactorReconciler{
		DefaulterValidator: monitoring.CreateCompactorDefaulterValidator(*cmOptions.MonitoringOptions),
		Client:             mgr.GetClient(),
		Scheme:             mgr.GetScheme(),
		Context:            ctx,
		Options:            cmOptions.MonitoringOptions.Compactor,
	}).SetupWithManager(mgr); err != nil {
		klog.Errorf("Unable to create Compactor controller: %v", err)
		return err
	}

	if err := (&monitoring.IngesterReconciler{
		DefaulterValidator: monitoring.CreateIngesterDefaulterValidator(*cmOptions.MonitoringOptions),
		Client:             mgr.GetClient(),
		Scheme:             mgr.GetScheme(),
		Context:            ctx,
	}).SetupWithManager(mgr); err != nil {
		klog.Errorf("Unable to create Ingester controller: %v", err)
		return err
	}

	if err := (&monitoring.RulerReconciler{
		DefaulterValidator:    monitoring.CreateRulerDefaulterValidator(*cmOptions.MonitoringOptions),
		ReloaderConfig:        cmOptions.MonitoringOptions.PrometheusConfigReloader,
		RulerQueryProxyConfig: cmOptions.MonitoringOptions.RulerQueryProxy,
		Client:                mgr.GetClient(),
		Scheme:                mgr.GetScheme(),
		Context:               ctx,
	}).SetupWithManager(mgr); err != nil {
		klog.Errorf("Unable to create Ruler controller: %v", err)
		return err
	}

	if err := (&monitoring.TenantReconciler{
		DefaulterValidator: monitoring.CreateTenantDefaulterValidator(*cmOptions.MonitoringOptions),
		Client:             mgr.GetClient(),
		Scheme:             mgr.GetScheme(),
		Context:            ctx,
		Options:            cmOptions.MonitoringOptions,
	}).SetupWithManager(mgr); err != nil {
		klog.Errorf("Unable to create Tenant controller: %v", err)
		return err
	}

	if err := (&monitoring.StorageReconciler{
		Client:  mgr.GetClient(),
		Scheme:  mgr.GetScheme(),
		Context: ctx,
	}).SetupWithManager(mgr); err != nil {
		klog.Errorf("Unable to create Storage controller: %v", err)
		return err
	}

	if cmOptions.MonitoringOptions.EnableKubeSphereAdapter {
		if err := (&monitoring.ClusterReconciler{
			Client:                          mgr.GetClient(),
			Scheme:                          mgr.GetScheme(),
			Context:                         ctx,
			KubesphereAdapterDefaultService: cmOptions.MonitoringOptions.KubeSphereAdapterService,
		}).SetupWithManager(mgr); err != nil {
			klog.Errorf("Unable to create Cluster controller: %v", err)
			return err
		}
	}
	return nil
}

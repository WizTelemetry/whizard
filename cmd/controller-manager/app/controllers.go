package app

import (
	"context"

	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	"github.com/kubesphere/paodin/cmd/controller-manager/app/options"
	"github.com/kubesphere/paodin/pkg/client/k8s"
	"github.com/kubesphere/paodin/pkg/controllers/monitoring"
	"github.com/kubesphere/paodin/pkg/informers"
)

func addControllers(mgr manager.Manager, client k8s.Client, informerFactory informers.InformerFactory,
	cmOptions *options.PaodinControllerManagerOptions, ctx context.Context) error {

	if err := (&monitoring.ServiceReconciler{
		DefaulterValidator: monitoring.CreateServiceDefaulterValidator(*cmOptions.MonitoringOptions),
		Client:             mgr.GetClient(),
		Scheme:             mgr.GetScheme(),
		Context:            ctx,
	}).SetupWithManager(mgr); err != nil {
		klog.Errorf("Unable to create Service controller: %v", err)
		return err
	}

	if err := (&monitoring.StoreReconciler{
		DefaulterValidator: monitoring.CreateStoreDefaulterValidator(*cmOptions.MonitoringOptions),
		Client:             mgr.GetClient(),
		Scheme:             mgr.GetScheme(),
		Context:            ctx,
	}).SetupWithManager(mgr); err != nil {
		klog.Errorf("Unable to create Store controller: %v", err)
		return err
	}

	if err := (&monitoring.ThanosReceiveIngestorReconciler{
		DefaulterValidator: monitoring.CreateThanosReceiveIngestorDefaulterValidator(*cmOptions.MonitoringOptions),
		Client:             mgr.GetClient(),
		Scheme:             mgr.GetScheme(),
		Context:            ctx,
	}).SetupWithManager(mgr); err != nil {
		klog.Errorf("Unable to create ThanosReceiveIngestor controller: %v", err)
		return err
	}

	if err := (&monitoring.ThanosRulerReconciler{
		DefaulterValidator: monitoring.CreateThanosRulerDefaulterValidator(*cmOptions.MonitoringOptions),
		ReloaderConfig:     cmOptions.MonitoringOptions.PrometheusConfigReloader,
		Client:             mgr.GetClient(),
		Scheme:             mgr.GetScheme(),
		Context:            ctx,
	}).SetupWithManager(mgr); err != nil {
		klog.Errorf("Unable to create ThanosRuler controller: %v", err)
		return err
	}

	return nil
}

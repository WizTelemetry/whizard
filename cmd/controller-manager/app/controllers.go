package app

import (
	"context"

	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	"github.com/WhizardTelemetry/whizard/pkg/client/k8s"
	controllers "github.com/WhizardTelemetry/whizard/pkg/controllers"
)

func addControllers(mgr manager.Manager, client k8s.Client, ctx context.Context) error {

	if err := (&controllers.GatewayReconciler{
		Client:  mgr.GetClient(),
		Scheme:  mgr.GetScheme(),
		Context: ctx,
	}).SetupWithManager(mgr); err != nil {
		klog.Errorf("Unable to create Gateway controller: %v", err)
		return err
	}

	if err := (&controllers.QueryFrontendReconciler{
		Client:  mgr.GetClient(),
		Scheme:  mgr.GetScheme(),
		Context: ctx,
	}).SetupWithManager(mgr); err != nil {
		klog.Errorf("Unable to create Query Frontend controller: %v", err)
		return err
	}

	if err := (&controllers.QueryReconciler{
		Client:  mgr.GetClient(),
		Scheme:  mgr.GetScheme(),
		Context: ctx,
	}).SetupWithManager(mgr); err != nil {
		klog.Errorf("Unable to create Query controller: %v", err)
		return err
	}

	if err := (&controllers.RouterReconciler{
		Client:  mgr.GetClient(),
		Scheme:  mgr.GetScheme(),
		Context: ctx,
	}).SetupWithManager(mgr); err != nil {
		klog.Errorf("Unable to create Router controller: %v", err)
		return err
	}

	if err := (&controllers.StoreReconciler{
		Client:  mgr.GetClient(),
		Scheme:  mgr.GetScheme(),
		Context: ctx,
	}).SetupWithManager(mgr); err != nil {
		klog.Errorf("Unable to create Store controller: %v", err)
		return err
	}

	if err := (&controllers.CompactorReconciler{
		Client:  mgr.GetClient(),
		Scheme:  mgr.GetScheme(),
		Context: ctx,
	}).SetupWithManager(mgr); err != nil {
		klog.Errorf("Unable to create Compactor controller: %v", err)
		return err
	}

	if err := (&controllers.IngesterReconciler{
		Client:  mgr.GetClient(),
		Scheme:  mgr.GetScheme(),
		Context: ctx,
	}).SetupWithManager(mgr); err != nil {
		klog.Errorf("Unable to create Ingester controller: %v", err)
		return err
	}

	if err := (&controllers.RulerReconciler{
		Client:  mgr.GetClient(),
		Scheme:  mgr.GetScheme(),
		Context: ctx,
	}).SetupWithManager(mgr); err != nil {
		klog.Errorf("Unable to create Ruler controller: %v", err)
		return err
	}

	if err := (&controllers.TenantReconciler{
		Client:  mgr.GetClient(),
		Scheme:  mgr.GetScheme(),
		Context: ctx,
	}).SetupWithManager(mgr); err != nil {
		klog.Errorf("Unable to create Tenant controller: %v", err)
		return err
	}

	if err := (&controllers.StorageReconciler{
		Client:  mgr.GetClient(),
		Scheme:  mgr.GetScheme(),
		Context: ctx,
	}).SetupWithManager(mgr); err != nil {
		klog.Errorf("Unable to create Storage controller: %v", err)
		return err
	}

	return nil
}

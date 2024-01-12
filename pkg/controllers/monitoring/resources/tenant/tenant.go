package tenant

import (
	monitoringv1alpha1 "github.com/kubesphere/whizard/pkg/api/monitoring/v1alpha1"
	"github.com/kubesphere/whizard/pkg/controllers/monitoring/resources"
)

type Tenant struct {
	tenant *monitoringv1alpha1.Tenant
	resources.BaseReconciler
}

func New(reconciler resources.BaseReconciler, tenant *monitoringv1alpha1.Tenant) (*Tenant, error) {
	if err := reconciler.SetService(tenant); err != nil {
		return nil, err
	}
	return &Tenant{
		tenant:         tenant,
		BaseReconciler: reconciler,
	}, nil
}

func (t *Tenant) Reconcile() error {
	if err := t.ingester(); err != nil {
		return err
	}
	if err := t.ruler(); err != nil {
		return err
	}
	if err := t.compactor(); err != nil {
		return err
	}
	if err := t.store(); err != nil {
		return err
	}
	return nil
}

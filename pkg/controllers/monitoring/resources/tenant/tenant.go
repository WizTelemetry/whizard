package tenant

import (
	monitoringv1alpha1 "github.com/kubesphere/paodin/pkg/api/monitoring/v1alpha1"
	"github.com/kubesphere/paodin/pkg/controllers/monitoring/options"
	"github.com/kubesphere/paodin/pkg/controllers/monitoring/resources"
)

type Tenant struct {
	tenant *monitoringv1alpha1.Tenant
	resources.BaseReconciler

	Options *options.Options
}

func New(reconciler resources.BaseReconciler, tenant *monitoringv1alpha1.Tenant, o *options.Options) *Tenant {
	return &Tenant{
		tenant:         tenant,
		BaseReconciler: reconciler,
		Options:        o,
	}
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

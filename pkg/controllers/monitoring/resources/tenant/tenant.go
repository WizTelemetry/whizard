package tenant

import (
	monitoringv1alpha1 "github.com/kubesphere/paodin/pkg/api/monitoring/v1alpha1"
	"github.com/kubesphere/paodin/pkg/controllers/monitoring/resources"
)

type Tenant struct {
	tenant *monitoringv1alpha1.Tenant
	resources.BaseReconciler
}

func New(reconciler resources.BaseReconciler, tenant *monitoringv1alpha1.Tenant) *Tenant {
	return &Tenant{
		tenant:         tenant,
		BaseReconciler: reconciler,
	}
}

func (t *Tenant) Reconcile() error {
	return t.ReconcileResources([]resources.Resource{
		t.ruler,
		t.receiveIngestor,
	})
}

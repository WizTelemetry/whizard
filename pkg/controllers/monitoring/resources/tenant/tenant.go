package tenant

import (
	"time"

	monitoringv1alpha1 "github.com/kubesphere/paodin/pkg/api/monitoring/v1alpha1"
	"github.com/kubesphere/paodin/pkg/controllers/monitoring/resources"
)

type Tenant struct {
	tenant *monitoringv1alpha1.Tenant
	resources.BaseReconciler

	DefaultTenantsPerIngester      int
	DefaultIngesterRetentionPeriod time.Duration
}

func New(reconciler resources.BaseReconciler, tenant *monitoringv1alpha1.Tenant, DefaultTenantsPerIngester int, defaultIngesterRetentionPeriod time.Duration) *Tenant {
	return &Tenant{
		tenant:                         tenant,
		BaseReconciler:                 reconciler,
		DefaultTenantsPerIngester:      DefaultTenantsPerIngester,
		DefaultIngesterRetentionPeriod: defaultIngesterRetentionPeriod,
	}
}

func (t *Tenant) Reconcile() error {
	if err := t.ingester(); err != nil {
		return err
	}
	if err := t.ruler(); err != nil {
		return err
	}
	return nil
}

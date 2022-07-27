package tenant

import (
	"time"

	monitoringv1alpha1 "github.com/kubesphere/paodin/pkg/api/monitoring/v1alpha1"
	"github.com/kubesphere/paodin/pkg/controllers/monitoring/resources"
)

type Tenant struct {
	tenant *monitoringv1alpha1.Tenant
	resources.BaseReconciler

	DefaultTenantsPerIngestor      int
	DefaultIngestorRetentionPeriod time.Duration
}

func New(reconciler resources.BaseReconciler, tenant *monitoringv1alpha1.Tenant, DefaultTenantsPerIngestor int, defaultIngestorRetentionPeriod time.Duration) *Tenant {
	return &Tenant{
		tenant:                         tenant,
		BaseReconciler:                 reconciler,
		DefaultTenantsPerIngestor:      DefaultTenantsPerIngestor,
		DefaultIngestorRetentionPeriod: defaultIngestorRetentionPeriod,
	}
}

func (t *Tenant) Reconcile() error {
	if err := t.receiveIngestor(); err != nil {
		return err
	}
	if err := t.ruler(); err != nil {
		return err
	}
	return nil
}

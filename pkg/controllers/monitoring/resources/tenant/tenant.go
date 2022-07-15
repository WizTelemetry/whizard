package tenant

import (
	"time"

	monitoringv1alpha1 "github.com/kubesphere/paodin/pkg/api/monitoring/v1alpha1"
	"github.com/kubesphere/paodin/pkg/controllers/monitoring/resources"
)

type Tenant struct {
	tenant *monitoringv1alpha1.Tenant
	resources.BaseReconciler

	DefaultTenantCountPerIngestor  int
	DefaultIngestorRetentionPeriod time.Duration
	DeleteIngestorEventChan        chan DeleteIngestorEvent
}

func New(reconciler resources.BaseReconciler, tenant *monitoringv1alpha1.Tenant, defaultTenantCountPerIngestor int, defaultIngestorRetentionPeriod time.Duration, deleteIngestorEventChan chan DeleteIngestorEvent) *Tenant {
	return &Tenant{
		tenant:                         tenant,
		BaseReconciler:                 reconciler,
		DefaultTenantCountPerIngestor:  defaultTenantCountPerIngestor,
		DefaultIngestorRetentionPeriod: defaultIngestorRetentionPeriod,
		DeleteIngestorEventChan:        deleteIngestorEventChan,
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

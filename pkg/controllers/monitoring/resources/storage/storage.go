package storage

import (
	monitoringv1alpha1 "github.com/kubesphere/whizard/pkg/api/monitoring/v1alpha1"
	"github.com/kubesphere/whizard/pkg/controllers/monitoring/resources"
)

type Storage struct {
	storage *monitoringv1alpha1.Storage
	resources.BaseReconciler
}

func New(reconciler resources.BaseReconciler, storage *monitoringv1alpha1.Storage) *Storage {
	return &Storage{
		storage:        storage,
		BaseReconciler: reconciler,
	}
}

func (s *Storage) Reconcile() error {
	return s.ReconcileResources([]resources.Resource{
		s.updateHashAnnotation,
	})
}

func (s *Storage) String() (string, error) {
	body, err := s.parseObjStoreConfig()
	return string(body), err
}

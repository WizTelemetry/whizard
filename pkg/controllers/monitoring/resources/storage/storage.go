package storage

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	monitoringv1alpha1 "github.com/kubesphere/whizard/pkg/api/monitoring/v1alpha1"
	"github.com/kubesphere/whizard/pkg/constants"
	"github.com/kubesphere/whizard/pkg/controllers/monitoring/resources"
	"github.com/kubesphere/whizard/pkg/util"
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

func (s *Storage) labels() map[string]string {
	labels := s.BaseLabels()
	labels[constants.LabelNameAppName] = constants.AppNameStorage
	labels[constants.LabelNameAppManagedBy] = s.storage.Name
	return labels
}

func (s *Storage) name(nameSuffix ...string) string {
	return s.QualifiedName(constants.AppNameBlockManager, s.storage.Name, nameSuffix...)
}

func (s *Storage) meta(name string) metav1.ObjectMeta {

	return metav1.ObjectMeta{
		Name:      name,
		Namespace: s.storage.Namespace,
		Labels:    s.labels(),
	}
}

// updateHashAnnotation generate the hash with objstoreConfig to write to the annotation. When the secret update hash changes, the storage update event is triggered.
func (s *Storage) updateHashAnnotation() (runtime.Object, resources.Operation, error) {

	if s.storage.Annotations == nil {
		s.storage.Annotations = make(map[string]string)
	}

	hashStr, err := s.GetStorageHash(util.Join(".", s.storage.Namespace, s.storage.Name))
	if err != nil {
		return nil, "", err
	}

	if v, ok := s.storage.Annotations[constants.LabelNameStorageHash]; !ok || v != hashStr {
		s.storage.Annotations[constants.LabelNameStorageHash] = hashStr
		return s.storage, resources.OperationCreateOrUpdate, nil
	}

	return nil, "", nil
}

func (s *Storage) Reconcile() error {
	return s.ReconcileResources([]resources.Resource{
		s.updateHashAnnotation,
		s.deployment,
		s.service,
	})
}

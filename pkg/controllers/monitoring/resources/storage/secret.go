package storage

import (
	"github.com/kubesphere/whizard/pkg/constants"
	"github.com/kubesphere/whizard/pkg/controllers/monitoring/resources"
	"github.com/kubesphere/whizard/pkg/util"
	"k8s.io/apimachinery/pkg/runtime"
)

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

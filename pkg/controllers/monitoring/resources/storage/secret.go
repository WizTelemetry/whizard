package storage

import (
	yaml "gopkg.in/yaml.v3"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/utils/pointer"

	monitoringv1alpha1 "github.com/kubesphere/paodin/pkg/api/monitoring/v1alpha1"
	"github.com/kubesphere/paodin/pkg/controllers/monitoring/resources"
)

func (s *Storage) secret() (runtime.Object, resources.Operation, error) {
	buff, err := parseObjStorageConfig(s.storage.Spec.Thanos)
	if err != nil {
		return nil, "", err
	}
	ls := make(map[string]string, 1)
	ls[monitoringv1alpha1.MonitoringPaodinStorage] = s.storage.Namespace + "." + s.storage.Name
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: s.storage.Namespace,
			Name:      "secret-" + s.storage.Name,
			Labels:    ls,
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion: s.storage.APIVersion,
					Kind:       s.storage.Kind,
					Name:       s.storage.Name,
					UID:        s.storage.UID,
					Controller: pointer.BoolPtr(true),
				},
			},
		},
		Type: corev1.SecretTypeOpaque,
		Data: map[string][]byte{
			resources.SecretThanosBucketKey: buff,
		},
	}

	return secret, resources.OperationCreateOrUpdate, nil
}

func parseObjStorageConfig(thanosStorageConfig *monitoringv1alpha1.ThanosStorage) ([]byte, error) {
	bucket := &BucketConfig{}

	if thanosStorageConfig.S3 != nil {
		bucket.Type = S3
		bucket.Config = *thanosStorageConfig.S3
	}

	return yaml.Marshal(bucket)
}

type ObjectStorageProvider string

const (
	FILESYSTEM ObjectStorageProvider = "FILESYSTEM"
	GCS        ObjectStorageProvider = "GCS"
	S3         ObjectStorageProvider = "S3"
	AZURE      ObjectStorageProvider = "AZURE"
	SWIFT      ObjectStorageProvider = "SWIFT"
	COS        ObjectStorageProvider = "COS"
	ALIYUNOSS  ObjectStorageProvider = "ALIYUNOSS"
	BOS        ObjectStorageProvider = "BOS"
)

type BucketConfig struct {
	Type   ObjectStorageProvider `yaml:"type"`
	Config interface{}           `yaml:"config"`
}

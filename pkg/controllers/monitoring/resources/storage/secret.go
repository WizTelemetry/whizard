package storage

import (
	"crypto/md5"
	"encoding/hex"

	yaml "gopkg.in/yaml.v3"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"

	"github.com/kubesphere/whizard/pkg/controllers/monitoring/resources"
)

var LabelNameStorageHash = "monitoring.whizard.io/storage-hash"

// updateHashAnnotation generate the hash with objstoreConfig to write to the annotation. When the secret update hash changes, the storage update event is triggered.
func (s *Storage) updateHashAnnotation() (runtime.Object, resources.Operation, error) {

	objStoreConfig, err := s.parseObjStoreConfig()
	if err != nil {
		return nil, "", err
	}
	c := md5.New()
	c.Write(objStoreConfig)
	hashStr := hex.EncodeToString(c.Sum(nil))

	if s.storage.Annotations == nil {
		s.storage.Annotations = make(map[string]string)
	}
	if v, ok := s.storage.Annotations[LabelNameStorageHash]; !ok || v != hashStr {
		s.storage.Annotations[LabelNameStorageHash] = hashStr
		return s.storage, resources.OperationCreateOrUpdate, nil
	}

	return nil, "", nil
}

func (s *Storage) parseObjStoreConfig() ([]byte, error) {
	bucket := &BucketConfig{}

	if s.storage.Spec.S3 != nil {
		akSecert := &corev1.Secret{}
		if err := s.Client.Get(s.Context, types.NamespacedName{
			Name:      s.storage.Spec.S3.AccessKeySecretRef.Name,
			Namespace: s.storage.Namespace,
		}, akSecert); err != nil {
			return nil, err
		}
		s.storage.Spec.S3.AccessKey = string(akSecert.Data[s.storage.Spec.S3.AccessKeySecretRef.Key])

		skSecret := &corev1.Secret{}
		if err := s.Client.Get(s.Context, types.NamespacedName{
			Name:      s.storage.Spec.S3.SecretKeySecretRef.Name,
			Namespace: s.storage.Namespace,
		}, skSecret); err != nil {
			return nil, err
		}
		s.storage.Spec.S3.SecretKey = string(skSecret.Data[s.storage.Spec.S3.SecretKeySecretRef.Key])

		bucket.Type = S3
		bucket.Config = *s.storage.Spec.S3
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

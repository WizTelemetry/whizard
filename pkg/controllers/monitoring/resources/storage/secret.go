package storage

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"

	"github.com/kubesphere/whizard/pkg/constants"
	"github.com/kubesphere/whizard/pkg/controllers/monitoring/resources"
	"gopkg.in/yaml.v3"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
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

	if s.storage.Spec.S3 != nil {
		b := &BucketConfig{
			S3,
			*s.storage.Spec.S3,
		}

		root := &yaml.Node{}
		bs, err := yaml.Marshal(b)
		if err != nil {
			return nil, err
		}

		if err := yaml.Unmarshal(bs, root); err != nil {
			return nil, err
		}

		akSecert := &corev1.Secret{}
		if err := s.Client.Get(s.Context, types.NamespacedName{
			Name:      s.storage.Spec.S3.AccessKey.Name,
			Namespace: s.storage.Namespace,
		}, akSecert); err != nil {
			return nil, err
		}
		if n := findNodeByKey(root, "access_key"); n != nil {
			n.SetString(string(akSecert.Data[s.storage.Spec.S3.AccessKey.Key]))
		}

		skSecret := &corev1.Secret{}
		if err := s.Client.Get(s.Context, types.NamespacedName{
			Name:      s.storage.Spec.S3.SecretKey.Name,
			Namespace: s.storage.Namespace,
		}, skSecret); err != nil {
			return nil, err
		}
		if n := findNodeByKey(root, "secret_key"); n != nil {
			n.SetString(string(skSecret.Data[s.storage.Spec.S3.SecretKey.Key]))
		}

		if ref := s.storage.Spec.S3.HTTPConfig.TLSConfig.CA; ref != nil {
			if n := findNodeByKey(root, "ca_file"); n != nil {
				n.SetString(fmt.Sprintf("%s%s/%s", constants.ConfigPath, ref.Name, ref.Key))
			}
		}

		if ref := s.storage.Spec.S3.HTTPConfig.TLSConfig.Cert; ref != nil {
			if n := findNodeByKey(root, "cert_file"); n != nil {
				n.SetString(fmt.Sprintf("%s%s/%s", constants.ConfigPath, ref.Name, ref.Key))
			}
		}

		if ref := s.storage.Spec.S3.HTTPConfig.TLSConfig.Key; ref != nil {
			if n := findNodeByKey(root, "key_file"); n != nil {
				n.SetString(fmt.Sprintf("%s%s/%s", constants.ConfigPath, ref.Name, ref.Key))
			}
		}

		return yaml.Marshal(root)
	}

	return nil, nil
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

func findNodeByKey(root *yaml.Node, key string) *yaml.Node {

	for i := 0; i < len(root.Content); i++ {
		if root.Content[i].Value == key && i+1 < len(root.Content) {
			return root.Content[i+1]
		}

		if n := findNodeByKey(root.Content[i], key); n != nil {
			return n
		}
	}
	return nil
}

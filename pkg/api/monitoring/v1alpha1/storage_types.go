/*
Copyright 2021 The KubeSphere authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	"strings"

	"github.com/prometheus/common/model"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

const (
	MonitoringPaodinService = "monitoring.paodin.io/service"
	MonitoringPaodinStorage = "monitoring.paodin.io/storage"
	MonitoringPaodinTenant  = "monitoring.paodin.io/tenant"

	FinalizerMonitoringPaodin = "finalizers.monitoring.paodin.io"
)

type StorageSpec struct {
	Thanos *ThanosStorage `json:"thanos"`
}

type ThanosStorage struct {
	S3 *S3 `json:"S3,omitempty"`
}

// Config stores the configuration for s3 bucket.
type S3 struct {
	Bucket             string            `yaml:"bucket,omitempty" json:"bucket"`
	Endpoint           string            `yaml:"endpoint,omitempty" json:"endpoint"`
	Region             string            `yaml:"region,omitempty" json:"region,omitempty"`
	AWSSDKAuth         bool              `yaml:"aws_sdk_auth,omitempty" json:"awsSdkAuth,omitempty"`
	AccessKey          string            `yaml:"access_key,omitempty" json:"accessKey"`
	Insecure           bool              `yaml:"insecure,omitempty" json:"insecure,omitempty"`
	SignatureV2        bool              `yaml:"signature_version2,omitempty" json:"signatureVersion2,omitempty"`
	SecretKey          string            `yaml:"secret_key,omitempty" json:"secretKey"`
	PutUserMetadata    map[string]string `yaml:"put_user_metadata,omitempty" json:"putUserMetadata,omitempty"`
	HTTPConfig         S3HTTPConfig      `yaml:"http_config,omitempty" json:"httpConfig,omitempty"`
	TraceConfig        S3TraceConfig     `yaml:"trace,omitempty" json:"trace,omitempty"`
	ListObjectsVersion string            `yaml:"list_objects_version,omitempty" json:"listObjectsVersion,omitempty"`
	// PartSize used for multipart upload. Only used if uploaded object size is known and larger than configured PartSize.
	// NOTE we need to make sure this number does not produce more parts than 10 000.
	PartSize    uint64      `yaml:"part_size,omitempty" json:"partSize,omitempty"`
	SSEConfig   S3SSEConfig `yaml:"sse_config,omitempty" json:"sseConfig,omitempty"`
	STSEndpoint string      `yaml:"sts_endpoint,omitempty" json:"stsEndpoint,omitempty"`
}

// S3SSEConfig deals with the configuration of SSE for Minio. The following options are valid:
// kmsencryptioncontext == https://docs.aws.amazon.com/kms/latest/developerguide/services-s3.html#s3-encryption-context
type S3SSEConfig struct {
	Type                 string            `yaml:"type,omitempty" json:"type,omitempty"`
	KMSKeyID             string            `yaml:"kms_key_id,omitempty" json:"kmsKeyId,omitempty"`
	KMSEncryptionContext map[string]string `yaml:"kms_encryption_context,omitempty" json:"kmsEncryptionContext,omitempty"`
	EncryptionKey        string            `yaml:"encryption_key,omitempty" json:"encryptionKey,omitempty"`
}

type S3TraceConfig struct {
	Enable bool `yaml:"enable,omitempty" json:"enable,omitempty"`
}

// S3HTTPConfig stores the http.Transport configuration for the s3 minio client.
type S3HTTPConfig struct {
	IdleConnTimeout       model.Duration `yaml:"idle_conn_timeout,omitempty" json:"idleConnTimeout,omitempty"`
	ResponseHeaderTimeout model.Duration `yaml:"response_header_timeout,omitempty" json:"responseHeaderTimeout,omitempty"`
	InsecureSkipVerify    bool           `yaml:"insecure_skip_verify,omitempty" json:"insecureSkipVerify,omitempty"`

	TLSHandshakeTimeout   model.Duration `yaml:"tls_handshake_timeout,omitempty" json:"tlsHandshakeTimeout,omitempty"`
	ExpectContinueTimeout model.Duration `yaml:"expect_continue_timeout,omitempty" json:"expectContinueTimeout,omitempty"`
	MaxIdleConns          int            `yaml:"max_idle_conns,omitempty" json:"maxIdleConns,omitempty"`
	MaxIdleConnsPerHost   int            `yaml:"max_idle_conns_per_host,omitempty" json:"maxIdleConnsPerHost,omitempty"`
	MaxConnsPerHost       int            `yaml:"max_conns_per_host,omitempty" json:"maxConnsPerHost,omitempty"`

	TLSConfig TLSConfig `yaml:"tls_config,omitempty" json:"tlsConfig,omitempty"`
}

// TLSConfig configures the options for TLS connections.
type TLSConfig struct {
	// The CA cert to use for the targets.
	CAFile string `yaml:"ca_file,omitempty" json:"caFile,omitempty"`
	// The client cert file for the targets.
	CertFile string `yaml:"cert_file,omitempty" json:"certFile,omitempty"`
	// The client key file for the targets.
	KeyFile string `yaml:"key_file,omitempty" json:"keyFile,omitempty"`
	// Used to verify the hostname for the targets.
	ServerName string `yaml:"server_name,omitempty" json:"serverName,omitempty"`
	// Disable target certificate validation.
	InsecureSkipVerify bool `yaml:"insecure_skip_verify,omitempty" json:"insecureSkipVerify,omitempty"`
}

type StorageStatus struct {
}

//+genclient
//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

type Storage struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   StorageSpec   `json:"spec,omitempty"`
	Status StorageStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

type StorageList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Storage `json:"items"`
}

func ManagedLabelByStorage(storage metav1.Object) map[string]string {
	return map[string]string{
		MonitoringPaodinStorage: storage.GetNamespace() + "." + storage.GetName(),
	}
}

func StorageNamespacedName(managedByStorage metav1.Object) *types.NamespacedName {
	ls := managedByStorage.GetLabels()
	if len(ls) == 0 {
		return nil
	}

	namespacedName := ls[MonitoringPaodinStorage]
	arr := strings.Split(namespacedName, ".")
	if len(arr) != 2 {
		return nil
	}

	return &types.NamespacedName{
		Namespace: arr[0],
		Name:      arr[1],
	}
}

func init() {
	SchemeBuilder = SchemeBuilder.
		Register(&Storage{}, &StorageList{})
}

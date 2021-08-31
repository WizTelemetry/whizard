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
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ThanosReceiveSpec defines the desired state of ThanosReceive
type ThanosReceiveSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// If specified, the pod's scheduling constraints.
	Affinity *corev1.Affinity `json:"affinity,omitempty"`
	// Define which Nodes the Pods are scheduled on.
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`
	// If specified, the pod's tolerations.
	Tolerations []corev1.Toleration `json:"tolerations,omitempty"`
	// Define resources requests and limits for single Pods.
	Resources corev1.ResourceRequirements `json:"resources,omitempty"`

	// Image is the thanos image with tag/version
	Image           string            `json:"image,omitempty"`
	ImagePullPolicy corev1.PullPolicy `json:"imagePullPolicy,omitempty"`
	// Number of replicas for a thanos receive component
	Replicas        *int32            `json:"replicas,omitempty"`

	// Router specifies the configs that thanos receive running in the router mode requires
	Router   *ReceiveRouterSpec   `json:"router,omitempty"`
	// Ingestor specifies the configs that thanos receive running in the ingestor mode requires
	Ingestor *ReceiveIngestorSpec `json:"ingestor,omitempty"`

	// TenantHeader configs the HTTP header specifying the replica number of a write request to thanos receive
	TenantHeader    string `json:"tenantHeader,omitempty"`
	// DefaultTenantId configs the default tenant ID to use when none is provided via a header
	DefaultTenantId string `json:"defaultTenantId,omitempty"`
	// TenantLabelName configs the label name through which the tenant will be announced.
	TenantLabelName string `json:"tenantLabelName,omitempty"`

	// LogLevel configs log filtering level. Possible options: error, warn, info, debug
	LogLevel  string `json:"level,omitempty"`
	// LogFormat configs log format to use. Possible options: logfmt or json
	LogFormat string `json:"format,omitempty"`
}

// ThanosReceiveStatus defines the observed state of ThanosReceive
type ThanosReceiveStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// ThanosReceive is the Schema for the thanosreceives API
type ThanosReceive struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ThanosReceiveSpec   `json:"spec,omitempty"`
	Status ThanosReceiveStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ThanosReceiveList contains a list of ThanosReceive
type ThanosReceiveList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ThanosReceive `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ThanosReceive{}, &ThanosReceiveList{})
}

// ReceiveRouterSpec defines the configs that thanos receive running in the router mode requires
type ReceiveRouterSpec struct {
	// HashringsRefreshInterval configs refresh interval to re-read the hashring configuration file
	HashringsRefreshInterval string `json:"hashringsRefreshInterval,omitempty"`

	// HardTenantHashrings are hashrings with non-empty tenants which match the tenant in the request
	HardTenantHashrings []*RouterHashringConfig `json:"hardTenantHashrings,omitempty"`
	// SoftTenantHashring is a hashring with empty tenants which is used when the tenant in the request
	// cannot be found in HardTenantHashrings.
	SoftTenantHashring *RouterHashringConfig `json:"softTenantHashring,omitempty"`

	ReplicationFactor *uint64 `json:"replicationFactor,omitempty"`

	// RemoteWriteIngress configs remote write request entry from services outside the cluster
	RemoteWriteIngress *IngressSpec            `json:"remoteWriteIngress,omitempty"`
}

// RouterHashringConfig defines the hashring config for a team of tenants
type RouterHashringConfig struct {
	Name    string   `json:"name,omitempty"`

	// Tenants if not empty indicates current config is for hard tenants; otherwise, it is for soft tenants.
	Tenants []string `json:"tenants,omitempty"`

	// Endpoints are statically configured endpoints which receive requests for the specified tenants
	Endpoints []string `json:"endpoints,omitempty"`

	// EndpointsNamespaceSelector and EndpointsSelector select endpoints which receive requests for the specified tenants
	// They only work when Endpoints is empty
	EndpointsNamespaceSelector *metav1.LabelSelector `json:"endpointsNamespaceSelector,omitempty"`
	EndpointsSelector          *metav1.LabelSelector `json:"endpointsSelector,omitempty"`
}

// ReceiveIngestorSpec defines the configs that thanos receive running in the ingestor mode requires
type ReceiveIngestorSpec struct {
	// LocalTSDBRetention configs how long to retain raw samples on local storage
	LocalTSDBRetention  string                    `json:"localTsdbRetention,omitempty"`
	// ObjectStorageConfig allows specifying a key of a Secret containing object store configuration
	ObjectStorageConfig *corev1.SecretKeySelector `json:"objectStorageConfig,omitempty"`
	// DataVolume specifies how volume shall be used
	DataVolume          *KubernetesVolume         `json:"dataVolume,omitempty"`
}

// KubernetesVolume defines the configured volume for a thanos receiver.
type KubernetesVolume struct {
	EmptyDir              *corev1.EmptyDirVolumeSource  `json:"emptyDir,omitempty"`
	PersistentVolumeClaim *corev1.PersistentVolumeClaim `json:"pvc,omitempty"`
}

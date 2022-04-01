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

// ServiceSpec defines the desired state of a Service
type ServiceSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// HTTP header to determine tenant for remote write requests.
	TenantHeader string `json:"tenantHeader,omitempty"`
	// Default tenant ID to use when none is provided via a header.
	DefaultTenantId string `json:"defaultTenantId,omitempty"`
	// Label name through which the tenant will be announced.
	TenantLabelName string `json:"tenantLabelName,omitempty"`

	// Define one Gateway instance to proxy and auth requests to thanos.
	Gateway *Gateway `json:"gateway,omitempty"`

	// Define one Thanos cluster.
	Thanos *Thanos `json:"thanos,omitempty"`
}

type Gateway struct {
	// Image is the gateway image with tag/version.
	Image string `json:"image,omitempty"`
	// Log filtering level. Possible options: error, warn, info, debug.
	LogLevel string `json:"logLevel,omitempty"`
	// Log format to use. Possible options: logfmt or json.
	LogFormat string `json:"logFormat,omitempty"`

	// If specified, the pod's scheduling constraints.
	Affinity *corev1.Affinity `json:"affinity,omitempty"`
	// Define which Nodes the Pods are scheduled on.
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`
	// If specified, the pod's tolerations.
	Tolerations []corev1.Toleration `json:"tolerations,omitempty"`
	// Define resources requests and limits for main container.
	Resources corev1.ResourceRequirements `json:"resources,omitempty"`
	// Number of replicas for a thanos component
	Replicas *int32 `json:"replicas,omitempty"`

	// Secret name for HTTP Server certificate (Kubernetes TLS secret type)
	ServerCertificate string `json:"serverCertificate,omitempty"`
	// Secret name for HTTP Client CA certificate (Kubernetes TLS secret type)
	ClientCACertificate string `json:"clientCaCertificate,omitempty"`
}

type Thanos struct {
	DefaultFields CommonThanosFields `json:"defaultFields,omitempty"`

	// ObjectStorageConfig allows specifying a key of a Secret containing object store configuration
	ObjectStorageConfig *corev1.SecretKeySelector `json:"objectStorageConfig,omitempty"`

	Query *Query `json:"query,omitempty"`

	Receive *Receive `json:"receive,omitempty"`

	StoreGateway *StoreGateway `json:"storeGateway,omitempty"`

	Compact *Compact `json:"compact,omitempty"`
}

type CommonThanosFields struct {
	// Image is the thanos image with tag/version
	Image string `json:"image,omitempty"`
	// Log filtering level. Possible options: error, warn, info, debug
	LogLevel string `json:"logLevel,omitempty"`
	// Log format to use. Possible options: logfmt or json
	LogFormat string `json:"logFormat,omitempty"`
}

type Query struct {
	CommonThanosFields `json:",inline"`

	// If specified, the pod's scheduling constraints.
	Affinity *corev1.Affinity `json:"affinity,omitempty"`
	// Define which Nodes the Pods are scheduled on.
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`
	// If specified, the pod's tolerations.
	Tolerations []corev1.Toleration `json:"tolerations,omitempty"`
	// Define resources requests and limits for main container.
	Resources corev1.ResourceRequirements `json:"resources,omitempty"`
	// Number of replicas for a thanos component
	Replicas *int32 `json:"replicas,omitempty"`

	// Additional StoreApi servers from which Thanos Query component queries from
	Stores []QueryStores `json:"stores,omitempty"`
	// Selector labels that will be exposed in info endpoint.
	SelectorLabels map[string]string `json:"selectorLabels,omitempty"`

	// Envoy is used to config sidecar which proxies requests requiring auth to the secure stores
	Envoy EnvoySpec `json:"envoy,omitempty"`
}

type QueryStores struct {
	// Address is the addresses of StoreApi server, which may be prefixed with 'dns+' or 'dnssrv+' to detect StoreAPI servers through respective DNS lookups.
	// For more info, see https://thanos.io/tip/thanos/service-discovery.md/#dns-service-discovery
	Addresses []string `json:"addresses,omitempty"`
	// Secret containing the CA cert to use for StoreApi connections
	CASecret *corev1.SecretKeySelector `json:"caSecret,omitempty"`
}

// EnvoySpec defines the desired state of envoy proxy sidecar which delegates requests to the secure thanos stores
type EnvoySpec struct {
	// Image is the thanos image with tag/version
	Image string `json:"image,omitempty"`
	// Define resources requests and limits for envoy container.
	Resources corev1.ResourceRequirements `json:"resources,omitempty"`
}

type Receive struct {
	Router    ReceiveRouter     `json:"router,omitempty"`
	Ingestors []ReceiveIngestor `json:"ingestors,omitempty"`
}

type ReceiveRouter struct {
	CommonThanosFields `json:",inline"`

	// If specified, the pod's scheduling constraints.
	Affinity *corev1.Affinity `json:"affinity,omitempty"`
	// Define which Nodes the Pods are scheduled on.
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`
	// If specified, the pod's tolerations.
	Tolerations []corev1.Toleration `json:"tolerations,omitempty"`
	// Define resources requests and limits for main container.
	Resources corev1.ResourceRequirements `json:"resources,omitempty"`
	// Number of replicas for a thanos component.
	Replicas *int32 `json:"replicas,omitempty"`

	// How many times to replicate incoming write requests
	ReplicationFactor *uint64 `json:"replicationFactor,omitempty"`
}

type ReceiveIngestor struct {
	CommonThanosFields `json:",inline"`

	// Ingestor name must be unique within current thanos cluster, which follows the regulation for k8s resource name.
	Name string `json:"name,omitempty"`
	// Tenants if not empty indicates current config is for hard tenants; otherwise, it is for soft tenants.
	Tenants []string `json:"tenants,omitempty"`

	// If specified, the pod's scheduling constraints.
	Affinity *corev1.Affinity `json:"affinity,omitempty"`
	// Define which Nodes the Pods are scheduled on.
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`
	// If specified, the pod's tolerations.
	Tolerations []corev1.Toleration `json:"tolerations,omitempty"`
	// Define resources requests and limits for main container.
	Resources corev1.ResourceRequirements `json:"resources,omitempty"`
	// Number of replicas for a thanos component
	Replicas *int32 `json:"replicas,omitempty"`

	// LocalTsdbRetention configs how long to retain raw samples on local storage
	LocalTsdbRetention string `json:"localTsdbRetention,omitempty"`

	// ObjectStorageConfig allows specifying a key of a Secret containing object store configuration
	ObjectStorageConfig *corev1.SecretKeySelector `json:"objectStorageConfig,omitempty"`
	// DataVolume specifies how volume shall be used
	DataVolume *KubernetesVolume `json:"dataVolume,omitempty"`
}

type StoreGateway struct {
	CommonThanosFields `json:",inline"`

	// If specified, the pod's scheduling constraints.
	Affinity *corev1.Affinity `json:"affinity,omitempty"`
	// Define which Nodes the Pods are scheduled on.
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`
	// If specified, the pod's tolerations.
	Tolerations []corev1.Toleration `json:"tolerations,omitempty"`
	// Define resources requests and limits for main container.
	Resources corev1.ResourceRequirements `json:"resources,omitempty"`
	// Number of replicas for a thanos component
	Replicas *int32 `json:"replicas,omitempty"`

	// MinTime specifies start of time range limit to serve
	MinTime string `json:"minTime,omitempty"`
	// MaxTime specifies end of time range limit to serve
	MaxTime string `json:"maxTime,omitempty"`

	// ObjectStorageConfig allows specifying a key of a Secret containing object store configuration
	ObjectStorageConfig *corev1.SecretKeySelector `json:"objectStorageConfig,omitempty"`
	// DataVolume specifies how volume shall be used
	DataVolume *KubernetesVolume `json:"dataVolume,omitempty"`
}

type Compact struct {
	CommonThanosFields `json:",inline"`

	// If specified, the pod's scheduling constraints.
	Affinity *corev1.Affinity `json:"affinity,omitempty"`
	// Define which Nodes the Pods are scheduled on.
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`
	// If specified, the pod's tolerations.
	Tolerations []corev1.Toleration `json:"tolerations,omitempty"`
	// Define resources requests and limits for main container.
	Resources corev1.ResourceRequirements `json:"resources,omitempty"`
	// Number of replicas for a thanos component
	Replicas *int32 `json:"replicas,omitempty"`

	// DownsamplingDisable specifies whether to disable downsampling
	DownsamplingDisable *bool `json:"downsamplingDisable,omitempty"`
	// Retention configs how long to retain samples
	Retention *Retention `json:"retention,omitempty"`

	// ObjectStorageConfig allows specifying a key of a Secret containing object store configuration
	ObjectStorageConfig *corev1.SecretKeySelector `json:"objectStorageConfig,omitempty"`
	// DataVolume specifies how volume shall be used
	DataVolume *KubernetesVolume `json:"dataVolume,omitempty"`
}

// KubernetesVolume defines the configured volume for a thanos instance.
type KubernetesVolume struct {
	EmptyDir              *corev1.EmptyDirVolumeSource  `json:"emptyDir,omitempty"`
	PersistentVolumeClaim *corev1.PersistentVolumeClaim `json:"pvc,omitempty"`
}

// Retention defines the config for retaining samples
type Retention struct {
	// RetentionRaw specifies how long to retain raw samples in bucket
	RetentionRaw string `json:"retentionRaw,omitempty"`
	// Retention5m specifies how long to retain samples of 5m resolution in bucket
	Retention5m string `json:"retention5m,omitempty"`
	// Retention1h specifies how long to retain samples of 1h resolution in bucket
	Retention1h string `json:"retention1h,omitempty"`
}

// ServiceStatus defines the observed state of Service
type ServiceStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
// +genclient

// Service is the Schema for the monitoring service API
type Service struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ServiceSpec   `json:"spec,omitempty"`
	Status ServiceStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ServiceList contains a list of Service
type ServiceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Service `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Service{}, &ServiceList{})
}

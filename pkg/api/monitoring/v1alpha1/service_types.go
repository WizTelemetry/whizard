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

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
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

	// Gateway to proxy and auth requests to Thanos Query and Thanos Receive Router defined in Thanos.
	Gateway *Gateway `json:"gateway,omitempty"`

	// Thanos cluster contains explicit Thanos Query and Thanos Receive Router,
	// and implicit Thanos Receive Ingestor and Thanos Store Gateway and Thanos Compact
	// which are selected by label selector `monitoring.paodin.io/service=<service_namespace>.<service_name>`.
	Thanos *Thanos `json:"thanos,omitempty"`
}

type Gateway struct {
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

	// Image is the gateway image with tag/version.
	Image string `json:"image,omitempty"`
	// Log filtering level. Possible options: error, warn, info, debug.
	LogLevel string `json:"logLevel,omitempty"`
	// Log format to use. Possible options: logfmt or json.
	LogFormat string `json:"logFormat,omitempty"`

	// Secret name for HTTP Server certificate (Kubernetes TLS secret type)
	ServerCertificate string `json:"serverCertificate,omitempty"`
	// Secret name for HTTP Client CA certificate (Kubernetes TLS secret type)
	ClientCACertificate string `json:"clientCaCertificate,omitempty"`
}

type Thanos struct {
	// Thanos Query component querys from the backends such as Thanos Receive Ingestor and Thanos Store Gateway by automated discovery.
	Query *Query `json:"query,omitempty"`
	// Thanos Receive Router component routes to the backends such as Thanos Receive Ingestor by automated discovery.
	ReceiveRouter *ThanosReceiveRouter `json:"receiveRouter,omitempty"`
	// Thanos Query frontend component implements a service deployed in front of queriers to improve query parallelization and caching.
	QueryFrontend *ThanosQueryFrontend `json:"queryFrontend,omitempty"`
}

type Query struct {

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

	// Image is the thanos image with tag/version
	Image string `json:"image,omitempty"`
	// Log filtering level. Possible options: error, warn, info, debug
	LogLevel string `json:"logLevel,omitempty"`
	// Log format to use. Possible options: logfmt or json
	LogFormat string `json:"logFormat,omitempty"`

	// Additional StoreApi servers from which Thanos Query component queries from
	Stores []QueryStores `json:"stores,omitempty"`
	// Selector labels that will be exposed in info endpoint.
	SelectorLabels map[string]string `json:"selectorLabels,omitempty"`
	// Labels to treat as a replica indicator along which data is deduplicated.
	ReplicaLabelNames []string `json:"replicaLabelNames,omitempty"`

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

type ThanosReceiveRouter struct {
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

	// Image is the thanos image with tag/version
	Image string `json:"image,omitempty"`
	// Log filtering level. Possible options: error, warn, info, debug
	LogLevel string `json:"logLevel,omitempty"`
	// Log format to use. Possible options: logfmt or json
	LogFormat string `json:"logFormat,omitempty"`

	// How many times to replicate incoming write requests
	ReplicationFactor *uint64 `json:"replicationFactor,omitempty"`
}

type ThanosQueryFrontend struct {
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

	// Image is the thanos image with tag/version
	Image string `json:"image,omitempty"`
	// Log filtering level. Possible options: error, warn, info, debug
	LogLevel string `json:"logLevel,omitempty"`
	// Log format to use. Possible options: logfmt or json
	LogFormat string `json:"logFormat,omitempty"`

	// MaxSizeInMemoryCacheConfig represents overall maximum number of bytes cache can contain. A unit suffix (KB, MB, GB) may be applied.
	MaxSizeInMemoryCacheConfig string `json:"maxSize,omitempty"`
	// MaxSizeItemsInMemoryCacheConfig represents the maximum number of entries in the cache.
	MaxSizeItemsInMemoryCacheConfig int32 `json:"maxSizeItems,omitempty"`
	// ValidityInMemoryCacheConfig represents the expiry duration for the cache.
	ValidityInMemoryCacheConfig int64 `json:"validity,omitempty"`
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

type ThanosStoreGateway struct {
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

	// Image is the thanos image with tag/version
	Image string `json:"image,omitempty"`
	// Log filtering level. Possible options: error, warn, info, debug
	LogLevel string `json:"logLevel,omitempty"`
	// Log format to use. Possible options: logfmt or json
	LogFormat string `json:"logFormat,omitempty"`

	// MinTime specifies start of time range limit to serve
	MinTime string `json:"minTime,omitempty"`
	// MaxTime specifies end of time range limit to serve
	MaxTime string `json:"maxTime,omitempty"`

	// DataVolume specifies how volume shall be used
	DataVolume *KubernetesVolume `json:"dataVolume,omitempty"`
}

type ThanosCompact struct {
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

	// Image is the thanos image with tag/version
	Image string `json:"image,omitempty"`
	// Log filtering level. Possible options: error, warn, info, debug
	LogLevel string `json:"logLevel,omitempty"`
	// Log format to use. Possible options: logfmt or json
	LogFormat string `json:"logFormat,omitempty"`

	// DownsamplingDisable specifies whether to disable downsampling
	DownsamplingDisable *bool `json:"downsamplingDisable,omitempty"`
	// Retention configs how long to retain samples
	Retention *Retention `json:"retention,omitempty"`

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
	RetentionRaw Duration `json:"retentionRaw,omitempty"`
	// Retention5m specifies how long to retain samples of 5m resolution in bucket
	Retention5m Duration `json:"retention5m,omitempty"`
	// Retention1h specifies how long to retain samples of 1h resolution in bucket
	Retention1h Duration `json:"retention1h,omitempty"`
}

// StoreSpec defines the desired state of a Store
type StoreSpec struct {
	// ObjectStorageConfig allows specifying a key of a Secret containing object store configuration
	ObjectStorageConfig *corev1.SecretKeySelector `json:"objectStorageConfig,omitempty"`
	// Thanos contains Thanos Store Gateway and Thanos Compact.
	Thanos *ThanosStore `json:"thanos,omitempty"`
}

type ThanosStore struct {
	// Thanos Store Gateway will be selected as query backends by Service.
	StoreGateway *ThanosStoreGateway `json:"storeGateway,omitempty"`
	// Thanos Compact as object storage data compactor and lifecycle manager.
	Compact *ThanosCompact `json:"compact,omitempty"`
}

// StoreStatus defines the observed state of Store
type StoreStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
// +genclient

// Store is the Schema for the Store API
type Store struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   StoreSpec   `json:"spec,omitempty"`
	Status StoreStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// StoreList contains a list of Store
type StoreList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Store `json:"items"`
}

// ThanosReceiveIngestorSpec defines the desired state of a ThanosReceiveIngestor
type ThanosReceiveIngestorSpec struct {
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
	// Number of replicas for a thanos component.
	Replicas *int32 `json:"replicas,omitempty"`

	// Image is the thanos image with tag/version
	Image string `json:"image,omitempty"`
	// Log filtering level. Possible options: error, warn, info, debug
	LogLevel string `json:"logLevel,omitempty"`
	// Log format to use. Possible options: logfmt or json
	LogFormat string `json:"logFormat,omitempty"`

	// LocalTsdbRetention configs how long to retain raw samples on local storage.
	LocalTsdbRetention string `json:"localTsdbRetention,omitempty"`

	// If specified, the object key of Store for long term storage.
	LongTermStore *ObjectReference `json:"longTermStore,omitempty"`

	// DataVolume specifies how volume shall be used
	DataVolume *KubernetesVolume `json:"dataVolume,omitempty"`
}

type ObjectReference struct {
	Namespace string `json:"namespace,omitempty"`
	Name      string `json:"name,omitempty"`
}

// ThanosReceiveIngestorStatus defines the observed state of Store
type ThanosReceiveIngestorStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
// +genclient

// ThanosReceiveIngestor is the Schema for the ThanosReceiveIngestor API
type ThanosReceiveIngestor struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ThanosReceiveIngestorSpec   `json:"spec,omitempty"`
	Status ThanosReceiveIngestorStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ThanosReceiveIngestorList contains a list of Store
type ThanosReceiveIngestorList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ThanosReceiveIngestor `json:"items"`
}

// ThanosRulerSpec defines the desired state of a ThanosRuler
type ThanosRulerSpec struct {
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

	// Image is the thanos image with tag/version
	Image string `json:"image,omitempty"`
	// Log filtering level. Possible options: error, warn, info, debug
	LogLevel string `json:"logLevel,omitempty"`
	// Log format to use. Possible options: logfmt or json
	LogFormat string `json:"logFormat,omitempty"`

	// AlertingRules to be selected for alerting.
	AlertingRuleSelector *metav1.LabelSelector `json:"alertingRuleSelector,omitempty"`
	// Namespaces to be selected for AlertingRules discovery. If nil, only
	// check own namespace.
	AlertingRuleNamespaceSelector *metav1.LabelSelector `json:"alertingRuleNamespaceSelector,omitempty"`

	// A label selector to select which PrometheusRules to mount for alerting and
	// recording.
	RuleSelector *metav1.LabelSelector `json:"ruleSelector,omitempty"`
	// Namespaces to be selected for Rules discovery. If unspecified, only
	// the same namespace as the ThanosRuler object is in is used.
	RuleNamespaceSelector *metav1.LabelSelector `json:"ruleNamespaceSelector,omitempty"`

	// Labels configure the external label pairs to ThanosRuler. A default replica label
	// `thanos_ruler_replica` will be always added  as a label with the value of the pod's name and it will be dropped in the alerts.
	Labels map[string]string `json:"labels,omitempty"`
	// AlertDropLabels configure the label names which should be dropped in ThanosRuler alerts.
	// The replica label `thanos_ruler_replica` will always be dropped in alerts.
	AlertDropLabels []string `json:"alertDropLabels,omitempty"`
	// Define URLs to send alerts to Alertmanager.  For Thanos v0.10.0 and higher,
	// AlertManagersConfig should be used instead.  Note: this field will be ignored
	// if AlertManagersConfig is specified.
	// Maps to the `alertmanagers.url` arg.
	AlertManagersURL []string `json:"alertmanagersUrl,omitempty"`
	// Define configuration for connecting to alertmanager.  Only available with thanos v0.10.0
	// and higher.  Maps to the `alertmanagers.config` arg.
	AlertManagersConfig *corev1.SecretKeySelector `json:"alertmanagersConfig,omitempty"`
	// Interval between consecutive evaluations. Default: `30s`
	// +kubebuilder:default:="30s"
	EvaluationInterval Duration `json:"evaluationInterval,omitempty"`

	// DataVolume specifies how volume shall be used
	DataVolume *KubernetesVolume `json:"dataVolume,omitempty"`
}

// ThanosRulerStatus defines the observed state of ThanosRuler
type ThanosRulerStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
// +genclient

// ThanosRuler is the Schema for the ThanosRuler API
type ThanosRuler struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ThanosRulerSpec   `json:"spec,omitempty"`
	Status ThanosRulerStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ThanosRulerList contains a list of ThanosRuler
type ThanosRulerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ThanosRuler `json:"items"`
}

// AlertingRuleSpec defines the desired state of a AlertingRule
type AlertingRuleSpec struct {
	Expr        intstr.IntOrString `json:"expr"`
	For         string             `json:"for,omitempty"`
	Labels      map[string]string  `json:"labels,omitempty"`
	Annotations map[string]string  `json:"annotations,omitempty"`
}

// AlertingRuleStatus defines the observed state of AlertingRule
type AlertingRuleStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
// +genclient

// AlertingRule is the Schema for the AlertingRule API
type AlertingRule struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AlertingRuleSpec   `json:"spec,omitempty"`
	Status AlertingRuleStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// AlertingRuleList contains a list of AlertingRule
type AlertingRuleList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AlertingRule `json:"items"`
}

// RuleGroupSpec defines the desired state of a RuleGroup
type RuleGroupSpec struct {
	Interval                string `json:"interval,omitempty"`
	PartialResponseStrategy string `json:"partial_response_strategy,omitempty"`
}

// RuleGroupStatus defines the observed state of RuleGroup
type RuleGroupStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
// +genclient

// RuleGroup is the Schema for the RuleGroup API
type RuleGroup struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RuleGroupSpec   `json:"spec,omitempty"`
	Status RuleGroupStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// RuleGroupList contains a list of RuleGroup
type RuleGroupList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []RuleGroup `json:"items"`
}

// Duration is a valid time unit
// Supported units: y, w, d, h, m, s, ms Examples: `30s`, `1m`, `1h20m15s`
// +kubebuilder:validation:Pattern:="^(0|(([0-9]+)y)?(([0-9]+)w)?(([0-9]+)d)?(([0-9]+)h)?(([0-9]+)m)?(([0-9]+)s)?(([0-9]+)ms)?)$"
type Duration string

func init() {
	SchemeBuilder = SchemeBuilder.
		Register(&Service{}, &ServiceList{}).
		Register(&ThanosReceiveIngestor{}, &ThanosReceiveIngestorList{}).
		Register(&Store{}, &StoreList{}).
		Register(&ThanosRuler{}, &ThanosRulerList{}).
		Register(&AlertingRule{}, &AlertingRuleList{}).
		Register(&RuleGroup{}, &RuleGroupList{})
}

func ManagedLabelByService(service metav1.Object) map[string]string {
	return map[string]string{
		"monitoring.paodin.io/service": service.GetNamespace() + "." + service.GetName(),
	}
}

func ServiceNamespacedName(managedByService metav1.Object) *types.NamespacedName {
	ls := managedByService.GetLabels()
	if len(ls) == 0 {
		return nil
	}

	namespacedName := ls["monitoring.paodin.io/service"]
	arr := strings.Split(namespacedName, ".")
	if len(arr) != 2 {
		return nil
	}

	return &types.NamespacedName{
		Namespace: arr[0],
		Name:      arr[1],
	}
}

func ManagedLabelByRuleGroup(ruleGroup metav1.Object) map[string]string {
	return map[string]string{
		"monitoring.paodin.io/rule-group": ruleGroup.GetName(),
	}
}

func RuleGroupName(managedByRuleGroup metav1.Object) string {
	ls := managedByRuleGroup.GetLabels()
	if ls == nil {
		return ""
	}
	return ls["monitoring.paodin.io/rule-group"]
}

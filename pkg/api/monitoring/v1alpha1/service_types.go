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
	"time"

	"k8s.io/api/autoscaling/v2beta2"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

const (
	FinalizerMonitoringPaodinDeletePVC = "finalizers.monitoring.paodin.io/deletePVC"

	DefaultTenantHeader    = "PAODIN-TENANT"
	DefaultTenantId        = "default-tenant"
	DefaultTenantLabelName = "tenant_id"
)

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

	Storage *ObjectReference `json:"storage,omitempty"`

	// Gateway to proxy and auth requests to Query and Router defined in Thanos.
	Gateway *Gateway `json:"gateway,omitempty"`

	// Query component querys from the backends such as Ingester and Store by automated discovery.
	Query *Query `json:"query,omitempty"`

	// Receive Router component routes to the backends such as Ingester by automated discovery.
	Router *Router `json:"router,omitempty"`

	// QueryFrontend component implements a service deployed in front of queriers to improve query parallelization and caching.
	QueryFrontend *QueryFrontend `json:"queryFrontend,omitempty"`
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
	// Number of replicas for a component
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

type Query struct {

	// If specified, the pod's scheduling constraints.
	Affinity *corev1.Affinity `json:"affinity,omitempty"`
	// Define which Nodes the Pods are scheduled on.
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`
	// If specified, the pod's tolerations.
	Tolerations []corev1.Toleration `json:"tolerations,omitempty"`
	// Define resources requests and limits for main container.
	Resources corev1.ResourceRequirements `json:"resources,omitempty"`
	// Number of replicas for a component
	Replicas *int32 `json:"replicas,omitempty"`

	// Image is the thanos image with tag/version
	Image string `json:"image,omitempty"`
	// Log filtering level. Possible options: error, warn, info, debug
	LogLevel string `json:"logLevel,omitempty"`
	// Log format to use. Possible options: logfmt or json
	LogFormat string `json:"logFormat,omitempty"`

	// Additional StoreApi servers from which Query component queries from
	Stores []QueryStores `json:"stores,omitempty"`
	// Selector labels that will be exposed in info endpoint.
	SelectorLabels map[string]string `json:"selectorLabels,omitempty"`
	// Labels to treat as a replica indicator along which data is deduplicated.
	ReplicaLabelNames []string `json:"replicaLabelNames,omitempty"`

	// Flags is a list of key/value that could be used to set strategy parameters.
	Flags map[string]string `json:"flags,omitempty"`

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
	// Image is the envoy image with tag/version
	Image string `json:"image,omitempty"`
	// Define resources requests and limits for envoy container.
	Resources corev1.ResourceRequirements `json:"resources,omitempty"`
}

type Router struct {
	// If specified, the pod's scheduling constraints.
	Affinity *corev1.Affinity `json:"affinity,omitempty"`
	// Define which Nodes the Pods are scheduled on.
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`
	// If specified, the pod's tolerations.
	Tolerations []corev1.Toleration `json:"tolerations,omitempty"`
	// Define resources requests and limits for main container.
	Resources corev1.ResourceRequirements `json:"resources,omitempty"`
	// Number of replicas for a component.
	Replicas *int32 `json:"replicas,omitempty"`

	// Image is the thanos image with tag/version
	Image string `json:"image,omitempty"`
	// Log filtering level. Possible options: error, warn, info, debug
	LogLevel string `json:"logLevel,omitempty"`
	// Log format to use. Possible options: logfmt or json
	LogFormat string `json:"logFormat,omitempty"`

	// How many times to replicate incoming write requests
	ReplicationFactor *uint64 `json:"replicationFactor,omitempty"`

	// Flags is a list of key/value that could be used to set strategy parameters.
	Flags map[string]string `json:"flags,omitempty"`
}

type QueryFrontend struct {
	// If specified, the pod's scheduling constraints.
	Affinity *corev1.Affinity `json:"affinity,omitempty"`
	// Define which Nodes the Pods are scheduled on.
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`
	// If specified, the pod's tolerations.
	Tolerations []corev1.Toleration `json:"tolerations,omitempty"`
	// Define resources requests and limits for main container.
	Resources corev1.ResourceRequirements `json:"resources,omitempty"`
	// Number of replicas for a component
	Replicas *int32 `json:"replicas,omitempty"`

	// Image is the thanos image with tag/version
	Image string `json:"image,omitempty"`
	// Log filtering level. Possible options: error, warn, info, debug
	LogLevel string `json:"logLevel,omitempty"`
	// Log format to use. Possible options: logfmt or json
	LogFormat string `json:"logFormat,omitempty"`
	// Flags is a list of key/value that could be used to set strategy parameters.
	Flags map[string]string `json:"flags,omitempty"`

	// CacheProviderConfig ...
	CacheConfig *ResponseCacheProviderConfig `json:"cacheConfig,omitempty"`
}

type CacheProvider string

const (
	INMEMORY  CacheProvider = "IN-MEMORY"
	MEMCACHED CacheProvider = "MEMCACHED"
	REDIS     CacheProvider = "REDIS"
)

// ResponseCacheProviderConfig is the initial ResponseCacheProviderConfig struct holder before parsing it into a specific cache provider.
// Based on the config type the config is then parsed into a specific cache provider.
type ResponseCacheProviderConfig struct {
	Type                        CacheProvider                `json:"type"`
	InMemoryResponseCacheConfig *InMemoryResponseCacheConfig `json:"inMemory,omitempty"`
}

// InMemoryResponseCacheConfig holds the configs for the in-memory cache provider.
type InMemoryResponseCacheConfig struct {
	// MaxSize represents overall maximum number of bytes cache can contain.
	MaxSize string `json:"maxSize" yaml:"max_size"`
	// MaxSizeItems represents the maximum number of entries in the cache.
	MaxSizeItems int `json:"maxSizeItems" yaml:"max_size_items"`
	// Validity represents the expiry duration for the cache.
	Validity time.Duration `json:"validity" yaml:"validity"`
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

// KubernetesVolume defines the configured volume for a instance.
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

// IndexCacheConfig specifies the index cache config.
type IndexCacheConfig struct {
	*InMemoryIndexCacheConfig `json:"inMemory,omitempty" yaml:"inMemory,omitempty"`
}

type InMemoryIndexCacheConfig struct {
	// MaxSize represents overall maximum number of bytes cache can contain.
	MaxSize string `json:"maxSize" yaml:"maxSize"`
}

type AutoScaler struct {
	// minReplicas is the lower limit for the number of replicas to which the autoscaler
	// can scale down.  It defaults to 1 pod.  minReplicas is allowed to be 0 if the
	// alpha feature gate HPAScaleToZero is enabled and at least one Object or External
	// metric is configured.  Scaling is active as long as at least one metric value is
	// available.
	// +optional
	MinReplicas *int32 `json:"minReplicas,omitempty" yaml:"minReplicas,omitempty"`
	// maxReplicas is the upper limit for the number of replicas to which the autoscaler can scale up.
	// It cannot be less that minReplicas.
	MaxReplicas int32 `json:"maxReplicas" yaml:"maxReplicas"`
	// metrics contains the specifications for which to use to calculate the
	// desired replica count (the maximum replica count across all metrics will
	// be used).  The desired replica count is calculated multiplying the
	// ratio between the target value and the current value by the current
	// number of pods.  Ergo, metrics used must decrease as the pod count is
	// increased, and vice-versa.  See the individual metric source types for
	// more information about how each type of metric must respond.
	// If not set, the default metric will be set to 80% average CPU utilization.
	// +optional
	Metrics []v2beta2.MetricSpec `json:"metrics,omitempty" yaml:"metrics,omitempty"`
	// behavior configures the scaling behavior of the target
	// in both Up and Down directions (scaleUp and scaleDown fields respectively).
	// If not set, the default HPAScalingRules for scale up and scale down are used.
	// +optional
	Behavior *v2beta2.HorizontalPodAutoscalerBehavior `json:"behavior,omitempty" yaml:"behavior,omitempty"`
}

type StoreSpec struct {
	// If specified, the pod's scheduling constraints.
	Affinity *corev1.Affinity `json:"affinity,omitempty"`
	// Define which Nodes the Pods are scheduled on.
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`
	// If specified, the pod's tolerations.
	Tolerations []corev1.Toleration `json:"tolerations,omitempty"`
	// Define resources requests and limits for main container.
	Resources corev1.ResourceRequirements `json:"resources,omitempty"`
	// Number of replicas for a component
	Replicas *int32 `json:"replicas,omitempty"`

	// Image is the thanos image with tag/version
	Image string `json:"image,omitempty"`
	// Image pull policy.
	// One of Always, Never, IfNotPresent.
	// Defaults to Always if :latest tag is specified, or IfNotPresent otherwise.
	// Cannot be updated.
	// +optional
	ImagePullPolicy corev1.PullPolicy `json:"imagePullPolicy,omitempty"`
	// Log filtering level. Possible options: error, warn, info, debug
	LogLevel string `json:"logLevel,omitempty"`
	// Log format to use. Possible options: logfmt or json
	LogFormat string `json:"logFormat,omitempty"`

	Storage *ObjectReference `json:"storage,omitempty"`

	// MinTime specifies start of time range limit to serve
	MinTime string `json:"minTime,omitempty"`
	// MaxTime specifies end of time range limit to serve
	MaxTime string `json:"maxTime,omitempty"`

	// IndexCacheConfig contains index cache configuration.
	IndexCacheConfig *IndexCacheConfig `json:"indexCacheConfig,omitempty"`

	// Flags is a list of key/value that could be used to set strategy parameters.
	Flags map[string]string `json:"flags,omitempty"`

	// DataVolume specifies how volume shall be used
	DataVolume *KubernetesVolume `json:"dataVolume,omitempty"`

	Scaler *AutoScaler `json:"scaler,omitempty"`
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

type CompactorSpec struct {
	// If specified, the pod's scheduling constraints.
	Affinity *corev1.Affinity `json:"affinity,omitempty"`
	// Define which Nodes the Pods are scheduled on.
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`
	// If specified, the pod's tolerations.
	Tolerations []corev1.Toleration `json:"tolerations,omitempty"`
	// Define resources requests and limits for main container.
	Resources corev1.ResourceRequirements `json:"resources,omitempty"`
	// Number of replicas for a component
	Replicas *int32 `json:"replicas,omitempty"`

	// Image is the thanos image with tag/version
	Image string `json:"image,omitempty"`
	// Image pull policy.
	// One of Always, Never, IfNotPresent.
	// Defaults to Always if :latest tag is specified, or IfNotPresent otherwise.
	// Cannot be updated.
	// +optional
	ImagePullPolicy corev1.PullPolicy `json:"imagePullPolicy,omitempty"`
	// Log filtering level. Possible options: error, warn, info, debug
	LogLevel string `json:"logLevel,omitempty"`
	// Log format to use. Possible options: logfmt or json
	LogFormat string `json:"logFormat,omitempty"`

	// DownsamplingDisable specifies whether to disable downsampling
	DownsamplingDisable *bool `json:"downsamplingDisable,omitempty"`
	// Retention configs how long to retain samples
	Retention *Retention `json:"retention,omitempty"`

	Storage *ObjectReference `json:"storage,omitempty"`

	// Flags is a list of key/value that could be used to set strategy parameters.
	Flags map[string]string `json:"flags,omitempty"`

	// DataVolume specifies how volume shall be used
	DataVolume *KubernetesVolume `json:"dataVolume,omitempty"`

	// Tenants if not empty indicates current config is for hard tenants; otherwise, it is for soft tenants.
	Tenants []string `json:"tenants,omitempty"`
}

// CompactorStatus defines the observed state of Compactor
type CompactorStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
// +genclient

// Compactor is the Schema for the Compactor API
type Compactor struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CompactorSpec   `json:"spec,omitempty"`
	Status CompactorStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// CompactorList contains a list of Compactor
type CompactorList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Compactor `json:"items"`
}

// IngesterSpec defines the desired state of a Ingester
type IngesterSpec struct {
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
	// Number of replicas for a component.
	Replicas *int32 `json:"replicas,omitempty"`

	// Image is the thanos image with tag/version
	Image string `json:"image,omitempty"`
	// Log filtering level. Possible options: error, warn, info, debug
	LogLevel string `json:"logLevel,omitempty"`
	// Log format to use. Possible options: logfmt or json
	LogFormat string `json:"logFormat,omitempty"`

	// LocalTsdbRetention configs how long to retain raw samples on local storage.
	LocalTsdbRetention string `json:"localTsdbRetention,omitempty"`

	// Flags is a list of key/value that could be used to set strategy parameters.
	Flags map[string]string `json:"flags,omitempty"`

	// If specified, the object key of Storage for long term storage.
	Storage *ObjectReference `json:"storage,omitempty"`

	// DataVolume specifies how volume shall be used
	DataVolume *KubernetesVolume `json:"dataVolume,omitempty"`
}

type ObjectReference struct {
	Namespace string `json:"namespace,omitempty"`
	Name      string `json:"name,omitempty"`
}

// IngesterStatus defines the observed state of Ingester
type IngesterStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
// +genclient

// Ingester is the Schema for the Ingester API
type Ingester struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   IngesterSpec   `json:"spec,omitempty"`
	Status IngesterStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// IngesterList contains a list of Ingester
type IngesterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Ingester `json:"items"`
}

// RulerSpec defines the desired state of a Ruler
type RulerSpec struct {
	// If specified, the pod's scheduling constraints.
	Affinity *corev1.Affinity `json:"affinity,omitempty"`
	// Define which Nodes the Pods are scheduled on.
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`
	// If specified, the pod's tolerations.
	Tolerations []corev1.Toleration `json:"tolerations,omitempty"`
	// Define resources requests and limits for main container.
	Resources corev1.ResourceRequirements `json:"resources,omitempty"`
	// Number of replicas for a component.
	Replicas *int32 `json:"replicas,omitempty"`

	// Image is the thanos image with tag/version
	Image string `json:"image,omitempty"`
	// Log filtering level. Possible options: error, warn, info, debug
	LogLevel string `json:"logLevel,omitempty"`
	// Log format to use. Possible options: logfmt or json
	LogFormat string `json:"logFormat,omitempty"`

	// A label selector to select which Rules to mount for alerting and
	// recording.
	RuleSelector *metav1.LabelSelector `json:"ruleSelector,omitempty"`
	// Namespaces to be selected for Rules discovery. If nil, only
	// the same namespace as the Ruler object is in is used.
	RuleNamespaceSelector *metav1.LabelSelector `json:"ruleNamespaceSelector,omitempty"`
	// A label selector to select which PrometheusRules to mount for alerting and
	// recording.
	PrometheusRuleSelector *metav1.LabelSelector `json:"prometheusRuleSelector,omitempty"`
	// Namespaces to be selected for PrometheusRules discovery. If unspecified, only
	// the same namespace as the Ruler object is in is used.
	PrometheusRuleNamespaceSelector *metav1.LabelSelector `json:"prometheusRuleNamespaceSelector,omitempty"`

	// Tenant if not empty indicates which tenant's data is evaluated for the selected rules;
	// otherwise, it is for all tenants.
	Tenant string `json:"tenant,omitempty"`

	// Labels configure the external label pairs to Ruler. A default replica label
	// `thanos_ruler_replica` will be always added  as a label with the value of the pod's name and it will be dropped in the alerts.
	Labels map[string]string `json:"labels,omitempty"`
	// AlertDropLabels configure the label names which should be dropped in Ruler alerts.
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

	// Flags is a list of key/value that could be used to set strategy parameters.
	Flags map[string]string `json:"flags,omitempty"`

	// DataVolume specifies how volume shall be used
	DataVolume *KubernetesVolume `json:"dataVolume,omitempty"`
}

// RulerStatus defines the observed state of Ruler
type RulerStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
// +genclient

// Ruler is the Schema for the Ruler API
type Ruler struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RulerSpec   `json:"spec,omitempty"`
	Status RulerStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// RulerList contains a list of Ruler
type RulerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Ruler `json:"items"`
}

// RuleSpec defines the desired state of a Rule
type RuleSpec struct {
	Alert       string             `json:"alert,omitempty"`
	Record      string             `json:"record,omitempty"`
	Expr        intstr.IntOrString `json:"expr"`
	For         Duration           `json:"for,omitempty"`
	Labels      map[string]string  `json:"labels,omitempty"`
	Annotations map[string]string  `json:"annotations,omitempty"`
}

// RuleStatus defines the observed state of Rule
type RuleStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
// +genclient

// Rule is the Schema for the Rule API
type Rule struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RuleSpec   `json:"spec,omitempty"`
	Status RuleStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// RuleList contains a list of Rule
type RuleList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Rule `json:"items"`
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
		Register(&Ingester{}, &IngesterList{}).
		Register(&Ruler{}, &RulerList{}).
		Register(&Store{}, &StoreList{}).
		Register(&Compactor{}, &CompactorList{}).
		Register(&Rule{}, &RuleList{}).
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

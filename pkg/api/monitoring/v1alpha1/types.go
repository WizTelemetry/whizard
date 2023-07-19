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
	"time"

	"k8s.io/api/autoscaling/v2beta2"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// Duration is a valid time unit
// Supported units: y, w, d, h, m, s, ms Examples: `30s`, `1m`, `1h20m15s`
// +kubebuilder:validation:Pattern:="^(0|(([0-9]+)y)?(([0-9]+)w)?(([0-9]+)d)?(([0-9]+)h)?(([0-9]+)m)?(([0-9]+)s)?(([0-9]+)ms)?)$"
type Duration string

// Retention defines the config for retaining samples
type Retention struct {
	// RetentionRaw specifies how long to retain raw samples in bucket
	RetentionRaw Duration `json:"retentionRaw,omitempty"`
	// Retention5m specifies how long to retain samples of 5m resolution in bucket
	Retention5m Duration `json:"retention5m,omitempty"`
	// Retention1h specifies how long to retain samples of 1h resolution in bucket
	Retention1h Duration `json:"retention1h,omitempty"`
}

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

	// Retention configs how long to retain samples
	Retention *Retention `json:"retention,omitempty"`

	Storage *ObjectReference `json:"storage,omitempty"`

	// RemoteWrites is the list of remote write configurations.
	// If it is configured, its targets will receive write requests from the Gateway and the Ruler.
	RemoteWrites []RemoteWriteSpec `json:"remoteWrites,omitempty"`
	// RemoteQuery is the remote query configuration and the remote target should have prometheus-compatible Query APIs.
	// If not configured, the Gateway will proxy all read requests through the QueryFrontend to the Query,
	// If configured, the Gateway will proxy metrics read requests through the QueryFrontend to the remote target,
	// but proxy rules read requests directly to the Query.
	RemoteQuery *RemoteQuerySpec `json:"remoteQuery,omitempty"`
}

// RemoteQuerySpec defines the configuration to query from remote service
// which should have prometheus-compatible Query APIs.
type RemoteQuerySpec struct {
	Name string `json:"name,omitempty"`
	URL  string `json:"url"`
}

// RemoteWriteSpec defines the remote write configuration.
type RemoteWriteSpec struct {
	Name string `json:"name,omitempty"`
	URL  string `json:"url"`
	// Custom HTTP headers to be sent along with each remote write request.
	Headers map[string]string `json:"headers,omitempty"`
	// Timeout for requests to the remote write endpoint.
	RemoteTimeout Duration `json:"remoteTimeout,omitempty"`
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

type CommonSpec struct {
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

	// Image is the component image with tag/version.
	Image string `json:"image,omitempty"`
	// Image pull policy.
	// One of Always, Never, IfNotPresent.
	// Defaults to Always if :latest tag is specified, or IfNotPresent otherwise.
	// Cannot be updated.
	// +optional
	ImagePullPolicy corev1.PullPolicy `json:"imagePullPolicy,omitempty"`
	// An optional list of references to secrets in the same namespace
	// to use for pulling images from registries
	ImagePullSecrets []corev1.LocalObjectReference `json:"imagePullSecrets,omitempty"`

	// Log filtering level. Possible options: error, warn, info, debug.
	LogLevel string `json:"logLevel,omitempty"`
	// Log format to use. Possible options: logfmt or json.
	LogFormat string `json:"logFormat,omitempty"`
	// Flags is the flags of component.
	Flags []string `json:"flags,omitempty"`
}

type GatewaySpec struct {
	CommonSpec `json:",inline"`

	// Secret name for HTTP Server certificate (Kubernetes TLS secret type)
	ServerCertificate string `json:"serverCertificate,omitempty"`
	// Secret name for HTTP Client CA certificate (Kubernetes TLS secret type)
	ClientCACertificate string `json:"clientCaCertificate,omitempty"`

	// NodePort is the port used to expose the gateway service.
	// If this is a valid node port, the gateway service type will be set to NodePort accordingly.
	NodePort int32 `json:"nodePort,omitempty"`
}

// GatewayStatus defines the observed state of Gateway
type GatewayStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
// +genclient

// Gateway is the Schema for the monitoring gateway API
type Gateway struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   GatewaySpec   `json:"spec,omitempty"`
	Status GatewayStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// GatewayList contains a list of Gateway
type GatewayList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Gateway `json:"items"`
}

type HTTPServerTLSConfig struct {
	// Secret containing the TLS key for the server.
	KeySecret corev1.SecretKeySelector `json:"keySecret"`
	// Contains the TLS certificate for the server.
	CertSecret corev1.SecretKeySelector `json:"certSecret"`
	// Contains the CA certificate for client certificate authentication to the server.
	ClientCASecret corev1.SecretKeySelector `json:"clientCASecret,omitempty"`

	/*
		// Server policy for client authentication. Maps to ClientAuth Policies.
		// For more detail on clientAuth options:
		// https://golang.org/pkg/crypto/tls/#ClientAuthType
		ClientAuthType string `json:"clientAuthType,omitempty"`
		// Minimum TLS version that is acceptable. Defaults to TLS12.
		MinVersion string `json:"minVersion,omitempty"`
		// Maximum TLS version that is acceptable. Defaults to TLS13.
		MaxVersion string `json:"maxVersion,omitempty"`
		// List of supported cipher suites for TLS versions up to TLS 1.2. If empty,
		// Go default cipher suites are used. Available cipher suites are documented
		// in the go documentation: https://golang.org/pkg/crypto/tls/#pkg-constants
		CipherSuites []string `json:"cipherSuites,omitempty"`
		// Controls whether the server selects the
		// client's most preferred cipher suite, or the server's most preferred
		// cipher suite. If true then the server's preference, as expressed in
		// the order of elements in cipherSuites, is used.
		PreferServerCipherSuites *bool `json:"preferServerCipherSuites,omitempty"`
		// Elliptic curves that will be used in an ECDHE handshake, in preference
		// order. Available curves are documented in the go documentation:
		// https://golang.org/pkg/crypto/tls/#CurveID
		CurvePreferences []string `json:"curvePreferences,omitempty"`
	*/
}

type QuerySpec struct {
	CommonSpec `json:",inline"`

	HTTPServerTLSConfig *HTTPServerTLSConfig `json:"httpServerTLSConfig,omitempty"`

	// Additional StoreApi servers from which Query component queries from
	Stores []QueryStores `json:"stores,omitempty"`
	// Selector labels that will be exposed in info endpoint.
	SelectorLabels map[string]string `json:"selectorLabels,omitempty"`
	// Labels to treat as a replica indicator along which data is deduplicated.
	ReplicaLabelNames []string `json:"replicaLabelNames,omitempty"`

	// Envoy is used to config sidecar which proxies requests requiring auth to the secure stores
	Envoy SidecarSpec `json:"envoy,omitempty"`
}

type QueryStores struct {
	// Address is the addresses of StoreApi server, which may be prefixed with 'dns+' or 'dnssrv+' to detect StoreAPI servers through respective DNS lookups.
	Addresses []string `json:"addresses,omitempty"`
	// Secret containing the CA cert to use for StoreApi connections
	CASecret *corev1.SecretKeySelector `json:"caSecret,omitempty"`
}

type SidecarSpec struct {
	// Image is the envoy image with tag/version
	Image string `json:"image,omitempty" yaml:"image,omitempty"`
	// Define resources requests and limits for sidecar container.
	Resources corev1.ResourceRequirements `json:"resources,omitempty" yaml:"resources,omitempty"`
}

// QueryStatus defines the observed state of Query
type QueryStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
// +genclient

// Query is the Schema for the monitoring query API
type Query struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   QuerySpec   `json:"spec,omitempty"`
	Status QueryStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// QueryList contains a list of Query
type QueryList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Query `json:"items"`
}

type RouterSpec struct {
	CommonSpec `json:",inline"`

	HTTPServerTLSConfig *HTTPServerTLSConfig `json:"httpServerTLSConfig,omitempty"`

	Envoy SidecarSpec `json:"envoy,omitempty"`
	// How many times to replicate incoming write requests
	ReplicationFactor *uint64 `json:"replicationFactor,omitempty"`
}

// RouterStatus defines the observed state of Query
type RouterStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
// +genclient

// Router is the Schema for the monitoring router API
type Router struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RouterSpec   `json:"spec,omitempty"`
	Status RouterStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// RouterList contains a list of Router
type RouterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Router `json:"items"`
}

type QueryFrontendSpec struct {
	CommonSpec `json:",inline"`

	HTTPServerTLSConfig *HTTPServerTLSConfig `json:"httpServerTLSConfig,omitempty"`

	Envoy SidecarSpec `json:"envoy,omitempty"`
	// CacheProviderConfig ...
	CacheConfig *ResponseCacheProviderConfig `json:"cacheConfig,omitempty"`
}

// QueryFrontendStatus defines the observed state of QueryFrontend
type QueryFrontendStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
// +genclient

// QueryFrontend is the Schema for the monitoring queryfrontend API
type QueryFrontend struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   QueryFrontendSpec   `json:"spec,omitempty"`
	Status QueryFrontendStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// QueryFrontendList contains a list of QueryFrontend
type QueryFrontendList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []QueryFrontend `json:"items"`
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
	Type                        CacheProvider                `json:"type" yaml:"type"`
	InMemoryResponseCacheConfig *InMemoryResponseCacheConfig `json:"inMemory,omitempty" yaml:"inMemory,omitempty"`
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

// KubernetesVolume defines the configured volume for a instance.
type KubernetesVolume struct {
	EmptyDir              *corev1.EmptyDirVolumeSource  `json:"emptyDir,omitempty"`
	PersistentVolumeClaim *corev1.PersistentVolumeClaim `json:"persistentVolumeClaim,omitempty"`
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
	CommonSpec `json:",inline"`

	// MinTime specifies start of time range limit to serve
	MinTime string `json:"minTime,omitempty"`
	// MaxTime specifies end of time range limit to serve
	MaxTime string `json:"maxTime,omitempty"`

	// IndexCacheConfig contains index cache configuration.
	IndexCacheConfig *IndexCacheConfig `json:"indexCacheConfig,omitempty"`

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
	CommonSpec `json:",inline"`

	// DisableDownsampling specifies whether to disable downsampling
	DisableDownsampling *bool `json:"disableDownsampling,omitempty"`

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
	CommonSpec `json:",inline"`

	// Tenants if not empty indicates current config is for hard tenants; otherwise, it is for soft tenants.
	Tenants []string `json:"tenants,omitempty"`

	// LocalTsdbRetention configs how long to retain raw samples on local storage.
	LocalTsdbRetention string `json:"localTsdbRetention,omitempty"`

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
	CommonSpec `json:",inline"`

	RulerQueryProxy SidecarSpec `json:"rulerQueryProxy,omitempty"`

	Envoy SidecarSpec `json:"envoy,omitempty"`

	PrometheusConfigReloader SidecarSpec `json:"prometheusConfigReloader,omitempty"`

	// Label selectors to select which PrometheusRules to mount for alerting and recording.
	// The result of multiple selectors are ORed.
	RuleSelectors []*metav1.LabelSelector `json:"ruleSelectors,omitempty"`
	// Namespaces to be selected for PrometheusRules discovery. If unspecified, only
	// the same namespace as the Ruler object is in is used.
	RuleNamespaceSelector *metav1.LabelSelector `json:"ruleNamespaceSelector,omitempty"`

	// Number of shards to take the hash of fully qualified name of the rule group in order to split rules.
	// Each shard of rules will be bound to one separate statefulset.
	Shards *int32 `json:"shards,omitempty"`

	// Tenant if not empty indicates which tenant's data is evaluated for the selected rules;
	// otherwise, it is for all tenants.
	Tenant string `json:"tenant,omitempty"`

	// Labels configure the external label pairs to Ruler. A default replica label
	// `ruler_replica` will be always added  as a label with the value of the pod's name and it will be dropped in the alerts.
	Labels map[string]string `json:"labels,omitempty"`
	// AlertDropLabels configure the label names which should be dropped in Ruler alerts.
	// The replica label `ruler_replica` will always be dropped in alerts.
	AlertDropLabels []string `json:"alertDropLabels,omitempty"`
	// Define URLs to send alerts to Alertmanager.
	// Note: this field will be ignored if AlertmanagersConfig is specified.
	// Maps to the `alertmanagers.url` arg.
	AlertmanagersURL []string `json:"alertmanagersUrl,omitempty"`
	// Define configuration for connecting to alertmanager. Maps to the `alertmanagers.config` arg.
	AlertmanagersConfig *corev1.SecretKeySelector `json:"alertmanagersConfig,omitempty"`
	// Interval between consecutive evaluations.
	EvaluationInterval Duration `json:"evaluationInterval,omitempty"`

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

func init() {
	SchemeBuilder = SchemeBuilder.
		Register(&Service{}, &ServiceList{}).
		Register(&Gateway{}, &GatewayList{}).
		Register(&Query{}, &QueryList{}).
		Register(&QueryFrontend{}, &QueryFrontendList{}).
		Register(&Router{}, &RouterList{}).
		Register(&Ingester{}, &IngesterList{}).
		Register(&Ruler{}, &RulerList{}).
		Register(&Store{}, &StoreList{}).
		Register(&Compactor{}, &CompactorList{})
}

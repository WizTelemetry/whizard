/*
Copyright 2024 the Whizard Authors.

Licensed under Apache License, Version 2.0 with a few additional conditions.

You may obtain a copy of the License at

    https://github.com/WhizardTelemetry/whizard/blob/main/LICENSE
*/

package v1alpha1

import (
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// CompactorSpec defines the desired state of Compactor
type CompactorSpec struct {
	// The tenants whose data is being compacted by the Compactor.
	Tenants []string `json:"tenants,omitempty"`

	// Disables downsampling.
	// This is not recommended, as querying long time ranges without non-downsampled data is not efficient and useful.
	// default: false
	DisableDownsampling *bool `json:"disableDownsampling,omitempty"`

	// Retention configs how long to retain samples
	Retention *Retention `json:"retention,omitempty"`

	// DataVolume specifies how volume shall be used
	DataVolume *KubernetesVolume `json:"dataVolume,omitempty"`

	CommonSpec `json:",inline"`
}

// Retention defines the config for retaining samples
type Retention struct {
	// How long to retain raw samples in bucket. Setting this to 0d will retain samples of this resolution forever
	// default: 0d
	RetentionRaw Duration `json:"retentionRaw,omitempty"`
	// How long to retain samples of resolution 1 (5 minutes) in bucket. Setting this to 0d will retain samples of this resolution forever
	// default: 0d
	Retention5m Duration `json:"retention5m,omitempty"`
	// How long to retain samples of resolution 2 (1 hour) in bucket. Setting this to 0d will retain samples of this resolution forever
	// default: 0d
	Retention1h Duration `json:"retention1h,omitempty"`
}

// CompactorStatus defines the observed state of Compactor
type CompactorStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +genclient
// +kubebuilder:object:root=true
// +kubebuilder:printcolumn:name="Replicas",type="integer",JSONPath=".spec.replicas",description="The number of desired replicas"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:subresource:status

// The `Compactor` custom resource definition (CRD) defines a desired [Compactor](https://thanos.io/tip/components/compactor.md/) setup to run in a Kubernetes cluster. It allows to specify many options such as the number of replicas, persistent storage and many more.
//
// For each `Compactor` resource, the Operator deploys a `StatefulSet` in the same namespace.
type Compactor struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CompactorSpec   `json:"spec,omitempty"`
	Status CompactorStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// CompactorList contains a list of Compactor
type CompactorList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Compactor `json:"items"`
}

// GatewaySpec defines the desired state of Gateway
type GatewaySpec struct {
	// Defines the configuration of the Gatewat web server.
	WebConfig *WebConfig `json:"webConfig,omitempty"`

	// If debug mode is on, gateway will proxy Query UI
	//
	// This is an *experimental feature*, it may change in any upcoming release in a breaking way.
	//
	DebugMode bool `json:"debug,omitempty"`

	// Deny unknown tenant data remote-write and query if enabled
	EnabledTenantsAdmission bool `json:"enabledTenantsAdmission,omitempty"`

	// NodePort is the port used to expose the gateway service.
	// If this is a valid node port, the gateway service type will be set to NodePort accordingly.
	NodePort int32 `json:"nodePort,omitempty"`

	CommonSpec `json:",inline"`
}

// GatewayStatus defines the observed state of Gateway
type GatewayStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +genclient
// +kubebuilder:object:root=true
// +kubebuilder:printcolumn:name="Replicas",type="integer",JSONPath=".spec.replicas",description="The number of desired replicas"
// +kubebuilder:printcolumn:name="NodePort",type="integer",JSONPath=".spec.nodePort",description="The nodePort of Gateway service"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:subresource:status

// The `Gateway` custom resource definition (CRD) defines a desired [Gateway](https://github.com/WhizardTelemetry/whizard-docs/blob/main/Architecture/components/whizard-monitoring-gateway.md) setup to run in a Kubernetes cluster. It allows to specify many options such as the number of replicas, and many more.
//
// For each `Gateway` resource, the Operator deploys a `Deployment` in the same namespace.
type Gateway struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   GatewaySpec   `json:"spec,omitempty"`
	Status GatewayStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// GatewayList contains a list of Gateway
type GatewayList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Gateway `json:"items"`
}

// IngesterSpec defines the desired state of Ingester
type IngesterSpec struct {
	// The tenants whose data is being ingested by the Ingester(ingesting receiver).
	Tenants []string `json:"tenants,omitempty"`

	// Enables target information in OTLP metrics ingested by Receive. If enabled, it converts the resource to the target info metric
	OtlpEnableTargetInfo *bool `json:"otlpEnableTargetInfo,omitempty"`

	// Resource attributes to include in OTLP metrics ingested by Receive.
	OtlpResourceAttributes []string `json:"otlpResourceAttributes,omitempty"`

	// LocalTsdbRetention configs how long to retain raw samples on local storage.
	LocalTsdbRetention string `json:"localTsdbRetention,omitempty"`

	// DataVolume specifies how volume shall be used
	DataVolume *KubernetesVolume `json:"dataVolume,omitempty"`

	CommonSpec `json:",inline"`

	IngesterTSDBCleanUp SidecarSpec `json:"ingesterTsdbCleanup,omitempty"`
}

// IngesterStatus defines the observed state of Ingester
type IngesterStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Tenants contain all tenants that have been configured for this Ingester object,
	// except those Tenant objects that have been deleted.
	Tenants []IngesterTenantStatus `json:"tenants,omitempty"`
}

type IngesterTenantStatus struct {
	Name string `json:"name"`
	// true represents that the tenant has been moved to other ingester but may left tsdb data in this ingester.
	Obsolete bool `json:"obsolete"`
}

// +genclient
// +kubebuilder:object:root=true
// +kubebuilder:printcolumn:name="Replicas",type="integer",JSONPath=".spec.replicas",description="The number of desired replicas"
// +kubebuilder:printcolumn:name="LocalTsdbRetention",type="string",JSONPath=".spec.localTsdbRetention",description="How long to retain raw samples on local storage."
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:subresource:status

// The `Ingester` custom resource definition (CRD) defines a desired [Ingesting Receive](https://thanos.io/tip/components/receive.md/) setup to run in a Kubernetes cluster. It allows to specify many options such as the number of replicas, persistent storage and many more.
//
// For each `Ingester` resource, the Operator deploys a `StatefulSet` in the same namespace.
type Ingester struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   IngesterSpec   `json:"spec,omitempty"`
	Status IngesterStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// IngesterList contains a list of Ingester
type IngesterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Ingester `json:"items"`
}

// QuerySpec defines the desired state of Query
type QuerySpec struct {
	// experimental PromQL engine, more info thanos.io/tip/components/query.md#promql-engine
	// default: prometheus
	// +kubebuilder:validation:Enum="";prometheus;thanos
	PromqlEngine string `json:"promqlEngine,omitempty"`

	// Selector labels that will be exposed in info endpoint.
	SelectorLabels map[string]string `json:"selectorLabels,omitempty"`

	// Labels to treat as a replica indicator along which data is deduplicated.
	ReplicaLabelNames []string `json:"replicaLabelNames,omitempty"`

	// Defines the configuration of the Thanos Query web server.
	WebConfig *WebConfig `json:"webConfig,omitempty"`

	// Additional StoreApi servers from which Query component queries from
	Stores []QueryStores `json:"stores,omitempty"`
	// Envoy is used to config sidecar which proxies requests requiring auth to the secure stores
	Envoy SidecarSpec `json:"envoy,omitempty"`

	CommonSpec `json:",inline"`
}

type QueryStores struct {
	// Address is the addresses of StoreApi server, which may be prefixed with 'dns+' or 'dnssrv+' to detect StoreAPI servers through respective DNS lookups.
	Addresses []string `json:"addresses,omitempty"`
	// Secret containing the CA cert to use for StoreApi connections
	CASecret *corev1.SecretKeySelector `json:"caSecret,omitempty"`
}

// QueryStatus defines the observed state of Query
type QueryStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +genclient
// +kubebuilder:object:root=true
// +kubebuilder:printcolumn:name="Replicas",type="integer",JSONPath=".spec.replicas",description="The number of desired replicas"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:subresource:status

// The `Query` custom resource definition (CRD) defines a desired [Query](https://thanos.io/tip/components/query.md/) setup to run in a Kubernetes cluster. It allows to specify many options such as the number of replicas, and many more.
//
// For each `Query` resource, the Operator deploys a `Deployment` in the same namespace.
type Query struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   QuerySpec   `json:"spec,omitempty"`
	Status QueryStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// QueryList contains a list of Query
type QueryList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Query `json:"items"`
}

// QueryFrontendSpec defines the desired state of QueryFrontend
type QueryFrontendSpec struct {
	// CacheProviderConfig specifies response cache configuration.
	CacheConfig *ResponseCacheProviderConfig `json:"cacheConfig,omitempty"`

	// Defines the configuration of the Thanos QueryFrontend web server.
	WebConfig *WebConfig `json:"webConfig,omitempty"`

	CommonSpec `json:",inline"`
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

// QueryFrontendStatus defines the observed state of QueryFrontend
type QueryFrontendStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +genclient
// +kubebuilder:object:root=true
// +kubebuilder:printcolumn:name="Replicas",type="integer",JSONPath=".spec.replicas",description="The number of desired replicas"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:subresource:status

// The `QueryFrontend` custom resource definition (CRD) defines a desired [QueryFrontend](https://thanos.io/tip/components/query-frontend.md/) setup to run in a Kubernetes cluster. It allows to specify many options such as the number of replicas, and many more.
//
// For each `QueryFrontend` resource, the Operator deploys a `Deployment` in the same namespace.
type QueryFrontend struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   QueryFrontendSpec   `json:"spec,omitempty"`
	Status QueryFrontendStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// QueryFrontendList contains a list of QueryFrontend
type QueryFrontendList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []QueryFrontend `json:"items"`
}

// RouterSpec defines the desired state of Router
type RouterSpec struct {

	// How many times to replicate incoming write requests
	ReplicationFactor *uint64 `json:"replicationFactor,omitempty"`

	// The protocol to use for replicating remote-write requests. One of protobuf,capnproto
	// +kubebuilder:validation:Enum=protobuf;capnproto
	ReplicationProtocol ReplicationProtocol `json:"replicationProtocol,omitempty"`

	// Defines the configuration of the Route(routing receiver) web server.
	WebConfig *WebConfig `json:"webConfig,omitempty"`

	CommonSpec `json:",inline"`
}

type ReplicationProtocol string

const (
	ProtobufReplication  ReplicationProtocol = "protobuf"
	CapNProtoReplication ReplicationProtocol = "capnproto"
)

// RouterStatus defines the observed state of Router
type RouterStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +genclient
// +kubebuilder:object:root=true
// +kubebuilder:printcolumn:name="Replicas",type="integer",JSONPath=".spec.replicas",description="The number of desired replicas"
// +kubebuilder:printcolumn:name="ReplicationFactor",type="integer",JSONPath=".spec.replicationFactor",description="How many times to replicate incoming write requests"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:subresource:status

// The `Router` custom resource definition (CRD) defines a desired [Routing Receivers](https://thanos.io/tip/components/receive.md/) setup to run in a Kubernetes cluster. It allows to specify many options such as the number of replicas, and many more.
//
// For each `Router` resource, the Operator deploys a `Deployment` in the same namespace.
type Router struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RouterSpec   `json:"spec,omitempty"`
	Status RouterStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// RouterList contains a list of Router
type RouterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Router `json:"items"`
}

// RulerSpec defines the desired state of Ruler
type RulerSpec struct {

	// Label selectors to select which PrometheusRules to mount for alerting and recording.
	// The result of multiple selectors are ORed.
	RuleSelectors []*metav1.LabelSelector `json:"ruleSelectors,omitempty"`
	// Namespaces to be selected for PrometheusRules discovery. If unspecified, only
	// the same namespace as the Ruler object is in is used.
	RuleNamespaceSelector *metav1.LabelSelector `json:"ruleNamespaceSelector,omitempty"`

	// Number of shards to take the hash of fully qualified name of the rule group in order to split rules.
	// Each shard of rules will be bound to one separate statefulset.
	// Default: 1
	// +kubebuilder:default:=1
	Shards *int32 `json:"shards,omitempty"`

	// Tenant if not empty indicates which tenant's data is evaluated for the selected rules;
	// otherwise, it is for all tenants.
	Tenant string `json:"tenant,omitempty"`

	QueryConfig *corev1.SecretKeySelector `json:"queryConfig,omitempty"`

	RemoteWriteConfig *corev1.SecretKeySelector `json:"remoteWriteConfig,omitempty"`

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
	//
	// Default: "1m"
	// +kubebuilder:default:="1m"
	EvaluationInterval Duration `json:"evaluationInterval,omitempty"`

	// DataVolume specifies how volume shall be used
	DataVolume *KubernetesVolume `json:"dataVolume,omitempty"`

	RulerQueryProxy          SidecarSpec `json:"rulerQueryProxy,omitempty"`
	RulerWriteProxy          SidecarSpec `json:"rulerWriteProxy,omitempty"`
	PrometheusConfigReloader SidecarSpec `json:"prometheusConfigReloader,omitempty"`

	CommonSpec `json:",inline"`
}

// RulerStatus defines the observed state of Ruler
type RulerStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +genclient
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// Ruler is the Schema for the rulers API
type Ruler struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RulerSpec   `json:"spec,omitempty"`
	Status RulerStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// RulerList contains a list of Ruler
type RulerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Ruler `json:"items"`
}

// StoreSpec defines the desired state of Store
type StoreSpec struct {

	// MinTime specifies start of time range limit to serve
	MinTime string `json:"minTime,omitempty"`
	// MaxTime specifies end of time range limit to serve
	MaxTime string `json:"maxTime,omitempty"`
	// TimeRanges is a list of TimeRange to partition Store.
	// If specified, the MinTime and MaxTime will be ignored.
	TimeRanges []TimeRange `json:"timeRanges,omitempty"`

	// IndexCacheConfig contains index cache configuration.
	IndexCacheConfig *IndexCacheConfig `json:"indexCacheConfig,omitempty"`

	// DataVolume specifies how volume shall be used
	DataVolume *KubernetesVolume `json:"dataVolume,omitempty"`

	CommonSpec `json:",inline"`
}

type TimeRange struct {
	// MinTime specifies start of time range limit to serve
	MinTime string `json:"minTime,omitempty" yaml:"minTime,omitempty"`
	// MaxTime specifies end of time range limit to serve
	MaxTime string `json:"maxTime,omitempty" yaml:"maxTime,omitempty"`
}

type IndexCacheConfig struct {
	*InMemoryIndexCacheConfig `json:"inMemory,omitempty" yaml:"inMemory,omitempty"`
}

type InMemoryIndexCacheConfig struct {
	// MaxSize represents overall maximum number of bytes cache can contain.
	MaxSize string `json:"maxSize" yaml:"maxSize"`
}

// StoreStatus defines the observed state of Store
type StoreStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +genclient
// +kubebuilder:object:root=true
// +kubebuilder:printcolumn:name="Replicas",type="integer",JSONPath=".spec.replicas",description="The number of desired replicas"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:subresource:status

// The `Store` custom resource definition (CRD) defines a desired [Compactor](https://thanos.io/tip/components/store.md/) setup to run in a Kubernetes cluster. It allows to specify many options such as the number of replicas, persistent storage and many more.
//
// For each `Store` resource, the Operator deploys a `StatefulSet` in the same namespace.
type Store struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   StoreSpec   `json:"spec,omitempty"`
	Status StoreStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// StoreList contains a list of Store
type StoreList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Store `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Compactor{}, &CompactorList{})
	SchemeBuilder.Register(&Gateway{}, &GatewayList{})
	SchemeBuilder.Register(&Ingester{}, &IngesterList{})
	SchemeBuilder.Register(&Query{}, &QueryList{})
	SchemeBuilder.Register(&QueryFrontend{}, &QueryFrontendList{})
	SchemeBuilder.Register(&Router{}, &RouterList{})
	SchemeBuilder.Register(&Ruler{}, &RulerList{})
	SchemeBuilder.Register(&Store{}, &StoreList{})
}

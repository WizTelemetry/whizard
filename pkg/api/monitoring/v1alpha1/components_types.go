/*
Copyright 2021 The WhizardTelemetry Authors.

This program is free software: you can redistribute it and/or modify
it under the terms of the Server Side Public License, version 1,
as published by MongoDB, Inc.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
Server Side Public License for more details.

You should have received a copy of the Server Side Public License
along with this program. If not, see
<http://www.mongodb.com/licensing/server-side-public-license>.

As a special exception, the copyright holders give permission to link the
code of portions of this program with the OpenSSL library under certain
conditions as described in each individual source file and distribute
linked combinations including the program with the OpenSSL library. You
must comply with the Server Side Public License in all respects for
all of the code used other than as permitted herein. If you modify file(s)
with this exception, you may extend this exception to your version of the
file(s), but you are not obligated to do so. If you do not wish to do so,
delete this exception statement from your version. If you delete this
exception statement from all source files in the program, then also delete
it in the license file.
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
	CommonSpec `json:",inline"`

	// DisableDownsampling specifies whether to disable downsampling
	DisableDownsampling *bool `json:"disableDownsampling,omitempty"`

	// Retention configs how long to retain samples
	Retention *Retention `json:"retention,omitempty"`

	// DataVolume specifies how volume shall be used
	DataVolume *KubernetesVolume `json:"dataVolume,omitempty"`

	// Tenants if not empty indicates current config is for hard tenants; otherwise, it is for soft tenants.
	Tenants []string `json:"tenants,omitempty"`
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

// CompactorStatus defines the observed state of Compactor
type CompactorStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +genclient
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// Compactor is the Schema for the compactors API
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
	CommonSpec `json:",inline"`

	WebConfig *WebConfig `json:"webConfig,omitempty"`

	// If debug mode is on, gateway will proxy Query UI
	DebugMode bool `json:"debug,omitempty"`

	// Deny unknown tenant data remote-write and query if enabled
	EnabledTenantsAdmission bool `json:"enabledTenantsAdmission,omitempty"`

	// NodePort is the port used to expose the gateway service.
	// If this is a valid node port, the gateway service type will be set to NodePort accordingly.
	NodePort int32 `json:"nodePort,omitempty"`
}

// GatewayStatus defines the observed state of Gateway
type GatewayStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +genclient
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// Gateway is the Schema for the gateways API
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
	CommonSpec `json:",inline"`

	IngesterTSDBCleanUp SidecarSpec `json:"ingesterTsdbCleanup,omitempty"`

	// Tenants if not empty indicates current config is for hard tenants; otherwise, it is for soft tenants.
	Tenants []string `json:"tenants,omitempty"`

	// LocalTsdbRetention configs how long to retain raw samples on local storage.
	LocalTsdbRetention string `json:"localTsdbRetention,omitempty"`

	// DataVolume specifies how volume shall be used
	DataVolume *KubernetesVolume `json:"dataVolume,omitempty"`
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
// +kubebuilder:subresource:status

// Ingester is the Schema for the ingesters API
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
	CommonSpec `json:",inline"`

	WebConfig *WebConfig `json:"webConfig,omitempty"`

	PromqlEngine string `json:"promqlEngine,omitempty"`

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

// QueryStatus defines the observed state of Query
type QueryStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +genclient
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// Query is the Schema for the queries API
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
	CommonSpec `json:",inline"`

	WebConfig *WebConfig `json:"webConfig,omitempty"`

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
// +kubebuilder:subresource:status

// QueryFrontend is the Schema for the queryfrontends API
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
	CommonSpec `json:",inline"`

	WebConfig *WebConfig `json:"webConfig,omitempty"`

	// How many times to replicate incoming write requests
	ReplicationFactor *uint64 `json:"replicationFactor,omitempty"`
}

// RouterStatus defines the observed state of Router
type RouterStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +genclient
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// Router is the Schema for the routers API
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
	CommonSpec `json:",inline"`

	RulerQueryProxy          SidecarSpec `json:"rulerQueryProxy,omitempty"`
	RulerWriteProxy          SidecarSpec `json:"rulerWriteProxy,omitempty"`
	PrometheusConfigReloader SidecarSpec `json:"prometheusConfigReloader,omitempty"`

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
	CommonSpec `json:",inline"`

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
// +kubebuilder:subresource:status

// Store is the Schema for the stores API
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

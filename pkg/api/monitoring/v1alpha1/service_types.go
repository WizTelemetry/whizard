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

// Package v1alpha1 contains API Schema definitions for the monitoring v1alpha1 API group
// +kubebuilder:object:generate=true
// +groupName=monitoring.whizard.io
package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ServiceSpec defines the desired state of a Service
type ServiceSpec struct {

	// HTTP header to determine tenant for remote write requests.
	TenantHeader string `json:"tenantHeader,omitempty"`
	// Default tenant ID to use when none is provided via a header.
	DefaultTenantId string `json:"defaultTenantId,omitempty"`
	// Label name through which the tenant will be announced.
	TenantLabelName string `json:"tenantLabelName,omitempty"`

	Storage *ObjectReference `json:"storage,omitempty"`

	// RemoteWrites is the list of remote write configurations.
	// If it is configured, its targets will receive write requests from the Gateway and the Ruler.
	RemoteWrites []RemoteWriteSpec `json:"remoteWrites,omitempty"`
	// RemoteQuery is the remote query configuration and the remote target should have prometheus-compatible Query APIs.
	// If not configured, the Gateway will proxy all read requests through the QueryFrontend to the Query,
	// If configured, the Gateway will proxy metrics read requests through the QueryFrontend to the remote target,
	// but proxy rules read requests directly to the Query.
	RemoteQuery *RemoteQuerySpec `json:"remoteQuery,omitempty"`

	GatewayTemplateSpec       GatewaySpec           `json:"gatewayTemplateSpec"`
	QueryFrontendTemplateSpec QueryFrontendSpec     `json:"queryFrontendTemplateSpec"`
	QueryTemplateSpec         QuerySpec             `json:"queryTemplateSpec"`
	RulerTemplateSpec         RulerTemplateSpec     `json:"rulerTemplateSpec"`
	RouterTemplateSpec        RouterSpec            `json:"routerTemplateSpec"`
	IngesterTemplateSpec      IngesterTemplateSpec  `json:"ingesterTemplateSpec"`
	StoreTemplateSpec         StoreSpec             `json:"storeTemplateSpec"`
	CompactorTemplateSpec     CompactorTemplateSpec `json:"compactorTemplateSpec"`
}

type IngesterTemplateSpec struct {
	IngesterSpec `json:",inline"`

	// DefaultTenantsPerIngester Whizard default tenant count per ingester.
	//
	// Default: 3
	// +kubebuilder:default:=3
	DefaultTenantsPerIngester int `json:"defaultTenantsPerIngester,omitempty"`
	// DefaultIngesterRetentionPeriod Whizard default ingester retention period when it has no tenant.
	//
	// Default: "3h"
	// +kubebuilder:default:="3h"
	DefaultIngesterRetentionPeriod Duration `json:"defaultIngesterRetentionPeriod,omitempty"`
	// DisableTSDBCleanup Disable the TSDB cleanup of ingester.
	// The cleanup will delete the blocks that belong to deleted tenants in the data directory of ingester TSDB.
	//
	// Default: true
	// +kubebuilder:default:=true
	DisableTSDBCleanup *bool `json:"disableTsdbCleanup,omitempty"`
}

type RulerTemplateSpec struct {
	RulerSpec `json:",inline"`

	// DisableAlertingRulesAutoSelection disable auto select alerting rules in tenant ruler
	//
	// Default: true
	// +kubebuilder:default:=true
	DisableAlertingRulesAutoSelection *bool `json:"disableAlertingRulesAutoSelection,omitempty"`
}

type CompactorTemplateSpec struct {
	CompactorSpec `json:",inline"`

	// DefaultTenantsPerIngester Whizard default tenant count per ingester.
	// Default: 10
	// +kubebuilder:default:=10
	DefaultTenantsPerCompactor int `json:"defaultTenantsPerCompactor,omitempty"`
}

// RemoteQuerySpec defines the configuration to query from remote service
// which should have prometheus-compatible Query APIs.
type RemoteQuerySpec struct {
	Name             string `json:"name,omitempty"`
	URL              string `json:"url"`
	HTTPClientConfig `json:",inline"`
}

// RemoteWriteSpec defines the remote write configuration.
type RemoteWriteSpec struct {
	Name string `json:"name,omitempty"`
	URL  string `json:"url"`
	// Custom HTTP headers to be sent along with each remote write request.
	Headers map[string]string `json:"headers,omitempty"`
	// Timeout for requests to the remote write endpoint.
	RemoteTimeout Duration `json:"remoteTimeout,omitempty"`

	HTTPClientConfig `json:",inline"`
}

// HTTPClientConfig configures an HTTP client.
type HTTPClientConfig struct {
	// The HTTP basic authentication credentials for the targets.
	BasicAuth BasicAuth `json:"basicAuth,omitempty"`
	// The bearer token for the targets.
	BearerToken string `json:"bearerToken,omitempty"`
}

// ServiceStatus defines the observed state of Service
type ServiceStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +genclient
// +k8s:openapi-gen=true
// +kubebuilder:object:root=true
// +kubebuilder:printcolumn:name="Storage",type="string",JSONPath=".spec.storage.name",description="The storage for service"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:subresource:status

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
	SchemeBuilder = SchemeBuilder.
		Register(&Service{}, &ServiceList{})
}

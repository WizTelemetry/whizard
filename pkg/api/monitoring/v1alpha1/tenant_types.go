/*
Copyright 2024 the Whizard Authors.

Licensed under Apache License, Version 2.0 with a few additional conditions.

You may obtain a copy of the License at

    https://github.com/WhizardTelemetry/whizard/blob/main/LICENSE
*/

package v1alpha1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// TenantSpec defines the desired state of Tenant
type TenantSpec struct {
	Tenant string `json:"tenant,omitempty"`
}

// TenantStatus defines the observed state of Tenant
type TenantStatus struct {
	Ruler     *ObjectReference `json:"ruler,omitempty"`
	Compactor *ObjectReference `json:"compactor,omitempty"`
	Ingester  *ObjectReference `json:"ingester,omitempty"`
}

// +genclient
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster

// The `Tenant` custom resource definition (CRD) defines the tenant configuration for multi-tenant data separation in Whizard.
// In Whizard, a tenant can represent various types of data sources, such as:
//
// - Monitoring data from a specific Kubernetes cluster
// - Monitoring data from a physical machine in a specific region
// - Monitoring data from a specific type of application
//
// When data is ingested, it will be tagged with the tenant label to ensure proper separation.
type Tenant struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   TenantSpec   `json:"spec,omitempty"`
	Status TenantStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// TenantList contains a list of Tenant
type TenantList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Tenant `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Tenant{}, &TenantList{})
}

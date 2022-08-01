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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type TenantSpec struct {
	Tenant  string           `json:"tenant,omitempty"`
	Storage *ObjectReference `json:"storage,omitempty"`
}

type TenantStatus struct {
	Store     *ObjectReference `json:"store,omitempty"`
	Ruler     *ObjectReference `json:"ruler,omitempty"`
	Compactor *ObjectReference `json:"compactor,omitempty"`
	Ingester  *ObjectReference `json:"ingester,omitempty"`
}

//+genclient
//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:scope=Cluster

// Tenant is the Schema for the monitoring Tenant API
type Tenant struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   TenantSpec   `json:"spec,omitempty"`
	Status TenantStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:resource:scope=Cluster

type TenantList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Tenant `json:"items"`
}

/*
	//+genclient
	//+kubebuilder:object:root=true
	//+kubebuilder:subresource:status
	//+kubebuilder:resource:scope=Cluster

	// Tenant is the Schema for the monitoring TenantGroup API
	type TenantGroup struct {
		metav1.TypeMeta   `json:",inline"`
		metav1.ObjectMeta `json:"metadata,omitempty"`

		Spec   TenantGroupSpec   `json:"spec,omitempty"`
		Status TenantGroupStatus `json:"status,omitempty"`
	}

	type TenantGroupSpec struct {
		TenantGroupId string           `json:"tenantGroupId"`
		PaodinService *ObjectReference `json:"paodinService"`
		PaodinStorage *ObjectReference `json:"paodinStorage"`
	}

	type TenantGroupStatus struct {
		Members []string
	}
*/
func init() {
	SchemeBuilder = SchemeBuilder.
		Register(&Tenant{}, &TenantList{})
}

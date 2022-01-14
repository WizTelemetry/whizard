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

// ThanosQuerySpec defines the desired state of ThanosQuery
type ThanosQuerySpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// If specified, the pod's scheduling constraints.
	Affinity *corev1.Affinity `json:"affinity,omitempty"`
	// Define which Nodes the Pods are scheduled on.
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`
	// If specified, the pod's tolerations.
	Tolerations []corev1.Toleration `json:"tolerations,omitempty"`
	// Define resources requests and limits for main container.
	Resources corev1.ResourceRequirements `json:"resources,omitempty"`

	// Image is the thanos image with tag/version
	Image           string            `json:"image,omitempty"`
	ImagePullPolicy corev1.PullPolicy `json:"imagePullPolicy,omitempty"`
	// Number of replicas for a thanos query component
	Replicas *int32 `json:"replicas,omitempty"`

	// Stores config store api servers from where series are queried
	Stores []QueryStore `json:"stores,omitempty"`
	// SelectorLabels config query selector labels that will be exposed in info endpoint
	SelectorLabels map[string]string `json:"selectorLabels,omitempty"`

	// HttpIngress configs http request entry from services outside the cluster
	HttpIngress *IngressSpec `json:"httpIngress,omitempty"`
	// HttpIngress configs grpc request entry from services outside the cluster
	GrpcIngress *IngressSpec `json:"grpcIngress,omitempty"`

	// LogLevel configs log filtering level. Possible options: error, warn, info, debug
	LogLevel string `json:"level,omitempty"`
	// LogFormat configs log format to use. Possible options: logfmt or json
	LogFormat string `json:"format,omitempty"`

	// Envoy is used to config envoy sidecar which proxies requests to the secure stores
	Envoy *EnvoySpec `json:"envoy,omitempty"`
}

// ThanosQueryStatus defines the observed state of ThanosQuery
type ThanosQueryStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// ThanosQuery is the Schema for the thanosqueries API
type ThanosQuery struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ThanosQuerySpec   `json:"spec,omitempty"`
	Status ThanosQueryStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ThanosQueryList contains a list of ThanosQuery
type ThanosQueryList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ThanosQuery `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ThanosQuery{}, &ThanosQueryList{})
}

type QueryStore struct {
	// Address is the address of a store api server, which may be prefixed with 'dns+' or 'dnssrv+' to detect
	// store API servers through respective DNS lookups. For more info, see https://thanos.io/tip/thanos/service-discovery.md/#dns-service-discovery
	Address    string `json:"address,omitempty"`
	SecretName string `json:"secretName,omitempty"`
}

type IngressSpec struct {
	Host       string `json:"host,omitempty"`
	Path       string `json:"path,omitempty"`
	SecretName string `json:"secretName,omitempty"`
}

// EnvoySpec defines the desired state of envoy proxy sidecar which delegates requests to the secure thanos stores
type EnvoySpec struct {
	// Image is the thanos image with tag/version
	Image           string            `json:"image,omitempty"`
	ImagePullPolicy corev1.PullPolicy `json:"imagePullPolicy,omitempty"`
	// Define resources requests and limits for envoy container.
	Resources corev1.ResourceRequirements `json:"resources,omitempty"`
}

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

// ThanosStorageSpec defines the desired state of ThanosStorage
type ThanosStorageSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// If specified, the pod's scheduling constraints.
	Affinity *corev1.Affinity `json:"affinity,omitempty"`
	// Define which Nodes the Pods are scheduled on.
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`
	// If specified, the pod's tolerations.
	Tolerations []corev1.Toleration `json:"tolerations,omitempty"`

	// Image is the thanos image with tag/version
	Image           string            `json:"image,omitempty"`
	ImagePullPolicy corev1.PullPolicy `json:"imagePullPolicy,omitempty"`

	// ObjectStorageConfig allows specifying a key of a Secret containing object store configuration
	ObjectStorageConfig *corev1.SecretKeySelector `json:"objectStorageConfig,omitempty"`

	// LogLevel configs log filtering level. Possible options: error, warn, info, debug
	LogLevel  string `json:"level,omitempty"`
	// LogFormat configs log format to use. Possible options: logfmt or json
	LogFormat string `json:"format,omitempty"`

	// Gateway specifies the configs for thanos store gateway
	Gateway *ThanosStoreGatewaySpec `json:"gateway,omitempty"`
	// Compact specifies the configs for thanos compact
	Compact *ThanosCompactSpec      `json:"compact,omitempty"`
}

// ThanosStorageStatus defines the observed state of ThanosStorage
type ThanosStorageStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// ThanosStorage is the Schema for the thanosstorages API
type ThanosStorage struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ThanosStorageSpec   `json:"spec,omitempty"`
	Status ThanosStorageStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ThanosStorageList contains a list of ThanosStorage
type ThanosStorageList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ThanosStorage `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ThanosStorage{}, &ThanosStorageList{})
}

// ThanosStoreGatewaySpec defines the desired state of thanos store gateway
type ThanosStoreGatewaySpec struct {
	// Define resources requests and limits for single Pods.
	Resources corev1.ResourceRequirements `json:"resources,omitempty"`
	Replicas   *int32            `json:"replicas,omitempty"`
	// DataVolume specifies how volume shall be used
	DataVolume *KubernetesVolume `json:"dataVolume,omitempty"`
	// MinTime specifies start of time range limit to serve
	MinTime    string            `json:"minTime,omitempty"`
	// MaxTime specifies end of time range limit to serve
	MaxTime    string            `json:"maxTime,omitempty"`
}

// ThanosCompactSpec defines the desired state of thanos compact
type ThanosCompactSpec struct {
	// Define resources requests and limits for single Pods.
	Resources corev1.ResourceRequirements `json:"resources,omitempty"`
	// DownsamplingDisable specifies whether to disable downsampling
	DownsamplingDisable *bool                   `json:"downsamplingDisable,omitempty"`
	// Retention configs how long to retain samples
	Retention           *ThanosStorageRetention `json:"retention,omitempty"`
	// DataVolume specifies how volume shall be used
	DataVolume          *KubernetesVolume       `json:"dataVolume,omitempty"`
}

// ThanosStorageRetention defines the config for retaining samples
type ThanosStorageRetention struct {
	// RetentionRaw specifies how long to retain raw samples in bucket
	RetentionRaw     string `json:"retentionRaw,omitempty"`
	// Retention5m specifies how long to retain samples of 5m resolution in bucket
	Retention5m string `json:"retention5m,omitempty"`
	// Retention1h specifies how long to retain samples of 1h resolution in bucket
	Retention1h string `json:"retention1h,omitempty"`
}

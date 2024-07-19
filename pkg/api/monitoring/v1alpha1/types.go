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
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// Duration is a valid time unit
// Supported units: y, w, d, h, m, s, ms Examples: `30s`, `1m`, `1h20m15s`
// +kubebuilder:validation:Pattern:="^(0|(([0-9]+)y)?(([0-9]+)w)?(([0-9]+)d)?(([0-9]+)h)?(([0-9]+)m)?(([0-9]+)s)?(([0-9]+)ms)?)$"
type Duration string

type CommonSpec struct {
	// Containers allows injecting additional containers or modifying operator generated containers.
	// Containers described here modify an operator generated
	// container if they share the same name and modifications are done via a
	// strategic merge patch.
	Containers runtime.RawExtension `json:"containers,omitempty"`

	// EmbeddedContainers
	EmbeddedContainers []corev1.Container `json:"-"`

	// PodMetadata configures labels and annotations which are propagated to the pods.
	//
	// * "kubectl.kubernetes.io/default-container" annotation, set to main pod.
	PodMetadata *EmbeddedObjectMetadata `json:"podMetadata,omitempty"`

	// Secrets is a list of Secrets in the same namespace as the component
	// object, which shall be mounted into the Prometheus Pods.
	// Each Secret is added to the StatefulSet/Deployment definition as a volume named `secret-<secret-name>`.
	// The Secrets are mounted into /etc/whizard/secrets/<secret-name> in the default container.
	Secrets []string `json:"secrets,omitempty"`
	// ConfigMaps is a list of ConfigMaps in the same namespace as the component
	// object, which shall be mounted into the default Pods.
	// Each ConfigMap is added to the StatefulSet/Deployment definition as a volume named `configmap-<configmap-name>`.
	// The ConfigMaps are mounted into /etc/whizard/configmaps/<configmap-name> in the default container.
	ConfigMaps []string `json:"configMaps,omitempty"`

	// If specified, the pod's scheduling constraints.
	Affinity *corev1.Affinity `json:"affinity,omitempty"`
	// Define which Nodes the Pods are scheduled on.
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`
	// If specified, the pod's tolerations.
	Tolerations []corev1.Toleration `json:"tolerations,omitempty"`
	// Define resources requests and limits for main container.
	Resources corev1.ResourceRequirements `json:"resources,omitempty"`
	// SecurityContext holds pod-level security attributes and common container settings.
	// This defaults to the default PodSecurityContext.
	// +optional
	SecurityContext *corev1.PodSecurityContext `json:"securityContext,omitempty"`
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

// KubernetesVolume defines the configured volume for a instance.
type KubernetesVolume struct {
	EmptyDir              *corev1.EmptyDirVolumeSource  `json:"emptyDir,omitempty"`
	PersistentVolumeClaim *corev1.PersistentVolumeClaim `json:"persistentVolumeClaim,omitempty"`
	// persistentVolumeClaimRetentionPolicy describes the lifecycle of persistent
	// volume claims created from persistentVolumeClaim.
	// This requires the kubernetes version >= 1.23 and its StatefulSetAutoDeletePVC feature gate to be enabled.
	// +optional
	PersistentVolumeClaimRetentionPolicy *appsv1.StatefulSetPersistentVolumeClaimRetentionPolicy `json:"persistentVolumeClaimRetentionPolicy,omitempty"`
}

// EmbeddedObjectMetadata contains a subset of the fields included in k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta
// Only fields which are relevant to embedded resources are included.
type EmbeddedObjectMetadata struct {
	// Name must be unique within a namespace. Is required when creating resources, although
	// some resources may allow a client to request the generation of an appropriate name
	// automatically. Name is primarily intended for creation idempotence and configuration
	// definition.
	// Cannot be updated.
	// More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names#names
	// +optional
	Name string `json:"name,omitempty" protobuf:"bytes,1,opt,name=name"`

	// Map of string keys and values that can be used to organize and categorize
	// (scope and select) objects. May match selectors of replication controllers
	// and services.
	// More info: http://kubernetes.io/docs/user-guide/labels
	// +optional
	Labels map[string]string `json:"labels,omitempty" protobuf:"bytes,11,rep,name=labels"`

	// Annotations is an unstructured key value map stored with a resource that may be
	// set by external tools to store and retrieve arbitrary metadata. They are not
	// queryable and should be preserved when modifying objects.
	// More info: http://kubernetes.io/docs/user-guide/annotations
	// +optional
	Annotations map[string]string `json:"annotations,omitempty" protobuf:"bytes,12,rep,name=annotations"`
}

type SidecarSpec struct {
	// Image is the envoy image with tag/version
	Image string `json:"image,omitempty" yaml:"image,omitempty"`
	// Define resources requests and limits for sidecar container.
	Resources corev1.ResourceRequirements `json:"resources,omitempty" yaml:"resources,omitempty"`
}

type ObjectReference struct {
	Namespace string `json:"namespace,omitempty"`
	Name      string `json:"name,omitempty"`
}

type WebConfig struct {
	HTTPServerTLSConfig *HTTPServerTLSConfig `json:"httpServerTLSConfig,omitempty"`
	HTTPServerConfig    *HTTPServerConfig    `json:"httpServerConfig,omitempty"`
	BasicAuthUsers      []BasicAuth          `json:"basicAuthUsers,omitempty"`
}

// HTTPClientConfig configures an HTTP client.
type HTTPClientConfig struct {
	// The HTTP basic authentication credentials for the targets.
	BasicAuth BasicAuth `json:"basicAuth,omitempty"`
	// The bearer token for the targets.
	BearerToken string `json:"bearerToken,omitempty"`
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

// BasicAuth allow an endpoint to authenticate over basic authentication
// +k8s:openapi-gen=true
type BasicAuth struct {
	// The secret in the service monitor namespace that contains the username
	// for authentication.
	Username corev1.SecretKeySelector `json:"username,omitempty"`
	// The secret in the service monitor namespace that contains the password
	// for authentication.
	Password corev1.SecretKeySelector `json:"password,omitempty"`
}

type HTTPServerConfig struct {
}

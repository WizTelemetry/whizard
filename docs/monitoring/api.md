---
title: "API reference"
description: "Prometheus operator generated API reference docs"
draft: false
images: []
menu: "operator"
weight: 151
toc: true
---
> This page is automatically generated with `gen-crd-api-reference-docs`.
<p>Packages:</p>
<ul>
<li>
<a href="#monitoring.whizard.io%2fv1alpha1">monitoring.whizard.io/v1alpha1</a>
</li>
</ul>
<h2 id="monitoring.whizard.io/v1alpha1">monitoring.whizard.io/v1alpha1</h2>
<div>
</div>
Resource Types:
<ul><li>
<a href="#monitoring.whizard.io/v1alpha1.Compactor">Compactor</a>
</li><li>
<a href="#monitoring.whizard.io/v1alpha1.Gateway">Gateway</a>
</li><li>
<a href="#monitoring.whizard.io/v1alpha1.Ingester">Ingester</a>
</li><li>
<a href="#monitoring.whizard.io/v1alpha1.Query">Query</a>
</li><li>
<a href="#monitoring.whizard.io/v1alpha1.QueryFrontend">QueryFrontend</a>
</li><li>
<a href="#monitoring.whizard.io/v1alpha1.Router">Router</a>
</li><li>
<a href="#monitoring.whizard.io/v1alpha1.Ruler">Ruler</a>
</li><li>
<a href="#monitoring.whizard.io/v1alpha1.Service">Service</a>
</li><li>
<a href="#monitoring.whizard.io/v1alpha1.Storage">Storage</a>
</li><li>
<a href="#monitoring.whizard.io/v1alpha1.Store">Store</a>
</li><li>
<a href="#monitoring.whizard.io/v1alpha1.Tenant">Tenant</a>
</li></ul>
<h3 id="monitoring.whizard.io/v1alpha1.Compactor">Compactor
</h3>
<div>
<p>The <code>Compactor</code> custom resource definition (CRD) defines a desired <a href="https://thanos.io/tip/components/compactor.md/">Compactor</a> setup to run in a Kubernetes cluster. It allows to specify many options such as the number of replicas, persistent storage and many more.</p>
<p>For each <code>Compactor</code> resource, the Operator deploys a <code>StatefulSet</code> in the same namespace.</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>apiVersion</code><br/>
string</td>
<td>
<code>
monitoring.whizard.io/v1alpha1
</code>
</td>
</tr>
<tr>
<td>
<code>kind</code><br/>
string
</td>
<td><code>Compactor</code></td>
</tr>
<tr>
<td>
<code>metadata</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#objectmeta-v1-meta">
Kubernetes meta/v1.ObjectMeta
</a>
</em>
</td>
<td>
Refer to the Kubernetes API documentation for the fields of the
<code>metadata</code> field.
</td>
</tr>
<tr>
<td>
<code>spec</code><br/>
<em>
<a href="#monitoring.whizard.io/v1alpha1.CompactorSpec">
CompactorSpec
</a>
</em>
</td>
<td>
<br/>
<br/>
<table>
<tr>
<td>
<code>tenants</code><br/>
<em>
[]string
</em>
</td>
<td>
<p>The tenants whose data is being compacted by the Compactor.</p>
</td>
</tr>
<tr>
<td>
<code>disableDownsampling</code><br/>
<em>
bool
</em>
</td>
<td>
<p>Disables downsampling.
This is not recommended, as querying long time ranges without non-downsampled data is not efficient and useful.
default: false</p>
</td>
</tr>
<tr>
<td>
<code>retention</code><br/>
<em>
<a href="#monitoring.whizard.io/v1alpha1.Retention">
Retention
</a>
</em>
</td>
<td>
<p>Retention configs how long to retain samples</p>
</td>
</tr>
<tr>
<td>
<code>dataVolume</code><br/>
<em>
<a href="#monitoring.whizard.io/v1alpha1.KubernetesVolume">
KubernetesVolume
</a>
</em>
</td>
<td>
<p>DataVolume specifies how volume shall be used</p>
</td>
</tr>
<tr>
<td>
<code>replicas</code><br/>
<em>
int32
</em>
</td>
<td>
<p>Number of component instances to deploy.</p>
</td>
</tr>
<tr>
<td>
<code>image</code><br/>
<em>
string
</em>
</td>
<td>
<p>Component container image URL.</p>
</td>
</tr>
<tr>
<td>
<code>imagePullPolicy</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#pullpolicy-v1-core">
Kubernetes core/v1.PullPolicy
</a>
</em>
</td>
<td>
<p>Image pull policy.</p>
</td>
</tr>
<tr>
<td>
<code>resources</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#resourcerequirements-v1-core">
Kubernetes core/v1.ResourceRequirements
</a>
</em>
</td>
<td>
<p>Resources defines the resource requirements for single Pods.</p>
</td>
</tr>
<tr>
<td>
<code>logLevel</code><br/>
<em>
string
</em>
</td>
<td>
<p>Log level for component to be configured with.</p>
</td>
</tr>
<tr>
<td>
<code>logFormat</code><br/>
<em>
string
</em>
</td>
<td>
<p>Log format for component to be configured with.</p>
</td>
</tr>
<tr>
<td>
<code>flags</code><br/>
<em>
[]string
</em>
</td>
<td>
<p>Flags allows setting additional flags for the component container.</p>
</td>
</tr>
<tr>
<td>
<code>podMetadata</code><br/>
<em>
<a href="#monitoring.whizard.io/v1alpha1.EmbeddedObjectMetadata">
EmbeddedObjectMetadata
</a>
</em>
</td>
<td>
<p>PodMetadata configures labels and annotations which are propagated to the pods.</p>
</td>
</tr>
<tr>
<td>
<code>configMaps</code><br/>
<em>
[]string
</em>
</td>
<td>
<p>ConfigMaps is a list of ConfigMaps in the same namespace as the component
object, which shall be mounted into the default Pods.
Each ConfigMap is added to the StatefulSet/Deployment definition as a volume named <code>configmap-&lt;configmap-name&gt;</code>.
The ConfigMaps are mounted into /etc/whizard/configmaps/<configmap-name> in the default container.</p>
</td>
</tr>
<tr>
<td>
<code>secrets</code><br/>
<em>
[]string
</em>
</td>
<td>
<p>Secrets is a list of Secrets in the same namespace as the component
object, which shall be mounted into the Prometheus Pods.
Each Secret is added to the StatefulSet/Deployment definition as a volume named <code>secret-&lt;secret-name&gt;</code>.
The Secrets are mounted into /etc/whizard/secrets/<secret-name> in the default container.</p>
</td>
</tr>
<tr>
<td>
<code>containers</code><br/>
<em>
<a href="https://pkg.go.dev/k8s.io/apimachinery/pkg/runtime.#RawExtension">
k8s.io/apimachinery/pkg/runtime.RawExtension
</a>
</em>
</td>
<td>
<p>Containers allows injecting additional containers or modifying operator generated containers.
Containers described here modify an operator generated
container if they share the same name and modifications are done via a
strategic merge patch.</p>
</td>
</tr>
<tr>
<td>
<code>imagePullSecrets</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#localobjectreference-v1-core">
[]Kubernetes core/v1.LocalObjectReference
</a>
</em>
</td>
<td>
<p>An optional list of references to secrets in the same namespace
to use for pulling images from registries</p>
</td>
</tr>
<tr>
<td>
<code>securityContext</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#podsecuritycontext-v1-core">
Kubernetes core/v1.PodSecurityContext
</a>
</em>
</td>
<td>
<p>SecurityContext holds pod-level security attributes and common container settings.
This defaults to the default PodSecurityContext.</p>
</td>
</tr>
<tr>
<td>
<code>affinity</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#affinity-v1-core">
Kubernetes core/v1.Affinity
</a>
</em>
</td>
<td>
<p>If specified, the pod&rsquo;s scheduling constraints.</p>
</td>
</tr>
<tr>
<td>
<code>nodeSelector</code><br/>
<em>
map[string]string
</em>
</td>
<td>
<p>Define which Nodes the Pods are scheduled on.</p>
</td>
</tr>
<tr>
<td>
<code>tolerations</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#toleration-v1-core">
[]Kubernetes core/v1.Toleration
</a>
</em>
</td>
<td>
<p>If specified, the pod&rsquo;s tolerations.</p>
</td>
</tr>
</table>
</td>
</tr>
<tr>
<td>
<code>status</code><br/>
<em>
<a href="#monitoring.whizard.io/v1alpha1.CompactorStatus">
CompactorStatus
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="monitoring.whizard.io/v1alpha1.Gateway">Gateway
</h3>
<div>
<p>The <code>Gateway</code> custom resource definition (CRD) defines a desired <a href="https://github.com/WhizardTelemetry/whizard-docs/blob/main/Architecture/components/whizard-monitoring-gateway.md">Gateway</a> setup to run in a Kubernetes cluster. It allows to specify many options such as the number of replicas, and many more.</p>
<p>For each <code>Gateway</code> resource, the Operator deploys a <code>Deployment</code> in the same namespace.</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>apiVersion</code><br/>
string</td>
<td>
<code>
monitoring.whizard.io/v1alpha1
</code>
</td>
</tr>
<tr>
<td>
<code>kind</code><br/>
string
</td>
<td><code>Gateway</code></td>
</tr>
<tr>
<td>
<code>metadata</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#objectmeta-v1-meta">
Kubernetes meta/v1.ObjectMeta
</a>
</em>
</td>
<td>
Refer to the Kubernetes API documentation for the fields of the
<code>metadata</code> field.
</td>
</tr>
<tr>
<td>
<code>spec</code><br/>
<em>
<a href="#monitoring.whizard.io/v1alpha1.GatewaySpec">
GatewaySpec
</a>
</em>
</td>
<td>
<br/>
<br/>
<table>
<tr>
<td>
<code>webConfig</code><br/>
<em>
<a href="#monitoring.whizard.io/v1alpha1.WebConfig">
WebConfig
</a>
</em>
</td>
<td>
<p>Defines the configuration of the Gatewat web server.</p>
</td>
</tr>
<tr>
<td>
<code>debug</code><br/>
<em>
bool
</em>
</td>
<td>
<p>If debug mode is on, gateway will proxy Query UI</p>
<p>This is an <em>experimental feature</em>, it may change in any upcoming release in a breaking way.</p>
</td>
</tr>
<tr>
<td>
<code>enabledTenantsAdmission</code><br/>
<em>
bool
</em>
</td>
<td>
<p>Deny unknown tenant data remote-write and query if enabled</p>
</td>
</tr>
<tr>
<td>
<code>nodePort</code><br/>
<em>
int32
</em>
</td>
<td>
<p>NodePort is the port used to expose the gateway service.
If this is a valid node port, the gateway service type will be set to NodePort accordingly.</p>
</td>
</tr>
<tr>
<td>
<code>replicas</code><br/>
<em>
int32
</em>
</td>
<td>
<p>Number of component instances to deploy.</p>
</td>
</tr>
<tr>
<td>
<code>image</code><br/>
<em>
string
</em>
</td>
<td>
<p>Component container image URL.</p>
</td>
</tr>
<tr>
<td>
<code>imagePullPolicy</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#pullpolicy-v1-core">
Kubernetes core/v1.PullPolicy
</a>
</em>
</td>
<td>
<p>Image pull policy.</p>
</td>
</tr>
<tr>
<td>
<code>resources</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#resourcerequirements-v1-core">
Kubernetes core/v1.ResourceRequirements
</a>
</em>
</td>
<td>
<p>Resources defines the resource requirements for single Pods.</p>
</td>
</tr>
<tr>
<td>
<code>logLevel</code><br/>
<em>
string
</em>
</td>
<td>
<p>Log level for component to be configured with.</p>
</td>
</tr>
<tr>
<td>
<code>logFormat</code><br/>
<em>
string
</em>
</td>
<td>
<p>Log format for component to be configured with.</p>
</td>
</tr>
<tr>
<td>
<code>flags</code><br/>
<em>
[]string
</em>
</td>
<td>
<p>Flags allows setting additional flags for the component container.</p>
</td>
</tr>
<tr>
<td>
<code>podMetadata</code><br/>
<em>
<a href="#monitoring.whizard.io/v1alpha1.EmbeddedObjectMetadata">
EmbeddedObjectMetadata
</a>
</em>
</td>
<td>
<p>PodMetadata configures labels and annotations which are propagated to the pods.</p>
</td>
</tr>
<tr>
<td>
<code>configMaps</code><br/>
<em>
[]string
</em>
</td>
<td>
<p>ConfigMaps is a list of ConfigMaps in the same namespace as the component
object, which shall be mounted into the default Pods.
Each ConfigMap is added to the StatefulSet/Deployment definition as a volume named <code>configmap-&lt;configmap-name&gt;</code>.
The ConfigMaps are mounted into /etc/whizard/configmaps/<configmap-name> in the default container.</p>
</td>
</tr>
<tr>
<td>
<code>secrets</code><br/>
<em>
[]string
</em>
</td>
<td>
<p>Secrets is a list of Secrets in the same namespace as the component
object, which shall be mounted into the Prometheus Pods.
Each Secret is added to the StatefulSet/Deployment definition as a volume named <code>secret-&lt;secret-name&gt;</code>.
The Secrets are mounted into /etc/whizard/secrets/<secret-name> in the default container.</p>
</td>
</tr>
<tr>
<td>
<code>containers</code><br/>
<em>
<a href="https://pkg.go.dev/k8s.io/apimachinery/pkg/runtime.#RawExtension">
k8s.io/apimachinery/pkg/runtime.RawExtension
</a>
</em>
</td>
<td>
<p>Containers allows injecting additional containers or modifying operator generated containers.
Containers described here modify an operator generated
container if they share the same name and modifications are done via a
strategic merge patch.</p>
</td>
</tr>
<tr>
<td>
<code>imagePullSecrets</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#localobjectreference-v1-core">
[]Kubernetes core/v1.LocalObjectReference
</a>
</em>
</td>
<td>
<p>An optional list of references to secrets in the same namespace
to use for pulling images from registries</p>
</td>
</tr>
<tr>
<td>
<code>securityContext</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#podsecuritycontext-v1-core">
Kubernetes core/v1.PodSecurityContext
</a>
</em>
</td>
<td>
<p>SecurityContext holds pod-level security attributes and common container settings.
This defaults to the default PodSecurityContext.</p>
</td>
</tr>
<tr>
<td>
<code>affinity</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#affinity-v1-core">
Kubernetes core/v1.Affinity
</a>
</em>
</td>
<td>
<p>If specified, the pod&rsquo;s scheduling constraints.</p>
</td>
</tr>
<tr>
<td>
<code>nodeSelector</code><br/>
<em>
map[string]string
</em>
</td>
<td>
<p>Define which Nodes the Pods are scheduled on.</p>
</td>
</tr>
<tr>
<td>
<code>tolerations</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#toleration-v1-core">
[]Kubernetes core/v1.Toleration
</a>
</em>
</td>
<td>
<p>If specified, the pod&rsquo;s tolerations.</p>
</td>
</tr>
</table>
</td>
</tr>
<tr>
<td>
<code>status</code><br/>
<em>
<a href="#monitoring.whizard.io/v1alpha1.GatewayStatus">
GatewayStatus
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="monitoring.whizard.io/v1alpha1.Ingester">Ingester
</h3>
<div>
<p>The <code>Ingester</code> custom resource definition (CRD) defines a desired <a href="https://thanos.io/tip/components/receive.md/">Ingesting Receive</a> setup to run in a Kubernetes cluster. It allows to specify many options such as the number of replicas, persistent storage and many more.</p>
<p>For each <code>Ingester</code> resource, the Operator deploys a <code>StatefulSet</code> in the same namespace.</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>apiVersion</code><br/>
string</td>
<td>
<code>
monitoring.whizard.io/v1alpha1
</code>
</td>
</tr>
<tr>
<td>
<code>kind</code><br/>
string
</td>
<td><code>Ingester</code></td>
</tr>
<tr>
<td>
<code>metadata</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#objectmeta-v1-meta">
Kubernetes meta/v1.ObjectMeta
</a>
</em>
</td>
<td>
Refer to the Kubernetes API documentation for the fields of the
<code>metadata</code> field.
</td>
</tr>
<tr>
<td>
<code>spec</code><br/>
<em>
<a href="#monitoring.whizard.io/v1alpha1.IngesterSpec">
IngesterSpec
</a>
</em>
</td>
<td>
<br/>
<br/>
<table>
<tr>
<td>
<code>tenants</code><br/>
<em>
[]string
</em>
</td>
<td>
<p>The tenants whose data is being ingested by the Ingester(ingesting receiver).</p>
</td>
</tr>
<tr>
<td>
<code>otlpEnableTargetInfo</code><br/>
<em>
bool
</em>
</td>
<td>
<p>Enables target information in OTLP metrics ingested by Receive. If enabled, it converts the resource to the target info metric</p>
</td>
</tr>
<tr>
<td>
<code>otlpResourceAttributes</code><br/>
<em>
[]string
</em>
</td>
<td>
<p>Resource attributes to include in OTLP metrics ingested by Receive.</p>
</td>
</tr>
<tr>
<td>
<code>localTsdbRetention</code><br/>
<em>
string
</em>
</td>
<td>
<p>LocalTsdbRetention configs how long to retain raw samples on local storage.</p>
</td>
</tr>
<tr>
<td>
<code>dataVolume</code><br/>
<em>
<a href="#monitoring.whizard.io/v1alpha1.KubernetesVolume">
KubernetesVolume
</a>
</em>
</td>
<td>
<p>DataVolume specifies how volume shall be used</p>
</td>
</tr>
<tr>
<td>
<code>replicas</code><br/>
<em>
int32
</em>
</td>
<td>
<p>Number of component instances to deploy.</p>
</td>
</tr>
<tr>
<td>
<code>image</code><br/>
<em>
string
</em>
</td>
<td>
<p>Component container image URL.</p>
</td>
</tr>
<tr>
<td>
<code>imagePullPolicy</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#pullpolicy-v1-core">
Kubernetes core/v1.PullPolicy
</a>
</em>
</td>
<td>
<p>Image pull policy.</p>
</td>
</tr>
<tr>
<td>
<code>resources</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#resourcerequirements-v1-core">
Kubernetes core/v1.ResourceRequirements
</a>
</em>
</td>
<td>
<p>Resources defines the resource requirements for single Pods.</p>
</td>
</tr>
<tr>
<td>
<code>logLevel</code><br/>
<em>
string
</em>
</td>
<td>
<p>Log level for component to be configured with.</p>
</td>
</tr>
<tr>
<td>
<code>logFormat</code><br/>
<em>
string
</em>
</td>
<td>
<p>Log format for component to be configured with.</p>
</td>
</tr>
<tr>
<td>
<code>flags</code><br/>
<em>
[]string
</em>
</td>
<td>
<p>Flags allows setting additional flags for the component container.</p>
</td>
</tr>
<tr>
<td>
<code>podMetadata</code><br/>
<em>
<a href="#monitoring.whizard.io/v1alpha1.EmbeddedObjectMetadata">
EmbeddedObjectMetadata
</a>
</em>
</td>
<td>
<p>PodMetadata configures labels and annotations which are propagated to the pods.</p>
</td>
</tr>
<tr>
<td>
<code>configMaps</code><br/>
<em>
[]string
</em>
</td>
<td>
<p>ConfigMaps is a list of ConfigMaps in the same namespace as the component
object, which shall be mounted into the default Pods.
Each ConfigMap is added to the StatefulSet/Deployment definition as a volume named <code>configmap-&lt;configmap-name&gt;</code>.
The ConfigMaps are mounted into /etc/whizard/configmaps/<configmap-name> in the default container.</p>
</td>
</tr>
<tr>
<td>
<code>secrets</code><br/>
<em>
[]string
</em>
</td>
<td>
<p>Secrets is a list of Secrets in the same namespace as the component
object, which shall be mounted into the Prometheus Pods.
Each Secret is added to the StatefulSet/Deployment definition as a volume named <code>secret-&lt;secret-name&gt;</code>.
The Secrets are mounted into /etc/whizard/secrets/<secret-name> in the default container.</p>
</td>
</tr>
<tr>
<td>
<code>containers</code><br/>
<em>
<a href="https://pkg.go.dev/k8s.io/apimachinery/pkg/runtime.#RawExtension">
k8s.io/apimachinery/pkg/runtime.RawExtension
</a>
</em>
</td>
<td>
<p>Containers allows injecting additional containers or modifying operator generated containers.
Containers described here modify an operator generated
container if they share the same name and modifications are done via a
strategic merge patch.</p>
</td>
</tr>
<tr>
<td>
<code>imagePullSecrets</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#localobjectreference-v1-core">
[]Kubernetes core/v1.LocalObjectReference
</a>
</em>
</td>
<td>
<p>An optional list of references to secrets in the same namespace
to use for pulling images from registries</p>
</td>
</tr>
<tr>
<td>
<code>securityContext</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#podsecuritycontext-v1-core">
Kubernetes core/v1.PodSecurityContext
</a>
</em>
</td>
<td>
<p>SecurityContext holds pod-level security attributes and common container settings.
This defaults to the default PodSecurityContext.</p>
</td>
</tr>
<tr>
<td>
<code>affinity</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#affinity-v1-core">
Kubernetes core/v1.Affinity
</a>
</em>
</td>
<td>
<p>If specified, the pod&rsquo;s scheduling constraints.</p>
</td>
</tr>
<tr>
<td>
<code>nodeSelector</code><br/>
<em>
map[string]string
</em>
</td>
<td>
<p>Define which Nodes the Pods are scheduled on.</p>
</td>
</tr>
<tr>
<td>
<code>tolerations</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#toleration-v1-core">
[]Kubernetes core/v1.Toleration
</a>
</em>
</td>
<td>
<p>If specified, the pod&rsquo;s tolerations.</p>
</td>
</tr>
<tr>
<td>
<code>ingesterTsdbCleanup</code><br/>
<em>
<a href="#monitoring.whizard.io/v1alpha1.SidecarSpec">
SidecarSpec
</a>
</em>
</td>
<td>
</td>
</tr>
</table>
</td>
</tr>
<tr>
<td>
<code>status</code><br/>
<em>
<a href="#monitoring.whizard.io/v1alpha1.IngesterStatus">
IngesterStatus
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="monitoring.whizard.io/v1alpha1.Query">Query
</h3>
<div>
<p>The <code>Query</code> custom resource definition (CRD) defines a desired <a href="https://thanos.io/tip/components/query.md/">Query</a> setup to run in a Kubernetes cluster. It allows to specify many options such as the number of replicas, and many more.</p>
<p>For each <code>Query</code> resource, the Operator deploys a <code>Deployment</code> in the same namespace.</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>apiVersion</code><br/>
string</td>
<td>
<code>
monitoring.whizard.io/v1alpha1
</code>
</td>
</tr>
<tr>
<td>
<code>kind</code><br/>
string
</td>
<td><code>Query</code></td>
</tr>
<tr>
<td>
<code>metadata</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#objectmeta-v1-meta">
Kubernetes meta/v1.ObjectMeta
</a>
</em>
</td>
<td>
Refer to the Kubernetes API documentation for the fields of the
<code>metadata</code> field.
</td>
</tr>
<tr>
<td>
<code>spec</code><br/>
<em>
<a href="#monitoring.whizard.io/v1alpha1.QuerySpec">
QuerySpec
</a>
</em>
</td>
<td>
<br/>
<br/>
<table>
<tr>
<td>
<code>promqlEngine</code><br/>
<em>
string
</em>
</td>
<td>
<p>experimental PromQL engine, more info thanos.io/tip/components/query.md#promql-engine
default: prometheus</p>
</td>
</tr>
<tr>
<td>
<code>selectorLabels</code><br/>
<em>
map[string]string
</em>
</td>
<td>
<p>Selector labels that will be exposed in info endpoint.</p>
</td>
</tr>
<tr>
<td>
<code>replicaLabelNames</code><br/>
<em>
[]string
</em>
</td>
<td>
<p>Labels to treat as a replica indicator along which data is deduplicated.</p>
</td>
</tr>
<tr>
<td>
<code>webConfig</code><br/>
<em>
<a href="#monitoring.whizard.io/v1alpha1.WebConfig">
WebConfig
</a>
</em>
</td>
<td>
<p>Defines the configuration of the Thanos Query web server.</p>
</td>
</tr>
<tr>
<td>
<code>stores</code><br/>
<em>
<a href="#monitoring.whizard.io/v1alpha1.QueryStores">
[]QueryStores
</a>
</em>
</td>
<td>
<p>Additional StoreApi servers from which Query component queries from</p>
</td>
</tr>
<tr>
<td>
<code>envoy</code><br/>
<em>
<a href="#monitoring.whizard.io/v1alpha1.SidecarSpec">
SidecarSpec
</a>
</em>
</td>
<td>
<p>Envoy is used to config sidecar which proxies requests requiring auth to the secure stores</p>
</td>
</tr>
<tr>
<td>
<code>replicas</code><br/>
<em>
int32
</em>
</td>
<td>
<p>Number of component instances to deploy.</p>
</td>
</tr>
<tr>
<td>
<code>image</code><br/>
<em>
string
</em>
</td>
<td>
<p>Component container image URL.</p>
</td>
</tr>
<tr>
<td>
<code>imagePullPolicy</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#pullpolicy-v1-core">
Kubernetes core/v1.PullPolicy
</a>
</em>
</td>
<td>
<p>Image pull policy.</p>
</td>
</tr>
<tr>
<td>
<code>resources</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#resourcerequirements-v1-core">
Kubernetes core/v1.ResourceRequirements
</a>
</em>
</td>
<td>
<p>Resources defines the resource requirements for single Pods.</p>
</td>
</tr>
<tr>
<td>
<code>logLevel</code><br/>
<em>
string
</em>
</td>
<td>
<p>Log level for component to be configured with.</p>
</td>
</tr>
<tr>
<td>
<code>logFormat</code><br/>
<em>
string
</em>
</td>
<td>
<p>Log format for component to be configured with.</p>
</td>
</tr>
<tr>
<td>
<code>flags</code><br/>
<em>
[]string
</em>
</td>
<td>
<p>Flags allows setting additional flags for the component container.</p>
</td>
</tr>
<tr>
<td>
<code>podMetadata</code><br/>
<em>
<a href="#monitoring.whizard.io/v1alpha1.EmbeddedObjectMetadata">
EmbeddedObjectMetadata
</a>
</em>
</td>
<td>
<p>PodMetadata configures labels and annotations which are propagated to the pods.</p>
</td>
</tr>
<tr>
<td>
<code>configMaps</code><br/>
<em>
[]string
</em>
</td>
<td>
<p>ConfigMaps is a list of ConfigMaps in the same namespace as the component
object, which shall be mounted into the default Pods.
Each ConfigMap is added to the StatefulSet/Deployment definition as a volume named <code>configmap-&lt;configmap-name&gt;</code>.
The ConfigMaps are mounted into /etc/whizard/configmaps/<configmap-name> in the default container.</p>
</td>
</tr>
<tr>
<td>
<code>secrets</code><br/>
<em>
[]string
</em>
</td>
<td>
<p>Secrets is a list of Secrets in the same namespace as the component
object, which shall be mounted into the Prometheus Pods.
Each Secret is added to the StatefulSet/Deployment definition as a volume named <code>secret-&lt;secret-name&gt;</code>.
The Secrets are mounted into /etc/whizard/secrets/<secret-name> in the default container.</p>
</td>
</tr>
<tr>
<td>
<code>containers</code><br/>
<em>
<a href="https://pkg.go.dev/k8s.io/apimachinery/pkg/runtime.#RawExtension">
k8s.io/apimachinery/pkg/runtime.RawExtension
</a>
</em>
</td>
<td>
<p>Containers allows injecting additional containers or modifying operator generated containers.
Containers described here modify an operator generated
container if they share the same name and modifications are done via a
strategic merge patch.</p>
</td>
</tr>
<tr>
<td>
<code>imagePullSecrets</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#localobjectreference-v1-core">
[]Kubernetes core/v1.LocalObjectReference
</a>
</em>
</td>
<td>
<p>An optional list of references to secrets in the same namespace
to use for pulling images from registries</p>
</td>
</tr>
<tr>
<td>
<code>securityContext</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#podsecuritycontext-v1-core">
Kubernetes core/v1.PodSecurityContext
</a>
</em>
</td>
<td>
<p>SecurityContext holds pod-level security attributes and common container settings.
This defaults to the default PodSecurityContext.</p>
</td>
</tr>
<tr>
<td>
<code>affinity</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#affinity-v1-core">
Kubernetes core/v1.Affinity
</a>
</em>
</td>
<td>
<p>If specified, the pod&rsquo;s scheduling constraints.</p>
</td>
</tr>
<tr>
<td>
<code>nodeSelector</code><br/>
<em>
map[string]string
</em>
</td>
<td>
<p>Define which Nodes the Pods are scheduled on.</p>
</td>
</tr>
<tr>
<td>
<code>tolerations</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#toleration-v1-core">
[]Kubernetes core/v1.Toleration
</a>
</em>
</td>
<td>
<p>If specified, the pod&rsquo;s tolerations.</p>
</td>
</tr>
</table>
</td>
</tr>
<tr>
<td>
<code>status</code><br/>
<em>
<a href="#monitoring.whizard.io/v1alpha1.QueryStatus">
QueryStatus
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="monitoring.whizard.io/v1alpha1.QueryFrontend">QueryFrontend
</h3>
<div>
<p>The <code>QueryFrontend</code> custom resource definition (CRD) defines a desired <a href="https://thanos.io/tip/components/query-frontend.md/">QueryFrontend</a> setup to run in a Kubernetes cluster. It allows to specify many options such as the number of replicas, and many more.</p>
<p>For each <code>QueryFrontend</code> resource, the Operator deploys a <code>Deployment</code> in the same namespace.</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>apiVersion</code><br/>
string</td>
<td>
<code>
monitoring.whizard.io/v1alpha1
</code>
</td>
</tr>
<tr>
<td>
<code>kind</code><br/>
string
</td>
<td><code>QueryFrontend</code></td>
</tr>
<tr>
<td>
<code>metadata</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#objectmeta-v1-meta">
Kubernetes meta/v1.ObjectMeta
</a>
</em>
</td>
<td>
Refer to the Kubernetes API documentation for the fields of the
<code>metadata</code> field.
</td>
</tr>
<tr>
<td>
<code>spec</code><br/>
<em>
<a href="#monitoring.whizard.io/v1alpha1.QueryFrontendSpec">
QueryFrontendSpec
</a>
</em>
</td>
<td>
<br/>
<br/>
<table>
<tr>
<td>
<code>cacheConfig</code><br/>
<em>
<a href="#monitoring.whizard.io/v1alpha1.ResponseCacheProviderConfig">
ResponseCacheProviderConfig
</a>
</em>
</td>
<td>
<p>CacheProviderConfig specifies response cache configuration.</p>
</td>
</tr>
<tr>
<td>
<code>webConfig</code><br/>
<em>
<a href="#monitoring.whizard.io/v1alpha1.WebConfig">
WebConfig
</a>
</em>
</td>
<td>
<p>Defines the configuration of the Thanos QueryFrontend web server.</p>
</td>
</tr>
<tr>
<td>
<code>replicas</code><br/>
<em>
int32
</em>
</td>
<td>
<p>Number of component instances to deploy.</p>
</td>
</tr>
<tr>
<td>
<code>image</code><br/>
<em>
string
</em>
</td>
<td>
<p>Component container image URL.</p>
</td>
</tr>
<tr>
<td>
<code>imagePullPolicy</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#pullpolicy-v1-core">
Kubernetes core/v1.PullPolicy
</a>
</em>
</td>
<td>
<p>Image pull policy.</p>
</td>
</tr>
<tr>
<td>
<code>resources</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#resourcerequirements-v1-core">
Kubernetes core/v1.ResourceRequirements
</a>
</em>
</td>
<td>
<p>Resources defines the resource requirements for single Pods.</p>
</td>
</tr>
<tr>
<td>
<code>logLevel</code><br/>
<em>
string
</em>
</td>
<td>
<p>Log level for component to be configured with.</p>
</td>
</tr>
<tr>
<td>
<code>logFormat</code><br/>
<em>
string
</em>
</td>
<td>
<p>Log format for component to be configured with.</p>
</td>
</tr>
<tr>
<td>
<code>flags</code><br/>
<em>
[]string
</em>
</td>
<td>
<p>Flags allows setting additional flags for the component container.</p>
</td>
</tr>
<tr>
<td>
<code>podMetadata</code><br/>
<em>
<a href="#monitoring.whizard.io/v1alpha1.EmbeddedObjectMetadata">
EmbeddedObjectMetadata
</a>
</em>
</td>
<td>
<p>PodMetadata configures labels and annotations which are propagated to the pods.</p>
</td>
</tr>
<tr>
<td>
<code>configMaps</code><br/>
<em>
[]string
</em>
</td>
<td>
<p>ConfigMaps is a list of ConfigMaps in the same namespace as the component
object, which shall be mounted into the default Pods.
Each ConfigMap is added to the StatefulSet/Deployment definition as a volume named <code>configmap-&lt;configmap-name&gt;</code>.
The ConfigMaps are mounted into /etc/whizard/configmaps/<configmap-name> in the default container.</p>
</td>
</tr>
<tr>
<td>
<code>secrets</code><br/>
<em>
[]string
</em>
</td>
<td>
<p>Secrets is a list of Secrets in the same namespace as the component
object, which shall be mounted into the Prometheus Pods.
Each Secret is added to the StatefulSet/Deployment definition as a volume named <code>secret-&lt;secret-name&gt;</code>.
The Secrets are mounted into /etc/whizard/secrets/<secret-name> in the default container.</p>
</td>
</tr>
<tr>
<td>
<code>containers</code><br/>
<em>
<a href="https://pkg.go.dev/k8s.io/apimachinery/pkg/runtime.#RawExtension">
k8s.io/apimachinery/pkg/runtime.RawExtension
</a>
</em>
</td>
<td>
<p>Containers allows injecting additional containers or modifying operator generated containers.
Containers described here modify an operator generated
container if they share the same name and modifications are done via a
strategic merge patch.</p>
</td>
</tr>
<tr>
<td>
<code>imagePullSecrets</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#localobjectreference-v1-core">
[]Kubernetes core/v1.LocalObjectReference
</a>
</em>
</td>
<td>
<p>An optional list of references to secrets in the same namespace
to use for pulling images from registries</p>
</td>
</tr>
<tr>
<td>
<code>securityContext</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#podsecuritycontext-v1-core">
Kubernetes core/v1.PodSecurityContext
</a>
</em>
</td>
<td>
<p>SecurityContext holds pod-level security attributes and common container settings.
This defaults to the default PodSecurityContext.</p>
</td>
</tr>
<tr>
<td>
<code>affinity</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#affinity-v1-core">
Kubernetes core/v1.Affinity
</a>
</em>
</td>
<td>
<p>If specified, the pod&rsquo;s scheduling constraints.</p>
</td>
</tr>
<tr>
<td>
<code>nodeSelector</code><br/>
<em>
map[string]string
</em>
</td>
<td>
<p>Define which Nodes the Pods are scheduled on.</p>
</td>
</tr>
<tr>
<td>
<code>tolerations</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#toleration-v1-core">
[]Kubernetes core/v1.Toleration
</a>
</em>
</td>
<td>
<p>If specified, the pod&rsquo;s tolerations.</p>
</td>
</tr>
</table>
</td>
</tr>
<tr>
<td>
<code>status</code><br/>
<em>
<a href="#monitoring.whizard.io/v1alpha1.QueryFrontendStatus">
QueryFrontendStatus
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="monitoring.whizard.io/v1alpha1.Router">Router
</h3>
<div>
<p>The <code>Router</code> custom resource definition (CRD) defines a desired <a href="https://thanos.io/tip/components/receive.md/">Routing Receivers</a> setup to run in a Kubernetes cluster. It allows to specify many options such as the number of replicas, and many more.</p>
<p>For each <code>Router</code> resource, the Operator deploys a <code>Deployment</code> in the same namespace.</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>apiVersion</code><br/>
string</td>
<td>
<code>
monitoring.whizard.io/v1alpha1
</code>
</td>
</tr>
<tr>
<td>
<code>kind</code><br/>
string
</td>
<td><code>Router</code></td>
</tr>
<tr>
<td>
<code>metadata</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#objectmeta-v1-meta">
Kubernetes meta/v1.ObjectMeta
</a>
</em>
</td>
<td>
Refer to the Kubernetes API documentation for the fields of the
<code>metadata</code> field.
</td>
</tr>
<tr>
<td>
<code>spec</code><br/>
<em>
<a href="#monitoring.whizard.io/v1alpha1.RouterSpec">
RouterSpec
</a>
</em>
</td>
<td>
<br/>
<br/>
<table>
<tr>
<td>
<code>replicationFactor</code><br/>
<em>
uint64
</em>
</td>
<td>
<p>How many times to replicate incoming write requests</p>
</td>
</tr>
<tr>
<td>
<code>replicationProtocol</code><br/>
<em>
<a href="#monitoring.whizard.io/v1alpha1.ReplicationProtocol">
ReplicationProtocol
</a>
</em>
</td>
<td>
<p>The protocol to use for replicating remote-write requests. One of protobuf,capnproto</p>
</td>
</tr>
<tr>
<td>
<code>webConfig</code><br/>
<em>
<a href="#monitoring.whizard.io/v1alpha1.WebConfig">
WebConfig
</a>
</em>
</td>
<td>
<p>Defines the configuration of the Route(routing receiver) web server.</p>
</td>
</tr>
<tr>
<td>
<code>replicas</code><br/>
<em>
int32
</em>
</td>
<td>
<p>Number of component instances to deploy.</p>
</td>
</tr>
<tr>
<td>
<code>image</code><br/>
<em>
string
</em>
</td>
<td>
<p>Component container image URL.</p>
</td>
</tr>
<tr>
<td>
<code>imagePullPolicy</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#pullpolicy-v1-core">
Kubernetes core/v1.PullPolicy
</a>
</em>
</td>
<td>
<p>Image pull policy.</p>
</td>
</tr>
<tr>
<td>
<code>resources</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#resourcerequirements-v1-core">
Kubernetes core/v1.ResourceRequirements
</a>
</em>
</td>
<td>
<p>Resources defines the resource requirements for single Pods.</p>
</td>
</tr>
<tr>
<td>
<code>logLevel</code><br/>
<em>
string
</em>
</td>
<td>
<p>Log level for component to be configured with.</p>
</td>
</tr>
<tr>
<td>
<code>logFormat</code><br/>
<em>
string
</em>
</td>
<td>
<p>Log format for component to be configured with.</p>
</td>
</tr>
<tr>
<td>
<code>flags</code><br/>
<em>
[]string
</em>
</td>
<td>
<p>Flags allows setting additional flags for the component container.</p>
</td>
</tr>
<tr>
<td>
<code>podMetadata</code><br/>
<em>
<a href="#monitoring.whizard.io/v1alpha1.EmbeddedObjectMetadata">
EmbeddedObjectMetadata
</a>
</em>
</td>
<td>
<p>PodMetadata configures labels and annotations which are propagated to the pods.</p>
</td>
</tr>
<tr>
<td>
<code>configMaps</code><br/>
<em>
[]string
</em>
</td>
<td>
<p>ConfigMaps is a list of ConfigMaps in the same namespace as the component
object, which shall be mounted into the default Pods.
Each ConfigMap is added to the StatefulSet/Deployment definition as a volume named <code>configmap-&lt;configmap-name&gt;</code>.
The ConfigMaps are mounted into /etc/whizard/configmaps/<configmap-name> in the default container.</p>
</td>
</tr>
<tr>
<td>
<code>secrets</code><br/>
<em>
[]string
</em>
</td>
<td>
<p>Secrets is a list of Secrets in the same namespace as the component
object, which shall be mounted into the Prometheus Pods.
Each Secret is added to the StatefulSet/Deployment definition as a volume named <code>secret-&lt;secret-name&gt;</code>.
The Secrets are mounted into /etc/whizard/secrets/<secret-name> in the default container.</p>
</td>
</tr>
<tr>
<td>
<code>containers</code><br/>
<em>
<a href="https://pkg.go.dev/k8s.io/apimachinery/pkg/runtime.#RawExtension">
k8s.io/apimachinery/pkg/runtime.RawExtension
</a>
</em>
</td>
<td>
<p>Containers allows injecting additional containers or modifying operator generated containers.
Containers described here modify an operator generated
container if they share the same name and modifications are done via a
strategic merge patch.</p>
</td>
</tr>
<tr>
<td>
<code>imagePullSecrets</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#localobjectreference-v1-core">
[]Kubernetes core/v1.LocalObjectReference
</a>
</em>
</td>
<td>
<p>An optional list of references to secrets in the same namespace
to use for pulling images from registries</p>
</td>
</tr>
<tr>
<td>
<code>securityContext</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#podsecuritycontext-v1-core">
Kubernetes core/v1.PodSecurityContext
</a>
</em>
</td>
<td>
<p>SecurityContext holds pod-level security attributes and common container settings.
This defaults to the default PodSecurityContext.</p>
</td>
</tr>
<tr>
<td>
<code>affinity</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#affinity-v1-core">
Kubernetes core/v1.Affinity
</a>
</em>
</td>
<td>
<p>If specified, the pod&rsquo;s scheduling constraints.</p>
</td>
</tr>
<tr>
<td>
<code>nodeSelector</code><br/>
<em>
map[string]string
</em>
</td>
<td>
<p>Define which Nodes the Pods are scheduled on.</p>
</td>
</tr>
<tr>
<td>
<code>tolerations</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#toleration-v1-core">
[]Kubernetes core/v1.Toleration
</a>
</em>
</td>
<td>
<p>If specified, the pod&rsquo;s tolerations.</p>
</td>
</tr>
</table>
</td>
</tr>
<tr>
<td>
<code>status</code><br/>
<em>
<a href="#monitoring.whizard.io/v1alpha1.RouterStatus">
RouterStatus
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="monitoring.whizard.io/v1alpha1.Ruler">Ruler
</h3>
<div>
<p>Ruler is the Schema for the rulers API</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>apiVersion</code><br/>
string</td>
<td>
<code>
monitoring.whizard.io/v1alpha1
</code>
</td>
</tr>
<tr>
<td>
<code>kind</code><br/>
string
</td>
<td><code>Ruler</code></td>
</tr>
<tr>
<td>
<code>metadata</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#objectmeta-v1-meta">
Kubernetes meta/v1.ObjectMeta
</a>
</em>
</td>
<td>
Refer to the Kubernetes API documentation for the fields of the
<code>metadata</code> field.
</td>
</tr>
<tr>
<td>
<code>spec</code><br/>
<em>
<a href="#monitoring.whizard.io/v1alpha1.RulerSpec">
RulerSpec
</a>
</em>
</td>
<td>
<br/>
<br/>
<table>
<tr>
<td>
<code>ruleSelectors</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#*k8s.io/apimachinery/pkg/apis/meta/v1.labelselector--">
[]*k8s.io/apimachinery/pkg/apis/meta/v1.LabelSelector
</a>
</em>
</td>
<td>
<p>Label selectors to select which PrometheusRules to mount for alerting and recording.
The result of multiple selectors are ORed.</p>
</td>
</tr>
<tr>
<td>
<code>ruleNamespaceSelector</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#labelselector-v1-meta">
Kubernetes meta/v1.LabelSelector
</a>
</em>
</td>
<td>
<p>Namespaces to be selected for PrometheusRules discovery. If unspecified, only
the same namespace as the Ruler object is in is used.</p>
</td>
</tr>
<tr>
<td>
<code>shards</code><br/>
<em>
int32
</em>
</td>
<td>
<p>Number of shards to take the hash of fully qualified name of the rule group in order to split rules.
Each shard of rules will be bound to one separate statefulset.
Default: 1</p>
</td>
</tr>
<tr>
<td>
<code>tenant</code><br/>
<em>
string
</em>
</td>
<td>
<p>Tenant if not empty indicates which tenant&rsquo;s data is evaluated for the selected rules;
otherwise, it is for all tenants.</p>
</td>
</tr>
<tr>
<td>
<code>queryConfig</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#secretkeyselector-v1-core">
Kubernetes core/v1.SecretKeySelector
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>remoteWriteConfig</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#secretkeyselector-v1-core">
Kubernetes core/v1.SecretKeySelector
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>labels</code><br/>
<em>
map[string]string
</em>
</td>
<td>
<p>Labels configure the external label pairs to Ruler. A default replica label
<code>ruler_replica</code> will be always added  as a label with the value of the pod&rsquo;s name and it will be dropped in the alerts.</p>
</td>
</tr>
<tr>
<td>
<code>alertDropLabels</code><br/>
<em>
[]string
</em>
</td>
<td>
<p>AlertDropLabels configure the label names which should be dropped in Ruler alerts.
The replica label <code>ruler_replica</code> will always be dropped in alerts.</p>
</td>
</tr>
<tr>
<td>
<code>alertmanagersUrl</code><br/>
<em>
[]string
</em>
</td>
<td>
<p>Define URLs to send alerts to Alertmanager.
Note: this field will be ignored if AlertmanagersConfig is specified.
Maps to the <code>alertmanagers.url</code> arg.</p>
</td>
</tr>
<tr>
<td>
<code>alertmanagersConfig</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#secretkeyselector-v1-core">
Kubernetes core/v1.SecretKeySelector
</a>
</em>
</td>
<td>
<p>Define configuration for connecting to alertmanager. Maps to the <code>alertmanagers.config</code> arg.</p>
</td>
</tr>
<tr>
<td>
<code>evaluationInterval</code><br/>
<em>
<a href="#monitoring.whizard.io/v1alpha1.Duration">
Duration
</a>
</em>
</td>
<td>
<p>Interval between consecutive evaluations.</p>
<p>Default: &ldquo;1m&rdquo;</p>
</td>
</tr>
<tr>
<td>
<code>dataVolume</code><br/>
<em>
<a href="#monitoring.whizard.io/v1alpha1.KubernetesVolume">
KubernetesVolume
</a>
</em>
</td>
<td>
<p>DataVolume specifies how volume shall be used</p>
</td>
</tr>
<tr>
<td>
<code>rulerQueryProxy</code><br/>
<em>
<a href="#monitoring.whizard.io/v1alpha1.SidecarSpec">
SidecarSpec
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>rulerWriteProxy</code><br/>
<em>
<a href="#monitoring.whizard.io/v1alpha1.SidecarSpec">
SidecarSpec
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>prometheusConfigReloader</code><br/>
<em>
<a href="#monitoring.whizard.io/v1alpha1.SidecarSpec">
SidecarSpec
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>replicas</code><br/>
<em>
int32
</em>
</td>
<td>
<p>Number of component instances to deploy.</p>
</td>
</tr>
<tr>
<td>
<code>image</code><br/>
<em>
string
</em>
</td>
<td>
<p>Component container image URL.</p>
</td>
</tr>
<tr>
<td>
<code>imagePullPolicy</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#pullpolicy-v1-core">
Kubernetes core/v1.PullPolicy
</a>
</em>
</td>
<td>
<p>Image pull policy.</p>
</td>
</tr>
<tr>
<td>
<code>resources</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#resourcerequirements-v1-core">
Kubernetes core/v1.ResourceRequirements
</a>
</em>
</td>
<td>
<p>Resources defines the resource requirements for single Pods.</p>
</td>
</tr>
<tr>
<td>
<code>logLevel</code><br/>
<em>
string
</em>
</td>
<td>
<p>Log level for component to be configured with.</p>
</td>
</tr>
<tr>
<td>
<code>logFormat</code><br/>
<em>
string
</em>
</td>
<td>
<p>Log format for component to be configured with.</p>
</td>
</tr>
<tr>
<td>
<code>flags</code><br/>
<em>
[]string
</em>
</td>
<td>
<p>Flags allows setting additional flags for the component container.</p>
</td>
</tr>
<tr>
<td>
<code>podMetadata</code><br/>
<em>
<a href="#monitoring.whizard.io/v1alpha1.EmbeddedObjectMetadata">
EmbeddedObjectMetadata
</a>
</em>
</td>
<td>
<p>PodMetadata configures labels and annotations which are propagated to the pods.</p>
</td>
</tr>
<tr>
<td>
<code>configMaps</code><br/>
<em>
[]string
</em>
</td>
<td>
<p>ConfigMaps is a list of ConfigMaps in the same namespace as the component
object, which shall be mounted into the default Pods.
Each ConfigMap is added to the StatefulSet/Deployment definition as a volume named <code>configmap-&lt;configmap-name&gt;</code>.
The ConfigMaps are mounted into /etc/whizard/configmaps/<configmap-name> in the default container.</p>
</td>
</tr>
<tr>
<td>
<code>secrets</code><br/>
<em>
[]string
</em>
</td>
<td>
<p>Secrets is a list of Secrets in the same namespace as the component
object, which shall be mounted into the Prometheus Pods.
Each Secret is added to the StatefulSet/Deployment definition as a volume named <code>secret-&lt;secret-name&gt;</code>.
The Secrets are mounted into /etc/whizard/secrets/<secret-name> in the default container.</p>
</td>
</tr>
<tr>
<td>
<code>containers</code><br/>
<em>
<a href="https://pkg.go.dev/k8s.io/apimachinery/pkg/runtime.#RawExtension">
k8s.io/apimachinery/pkg/runtime.RawExtension
</a>
</em>
</td>
<td>
<p>Containers allows injecting additional containers or modifying operator generated containers.
Containers described here modify an operator generated
container if they share the same name and modifications are done via a
strategic merge patch.</p>
</td>
</tr>
<tr>
<td>
<code>imagePullSecrets</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#localobjectreference-v1-core">
[]Kubernetes core/v1.LocalObjectReference
</a>
</em>
</td>
<td>
<p>An optional list of references to secrets in the same namespace
to use for pulling images from registries</p>
</td>
</tr>
<tr>
<td>
<code>securityContext</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#podsecuritycontext-v1-core">
Kubernetes core/v1.PodSecurityContext
</a>
</em>
</td>
<td>
<p>SecurityContext holds pod-level security attributes and common container settings.
This defaults to the default PodSecurityContext.</p>
</td>
</tr>
<tr>
<td>
<code>affinity</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#affinity-v1-core">
Kubernetes core/v1.Affinity
</a>
</em>
</td>
<td>
<p>If specified, the pod&rsquo;s scheduling constraints.</p>
</td>
</tr>
<tr>
<td>
<code>nodeSelector</code><br/>
<em>
map[string]string
</em>
</td>
<td>
<p>Define which Nodes the Pods are scheduled on.</p>
</td>
</tr>
<tr>
<td>
<code>tolerations</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#toleration-v1-core">
[]Kubernetes core/v1.Toleration
</a>
</em>
</td>
<td>
<p>If specified, the pod&rsquo;s tolerations.</p>
</td>
</tr>
</table>
</td>
</tr>
<tr>
<td>
<code>status</code><br/>
<em>
<a href="#monitoring.whizard.io/v1alpha1.RulerStatus">
RulerStatus
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="monitoring.whizard.io/v1alpha1.Service">Service
</h3>
<div>
<p>The <code>Service</code> custom resource definition (CRD) defines the Whizard service configuration.
The `ServiceSpec“ has component configuration templates. Some components scale based on the number of tenants and load service configurations</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>apiVersion</code><br/>
string</td>
<td>
<code>
monitoring.whizard.io/v1alpha1
</code>
</td>
</tr>
<tr>
<td>
<code>kind</code><br/>
string
</td>
<td><code>Service</code></td>
</tr>
<tr>
<td>
<code>metadata</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#objectmeta-v1-meta">
Kubernetes meta/v1.ObjectMeta
</a>
</em>
</td>
<td>
Refer to the Kubernetes API documentation for the fields of the
<code>metadata</code> field.
</td>
</tr>
<tr>
<td>
<code>spec</code><br/>
<em>
<a href="#monitoring.whizard.io/v1alpha1.ServiceSpec">
ServiceSpec
</a>
</em>
</td>
<td>
<br/>
<br/>
<table>
<tr>
<td>
<code>tenantHeader</code><br/>
<em>
string
</em>
</td>
<td>
<p>HTTP header to determine tenant for remote write requests.</p>
</td>
</tr>
<tr>
<td>
<code>defaultTenantId</code><br/>
<em>
string
</em>
</td>
<td>
<p>Default tenant ID to use when none is provided via a header.</p>
</td>
</tr>
<tr>
<td>
<code>tenantLabelName</code><br/>
<em>
string
</em>
</td>
<td>
<p>Label name through which the tenant will be announced.</p>
</td>
</tr>
<tr>
<td>
<code>storage</code><br/>
<em>
<a href="#monitoring.whizard.io/v1alpha1.ObjectReference">
ObjectReference
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>remoteWrites</code><br/>
<em>
<a href="#monitoring.whizard.io/v1alpha1.RemoteWriteSpec">
[]RemoteWriteSpec
</a>
</em>
</td>
<td>
<p>RemoteWrites is the list of remote write configurations.
If it is configured, its targets will receive write requests from the Gateway and the Ruler.</p>
</td>
</tr>
<tr>
<td>
<code>remoteQuery</code><br/>
<em>
<a href="#monitoring.whizard.io/v1alpha1.RemoteQuerySpec">
RemoteQuerySpec
</a>
</em>
</td>
<td>
<p>RemoteQuery is the remote query configuration and the remote target should have prometheus-compatible Query APIs.
If not configured, the Gateway will proxy all read requests through the QueryFrontend to the Query,
If configured, the Gateway will proxy metrics read requests through the QueryFrontend to the remote target,
but proxy rules read requests directly to the Query.</p>
</td>
</tr>
<tr>
<td>
<code>gatewayTemplateSpec</code><br/>
<em>
<a href="#monitoring.whizard.io/v1alpha1.GatewaySpec">
GatewaySpec
</a>
</em>
</td>
<td>
<p>GatewayTemplateSpec defines the Gateway configuration template.</p>
</td>
</tr>
<tr>
<td>
<code>queryFrontendTemplateSpec</code><br/>
<em>
<a href="#monitoring.whizard.io/v1alpha1.QueryFrontendSpec">
QueryFrontendSpec
</a>
</em>
</td>
<td>
<p>QueryFrontendTemplateSpec defines the QueryFrontend configuration template.</p>
</td>
</tr>
<tr>
<td>
<code>queryTemplateSpec</code><br/>
<em>
<a href="#monitoring.whizard.io/v1alpha1.QuerySpec">
QuerySpec
</a>
</em>
</td>
<td>
<p>QueryTemplateSpec defines the Query configuration template.</p>
</td>
</tr>
<tr>
<td>
<code>rulerTemplateSpec</code><br/>
<em>
<a href="#monitoring.whizard.io/v1alpha1.RulerTemplateSpec">
RulerTemplateSpec
</a>
</em>
</td>
<td>
<p>RulerTemplateSpec defines the Ruler configuration template.</p>
</td>
</tr>
<tr>
<td>
<code>routerTemplateSpec</code><br/>
<em>
<a href="#monitoring.whizard.io/v1alpha1.RouterSpec">
RouterSpec
</a>
</em>
</td>
<td>
<p>RouterTemplateSpec defines the Router configuration template.</p>
</td>
</tr>
<tr>
<td>
<code>ingesterTemplateSpec</code><br/>
<em>
<a href="#monitoring.whizard.io/v1alpha1.IngesterTemplateSpec">
IngesterTemplateSpec
</a>
</em>
</td>
<td>
<p>IngesterTemplateSpec defines the Ingester configuration template.</p>
</td>
</tr>
<tr>
<td>
<code>storeTemplateSpec</code><br/>
<em>
<a href="#monitoring.whizard.io/v1alpha1.StoreSpec">
StoreSpec
</a>
</em>
</td>
<td>
<p>StoreTemplateSpec defines the Store configuration template.</p>
</td>
</tr>
<tr>
<td>
<code>compactorTemplateSpec</code><br/>
<em>
<a href="#monitoring.whizard.io/v1alpha1.CompactorTemplateSpec">
CompactorTemplateSpec
</a>
</em>
</td>
<td>
<p>CompactorTemplateSpec defines the Compactor configuration template.</p>
</td>
</tr>
</table>
</td>
</tr>
<tr>
<td>
<code>status</code><br/>
<em>
<a href="#monitoring.whizard.io/v1alpha1.ServiceStatus">
ServiceStatus
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="monitoring.whizard.io/v1alpha1.Storage">Storage
</h3>
<div>
<p>The <code>Storage</code> custom resource definition (CRD) defines how to configure access to object storage.
More info <a href="https://thanos.io/tip/thanos/storage.md/">https://thanos.io/tip/thanos/storage.md/</a>
Current object storage client implementations: S3, other in progress.</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>apiVersion</code><br/>
string</td>
<td>
<code>
monitoring.whizard.io/v1alpha1
</code>
</td>
</tr>
<tr>
<td>
<code>kind</code><br/>
string
</td>
<td><code>Storage</code></td>
</tr>
<tr>
<td>
<code>metadata</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#objectmeta-v1-meta">
Kubernetes meta/v1.ObjectMeta
</a>
</em>
</td>
<td>
Refer to the Kubernetes API documentation for the fields of the
<code>metadata</code> field.
</td>
</tr>
<tr>
<td>
<code>spec</code><br/>
<em>
<a href="#monitoring.whizard.io/v1alpha1.StorageSpec">
StorageSpec
</a>
</em>
</td>
<td>
<br/>
<br/>
<table>
<tr>
<td>
<code>blockManager</code><br/>
<em>
<a href="#monitoring.whizard.io/v1alpha1.BlockManager">
BlockManager
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>S3</code><br/>
<em>
<a href="#monitoring.whizard.io/v1alpha1.S3">
S3
</a>
</em>
</td>
<td>
</td>
</tr>
</table>
</td>
</tr>
<tr>
<td>
<code>status</code><br/>
<em>
<a href="#monitoring.whizard.io/v1alpha1.StorageStatus">
StorageStatus
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="monitoring.whizard.io/v1alpha1.Store">Store
</h3>
<div>
<p>The <code>Store</code> custom resource definition (CRD) defines a desired <a href="https://thanos.io/tip/components/store.md/">Compactor</a> setup to run in a Kubernetes cluster. It allows to specify many options such as the number of replicas, persistent storage and many more.</p>
<p>For each <code>Store</code> resource, the Operator deploys a <code>StatefulSet</code> in the same namespace.</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>apiVersion</code><br/>
string</td>
<td>
<code>
monitoring.whizard.io/v1alpha1
</code>
</td>
</tr>
<tr>
<td>
<code>kind</code><br/>
string
</td>
<td><code>Store</code></td>
</tr>
<tr>
<td>
<code>metadata</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#objectmeta-v1-meta">
Kubernetes meta/v1.ObjectMeta
</a>
</em>
</td>
<td>
Refer to the Kubernetes API documentation for the fields of the
<code>metadata</code> field.
</td>
</tr>
<tr>
<td>
<code>spec</code><br/>
<em>
<a href="#monitoring.whizard.io/v1alpha1.StoreSpec">
StoreSpec
</a>
</em>
</td>
<td>
<br/>
<br/>
<table>
<tr>
<td>
<code>minTime</code><br/>
<em>
string
</em>
</td>
<td>
<p>MinTime specifies start of time range limit to serve</p>
</td>
</tr>
<tr>
<td>
<code>maxTime</code><br/>
<em>
string
</em>
</td>
<td>
<p>MaxTime specifies end of time range limit to serve</p>
</td>
</tr>
<tr>
<td>
<code>timeRanges</code><br/>
<em>
<a href="#monitoring.whizard.io/v1alpha1.TimeRange">
[]TimeRange
</a>
</em>
</td>
<td>
<p>TimeRanges is a list of TimeRange to partition Store.
If specified, the MinTime and MaxTime will be ignored.</p>
</td>
</tr>
<tr>
<td>
<code>indexCacheConfig</code><br/>
<em>
<a href="#monitoring.whizard.io/v1alpha1.IndexCacheConfig">
IndexCacheConfig
</a>
</em>
</td>
<td>
<p>IndexCacheConfig contains index cache configuration.</p>
</td>
</tr>
<tr>
<td>
<code>dataVolume</code><br/>
<em>
<a href="#monitoring.whizard.io/v1alpha1.KubernetesVolume">
KubernetesVolume
</a>
</em>
</td>
<td>
<p>DataVolume specifies how volume shall be used</p>
</td>
</tr>
<tr>
<td>
<code>replicas</code><br/>
<em>
int32
</em>
</td>
<td>
<p>Number of component instances to deploy.</p>
</td>
</tr>
<tr>
<td>
<code>image</code><br/>
<em>
string
</em>
</td>
<td>
<p>Component container image URL.</p>
</td>
</tr>
<tr>
<td>
<code>imagePullPolicy</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#pullpolicy-v1-core">
Kubernetes core/v1.PullPolicy
</a>
</em>
</td>
<td>
<p>Image pull policy.</p>
</td>
</tr>
<tr>
<td>
<code>resources</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#resourcerequirements-v1-core">
Kubernetes core/v1.ResourceRequirements
</a>
</em>
</td>
<td>
<p>Resources defines the resource requirements for single Pods.</p>
</td>
</tr>
<tr>
<td>
<code>logLevel</code><br/>
<em>
string
</em>
</td>
<td>
<p>Log level for component to be configured with.</p>
</td>
</tr>
<tr>
<td>
<code>logFormat</code><br/>
<em>
string
</em>
</td>
<td>
<p>Log format for component to be configured with.</p>
</td>
</tr>
<tr>
<td>
<code>flags</code><br/>
<em>
[]string
</em>
</td>
<td>
<p>Flags allows setting additional flags for the component container.</p>
</td>
</tr>
<tr>
<td>
<code>podMetadata</code><br/>
<em>
<a href="#monitoring.whizard.io/v1alpha1.EmbeddedObjectMetadata">
EmbeddedObjectMetadata
</a>
</em>
</td>
<td>
<p>PodMetadata configures labels and annotations which are propagated to the pods.</p>
</td>
</tr>
<tr>
<td>
<code>configMaps</code><br/>
<em>
[]string
</em>
</td>
<td>
<p>ConfigMaps is a list of ConfigMaps in the same namespace as the component
object, which shall be mounted into the default Pods.
Each ConfigMap is added to the StatefulSet/Deployment definition as a volume named <code>configmap-&lt;configmap-name&gt;</code>.
The ConfigMaps are mounted into /etc/whizard/configmaps/<configmap-name> in the default container.</p>
</td>
</tr>
<tr>
<td>
<code>secrets</code><br/>
<em>
[]string
</em>
</td>
<td>
<p>Secrets is a list of Secrets in the same namespace as the component
object, which shall be mounted into the Prometheus Pods.
Each Secret is added to the StatefulSet/Deployment definition as a volume named <code>secret-&lt;secret-name&gt;</code>.
The Secrets are mounted into /etc/whizard/secrets/<secret-name> in the default container.</p>
</td>
</tr>
<tr>
<td>
<code>containers</code><br/>
<em>
<a href="https://pkg.go.dev/k8s.io/apimachinery/pkg/runtime.#RawExtension">
k8s.io/apimachinery/pkg/runtime.RawExtension
</a>
</em>
</td>
<td>
<p>Containers allows injecting additional containers or modifying operator generated containers.
Containers described here modify an operator generated
container if they share the same name and modifications are done via a
strategic merge patch.</p>
</td>
</tr>
<tr>
<td>
<code>imagePullSecrets</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#localobjectreference-v1-core">
[]Kubernetes core/v1.LocalObjectReference
</a>
</em>
</td>
<td>
<p>An optional list of references to secrets in the same namespace
to use for pulling images from registries</p>
</td>
</tr>
<tr>
<td>
<code>securityContext</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#podsecuritycontext-v1-core">
Kubernetes core/v1.PodSecurityContext
</a>
</em>
</td>
<td>
<p>SecurityContext holds pod-level security attributes and common container settings.
This defaults to the default PodSecurityContext.</p>
</td>
</tr>
<tr>
<td>
<code>affinity</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#affinity-v1-core">
Kubernetes core/v1.Affinity
</a>
</em>
</td>
<td>
<p>If specified, the pod&rsquo;s scheduling constraints.</p>
</td>
</tr>
<tr>
<td>
<code>nodeSelector</code><br/>
<em>
map[string]string
</em>
</td>
<td>
<p>Define which Nodes the Pods are scheduled on.</p>
</td>
</tr>
<tr>
<td>
<code>tolerations</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#toleration-v1-core">
[]Kubernetes core/v1.Toleration
</a>
</em>
</td>
<td>
<p>If specified, the pod&rsquo;s tolerations.</p>
</td>
</tr>
</table>
</td>
</tr>
<tr>
<td>
<code>status</code><br/>
<em>
<a href="#monitoring.whizard.io/v1alpha1.StoreStatus">
StoreStatus
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="monitoring.whizard.io/v1alpha1.Tenant">Tenant
</h3>
<div>
<p>The <code>Tenant</code> custom resource definition (CRD) defines the tenant configuration for multi-tenant data separation in Whizard.
In Whizard, a tenant can represent various types of data sources, such as:</p>
<ul>
<li>Monitoring data from a specific Kubernetes cluster</li>
<li>Monitoring data from a physical machine in a specific region</li>
<li>Monitoring data from a specific type of application</li>
</ul>
<p>When data is ingested, it will be tagged with the tenant label to ensure proper separation.</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>apiVersion</code><br/>
string</td>
<td>
<code>
monitoring.whizard.io/v1alpha1
</code>
</td>
</tr>
<tr>
<td>
<code>kind</code><br/>
string
</td>
<td><code>Tenant</code></td>
</tr>
<tr>
<td>
<code>metadata</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#objectmeta-v1-meta">
Kubernetes meta/v1.ObjectMeta
</a>
</em>
</td>
<td>
Refer to the Kubernetes API documentation for the fields of the
<code>metadata</code> field.
</td>
</tr>
<tr>
<td>
<code>spec</code><br/>
<em>
<a href="#monitoring.whizard.io/v1alpha1.TenantSpec">
TenantSpec
</a>
</em>
</td>
<td>
<br/>
<br/>
<table>
<tr>
<td>
<code>tenant</code><br/>
<em>
string
</em>
</td>
<td>
</td>
</tr>
</table>
</td>
</tr>
<tr>
<td>
<code>status</code><br/>
<em>
<a href="#monitoring.whizard.io/v1alpha1.TenantStatus">
TenantStatus
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="monitoring.whizard.io/v1alpha1.BasicAuth">BasicAuth
</h3>
<p>
(<em>Appears on:</em><a href="#monitoring.whizard.io/v1alpha1.HTTPClientConfig">HTTPClientConfig</a>, <a href="#monitoring.whizard.io/v1alpha1.WebConfig">WebConfig</a>)
</p>
<div>
<p>BasicAuth allow an endpoint to authenticate over basic authentication</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>username</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#secretkeyselector-v1-core">
Kubernetes core/v1.SecretKeySelector
</a>
</em>
</td>
<td>
<p>The secret in the service monitor namespace that contains the username
for authentication.</p>
</td>
</tr>
<tr>
<td>
<code>password</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#secretkeyselector-v1-core">
Kubernetes core/v1.SecretKeySelector
</a>
</em>
</td>
<td>
<p>The secret in the service monitor namespace that contains the password
for authentication.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="monitoring.whizard.io/v1alpha1.BlockGC">BlockGC
</h3>
<p>
(<em>Appears on:</em><a href="#monitoring.whizard.io/v1alpha1.BlockManager">BlockManager</a>)
</p>
<div>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>enable</code><br/>
<em>
bool
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>resources</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#resourcerequirements-v1-core">
Kubernetes core/v1.ResourceRequirements
</a>
</em>
</td>
<td>
<p>Define resources requests and limits for main container.</p>
</td>
</tr>
<tr>
<td>
<code>image</code><br/>
<em>
string
</em>
</td>
<td>
<p>Image is the component image with tag/version.</p>
</td>
</tr>
<tr>
<td>
<code>imagePullPolicy</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#pullpolicy-v1-core">
Kubernetes core/v1.PullPolicy
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Image pull policy.
One of Always, Never, IfNotPresent.
Defaults to Always if :latest tag is specified, or IfNotPresent otherwise.
Cannot be updated.</p>
</td>
</tr>
<tr>
<td>
<code>gcInterval</code><br/>
<em>
<a href="https://pkg.go.dev/k8s.io/apimachinery/pkg/apis/meta/v1#Duration">
Kubernetes meta/v1.Duration
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>cleanupTimeout</code><br/>
<em>
<a href="https://pkg.go.dev/k8s.io/apimachinery/pkg/apis/meta/v1#Duration">
Kubernetes meta/v1.Duration
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>defaultTenantId</code><br/>
<em>
string
</em>
</td>
<td>
<p>Default tenant ID to use when none is provided via a header.</p>
</td>
</tr>
<tr>
<td>
<code>tenantLabelName</code><br/>
<em>
string
</em>
</td>
<td>
<p>Label name through which the tenant will be announced.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="monitoring.whizard.io/v1alpha1.BlockManager">BlockManager
</h3>
<p>
(<em>Appears on:</em><a href="#monitoring.whizard.io/v1alpha1.StorageSpec">StorageSpec</a>)
</p>
<div>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>enable</code><br/>
<em>
bool
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>replicas</code><br/>
<em>
int32
</em>
</td>
<td>
<p>Number of component instances to deploy.</p>
</td>
</tr>
<tr>
<td>
<code>image</code><br/>
<em>
string
</em>
</td>
<td>
<p>Component container image URL.</p>
</td>
</tr>
<tr>
<td>
<code>imagePullPolicy</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#pullpolicy-v1-core">
Kubernetes core/v1.PullPolicy
</a>
</em>
</td>
<td>
<p>Image pull policy.</p>
</td>
</tr>
<tr>
<td>
<code>resources</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#resourcerequirements-v1-core">
Kubernetes core/v1.ResourceRequirements
</a>
</em>
</td>
<td>
<p>Resources defines the resource requirements for single Pods.</p>
</td>
</tr>
<tr>
<td>
<code>logLevel</code><br/>
<em>
string
</em>
</td>
<td>
<p>Log level for component to be configured with.</p>
</td>
</tr>
<tr>
<td>
<code>logFormat</code><br/>
<em>
string
</em>
</td>
<td>
<p>Log format for component to be configured with.</p>
</td>
</tr>
<tr>
<td>
<code>flags</code><br/>
<em>
[]string
</em>
</td>
<td>
<p>Flags allows setting additional flags for the component container.</p>
</td>
</tr>
<tr>
<td>
<code>podMetadata</code><br/>
<em>
<a href="#monitoring.whizard.io/v1alpha1.EmbeddedObjectMetadata">
EmbeddedObjectMetadata
</a>
</em>
</td>
<td>
<p>PodMetadata configures labels and annotations which are propagated to the pods.</p>
</td>
</tr>
<tr>
<td>
<code>configMaps</code><br/>
<em>
[]string
</em>
</td>
<td>
<p>ConfigMaps is a list of ConfigMaps in the same namespace as the component
object, which shall be mounted into the default Pods.
Each ConfigMap is added to the StatefulSet/Deployment definition as a volume named <code>configmap-&lt;configmap-name&gt;</code>.
The ConfigMaps are mounted into /etc/whizard/configmaps/<configmap-name> in the default container.</p>
</td>
</tr>
<tr>
<td>
<code>secrets</code><br/>
<em>
[]string
</em>
</td>
<td>
<p>Secrets is a list of Secrets in the same namespace as the component
object, which shall be mounted into the Prometheus Pods.
Each Secret is added to the StatefulSet/Deployment definition as a volume named <code>secret-&lt;secret-name&gt;</code>.
The Secrets are mounted into /etc/whizard/secrets/<secret-name> in the default container.</p>
</td>
</tr>
<tr>
<td>
<code>containers</code><br/>
<em>
<a href="https://pkg.go.dev/k8s.io/apimachinery/pkg/runtime.#RawExtension">
k8s.io/apimachinery/pkg/runtime.RawExtension
</a>
</em>
</td>
<td>
<p>Containers allows injecting additional containers or modifying operator generated containers.
Containers described here modify an operator generated
container if they share the same name and modifications are done via a
strategic merge patch.</p>
</td>
</tr>
<tr>
<td>
<code>imagePullSecrets</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#localobjectreference-v1-core">
[]Kubernetes core/v1.LocalObjectReference
</a>
</em>
</td>
<td>
<p>An optional list of references to secrets in the same namespace
to use for pulling images from registries</p>
</td>
</tr>
<tr>
<td>
<code>securityContext</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#podsecuritycontext-v1-core">
Kubernetes core/v1.PodSecurityContext
</a>
</em>
</td>
<td>
<p>SecurityContext holds pod-level security attributes and common container settings.
This defaults to the default PodSecurityContext.</p>
</td>
</tr>
<tr>
<td>
<code>affinity</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#affinity-v1-core">
Kubernetes core/v1.Affinity
</a>
</em>
</td>
<td>
<p>If specified, the pod&rsquo;s scheduling constraints.</p>
</td>
</tr>
<tr>
<td>
<code>nodeSelector</code><br/>
<em>
map[string]string
</em>
</td>
<td>
<p>Define which Nodes the Pods are scheduled on.</p>
</td>
</tr>
<tr>
<td>
<code>tolerations</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#toleration-v1-core">
[]Kubernetes core/v1.Toleration
</a>
</em>
</td>
<td>
<p>If specified, the pod&rsquo;s tolerations.</p>
</td>
</tr>
<tr>
<td>
<code>serviceAccountName</code><br/>
<em>
string
</em>
</td>
<td>
<p>ServiceAccountName is the name of the ServiceAccount to use to run bucket Pods.</p>
</td>
</tr>
<tr>
<td>
<code>nodePort</code><br/>
<em>
int32
</em>
</td>
<td>
<p>NodePort is the port used to expose the bucket service.
If this is a valid node port, the gateway service type will be set to NodePort accordingly.</p>
</td>
</tr>
<tr>
<td>
<code>blockSyncInterval</code><br/>
<em>
<a href="https://pkg.go.dev/k8s.io/apimachinery/pkg/apis/meta/v1#Duration">
Kubernetes meta/v1.Duration
</a>
</em>
</td>
<td>
<p>Interval to sync block metadata from object storage</p>
</td>
</tr>
<tr>
<td>
<code>gc</code><br/>
<em>
<a href="#monitoring.whizard.io/v1alpha1.BlockGC">
BlockGC
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="monitoring.whizard.io/v1alpha1.CacheProvider">CacheProvider
(<code>string</code> alias)</h3>
<p>
(<em>Appears on:</em><a href="#monitoring.whizard.io/v1alpha1.ResponseCacheProviderConfig">ResponseCacheProviderConfig</a>)
</p>
<div>
</div>
<table>
<thead>
<tr>
<th>Value</th>
<th>Description</th>
</tr>
</thead>
<tbody><tr><td><p>&#34;IN-MEMORY&#34;</p></td>
<td></td>
</tr><tr><td><p>&#34;MEMCACHED&#34;</p></td>
<td></td>
</tr><tr><td><p>&#34;REDIS&#34;</p></td>
<td></td>
</tr></tbody>
</table>
<h3 id="monitoring.whizard.io/v1alpha1.CommonSpec">CommonSpec
</h3>
<p>
(<em>Appears on:</em><a href="#monitoring.whizard.io/v1alpha1.BlockManager">BlockManager</a>, <a href="#monitoring.whizard.io/v1alpha1.CompactorSpec">CompactorSpec</a>, <a href="#monitoring.whizard.io/v1alpha1.GatewaySpec">GatewaySpec</a>, <a href="#monitoring.whizard.io/v1alpha1.IngesterSpec">IngesterSpec</a>, <a href="#monitoring.whizard.io/v1alpha1.QueryFrontendSpec">QueryFrontendSpec</a>, <a href="#monitoring.whizard.io/v1alpha1.QuerySpec">QuerySpec</a>, <a href="#monitoring.whizard.io/v1alpha1.RouterSpec">RouterSpec</a>, <a href="#monitoring.whizard.io/v1alpha1.RulerSpec">RulerSpec</a>, <a href="#monitoring.whizard.io/v1alpha1.StoreSpec">StoreSpec</a>)
</p>
<div>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>replicas</code><br/>
<em>
int32
</em>
</td>
<td>
<p>Number of component instances to deploy.</p>
</td>
</tr>
<tr>
<td>
<code>image</code><br/>
<em>
string
</em>
</td>
<td>
<p>Component container image URL.</p>
</td>
</tr>
<tr>
<td>
<code>imagePullPolicy</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#pullpolicy-v1-core">
Kubernetes core/v1.PullPolicy
</a>
</em>
</td>
<td>
<p>Image pull policy.</p>
</td>
</tr>
<tr>
<td>
<code>resources</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#resourcerequirements-v1-core">
Kubernetes core/v1.ResourceRequirements
</a>
</em>
</td>
<td>
<p>Resources defines the resource requirements for single Pods.</p>
</td>
</tr>
<tr>
<td>
<code>logLevel</code><br/>
<em>
string
</em>
</td>
<td>
<p>Log level for component to be configured with.</p>
</td>
</tr>
<tr>
<td>
<code>logFormat</code><br/>
<em>
string
</em>
</td>
<td>
<p>Log format for component to be configured with.</p>
</td>
</tr>
<tr>
<td>
<code>flags</code><br/>
<em>
[]string
</em>
</td>
<td>
<p>Flags allows setting additional flags for the component container.</p>
</td>
</tr>
<tr>
<td>
<code>podMetadata</code><br/>
<em>
<a href="#monitoring.whizard.io/v1alpha1.EmbeddedObjectMetadata">
EmbeddedObjectMetadata
</a>
</em>
</td>
<td>
<p>PodMetadata configures labels and annotations which are propagated to the pods.</p>
</td>
</tr>
<tr>
<td>
<code>configMaps</code><br/>
<em>
[]string
</em>
</td>
<td>
<p>ConfigMaps is a list of ConfigMaps in the same namespace as the component
object, which shall be mounted into the default Pods.
Each ConfigMap is added to the StatefulSet/Deployment definition as a volume named <code>configmap-&lt;configmap-name&gt;</code>.
The ConfigMaps are mounted into /etc/whizard/configmaps/<configmap-name> in the default container.</p>
</td>
</tr>
<tr>
<td>
<code>secrets</code><br/>
<em>
[]string
</em>
</td>
<td>
<p>Secrets is a list of Secrets in the same namespace as the component
object, which shall be mounted into the Prometheus Pods.
Each Secret is added to the StatefulSet/Deployment definition as a volume named <code>secret-&lt;secret-name&gt;</code>.
The Secrets are mounted into /etc/whizard/secrets/<secret-name> in the default container.</p>
</td>
</tr>
<tr>
<td>
<code>containers</code><br/>
<em>
<a href="https://pkg.go.dev/k8s.io/apimachinery/pkg/runtime.#RawExtension">
k8s.io/apimachinery/pkg/runtime.RawExtension
</a>
</em>
</td>
<td>
<p>Containers allows injecting additional containers or modifying operator generated containers.
Containers described here modify an operator generated
container if they share the same name and modifications are done via a
strategic merge patch.</p>
</td>
</tr>
<tr>
<td>
<code>imagePullSecrets</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#localobjectreference-v1-core">
[]Kubernetes core/v1.LocalObjectReference
</a>
</em>
</td>
<td>
<p>An optional list of references to secrets in the same namespace
to use for pulling images from registries</p>
</td>
</tr>
<tr>
<td>
<code>securityContext</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#podsecuritycontext-v1-core">
Kubernetes core/v1.PodSecurityContext
</a>
</em>
</td>
<td>
<p>SecurityContext holds pod-level security attributes and common container settings.
This defaults to the default PodSecurityContext.</p>
</td>
</tr>
<tr>
<td>
<code>affinity</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#affinity-v1-core">
Kubernetes core/v1.Affinity
</a>
</em>
</td>
<td>
<p>If specified, the pod&rsquo;s scheduling constraints.</p>
</td>
</tr>
<tr>
<td>
<code>nodeSelector</code><br/>
<em>
map[string]string
</em>
</td>
<td>
<p>Define which Nodes the Pods are scheduled on.</p>
</td>
</tr>
<tr>
<td>
<code>tolerations</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#toleration-v1-core">
[]Kubernetes core/v1.Toleration
</a>
</em>
</td>
<td>
<p>If specified, the pod&rsquo;s tolerations.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="monitoring.whizard.io/v1alpha1.CompactorSpec">CompactorSpec
</h3>
<p>
(<em>Appears on:</em><a href="#monitoring.whizard.io/v1alpha1.Compactor">Compactor</a>, <a href="#monitoring.whizard.io/v1alpha1.CompactorTemplateSpec">CompactorTemplateSpec</a>)
</p>
<div>
<p>CompactorSpec defines the desired state of Compactor</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>tenants</code><br/>
<em>
[]string
</em>
</td>
<td>
<p>The tenants whose data is being compacted by the Compactor.</p>
</td>
</tr>
<tr>
<td>
<code>disableDownsampling</code><br/>
<em>
bool
</em>
</td>
<td>
<p>Disables downsampling.
This is not recommended, as querying long time ranges without non-downsampled data is not efficient and useful.
default: false</p>
</td>
</tr>
<tr>
<td>
<code>retention</code><br/>
<em>
<a href="#monitoring.whizard.io/v1alpha1.Retention">
Retention
</a>
</em>
</td>
<td>
<p>Retention configs how long to retain samples</p>
</td>
</tr>
<tr>
<td>
<code>dataVolume</code><br/>
<em>
<a href="#monitoring.whizard.io/v1alpha1.KubernetesVolume">
KubernetesVolume
</a>
</em>
</td>
<td>
<p>DataVolume specifies how volume shall be used</p>
</td>
</tr>
<tr>
<td>
<code>replicas</code><br/>
<em>
int32
</em>
</td>
<td>
<p>Number of component instances to deploy.</p>
</td>
</tr>
<tr>
<td>
<code>image</code><br/>
<em>
string
</em>
</td>
<td>
<p>Component container image URL.</p>
</td>
</tr>
<tr>
<td>
<code>imagePullPolicy</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#pullpolicy-v1-core">
Kubernetes core/v1.PullPolicy
</a>
</em>
</td>
<td>
<p>Image pull policy.</p>
</td>
</tr>
<tr>
<td>
<code>resources</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#resourcerequirements-v1-core">
Kubernetes core/v1.ResourceRequirements
</a>
</em>
</td>
<td>
<p>Resources defines the resource requirements for single Pods.</p>
</td>
</tr>
<tr>
<td>
<code>logLevel</code><br/>
<em>
string
</em>
</td>
<td>
<p>Log level for component to be configured with.</p>
</td>
</tr>
<tr>
<td>
<code>logFormat</code><br/>
<em>
string
</em>
</td>
<td>
<p>Log format for component to be configured with.</p>
</td>
</tr>
<tr>
<td>
<code>flags</code><br/>
<em>
[]string
</em>
</td>
<td>
<p>Flags allows setting additional flags for the component container.</p>
</td>
</tr>
<tr>
<td>
<code>podMetadata</code><br/>
<em>
<a href="#monitoring.whizard.io/v1alpha1.EmbeddedObjectMetadata">
EmbeddedObjectMetadata
</a>
</em>
</td>
<td>
<p>PodMetadata configures labels and annotations which are propagated to the pods.</p>
</td>
</tr>
<tr>
<td>
<code>configMaps</code><br/>
<em>
[]string
</em>
</td>
<td>
<p>ConfigMaps is a list of ConfigMaps in the same namespace as the component
object, which shall be mounted into the default Pods.
Each ConfigMap is added to the StatefulSet/Deployment definition as a volume named <code>configmap-&lt;configmap-name&gt;</code>.
The ConfigMaps are mounted into /etc/whizard/configmaps/<configmap-name> in the default container.</p>
</td>
</tr>
<tr>
<td>
<code>secrets</code><br/>
<em>
[]string
</em>
</td>
<td>
<p>Secrets is a list of Secrets in the same namespace as the component
object, which shall be mounted into the Prometheus Pods.
Each Secret is added to the StatefulSet/Deployment definition as a volume named <code>secret-&lt;secret-name&gt;</code>.
The Secrets are mounted into /etc/whizard/secrets/<secret-name> in the default container.</p>
</td>
</tr>
<tr>
<td>
<code>containers</code><br/>
<em>
<a href="https://pkg.go.dev/k8s.io/apimachinery/pkg/runtime.#RawExtension">
k8s.io/apimachinery/pkg/runtime.RawExtension
</a>
</em>
</td>
<td>
<p>Containers allows injecting additional containers or modifying operator generated containers.
Containers described here modify an operator generated
container if they share the same name and modifications are done via a
strategic merge patch.</p>
</td>
</tr>
<tr>
<td>
<code>imagePullSecrets</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#localobjectreference-v1-core">
[]Kubernetes core/v1.LocalObjectReference
</a>
</em>
</td>
<td>
<p>An optional list of references to secrets in the same namespace
to use for pulling images from registries</p>
</td>
</tr>
<tr>
<td>
<code>securityContext</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#podsecuritycontext-v1-core">
Kubernetes core/v1.PodSecurityContext
</a>
</em>
</td>
<td>
<p>SecurityContext holds pod-level security attributes and common container settings.
This defaults to the default PodSecurityContext.</p>
</td>
</tr>
<tr>
<td>
<code>affinity</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#affinity-v1-core">
Kubernetes core/v1.Affinity
</a>
</em>
</td>
<td>
<p>If specified, the pod&rsquo;s scheduling constraints.</p>
</td>
</tr>
<tr>
<td>
<code>nodeSelector</code><br/>
<em>
map[string]string
</em>
</td>
<td>
<p>Define which Nodes the Pods are scheduled on.</p>
</td>
</tr>
<tr>
<td>
<code>tolerations</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#toleration-v1-core">
[]Kubernetes core/v1.Toleration
</a>
</em>
</td>
<td>
<p>If specified, the pod&rsquo;s tolerations.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="monitoring.whizard.io/v1alpha1.CompactorStatus">CompactorStatus
</h3>
<p>
(<em>Appears on:</em><a href="#monitoring.whizard.io/v1alpha1.Compactor">Compactor</a>)
</p>
<div>
<p>CompactorStatus defines the observed state of Compactor</p>
</div>
<h3 id="monitoring.whizard.io/v1alpha1.CompactorTemplateSpec">CompactorTemplateSpec
</h3>
<p>
(<em>Appears on:</em><a href="#monitoring.whizard.io/v1alpha1.ServiceSpec">ServiceSpec</a>)
</p>
<div>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>tenants</code><br/>
<em>
[]string
</em>
</td>
<td>
<p>The tenants whose data is being compacted by the Compactor.</p>
</td>
</tr>
<tr>
<td>
<code>disableDownsampling</code><br/>
<em>
bool
</em>
</td>
<td>
<p>Disables downsampling.
This is not recommended, as querying long time ranges without non-downsampled data is not efficient and useful.
default: false</p>
</td>
</tr>
<tr>
<td>
<code>retention</code><br/>
<em>
<a href="#monitoring.whizard.io/v1alpha1.Retention">
Retention
</a>
</em>
</td>
<td>
<p>Retention configs how long to retain samples</p>
</td>
</tr>
<tr>
<td>
<code>dataVolume</code><br/>
<em>
<a href="#monitoring.whizard.io/v1alpha1.KubernetesVolume">
KubernetesVolume
</a>
</em>
</td>
<td>
<p>DataVolume specifies how volume shall be used</p>
</td>
</tr>
<tr>
<td>
<code>replicas</code><br/>
<em>
int32
</em>
</td>
<td>
<p>Number of component instances to deploy.</p>
</td>
</tr>
<tr>
<td>
<code>image</code><br/>
<em>
string
</em>
</td>
<td>
<p>Component container image URL.</p>
</td>
</tr>
<tr>
<td>
<code>imagePullPolicy</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#pullpolicy-v1-core">
Kubernetes core/v1.PullPolicy
</a>
</em>
</td>
<td>
<p>Image pull policy.</p>
</td>
</tr>
<tr>
<td>
<code>resources</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#resourcerequirements-v1-core">
Kubernetes core/v1.ResourceRequirements
</a>
</em>
</td>
<td>
<p>Resources defines the resource requirements for single Pods.</p>
</td>
</tr>
<tr>
<td>
<code>logLevel</code><br/>
<em>
string
</em>
</td>
<td>
<p>Log level for component to be configured with.</p>
</td>
</tr>
<tr>
<td>
<code>logFormat</code><br/>
<em>
string
</em>
</td>
<td>
<p>Log format for component to be configured with.</p>
</td>
</tr>
<tr>
<td>
<code>flags</code><br/>
<em>
[]string
</em>
</td>
<td>
<p>Flags allows setting additional flags for the component container.</p>
</td>
</tr>
<tr>
<td>
<code>podMetadata</code><br/>
<em>
<a href="#monitoring.whizard.io/v1alpha1.EmbeddedObjectMetadata">
EmbeddedObjectMetadata
</a>
</em>
</td>
<td>
<p>PodMetadata configures labels and annotations which are propagated to the pods.</p>
</td>
</tr>
<tr>
<td>
<code>configMaps</code><br/>
<em>
[]string
</em>
</td>
<td>
<p>ConfigMaps is a list of ConfigMaps in the same namespace as the component
object, which shall be mounted into the default Pods.
Each ConfigMap is added to the StatefulSet/Deployment definition as a volume named <code>configmap-&lt;configmap-name&gt;</code>.
The ConfigMaps are mounted into /etc/whizard/configmaps/<configmap-name> in the default container.</p>
</td>
</tr>
<tr>
<td>
<code>secrets</code><br/>
<em>
[]string
</em>
</td>
<td>
<p>Secrets is a list of Secrets in the same namespace as the component
object, which shall be mounted into the Prometheus Pods.
Each Secret is added to the StatefulSet/Deployment definition as a volume named <code>secret-&lt;secret-name&gt;</code>.
The Secrets are mounted into /etc/whizard/secrets/<secret-name> in the default container.</p>
</td>
</tr>
<tr>
<td>
<code>containers</code><br/>
<em>
<a href="https://pkg.go.dev/k8s.io/apimachinery/pkg/runtime.#RawExtension">
k8s.io/apimachinery/pkg/runtime.RawExtension
</a>
</em>
</td>
<td>
<p>Containers allows injecting additional containers or modifying operator generated containers.
Containers described here modify an operator generated
container if they share the same name and modifications are done via a
strategic merge patch.</p>
</td>
</tr>
<tr>
<td>
<code>imagePullSecrets</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#localobjectreference-v1-core">
[]Kubernetes core/v1.LocalObjectReference
</a>
</em>
</td>
<td>
<p>An optional list of references to secrets in the same namespace
to use for pulling images from registries</p>
</td>
</tr>
<tr>
<td>
<code>securityContext</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#podsecuritycontext-v1-core">
Kubernetes core/v1.PodSecurityContext
</a>
</em>
</td>
<td>
<p>SecurityContext holds pod-level security attributes and common container settings.
This defaults to the default PodSecurityContext.</p>
</td>
</tr>
<tr>
<td>
<code>affinity</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#affinity-v1-core">
Kubernetes core/v1.Affinity
</a>
</em>
</td>
<td>
<p>If specified, the pod&rsquo;s scheduling constraints.</p>
</td>
</tr>
<tr>
<td>
<code>nodeSelector</code><br/>
<em>
map[string]string
</em>
</td>
<td>
<p>Define which Nodes the Pods are scheduled on.</p>
</td>
</tr>
<tr>
<td>
<code>tolerations</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#toleration-v1-core">
[]Kubernetes core/v1.Toleration
</a>
</em>
</td>
<td>
<p>If specified, the pod&rsquo;s tolerations.</p>
</td>
</tr>
<tr>
<td>
<code>defaultTenantsPerCompactor</code><br/>
<em>
int
</em>
</td>
<td>
<p>DefaultTenantsPerIngester Whizard default tenant count per ingester.
Default: 10</p>
</td>
</tr>
</tbody>
</table>
<h3 id="monitoring.whizard.io/v1alpha1.Duration">Duration
(<code>string</code> alias)</h3>
<p>
(<em>Appears on:</em><a href="#monitoring.whizard.io/v1alpha1.IngesterTemplateSpec">IngesterTemplateSpec</a>, <a href="#monitoring.whizard.io/v1alpha1.RemoteWriteSpec">RemoteWriteSpec</a>, <a href="#monitoring.whizard.io/v1alpha1.Retention">Retention</a>, <a href="#monitoring.whizard.io/v1alpha1.RulerSpec">RulerSpec</a>)
</p>
<div>
<p>Duration is a valid time unit
Supported units: y, w, d, h, m, s, ms Examples: <code>30s</code>, <code>1m</code>, <code>1h20m15s</code></p>
</div>
<h3 id="monitoring.whizard.io/v1alpha1.EmbeddedObjectMetadata">EmbeddedObjectMetadata
</h3>
<p>
(<em>Appears on:</em><a href="#monitoring.whizard.io/v1alpha1.CommonSpec">CommonSpec</a>)
</p>
<div>
<p>EmbeddedObjectMetadata contains a subset of the fields included in k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta
Only fields which are relevant to embedded resources are included.</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>name</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>Name must be unique within a namespace. Is required when creating resources, although
some resources may allow a client to request the generation of an appropriate name
automatically. Name is primarily intended for creation idempotence and configuration
definition.
Cannot be updated.
More info: <a href="https://kubernetes.io/docs/concepts/overview/working-with-objects/names#names">https://kubernetes.io/docs/concepts/overview/working-with-objects/names#names</a></p>
</td>
</tr>
<tr>
<td>
<code>labels</code><br/>
<em>
map[string]string
</em>
</td>
<td>
<em>(Optional)</em>
<p>Map of string keys and values that can be used to organize and categorize
(scope and select) objects. May match selectors of replication controllers
and services.
More info: <a href="http://kubernetes.io/docs/user-guide/labels">http://kubernetes.io/docs/user-guide/labels</a></p>
</td>
</tr>
<tr>
<td>
<code>annotations</code><br/>
<em>
map[string]string
</em>
</td>
<td>
<em>(Optional)</em>
<p>Annotations is an unstructured key value map stored with a resource that may be
set by external tools to store and retrieve arbitrary metadata. They are not
queryable and should be preserved when modifying objects.
More info: <a href="http://kubernetes.io/docs/user-guide/annotations">http://kubernetes.io/docs/user-guide/annotations</a></p>
</td>
</tr>
</tbody>
</table>
<h3 id="monitoring.whizard.io/v1alpha1.GatewaySpec">GatewaySpec
</h3>
<p>
(<em>Appears on:</em><a href="#monitoring.whizard.io/v1alpha1.Gateway">Gateway</a>, <a href="#monitoring.whizard.io/v1alpha1.ServiceSpec">ServiceSpec</a>)
</p>
<div>
<p>GatewaySpec defines the desired state of Gateway</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>webConfig</code><br/>
<em>
<a href="#monitoring.whizard.io/v1alpha1.WebConfig">
WebConfig
</a>
</em>
</td>
<td>
<p>Defines the configuration of the Gatewat web server.</p>
</td>
</tr>
<tr>
<td>
<code>debug</code><br/>
<em>
bool
</em>
</td>
<td>
<p>If debug mode is on, gateway will proxy Query UI</p>
<p>This is an <em>experimental feature</em>, it may change in any upcoming release in a breaking way.</p>
</td>
</tr>
<tr>
<td>
<code>enabledTenantsAdmission</code><br/>
<em>
bool
</em>
</td>
<td>
<p>Deny unknown tenant data remote-write and query if enabled</p>
</td>
</tr>
<tr>
<td>
<code>nodePort</code><br/>
<em>
int32
</em>
</td>
<td>
<p>NodePort is the port used to expose the gateway service.
If this is a valid node port, the gateway service type will be set to NodePort accordingly.</p>
</td>
</tr>
<tr>
<td>
<code>replicas</code><br/>
<em>
int32
</em>
</td>
<td>
<p>Number of component instances to deploy.</p>
</td>
</tr>
<tr>
<td>
<code>image</code><br/>
<em>
string
</em>
</td>
<td>
<p>Component container image URL.</p>
</td>
</tr>
<tr>
<td>
<code>imagePullPolicy</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#pullpolicy-v1-core">
Kubernetes core/v1.PullPolicy
</a>
</em>
</td>
<td>
<p>Image pull policy.</p>
</td>
</tr>
<tr>
<td>
<code>resources</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#resourcerequirements-v1-core">
Kubernetes core/v1.ResourceRequirements
</a>
</em>
</td>
<td>
<p>Resources defines the resource requirements for single Pods.</p>
</td>
</tr>
<tr>
<td>
<code>logLevel</code><br/>
<em>
string
</em>
</td>
<td>
<p>Log level for component to be configured with.</p>
</td>
</tr>
<tr>
<td>
<code>logFormat</code><br/>
<em>
string
</em>
</td>
<td>
<p>Log format for component to be configured with.</p>
</td>
</tr>
<tr>
<td>
<code>flags</code><br/>
<em>
[]string
</em>
</td>
<td>
<p>Flags allows setting additional flags for the component container.</p>
</td>
</tr>
<tr>
<td>
<code>podMetadata</code><br/>
<em>
<a href="#monitoring.whizard.io/v1alpha1.EmbeddedObjectMetadata">
EmbeddedObjectMetadata
</a>
</em>
</td>
<td>
<p>PodMetadata configures labels and annotations which are propagated to the pods.</p>
</td>
</tr>
<tr>
<td>
<code>configMaps</code><br/>
<em>
[]string
</em>
</td>
<td>
<p>ConfigMaps is a list of ConfigMaps in the same namespace as the component
object, which shall be mounted into the default Pods.
Each ConfigMap is added to the StatefulSet/Deployment definition as a volume named <code>configmap-&lt;configmap-name&gt;</code>.
The ConfigMaps are mounted into /etc/whizard/configmaps/<configmap-name> in the default container.</p>
</td>
</tr>
<tr>
<td>
<code>secrets</code><br/>
<em>
[]string
</em>
</td>
<td>
<p>Secrets is a list of Secrets in the same namespace as the component
object, which shall be mounted into the Prometheus Pods.
Each Secret is added to the StatefulSet/Deployment definition as a volume named <code>secret-&lt;secret-name&gt;</code>.
The Secrets are mounted into /etc/whizard/secrets/<secret-name> in the default container.</p>
</td>
</tr>
<tr>
<td>
<code>containers</code><br/>
<em>
<a href="https://pkg.go.dev/k8s.io/apimachinery/pkg/runtime.#RawExtension">
k8s.io/apimachinery/pkg/runtime.RawExtension
</a>
</em>
</td>
<td>
<p>Containers allows injecting additional containers or modifying operator generated containers.
Containers described here modify an operator generated
container if they share the same name and modifications are done via a
strategic merge patch.</p>
</td>
</tr>
<tr>
<td>
<code>imagePullSecrets</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#localobjectreference-v1-core">
[]Kubernetes core/v1.LocalObjectReference
</a>
</em>
</td>
<td>
<p>An optional list of references to secrets in the same namespace
to use for pulling images from registries</p>
</td>
</tr>
<tr>
<td>
<code>securityContext</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#podsecuritycontext-v1-core">
Kubernetes core/v1.PodSecurityContext
</a>
</em>
</td>
<td>
<p>SecurityContext holds pod-level security attributes and common container settings.
This defaults to the default PodSecurityContext.</p>
</td>
</tr>
<tr>
<td>
<code>affinity</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#affinity-v1-core">
Kubernetes core/v1.Affinity
</a>
</em>
</td>
<td>
<p>If specified, the pod&rsquo;s scheduling constraints.</p>
</td>
</tr>
<tr>
<td>
<code>nodeSelector</code><br/>
<em>
map[string]string
</em>
</td>
<td>
<p>Define which Nodes the Pods are scheduled on.</p>
</td>
</tr>
<tr>
<td>
<code>tolerations</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#toleration-v1-core">
[]Kubernetes core/v1.Toleration
</a>
</em>
</td>
<td>
<p>If specified, the pod&rsquo;s tolerations.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="monitoring.whizard.io/v1alpha1.GatewayStatus">GatewayStatus
</h3>
<p>
(<em>Appears on:</em><a href="#monitoring.whizard.io/v1alpha1.Gateway">Gateway</a>)
</p>
<div>
<p>GatewayStatus defines the observed state of Gateway</p>
</div>
<h3 id="monitoring.whizard.io/v1alpha1.HTTPClientConfig">HTTPClientConfig
</h3>
<p>
(<em>Appears on:</em><a href="#monitoring.whizard.io/v1alpha1.RemoteQuerySpec">RemoteQuerySpec</a>, <a href="#monitoring.whizard.io/v1alpha1.RemoteWriteSpec">RemoteWriteSpec</a>)
</p>
<div>
<p>HTTPClientConfig configures an HTTP client.</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>basicAuth</code><br/>
<em>
<a href="#monitoring.whizard.io/v1alpha1.BasicAuth">
BasicAuth
</a>
</em>
</td>
<td>
<p>The HTTP basic authentication credentials for the targets.</p>
</td>
</tr>
<tr>
<td>
<code>bearerToken</code><br/>
<em>
string
</em>
</td>
<td>
<p>The bearer token for the targets.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="monitoring.whizard.io/v1alpha1.HTTPServerConfig">HTTPServerConfig
</h3>
<p>
(<em>Appears on:</em><a href="#monitoring.whizard.io/v1alpha1.WebConfig">WebConfig</a>)
</p>
<div>
</div>
<h3 id="monitoring.whizard.io/v1alpha1.HTTPServerTLSConfig">HTTPServerTLSConfig
</h3>
<p>
(<em>Appears on:</em><a href="#monitoring.whizard.io/v1alpha1.WebConfig">WebConfig</a>)
</p>
<div>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>keySecret</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#secretkeyselector-v1-core">
Kubernetes core/v1.SecretKeySelector
</a>
</em>
</td>
<td>
<p>Secret containing the TLS key for the server.</p>
</td>
</tr>
<tr>
<td>
<code>certSecret</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#secretkeyselector-v1-core">
Kubernetes core/v1.SecretKeySelector
</a>
</em>
</td>
<td>
<p>Contains the TLS certificate for the server.</p>
</td>
</tr>
<tr>
<td>
<code>clientCASecret</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#secretkeyselector-v1-core">
Kubernetes core/v1.SecretKeySelector
</a>
</em>
</td>
<td>
<p>Contains the CA certificate for client certificate authentication to the server.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="monitoring.whizard.io/v1alpha1.InMemoryIndexCacheConfig">InMemoryIndexCacheConfig
</h3>
<p>
(<em>Appears on:</em><a href="#monitoring.whizard.io/v1alpha1.IndexCacheConfig">IndexCacheConfig</a>)
</p>
<div>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>maxSize</code><br/>
<em>
string
</em>
</td>
<td>
<p>MaxSize represents overall maximum number of bytes cache can contain.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="monitoring.whizard.io/v1alpha1.InMemoryResponseCacheConfig">InMemoryResponseCacheConfig
</h3>
<p>
(<em>Appears on:</em><a href="#monitoring.whizard.io/v1alpha1.ResponseCacheProviderConfig">ResponseCacheProviderConfig</a>)
</p>
<div>
<p>InMemoryResponseCacheConfig holds the configs for the in-memory cache provider.</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>maxSize</code><br/>
<em>
string
</em>
</td>
<td>
<p>MaxSize represents overall maximum number of bytes cache can contain.</p>
</td>
</tr>
<tr>
<td>
<code>maxSizeItems</code><br/>
<em>
int
</em>
</td>
<td>
<p>MaxSizeItems represents the maximum number of entries in the cache.</p>
</td>
</tr>
<tr>
<td>
<code>validity</code><br/>
<em>
time.Duration
</em>
</td>
<td>
<p>Validity represents the expiry duration for the cache.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="monitoring.whizard.io/v1alpha1.IndexCacheConfig">IndexCacheConfig
</h3>
<p>
(<em>Appears on:</em><a href="#monitoring.whizard.io/v1alpha1.StoreSpec">StoreSpec</a>)
</p>
<div>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>inMemory</code><br/>
<em>
<a href="#monitoring.whizard.io/v1alpha1.InMemoryIndexCacheConfig">
InMemoryIndexCacheConfig
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="monitoring.whizard.io/v1alpha1.IngesterSpec">IngesterSpec
</h3>
<p>
(<em>Appears on:</em><a href="#monitoring.whizard.io/v1alpha1.Ingester">Ingester</a>, <a href="#monitoring.whizard.io/v1alpha1.IngesterTemplateSpec">IngesterTemplateSpec</a>)
</p>
<div>
<p>IngesterSpec defines the desired state of Ingester</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>tenants</code><br/>
<em>
[]string
</em>
</td>
<td>
<p>The tenants whose data is being ingested by the Ingester(ingesting receiver).</p>
</td>
</tr>
<tr>
<td>
<code>otlpEnableTargetInfo</code><br/>
<em>
bool
</em>
</td>
<td>
<p>Enables target information in OTLP metrics ingested by Receive. If enabled, it converts the resource to the target info metric</p>
</td>
</tr>
<tr>
<td>
<code>otlpResourceAttributes</code><br/>
<em>
[]string
</em>
</td>
<td>
<p>Resource attributes to include in OTLP metrics ingested by Receive.</p>
</td>
</tr>
<tr>
<td>
<code>localTsdbRetention</code><br/>
<em>
string
</em>
</td>
<td>
<p>LocalTsdbRetention configs how long to retain raw samples on local storage.</p>
</td>
</tr>
<tr>
<td>
<code>dataVolume</code><br/>
<em>
<a href="#monitoring.whizard.io/v1alpha1.KubernetesVolume">
KubernetesVolume
</a>
</em>
</td>
<td>
<p>DataVolume specifies how volume shall be used</p>
</td>
</tr>
<tr>
<td>
<code>replicas</code><br/>
<em>
int32
</em>
</td>
<td>
<p>Number of component instances to deploy.</p>
</td>
</tr>
<tr>
<td>
<code>image</code><br/>
<em>
string
</em>
</td>
<td>
<p>Component container image URL.</p>
</td>
</tr>
<tr>
<td>
<code>imagePullPolicy</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#pullpolicy-v1-core">
Kubernetes core/v1.PullPolicy
</a>
</em>
</td>
<td>
<p>Image pull policy.</p>
</td>
</tr>
<tr>
<td>
<code>resources</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#resourcerequirements-v1-core">
Kubernetes core/v1.ResourceRequirements
</a>
</em>
</td>
<td>
<p>Resources defines the resource requirements for single Pods.</p>
</td>
</tr>
<tr>
<td>
<code>logLevel</code><br/>
<em>
string
</em>
</td>
<td>
<p>Log level for component to be configured with.</p>
</td>
</tr>
<tr>
<td>
<code>logFormat</code><br/>
<em>
string
</em>
</td>
<td>
<p>Log format for component to be configured with.</p>
</td>
</tr>
<tr>
<td>
<code>flags</code><br/>
<em>
[]string
</em>
</td>
<td>
<p>Flags allows setting additional flags for the component container.</p>
</td>
</tr>
<tr>
<td>
<code>podMetadata</code><br/>
<em>
<a href="#monitoring.whizard.io/v1alpha1.EmbeddedObjectMetadata">
EmbeddedObjectMetadata
</a>
</em>
</td>
<td>
<p>PodMetadata configures labels and annotations which are propagated to the pods.</p>
</td>
</tr>
<tr>
<td>
<code>configMaps</code><br/>
<em>
[]string
</em>
</td>
<td>
<p>ConfigMaps is a list of ConfigMaps in the same namespace as the component
object, which shall be mounted into the default Pods.
Each ConfigMap is added to the StatefulSet/Deployment definition as a volume named <code>configmap-&lt;configmap-name&gt;</code>.
The ConfigMaps are mounted into /etc/whizard/configmaps/<configmap-name> in the default container.</p>
</td>
</tr>
<tr>
<td>
<code>secrets</code><br/>
<em>
[]string
</em>
</td>
<td>
<p>Secrets is a list of Secrets in the same namespace as the component
object, which shall be mounted into the Prometheus Pods.
Each Secret is added to the StatefulSet/Deployment definition as a volume named <code>secret-&lt;secret-name&gt;</code>.
The Secrets are mounted into /etc/whizard/secrets/<secret-name> in the default container.</p>
</td>
</tr>
<tr>
<td>
<code>containers</code><br/>
<em>
<a href="https://pkg.go.dev/k8s.io/apimachinery/pkg/runtime.#RawExtension">
k8s.io/apimachinery/pkg/runtime.RawExtension
</a>
</em>
</td>
<td>
<p>Containers allows injecting additional containers or modifying operator generated containers.
Containers described here modify an operator generated
container if they share the same name and modifications are done via a
strategic merge patch.</p>
</td>
</tr>
<tr>
<td>
<code>imagePullSecrets</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#localobjectreference-v1-core">
[]Kubernetes core/v1.LocalObjectReference
</a>
</em>
</td>
<td>
<p>An optional list of references to secrets in the same namespace
to use for pulling images from registries</p>
</td>
</tr>
<tr>
<td>
<code>securityContext</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#podsecuritycontext-v1-core">
Kubernetes core/v1.PodSecurityContext
</a>
</em>
</td>
<td>
<p>SecurityContext holds pod-level security attributes and common container settings.
This defaults to the default PodSecurityContext.</p>
</td>
</tr>
<tr>
<td>
<code>affinity</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#affinity-v1-core">
Kubernetes core/v1.Affinity
</a>
</em>
</td>
<td>
<p>If specified, the pod&rsquo;s scheduling constraints.</p>
</td>
</tr>
<tr>
<td>
<code>nodeSelector</code><br/>
<em>
map[string]string
</em>
</td>
<td>
<p>Define which Nodes the Pods are scheduled on.</p>
</td>
</tr>
<tr>
<td>
<code>tolerations</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#toleration-v1-core">
[]Kubernetes core/v1.Toleration
</a>
</em>
</td>
<td>
<p>If specified, the pod&rsquo;s tolerations.</p>
</td>
</tr>
<tr>
<td>
<code>ingesterTsdbCleanup</code><br/>
<em>
<a href="#monitoring.whizard.io/v1alpha1.SidecarSpec">
SidecarSpec
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="monitoring.whizard.io/v1alpha1.IngesterStatus">IngesterStatus
</h3>
<p>
(<em>Appears on:</em><a href="#monitoring.whizard.io/v1alpha1.Ingester">Ingester</a>)
</p>
<div>
<p>IngesterStatus defines the observed state of Ingester</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>tenants</code><br/>
<em>
<a href="#monitoring.whizard.io/v1alpha1.IngesterTenantStatus">
[]IngesterTenantStatus
</a>
</em>
</td>
<td>
<p>Tenants contain all tenants that have been configured for this Ingester object,
except those Tenant objects that have been deleted.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="monitoring.whizard.io/v1alpha1.IngesterTemplateSpec">IngesterTemplateSpec
</h3>
<p>
(<em>Appears on:</em><a href="#monitoring.whizard.io/v1alpha1.ServiceSpec">ServiceSpec</a>)
</p>
<div>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>tenants</code><br/>
<em>
[]string
</em>
</td>
<td>
<p>The tenants whose data is being ingested by the Ingester(ingesting receiver).</p>
</td>
</tr>
<tr>
<td>
<code>otlpEnableTargetInfo</code><br/>
<em>
bool
</em>
</td>
<td>
<p>Enables target information in OTLP metrics ingested by Receive. If enabled, it converts the resource to the target info metric</p>
</td>
</tr>
<tr>
<td>
<code>otlpResourceAttributes</code><br/>
<em>
[]string
</em>
</td>
<td>
<p>Resource attributes to include in OTLP metrics ingested by Receive.</p>
</td>
</tr>
<tr>
<td>
<code>localTsdbRetention</code><br/>
<em>
string
</em>
</td>
<td>
<p>LocalTsdbRetention configs how long to retain raw samples on local storage.</p>
</td>
</tr>
<tr>
<td>
<code>dataVolume</code><br/>
<em>
<a href="#monitoring.whizard.io/v1alpha1.KubernetesVolume">
KubernetesVolume
</a>
</em>
</td>
<td>
<p>DataVolume specifies how volume shall be used</p>
</td>
</tr>
<tr>
<td>
<code>replicas</code><br/>
<em>
int32
</em>
</td>
<td>
<p>Number of component instances to deploy.</p>
</td>
</tr>
<tr>
<td>
<code>image</code><br/>
<em>
string
</em>
</td>
<td>
<p>Component container image URL.</p>
</td>
</tr>
<tr>
<td>
<code>imagePullPolicy</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#pullpolicy-v1-core">
Kubernetes core/v1.PullPolicy
</a>
</em>
</td>
<td>
<p>Image pull policy.</p>
</td>
</tr>
<tr>
<td>
<code>resources</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#resourcerequirements-v1-core">
Kubernetes core/v1.ResourceRequirements
</a>
</em>
</td>
<td>
<p>Resources defines the resource requirements for single Pods.</p>
</td>
</tr>
<tr>
<td>
<code>logLevel</code><br/>
<em>
string
</em>
</td>
<td>
<p>Log level for component to be configured with.</p>
</td>
</tr>
<tr>
<td>
<code>logFormat</code><br/>
<em>
string
</em>
</td>
<td>
<p>Log format for component to be configured with.</p>
</td>
</tr>
<tr>
<td>
<code>flags</code><br/>
<em>
[]string
</em>
</td>
<td>
<p>Flags allows setting additional flags for the component container.</p>
</td>
</tr>
<tr>
<td>
<code>podMetadata</code><br/>
<em>
<a href="#monitoring.whizard.io/v1alpha1.EmbeddedObjectMetadata">
EmbeddedObjectMetadata
</a>
</em>
</td>
<td>
<p>PodMetadata configures labels and annotations which are propagated to the pods.</p>
</td>
</tr>
<tr>
<td>
<code>configMaps</code><br/>
<em>
[]string
</em>
</td>
<td>
<p>ConfigMaps is a list of ConfigMaps in the same namespace as the component
object, which shall be mounted into the default Pods.
Each ConfigMap is added to the StatefulSet/Deployment definition as a volume named <code>configmap-&lt;configmap-name&gt;</code>.
The ConfigMaps are mounted into /etc/whizard/configmaps/<configmap-name> in the default container.</p>
</td>
</tr>
<tr>
<td>
<code>secrets</code><br/>
<em>
[]string
</em>
</td>
<td>
<p>Secrets is a list of Secrets in the same namespace as the component
object, which shall be mounted into the Prometheus Pods.
Each Secret is added to the StatefulSet/Deployment definition as a volume named <code>secret-&lt;secret-name&gt;</code>.
The Secrets are mounted into /etc/whizard/secrets/<secret-name> in the default container.</p>
</td>
</tr>
<tr>
<td>
<code>containers</code><br/>
<em>
<a href="https://pkg.go.dev/k8s.io/apimachinery/pkg/runtime.#RawExtension">
k8s.io/apimachinery/pkg/runtime.RawExtension
</a>
</em>
</td>
<td>
<p>Containers allows injecting additional containers or modifying operator generated containers.
Containers described here modify an operator generated
container if they share the same name and modifications are done via a
strategic merge patch.</p>
</td>
</tr>
<tr>
<td>
<code>imagePullSecrets</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#localobjectreference-v1-core">
[]Kubernetes core/v1.LocalObjectReference
</a>
</em>
</td>
<td>
<p>An optional list of references to secrets in the same namespace
to use for pulling images from registries</p>
</td>
</tr>
<tr>
<td>
<code>securityContext</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#podsecuritycontext-v1-core">
Kubernetes core/v1.PodSecurityContext
</a>
</em>
</td>
<td>
<p>SecurityContext holds pod-level security attributes and common container settings.
This defaults to the default PodSecurityContext.</p>
</td>
</tr>
<tr>
<td>
<code>affinity</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#affinity-v1-core">
Kubernetes core/v1.Affinity
</a>
</em>
</td>
<td>
<p>If specified, the pod&rsquo;s scheduling constraints.</p>
</td>
</tr>
<tr>
<td>
<code>nodeSelector</code><br/>
<em>
map[string]string
</em>
</td>
<td>
<p>Define which Nodes the Pods are scheduled on.</p>
</td>
</tr>
<tr>
<td>
<code>tolerations</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#toleration-v1-core">
[]Kubernetes core/v1.Toleration
</a>
</em>
</td>
<td>
<p>If specified, the pod&rsquo;s tolerations.</p>
</td>
</tr>
<tr>
<td>
<code>ingesterTsdbCleanup</code><br/>
<em>
<a href="#monitoring.whizard.io/v1alpha1.SidecarSpec">
SidecarSpec
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>defaultTenantsPerIngester</code><br/>
<em>
int
</em>
</td>
<td>
<p>DefaultTenantsPerIngester Whizard default tenant count per ingester.</p>
<p>Default: 3</p>
</td>
</tr>
<tr>
<td>
<code>defaultIngesterRetentionPeriod</code><br/>
<em>
<a href="#monitoring.whizard.io/v1alpha1.Duration">
Duration
</a>
</em>
</td>
<td>
<p>DefaultIngesterRetentionPeriod Whizard default ingester retention period when it has no tenant.</p>
<p>Default: &ldquo;3h&rdquo;</p>
</td>
</tr>
<tr>
<td>
<code>disableTsdbCleanup</code><br/>
<em>
bool
</em>
</td>
<td>
<p>DisableTSDBCleanup Disable the TSDB cleanup of ingester.
The cleanup will delete the blocks that belong to deleted tenants in the data directory of ingester TSDB.</p>
<p>Default: true</p>
</td>
</tr>
</tbody>
</table>
<h3 id="monitoring.whizard.io/v1alpha1.IngesterTenantStatus">IngesterTenantStatus
</h3>
<p>
(<em>Appears on:</em><a href="#monitoring.whizard.io/v1alpha1.IngesterStatus">IngesterStatus</a>)
</p>
<div>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>name</code><br/>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>obsolete</code><br/>
<em>
bool
</em>
</td>
<td>
<p>true represents that the tenant has been moved to other ingester but may left tsdb data in this ingester.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="monitoring.whizard.io/v1alpha1.KubernetesVolume">KubernetesVolume
</h3>
<p>
(<em>Appears on:</em><a href="#monitoring.whizard.io/v1alpha1.CompactorSpec">CompactorSpec</a>, <a href="#monitoring.whizard.io/v1alpha1.IngesterSpec">IngesterSpec</a>, <a href="#monitoring.whizard.io/v1alpha1.RulerSpec">RulerSpec</a>, <a href="#monitoring.whizard.io/v1alpha1.StoreSpec">StoreSpec</a>)
</p>
<div>
<p>KubernetesVolume defines the configured storage for component.
If no storage option is specified, then by default an <a href="https://kubernetes.io/docs/concepts/storage/volumes/#emptydir">EmptyDir</a> will be used.</p>
<p>If multiple storage options are specified, priority will be given as follows:
1. emptyDir
2. persistentVolumeClaim</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>emptyDir</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#emptydirvolumesource-v1-core">
Kubernetes core/v1.EmptyDirVolumeSource
</a>
</em>
</td>
<td>
<p>emptyDir represents a temporary directory that shares a pod&rsquo;s lifetime.</p>
</td>
</tr>
<tr>
<td>
<code>persistentVolumeClaim</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#persistentvolumeclaim-v1-core">
Kubernetes core/v1.PersistentVolumeClaim
</a>
</em>
</td>
<td>
<p>Defines the PVC spec to be used by the component StatefulSets.</p>
</td>
</tr>
<tr>
<td>
<code>persistentVolumeClaimRetentionPolicy</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#statefulsetpersistentvolumeclaimretentionpolicy-v1-apps">
Kubernetes apps/v1.StatefulSetPersistentVolumeClaimRetentionPolicy
</a>
</em>
</td>
<td>
<p>persistentVolumeClaimRetentionPolicy describes the lifecycle of persistent
volume claims created from persistentVolumeClaim.
This requires the kubernetes version &gt;= 1.23 and its StatefulSetAutoDeletePVC feature gate to be enabled.</p>
<p>This is an <em>experimental feature</em>, it may change in any upcoming release in a breaking way.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="monitoring.whizard.io/v1alpha1.ObjectReference">ObjectReference
</h3>
<p>
(<em>Appears on:</em><a href="#monitoring.whizard.io/v1alpha1.ServiceSpec">ServiceSpec</a>, <a href="#monitoring.whizard.io/v1alpha1.TenantStatus">TenantStatus</a>)
</p>
<div>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>namespace</code><br/>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>name</code><br/>
<em>
string
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="monitoring.whizard.io/v1alpha1.QueryFrontendSpec">QueryFrontendSpec
</h3>
<p>
(<em>Appears on:</em><a href="#monitoring.whizard.io/v1alpha1.QueryFrontend">QueryFrontend</a>, <a href="#monitoring.whizard.io/v1alpha1.ServiceSpec">ServiceSpec</a>)
</p>
<div>
<p>QueryFrontendSpec defines the desired state of QueryFrontend</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>cacheConfig</code><br/>
<em>
<a href="#monitoring.whizard.io/v1alpha1.ResponseCacheProviderConfig">
ResponseCacheProviderConfig
</a>
</em>
</td>
<td>
<p>CacheProviderConfig specifies response cache configuration.</p>
</td>
</tr>
<tr>
<td>
<code>webConfig</code><br/>
<em>
<a href="#monitoring.whizard.io/v1alpha1.WebConfig">
WebConfig
</a>
</em>
</td>
<td>
<p>Defines the configuration of the Thanos QueryFrontend web server.</p>
</td>
</tr>
<tr>
<td>
<code>replicas</code><br/>
<em>
int32
</em>
</td>
<td>
<p>Number of component instances to deploy.</p>
</td>
</tr>
<tr>
<td>
<code>image</code><br/>
<em>
string
</em>
</td>
<td>
<p>Component container image URL.</p>
</td>
</tr>
<tr>
<td>
<code>imagePullPolicy</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#pullpolicy-v1-core">
Kubernetes core/v1.PullPolicy
</a>
</em>
</td>
<td>
<p>Image pull policy.</p>
</td>
</tr>
<tr>
<td>
<code>resources</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#resourcerequirements-v1-core">
Kubernetes core/v1.ResourceRequirements
</a>
</em>
</td>
<td>
<p>Resources defines the resource requirements for single Pods.</p>
</td>
</tr>
<tr>
<td>
<code>logLevel</code><br/>
<em>
string
</em>
</td>
<td>
<p>Log level for component to be configured with.</p>
</td>
</tr>
<tr>
<td>
<code>logFormat</code><br/>
<em>
string
</em>
</td>
<td>
<p>Log format for component to be configured with.</p>
</td>
</tr>
<tr>
<td>
<code>flags</code><br/>
<em>
[]string
</em>
</td>
<td>
<p>Flags allows setting additional flags for the component container.</p>
</td>
</tr>
<tr>
<td>
<code>podMetadata</code><br/>
<em>
<a href="#monitoring.whizard.io/v1alpha1.EmbeddedObjectMetadata">
EmbeddedObjectMetadata
</a>
</em>
</td>
<td>
<p>PodMetadata configures labels and annotations which are propagated to the pods.</p>
</td>
</tr>
<tr>
<td>
<code>configMaps</code><br/>
<em>
[]string
</em>
</td>
<td>
<p>ConfigMaps is a list of ConfigMaps in the same namespace as the component
object, which shall be mounted into the default Pods.
Each ConfigMap is added to the StatefulSet/Deployment definition as a volume named <code>configmap-&lt;configmap-name&gt;</code>.
The ConfigMaps are mounted into /etc/whizard/configmaps/<configmap-name> in the default container.</p>
</td>
</tr>
<tr>
<td>
<code>secrets</code><br/>
<em>
[]string
</em>
</td>
<td>
<p>Secrets is a list of Secrets in the same namespace as the component
object, which shall be mounted into the Prometheus Pods.
Each Secret is added to the StatefulSet/Deployment definition as a volume named <code>secret-&lt;secret-name&gt;</code>.
The Secrets are mounted into /etc/whizard/secrets/<secret-name> in the default container.</p>
</td>
</tr>
<tr>
<td>
<code>containers</code><br/>
<em>
<a href="https://pkg.go.dev/k8s.io/apimachinery/pkg/runtime.#RawExtension">
k8s.io/apimachinery/pkg/runtime.RawExtension
</a>
</em>
</td>
<td>
<p>Containers allows injecting additional containers or modifying operator generated containers.
Containers described here modify an operator generated
container if they share the same name and modifications are done via a
strategic merge patch.</p>
</td>
</tr>
<tr>
<td>
<code>imagePullSecrets</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#localobjectreference-v1-core">
[]Kubernetes core/v1.LocalObjectReference
</a>
</em>
</td>
<td>
<p>An optional list of references to secrets in the same namespace
to use for pulling images from registries</p>
</td>
</tr>
<tr>
<td>
<code>securityContext</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#podsecuritycontext-v1-core">
Kubernetes core/v1.PodSecurityContext
</a>
</em>
</td>
<td>
<p>SecurityContext holds pod-level security attributes and common container settings.
This defaults to the default PodSecurityContext.</p>
</td>
</tr>
<tr>
<td>
<code>affinity</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#affinity-v1-core">
Kubernetes core/v1.Affinity
</a>
</em>
</td>
<td>
<p>If specified, the pod&rsquo;s scheduling constraints.</p>
</td>
</tr>
<tr>
<td>
<code>nodeSelector</code><br/>
<em>
map[string]string
</em>
</td>
<td>
<p>Define which Nodes the Pods are scheduled on.</p>
</td>
</tr>
<tr>
<td>
<code>tolerations</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#toleration-v1-core">
[]Kubernetes core/v1.Toleration
</a>
</em>
</td>
<td>
<p>If specified, the pod&rsquo;s tolerations.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="monitoring.whizard.io/v1alpha1.QueryFrontendStatus">QueryFrontendStatus
</h3>
<p>
(<em>Appears on:</em><a href="#monitoring.whizard.io/v1alpha1.QueryFrontend">QueryFrontend</a>)
</p>
<div>
<p>QueryFrontendStatus defines the observed state of QueryFrontend</p>
</div>
<h3 id="monitoring.whizard.io/v1alpha1.QuerySpec">QuerySpec
</h3>
<p>
(<em>Appears on:</em><a href="#monitoring.whizard.io/v1alpha1.Query">Query</a>, <a href="#monitoring.whizard.io/v1alpha1.ServiceSpec">ServiceSpec</a>)
</p>
<div>
<p>QuerySpec defines the desired state of Query</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>promqlEngine</code><br/>
<em>
string
</em>
</td>
<td>
<p>experimental PromQL engine, more info thanos.io/tip/components/query.md#promql-engine
default: prometheus</p>
</td>
</tr>
<tr>
<td>
<code>selectorLabels</code><br/>
<em>
map[string]string
</em>
</td>
<td>
<p>Selector labels that will be exposed in info endpoint.</p>
</td>
</tr>
<tr>
<td>
<code>replicaLabelNames</code><br/>
<em>
[]string
</em>
</td>
<td>
<p>Labels to treat as a replica indicator along which data is deduplicated.</p>
</td>
</tr>
<tr>
<td>
<code>webConfig</code><br/>
<em>
<a href="#monitoring.whizard.io/v1alpha1.WebConfig">
WebConfig
</a>
</em>
</td>
<td>
<p>Defines the configuration of the Thanos Query web server.</p>
</td>
</tr>
<tr>
<td>
<code>stores</code><br/>
<em>
<a href="#monitoring.whizard.io/v1alpha1.QueryStores">
[]QueryStores
</a>
</em>
</td>
<td>
<p>Additional StoreApi servers from which Query component queries from</p>
</td>
</tr>
<tr>
<td>
<code>envoy</code><br/>
<em>
<a href="#monitoring.whizard.io/v1alpha1.SidecarSpec">
SidecarSpec
</a>
</em>
</td>
<td>
<p>Envoy is used to config sidecar which proxies requests requiring auth to the secure stores</p>
</td>
</tr>
<tr>
<td>
<code>replicas</code><br/>
<em>
int32
</em>
</td>
<td>
<p>Number of component instances to deploy.</p>
</td>
</tr>
<tr>
<td>
<code>image</code><br/>
<em>
string
</em>
</td>
<td>
<p>Component container image URL.</p>
</td>
</tr>
<tr>
<td>
<code>imagePullPolicy</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#pullpolicy-v1-core">
Kubernetes core/v1.PullPolicy
</a>
</em>
</td>
<td>
<p>Image pull policy.</p>
</td>
</tr>
<tr>
<td>
<code>resources</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#resourcerequirements-v1-core">
Kubernetes core/v1.ResourceRequirements
</a>
</em>
</td>
<td>
<p>Resources defines the resource requirements for single Pods.</p>
</td>
</tr>
<tr>
<td>
<code>logLevel</code><br/>
<em>
string
</em>
</td>
<td>
<p>Log level for component to be configured with.</p>
</td>
</tr>
<tr>
<td>
<code>logFormat</code><br/>
<em>
string
</em>
</td>
<td>
<p>Log format for component to be configured with.</p>
</td>
</tr>
<tr>
<td>
<code>flags</code><br/>
<em>
[]string
</em>
</td>
<td>
<p>Flags allows setting additional flags for the component container.</p>
</td>
</tr>
<tr>
<td>
<code>podMetadata</code><br/>
<em>
<a href="#monitoring.whizard.io/v1alpha1.EmbeddedObjectMetadata">
EmbeddedObjectMetadata
</a>
</em>
</td>
<td>
<p>PodMetadata configures labels and annotations which are propagated to the pods.</p>
</td>
</tr>
<tr>
<td>
<code>configMaps</code><br/>
<em>
[]string
</em>
</td>
<td>
<p>ConfigMaps is a list of ConfigMaps in the same namespace as the component
object, which shall be mounted into the default Pods.
Each ConfigMap is added to the StatefulSet/Deployment definition as a volume named <code>configmap-&lt;configmap-name&gt;</code>.
The ConfigMaps are mounted into /etc/whizard/configmaps/<configmap-name> in the default container.</p>
</td>
</tr>
<tr>
<td>
<code>secrets</code><br/>
<em>
[]string
</em>
</td>
<td>
<p>Secrets is a list of Secrets in the same namespace as the component
object, which shall be mounted into the Prometheus Pods.
Each Secret is added to the StatefulSet/Deployment definition as a volume named <code>secret-&lt;secret-name&gt;</code>.
The Secrets are mounted into /etc/whizard/secrets/<secret-name> in the default container.</p>
</td>
</tr>
<tr>
<td>
<code>containers</code><br/>
<em>
<a href="https://pkg.go.dev/k8s.io/apimachinery/pkg/runtime.#RawExtension">
k8s.io/apimachinery/pkg/runtime.RawExtension
</a>
</em>
</td>
<td>
<p>Containers allows injecting additional containers or modifying operator generated containers.
Containers described here modify an operator generated
container if they share the same name and modifications are done via a
strategic merge patch.</p>
</td>
</tr>
<tr>
<td>
<code>imagePullSecrets</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#localobjectreference-v1-core">
[]Kubernetes core/v1.LocalObjectReference
</a>
</em>
</td>
<td>
<p>An optional list of references to secrets in the same namespace
to use for pulling images from registries</p>
</td>
</tr>
<tr>
<td>
<code>securityContext</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#podsecuritycontext-v1-core">
Kubernetes core/v1.PodSecurityContext
</a>
</em>
</td>
<td>
<p>SecurityContext holds pod-level security attributes and common container settings.
This defaults to the default PodSecurityContext.</p>
</td>
</tr>
<tr>
<td>
<code>affinity</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#affinity-v1-core">
Kubernetes core/v1.Affinity
</a>
</em>
</td>
<td>
<p>If specified, the pod&rsquo;s scheduling constraints.</p>
</td>
</tr>
<tr>
<td>
<code>nodeSelector</code><br/>
<em>
map[string]string
</em>
</td>
<td>
<p>Define which Nodes the Pods are scheduled on.</p>
</td>
</tr>
<tr>
<td>
<code>tolerations</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#toleration-v1-core">
[]Kubernetes core/v1.Toleration
</a>
</em>
</td>
<td>
<p>If specified, the pod&rsquo;s tolerations.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="monitoring.whizard.io/v1alpha1.QueryStatus">QueryStatus
</h3>
<p>
(<em>Appears on:</em><a href="#monitoring.whizard.io/v1alpha1.Query">Query</a>)
</p>
<div>
<p>QueryStatus defines the observed state of Query</p>
</div>
<h3 id="monitoring.whizard.io/v1alpha1.QueryStores">QueryStores
</h3>
<p>
(<em>Appears on:</em><a href="#monitoring.whizard.io/v1alpha1.QuerySpec">QuerySpec</a>)
</p>
<div>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>addresses</code><br/>
<em>
[]string
</em>
</td>
<td>
<p>Address is the addresses of StoreApi server, which may be prefixed with &lsquo;dns+&rsquo; or &lsquo;dnssrv+&rsquo; to detect StoreAPI servers through respective DNS lookups.</p>
</td>
</tr>
<tr>
<td>
<code>caSecret</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#secretkeyselector-v1-core">
Kubernetes core/v1.SecretKeySelector
</a>
</em>
</td>
<td>
<p>Secret containing the CA cert to use for StoreApi connections</p>
</td>
</tr>
</tbody>
</table>
<h3 id="monitoring.whizard.io/v1alpha1.RemoteQuerySpec">RemoteQuerySpec
</h3>
<p>
(<em>Appears on:</em><a href="#monitoring.whizard.io/v1alpha1.ServiceSpec">ServiceSpec</a>)
</p>
<div>
<p>RemoteQuerySpec defines the configuration to query from remote service
which should have prometheus-compatible Query APIs.</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>name</code><br/>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>url</code><br/>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>basicAuth</code><br/>
<em>
<a href="#monitoring.whizard.io/v1alpha1.BasicAuth">
BasicAuth
</a>
</em>
</td>
<td>
<p>The HTTP basic authentication credentials for the targets.</p>
</td>
</tr>
<tr>
<td>
<code>bearerToken</code><br/>
<em>
string
</em>
</td>
<td>
<p>The bearer token for the targets.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="monitoring.whizard.io/v1alpha1.RemoteWriteSpec">RemoteWriteSpec
</h3>
<p>
(<em>Appears on:</em><a href="#monitoring.whizard.io/v1alpha1.ServiceSpec">ServiceSpec</a>)
</p>
<div>
<p>RemoteWriteSpec defines the remote write configuration.</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>name</code><br/>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>url</code><br/>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>headers</code><br/>
<em>
map[string]string
</em>
</td>
<td>
<p>Custom HTTP headers to be sent along with each remote write request.</p>
</td>
</tr>
<tr>
<td>
<code>remoteTimeout</code><br/>
<em>
<a href="#monitoring.whizard.io/v1alpha1.Duration">
Duration
</a>
</em>
</td>
<td>
<p>Timeout for requests to the remote write endpoint.</p>
</td>
</tr>
<tr>
<td>
<code>basicAuth</code><br/>
<em>
<a href="#monitoring.whizard.io/v1alpha1.BasicAuth">
BasicAuth
</a>
</em>
</td>
<td>
<p>The HTTP basic authentication credentials for the targets.</p>
</td>
</tr>
<tr>
<td>
<code>bearerToken</code><br/>
<em>
string
</em>
</td>
<td>
<p>The bearer token for the targets.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="monitoring.whizard.io/v1alpha1.ReplicationProtocol">ReplicationProtocol
(<code>string</code> alias)</h3>
<p>
(<em>Appears on:</em><a href="#monitoring.whizard.io/v1alpha1.RouterSpec">RouterSpec</a>)
</p>
<div>
</div>
<table>
<thead>
<tr>
<th>Value</th>
<th>Description</th>
</tr>
</thead>
<tbody><tr><td><p>&#34;capnproto&#34;</p></td>
<td></td>
</tr><tr><td><p>&#34;protobuf&#34;</p></td>
<td></td>
</tr></tbody>
</table>
<h3 id="monitoring.whizard.io/v1alpha1.ResponseCacheProviderConfig">ResponseCacheProviderConfig
</h3>
<p>
(<em>Appears on:</em><a href="#monitoring.whizard.io/v1alpha1.QueryFrontendSpec">QueryFrontendSpec</a>)
</p>
<div>
<p>ResponseCacheProviderConfig is the initial ResponseCacheProviderConfig struct holder before parsing it into a specific cache provider.
Based on the config type the config is then parsed into a specific cache provider.</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>type</code><br/>
<em>
<a href="#monitoring.whizard.io/v1alpha1.CacheProvider">
CacheProvider
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>inMemory</code><br/>
<em>
<a href="#monitoring.whizard.io/v1alpha1.InMemoryResponseCacheConfig">
InMemoryResponseCacheConfig
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="monitoring.whizard.io/v1alpha1.Retention">Retention
</h3>
<p>
(<em>Appears on:</em><a href="#monitoring.whizard.io/v1alpha1.CompactorSpec">CompactorSpec</a>)
</p>
<div>
<p>Retention defines the config for retaining samples</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>retentionRaw</code><br/>
<em>
<a href="#monitoring.whizard.io/v1alpha1.Duration">
Duration
</a>
</em>
</td>
<td>
<p>How long to retain raw samples in bucket. Setting this to 0d will retain samples of this resolution forever
default: 0d</p>
</td>
</tr>
<tr>
<td>
<code>retention5m</code><br/>
<em>
<a href="#monitoring.whizard.io/v1alpha1.Duration">
Duration
</a>
</em>
</td>
<td>
<p>How long to retain samples of resolution 1 (5 minutes) in bucket. Setting this to 0d will retain samples of this resolution forever
default: 0d</p>
</td>
</tr>
<tr>
<td>
<code>retention1h</code><br/>
<em>
<a href="#monitoring.whizard.io/v1alpha1.Duration">
Duration
</a>
</em>
</td>
<td>
<p>How long to retain samples of resolution 2 (1 hour) in bucket. Setting this to 0d will retain samples of this resolution forever
default: 0d</p>
</td>
</tr>
</tbody>
</table>
<h3 id="monitoring.whizard.io/v1alpha1.RouterSpec">RouterSpec
</h3>
<p>
(<em>Appears on:</em><a href="#monitoring.whizard.io/v1alpha1.Router">Router</a>, <a href="#monitoring.whizard.io/v1alpha1.ServiceSpec">ServiceSpec</a>)
</p>
<div>
<p>RouterSpec defines the desired state of Router</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>replicationFactor</code><br/>
<em>
uint64
</em>
</td>
<td>
<p>How many times to replicate incoming write requests</p>
</td>
</tr>
<tr>
<td>
<code>replicationProtocol</code><br/>
<em>
<a href="#monitoring.whizard.io/v1alpha1.ReplicationProtocol">
ReplicationProtocol
</a>
</em>
</td>
<td>
<p>The protocol to use for replicating remote-write requests. One of protobuf,capnproto</p>
</td>
</tr>
<tr>
<td>
<code>webConfig</code><br/>
<em>
<a href="#monitoring.whizard.io/v1alpha1.WebConfig">
WebConfig
</a>
</em>
</td>
<td>
<p>Defines the configuration of the Route(routing receiver) web server.</p>
</td>
</tr>
<tr>
<td>
<code>replicas</code><br/>
<em>
int32
</em>
</td>
<td>
<p>Number of component instances to deploy.</p>
</td>
</tr>
<tr>
<td>
<code>image</code><br/>
<em>
string
</em>
</td>
<td>
<p>Component container image URL.</p>
</td>
</tr>
<tr>
<td>
<code>imagePullPolicy</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#pullpolicy-v1-core">
Kubernetes core/v1.PullPolicy
</a>
</em>
</td>
<td>
<p>Image pull policy.</p>
</td>
</tr>
<tr>
<td>
<code>resources</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#resourcerequirements-v1-core">
Kubernetes core/v1.ResourceRequirements
</a>
</em>
</td>
<td>
<p>Resources defines the resource requirements for single Pods.</p>
</td>
</tr>
<tr>
<td>
<code>logLevel</code><br/>
<em>
string
</em>
</td>
<td>
<p>Log level for component to be configured with.</p>
</td>
</tr>
<tr>
<td>
<code>logFormat</code><br/>
<em>
string
</em>
</td>
<td>
<p>Log format for component to be configured with.</p>
</td>
</tr>
<tr>
<td>
<code>flags</code><br/>
<em>
[]string
</em>
</td>
<td>
<p>Flags allows setting additional flags for the component container.</p>
</td>
</tr>
<tr>
<td>
<code>podMetadata</code><br/>
<em>
<a href="#monitoring.whizard.io/v1alpha1.EmbeddedObjectMetadata">
EmbeddedObjectMetadata
</a>
</em>
</td>
<td>
<p>PodMetadata configures labels and annotations which are propagated to the pods.</p>
</td>
</tr>
<tr>
<td>
<code>configMaps</code><br/>
<em>
[]string
</em>
</td>
<td>
<p>ConfigMaps is a list of ConfigMaps in the same namespace as the component
object, which shall be mounted into the default Pods.
Each ConfigMap is added to the StatefulSet/Deployment definition as a volume named <code>configmap-&lt;configmap-name&gt;</code>.
The ConfigMaps are mounted into /etc/whizard/configmaps/<configmap-name> in the default container.</p>
</td>
</tr>
<tr>
<td>
<code>secrets</code><br/>
<em>
[]string
</em>
</td>
<td>
<p>Secrets is a list of Secrets in the same namespace as the component
object, which shall be mounted into the Prometheus Pods.
Each Secret is added to the StatefulSet/Deployment definition as a volume named <code>secret-&lt;secret-name&gt;</code>.
The Secrets are mounted into /etc/whizard/secrets/<secret-name> in the default container.</p>
</td>
</tr>
<tr>
<td>
<code>containers</code><br/>
<em>
<a href="https://pkg.go.dev/k8s.io/apimachinery/pkg/runtime.#RawExtension">
k8s.io/apimachinery/pkg/runtime.RawExtension
</a>
</em>
</td>
<td>
<p>Containers allows injecting additional containers or modifying operator generated containers.
Containers described here modify an operator generated
container if they share the same name and modifications are done via a
strategic merge patch.</p>
</td>
</tr>
<tr>
<td>
<code>imagePullSecrets</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#localobjectreference-v1-core">
[]Kubernetes core/v1.LocalObjectReference
</a>
</em>
</td>
<td>
<p>An optional list of references to secrets in the same namespace
to use for pulling images from registries</p>
</td>
</tr>
<tr>
<td>
<code>securityContext</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#podsecuritycontext-v1-core">
Kubernetes core/v1.PodSecurityContext
</a>
</em>
</td>
<td>
<p>SecurityContext holds pod-level security attributes and common container settings.
This defaults to the default PodSecurityContext.</p>
</td>
</tr>
<tr>
<td>
<code>affinity</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#affinity-v1-core">
Kubernetes core/v1.Affinity
</a>
</em>
</td>
<td>
<p>If specified, the pod&rsquo;s scheduling constraints.</p>
</td>
</tr>
<tr>
<td>
<code>nodeSelector</code><br/>
<em>
map[string]string
</em>
</td>
<td>
<p>Define which Nodes the Pods are scheduled on.</p>
</td>
</tr>
<tr>
<td>
<code>tolerations</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#toleration-v1-core">
[]Kubernetes core/v1.Toleration
</a>
</em>
</td>
<td>
<p>If specified, the pod&rsquo;s tolerations.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="monitoring.whizard.io/v1alpha1.RouterStatus">RouterStatus
</h3>
<p>
(<em>Appears on:</em><a href="#monitoring.whizard.io/v1alpha1.Router">Router</a>)
</p>
<div>
<p>RouterStatus defines the observed state of Router</p>
</div>
<h3 id="monitoring.whizard.io/v1alpha1.RulerSpec">RulerSpec
</h3>
<p>
(<em>Appears on:</em><a href="#monitoring.whizard.io/v1alpha1.Ruler">Ruler</a>, <a href="#monitoring.whizard.io/v1alpha1.RulerTemplateSpec">RulerTemplateSpec</a>)
</p>
<div>
<p>RulerSpec defines the desired state of Ruler</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>ruleSelectors</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#*k8s.io/apimachinery/pkg/apis/meta/v1.labelselector--">
[]*k8s.io/apimachinery/pkg/apis/meta/v1.LabelSelector
</a>
</em>
</td>
<td>
<p>Label selectors to select which PrometheusRules to mount for alerting and recording.
The result of multiple selectors are ORed.</p>
</td>
</tr>
<tr>
<td>
<code>ruleNamespaceSelector</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#labelselector-v1-meta">
Kubernetes meta/v1.LabelSelector
</a>
</em>
</td>
<td>
<p>Namespaces to be selected for PrometheusRules discovery. If unspecified, only
the same namespace as the Ruler object is in is used.</p>
</td>
</tr>
<tr>
<td>
<code>shards</code><br/>
<em>
int32
</em>
</td>
<td>
<p>Number of shards to take the hash of fully qualified name of the rule group in order to split rules.
Each shard of rules will be bound to one separate statefulset.
Default: 1</p>
</td>
</tr>
<tr>
<td>
<code>tenant</code><br/>
<em>
string
</em>
</td>
<td>
<p>Tenant if not empty indicates which tenant&rsquo;s data is evaluated for the selected rules;
otherwise, it is for all tenants.</p>
</td>
</tr>
<tr>
<td>
<code>queryConfig</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#secretkeyselector-v1-core">
Kubernetes core/v1.SecretKeySelector
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>remoteWriteConfig</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#secretkeyselector-v1-core">
Kubernetes core/v1.SecretKeySelector
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>labels</code><br/>
<em>
map[string]string
</em>
</td>
<td>
<p>Labels configure the external label pairs to Ruler. A default replica label
<code>ruler_replica</code> will be always added  as a label with the value of the pod&rsquo;s name and it will be dropped in the alerts.</p>
</td>
</tr>
<tr>
<td>
<code>alertDropLabels</code><br/>
<em>
[]string
</em>
</td>
<td>
<p>AlertDropLabels configure the label names which should be dropped in Ruler alerts.
The replica label <code>ruler_replica</code> will always be dropped in alerts.</p>
</td>
</tr>
<tr>
<td>
<code>alertmanagersUrl</code><br/>
<em>
[]string
</em>
</td>
<td>
<p>Define URLs to send alerts to Alertmanager.
Note: this field will be ignored if AlertmanagersConfig is specified.
Maps to the <code>alertmanagers.url</code> arg.</p>
</td>
</tr>
<tr>
<td>
<code>alertmanagersConfig</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#secretkeyselector-v1-core">
Kubernetes core/v1.SecretKeySelector
</a>
</em>
</td>
<td>
<p>Define configuration for connecting to alertmanager. Maps to the <code>alertmanagers.config</code> arg.</p>
</td>
</tr>
<tr>
<td>
<code>evaluationInterval</code><br/>
<em>
<a href="#monitoring.whizard.io/v1alpha1.Duration">
Duration
</a>
</em>
</td>
<td>
<p>Interval between consecutive evaluations.</p>
<p>Default: &ldquo;1m&rdquo;</p>
</td>
</tr>
<tr>
<td>
<code>dataVolume</code><br/>
<em>
<a href="#monitoring.whizard.io/v1alpha1.KubernetesVolume">
KubernetesVolume
</a>
</em>
</td>
<td>
<p>DataVolume specifies how volume shall be used</p>
</td>
</tr>
<tr>
<td>
<code>rulerQueryProxy</code><br/>
<em>
<a href="#monitoring.whizard.io/v1alpha1.SidecarSpec">
SidecarSpec
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>rulerWriteProxy</code><br/>
<em>
<a href="#monitoring.whizard.io/v1alpha1.SidecarSpec">
SidecarSpec
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>prometheusConfigReloader</code><br/>
<em>
<a href="#monitoring.whizard.io/v1alpha1.SidecarSpec">
SidecarSpec
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>replicas</code><br/>
<em>
int32
</em>
</td>
<td>
<p>Number of component instances to deploy.</p>
</td>
</tr>
<tr>
<td>
<code>image</code><br/>
<em>
string
</em>
</td>
<td>
<p>Component container image URL.</p>
</td>
</tr>
<tr>
<td>
<code>imagePullPolicy</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#pullpolicy-v1-core">
Kubernetes core/v1.PullPolicy
</a>
</em>
</td>
<td>
<p>Image pull policy.</p>
</td>
</tr>
<tr>
<td>
<code>resources</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#resourcerequirements-v1-core">
Kubernetes core/v1.ResourceRequirements
</a>
</em>
</td>
<td>
<p>Resources defines the resource requirements for single Pods.</p>
</td>
</tr>
<tr>
<td>
<code>logLevel</code><br/>
<em>
string
</em>
</td>
<td>
<p>Log level for component to be configured with.</p>
</td>
</tr>
<tr>
<td>
<code>logFormat</code><br/>
<em>
string
</em>
</td>
<td>
<p>Log format for component to be configured with.</p>
</td>
</tr>
<tr>
<td>
<code>flags</code><br/>
<em>
[]string
</em>
</td>
<td>
<p>Flags allows setting additional flags for the component container.</p>
</td>
</tr>
<tr>
<td>
<code>podMetadata</code><br/>
<em>
<a href="#monitoring.whizard.io/v1alpha1.EmbeddedObjectMetadata">
EmbeddedObjectMetadata
</a>
</em>
</td>
<td>
<p>PodMetadata configures labels and annotations which are propagated to the pods.</p>
</td>
</tr>
<tr>
<td>
<code>configMaps</code><br/>
<em>
[]string
</em>
</td>
<td>
<p>ConfigMaps is a list of ConfigMaps in the same namespace as the component
object, which shall be mounted into the default Pods.
Each ConfigMap is added to the StatefulSet/Deployment definition as a volume named <code>configmap-&lt;configmap-name&gt;</code>.
The ConfigMaps are mounted into /etc/whizard/configmaps/<configmap-name> in the default container.</p>
</td>
</tr>
<tr>
<td>
<code>secrets</code><br/>
<em>
[]string
</em>
</td>
<td>
<p>Secrets is a list of Secrets in the same namespace as the component
object, which shall be mounted into the Prometheus Pods.
Each Secret is added to the StatefulSet/Deployment definition as a volume named <code>secret-&lt;secret-name&gt;</code>.
The Secrets are mounted into /etc/whizard/secrets/<secret-name> in the default container.</p>
</td>
</tr>
<tr>
<td>
<code>containers</code><br/>
<em>
<a href="https://pkg.go.dev/k8s.io/apimachinery/pkg/runtime.#RawExtension">
k8s.io/apimachinery/pkg/runtime.RawExtension
</a>
</em>
</td>
<td>
<p>Containers allows injecting additional containers or modifying operator generated containers.
Containers described here modify an operator generated
container if they share the same name and modifications are done via a
strategic merge patch.</p>
</td>
</tr>
<tr>
<td>
<code>imagePullSecrets</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#localobjectreference-v1-core">
[]Kubernetes core/v1.LocalObjectReference
</a>
</em>
</td>
<td>
<p>An optional list of references to secrets in the same namespace
to use for pulling images from registries</p>
</td>
</tr>
<tr>
<td>
<code>securityContext</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#podsecuritycontext-v1-core">
Kubernetes core/v1.PodSecurityContext
</a>
</em>
</td>
<td>
<p>SecurityContext holds pod-level security attributes and common container settings.
This defaults to the default PodSecurityContext.</p>
</td>
</tr>
<tr>
<td>
<code>affinity</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#affinity-v1-core">
Kubernetes core/v1.Affinity
</a>
</em>
</td>
<td>
<p>If specified, the pod&rsquo;s scheduling constraints.</p>
</td>
</tr>
<tr>
<td>
<code>nodeSelector</code><br/>
<em>
map[string]string
</em>
</td>
<td>
<p>Define which Nodes the Pods are scheduled on.</p>
</td>
</tr>
<tr>
<td>
<code>tolerations</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#toleration-v1-core">
[]Kubernetes core/v1.Toleration
</a>
</em>
</td>
<td>
<p>If specified, the pod&rsquo;s tolerations.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="monitoring.whizard.io/v1alpha1.RulerStatus">RulerStatus
</h3>
<p>
(<em>Appears on:</em><a href="#monitoring.whizard.io/v1alpha1.Ruler">Ruler</a>)
</p>
<div>
<p>RulerStatus defines the observed state of Ruler</p>
</div>
<h3 id="monitoring.whizard.io/v1alpha1.RulerTemplateSpec">RulerTemplateSpec
</h3>
<p>
(<em>Appears on:</em><a href="#monitoring.whizard.io/v1alpha1.ServiceSpec">ServiceSpec</a>)
</p>
<div>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>ruleSelectors</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#*k8s.io/apimachinery/pkg/apis/meta/v1.labelselector--">
[]*k8s.io/apimachinery/pkg/apis/meta/v1.LabelSelector
</a>
</em>
</td>
<td>
<p>Label selectors to select which PrometheusRules to mount for alerting and recording.
The result of multiple selectors are ORed.</p>
</td>
</tr>
<tr>
<td>
<code>ruleNamespaceSelector</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#labelselector-v1-meta">
Kubernetes meta/v1.LabelSelector
</a>
</em>
</td>
<td>
<p>Namespaces to be selected for PrometheusRules discovery. If unspecified, only
the same namespace as the Ruler object is in is used.</p>
</td>
</tr>
<tr>
<td>
<code>shards</code><br/>
<em>
int32
</em>
</td>
<td>
<p>Number of shards to take the hash of fully qualified name of the rule group in order to split rules.
Each shard of rules will be bound to one separate statefulset.
Default: 1</p>
</td>
</tr>
<tr>
<td>
<code>tenant</code><br/>
<em>
string
</em>
</td>
<td>
<p>Tenant if not empty indicates which tenant&rsquo;s data is evaluated for the selected rules;
otherwise, it is for all tenants.</p>
</td>
</tr>
<tr>
<td>
<code>queryConfig</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#secretkeyselector-v1-core">
Kubernetes core/v1.SecretKeySelector
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>remoteWriteConfig</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#secretkeyselector-v1-core">
Kubernetes core/v1.SecretKeySelector
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>labels</code><br/>
<em>
map[string]string
</em>
</td>
<td>
<p>Labels configure the external label pairs to Ruler. A default replica label
<code>ruler_replica</code> will be always added  as a label with the value of the pod&rsquo;s name and it will be dropped in the alerts.</p>
</td>
</tr>
<tr>
<td>
<code>alertDropLabels</code><br/>
<em>
[]string
</em>
</td>
<td>
<p>AlertDropLabels configure the label names which should be dropped in Ruler alerts.
The replica label <code>ruler_replica</code> will always be dropped in alerts.</p>
</td>
</tr>
<tr>
<td>
<code>alertmanagersUrl</code><br/>
<em>
[]string
</em>
</td>
<td>
<p>Define URLs to send alerts to Alertmanager.
Note: this field will be ignored if AlertmanagersConfig is specified.
Maps to the <code>alertmanagers.url</code> arg.</p>
</td>
</tr>
<tr>
<td>
<code>alertmanagersConfig</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#secretkeyselector-v1-core">
Kubernetes core/v1.SecretKeySelector
</a>
</em>
</td>
<td>
<p>Define configuration for connecting to alertmanager. Maps to the <code>alertmanagers.config</code> arg.</p>
</td>
</tr>
<tr>
<td>
<code>evaluationInterval</code><br/>
<em>
<a href="#monitoring.whizard.io/v1alpha1.Duration">
Duration
</a>
</em>
</td>
<td>
<p>Interval between consecutive evaluations.</p>
<p>Default: &ldquo;1m&rdquo;</p>
</td>
</tr>
<tr>
<td>
<code>dataVolume</code><br/>
<em>
<a href="#monitoring.whizard.io/v1alpha1.KubernetesVolume">
KubernetesVolume
</a>
</em>
</td>
<td>
<p>DataVolume specifies how volume shall be used</p>
</td>
</tr>
<tr>
<td>
<code>rulerQueryProxy</code><br/>
<em>
<a href="#monitoring.whizard.io/v1alpha1.SidecarSpec">
SidecarSpec
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>rulerWriteProxy</code><br/>
<em>
<a href="#monitoring.whizard.io/v1alpha1.SidecarSpec">
SidecarSpec
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>prometheusConfigReloader</code><br/>
<em>
<a href="#monitoring.whizard.io/v1alpha1.SidecarSpec">
SidecarSpec
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>replicas</code><br/>
<em>
int32
</em>
</td>
<td>
<p>Number of component instances to deploy.</p>
</td>
</tr>
<tr>
<td>
<code>image</code><br/>
<em>
string
</em>
</td>
<td>
<p>Component container image URL.</p>
</td>
</tr>
<tr>
<td>
<code>imagePullPolicy</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#pullpolicy-v1-core">
Kubernetes core/v1.PullPolicy
</a>
</em>
</td>
<td>
<p>Image pull policy.</p>
</td>
</tr>
<tr>
<td>
<code>resources</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#resourcerequirements-v1-core">
Kubernetes core/v1.ResourceRequirements
</a>
</em>
</td>
<td>
<p>Resources defines the resource requirements for single Pods.</p>
</td>
</tr>
<tr>
<td>
<code>logLevel</code><br/>
<em>
string
</em>
</td>
<td>
<p>Log level for component to be configured with.</p>
</td>
</tr>
<tr>
<td>
<code>logFormat</code><br/>
<em>
string
</em>
</td>
<td>
<p>Log format for component to be configured with.</p>
</td>
</tr>
<tr>
<td>
<code>flags</code><br/>
<em>
[]string
</em>
</td>
<td>
<p>Flags allows setting additional flags for the component container.</p>
</td>
</tr>
<tr>
<td>
<code>podMetadata</code><br/>
<em>
<a href="#monitoring.whizard.io/v1alpha1.EmbeddedObjectMetadata">
EmbeddedObjectMetadata
</a>
</em>
</td>
<td>
<p>PodMetadata configures labels and annotations which are propagated to the pods.</p>
</td>
</tr>
<tr>
<td>
<code>configMaps</code><br/>
<em>
[]string
</em>
</td>
<td>
<p>ConfigMaps is a list of ConfigMaps in the same namespace as the component
object, which shall be mounted into the default Pods.
Each ConfigMap is added to the StatefulSet/Deployment definition as a volume named <code>configmap-&lt;configmap-name&gt;</code>.
The ConfigMaps are mounted into /etc/whizard/configmaps/<configmap-name> in the default container.</p>
</td>
</tr>
<tr>
<td>
<code>secrets</code><br/>
<em>
[]string
</em>
</td>
<td>
<p>Secrets is a list of Secrets in the same namespace as the component
object, which shall be mounted into the Prometheus Pods.
Each Secret is added to the StatefulSet/Deployment definition as a volume named <code>secret-&lt;secret-name&gt;</code>.
The Secrets are mounted into /etc/whizard/secrets/<secret-name> in the default container.</p>
</td>
</tr>
<tr>
<td>
<code>containers</code><br/>
<em>
<a href="https://pkg.go.dev/k8s.io/apimachinery/pkg/runtime.#RawExtension">
k8s.io/apimachinery/pkg/runtime.RawExtension
</a>
</em>
</td>
<td>
<p>Containers allows injecting additional containers or modifying operator generated containers.
Containers described here modify an operator generated
container if they share the same name and modifications are done via a
strategic merge patch.</p>
</td>
</tr>
<tr>
<td>
<code>imagePullSecrets</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#localobjectreference-v1-core">
[]Kubernetes core/v1.LocalObjectReference
</a>
</em>
</td>
<td>
<p>An optional list of references to secrets in the same namespace
to use for pulling images from registries</p>
</td>
</tr>
<tr>
<td>
<code>securityContext</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#podsecuritycontext-v1-core">
Kubernetes core/v1.PodSecurityContext
</a>
</em>
</td>
<td>
<p>SecurityContext holds pod-level security attributes and common container settings.
This defaults to the default PodSecurityContext.</p>
</td>
</tr>
<tr>
<td>
<code>affinity</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#affinity-v1-core">
Kubernetes core/v1.Affinity
</a>
</em>
</td>
<td>
<p>If specified, the pod&rsquo;s scheduling constraints.</p>
</td>
</tr>
<tr>
<td>
<code>nodeSelector</code><br/>
<em>
map[string]string
</em>
</td>
<td>
<p>Define which Nodes the Pods are scheduled on.</p>
</td>
</tr>
<tr>
<td>
<code>tolerations</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#toleration-v1-core">
[]Kubernetes core/v1.Toleration
</a>
</em>
</td>
<td>
<p>If specified, the pod&rsquo;s tolerations.</p>
</td>
</tr>
<tr>
<td>
<code>disableAlertingRulesAutoSelection</code><br/>
<em>
bool
</em>
</td>
<td>
<p>DisableAlertingRulesAutoSelection disable auto select alerting rules in tenant ruler</p>
<p>Default: true</p>
</td>
</tr>
</tbody>
</table>
<h3 id="monitoring.whizard.io/v1alpha1.S3">S3
</h3>
<p>
(<em>Appears on:</em><a href="#monitoring.whizard.io/v1alpha1.StorageSpec">StorageSpec</a>)
</p>
<div>
<p>Config stores the configuration for s3 bucket.
<a href="https://github.com/thanos-io/objstore/blob/main/providers/s3">https://github.com/thanos-io/objstore/blob/main/providers/s3</a></p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>bucket</code><br/>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>endpoint</code><br/>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>region</code><br/>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>disableDualstack</code><br/>
<em>
bool
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>awsSdkAuth</code><br/>
<em>
bool
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>accessKey</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#secretkeyselector-v1-core">
Kubernetes core/v1.SecretKeySelector
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>insecure</code><br/>
<em>
bool
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>signatureVersion2</code><br/>
<em>
bool
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>secretKey</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#secretkeyselector-v1-core">
Kubernetes core/v1.SecretKeySelector
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>putUserMetadata</code><br/>
<em>
map[string]string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>httpConfig</code><br/>
<em>
<a href="#monitoring.whizard.io/v1alpha1.S3HTTPConfig">
S3HTTPConfig
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>trace</code><br/>
<em>
<a href="#monitoring.whizard.io/v1alpha1.S3TraceConfig">
S3TraceConfig
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>listObjectsVersion</code><br/>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>sendContentMd5</code><br/>
<em>
bool
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>disableMultipart</code><br/>
<em>
bool
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>partSize</code><br/>
<em>
uint64
</em>
</td>
<td>
<p>PartSize used for multipart upload. Only used if uploaded object size is known and larger than configured PartSize.
NOTE we need to make sure this number does not produce more parts than 10 000.</p>
</td>
</tr>
<tr>
<td>
<code>sseConfig</code><br/>
<em>
<a href="#monitoring.whizard.io/v1alpha1.S3SSEConfig">
S3SSEConfig
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>stsEndpoint</code><br/>
<em>
string
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="monitoring.whizard.io/v1alpha1.S3HTTPConfig">S3HTTPConfig
</h3>
<p>
(<em>Appears on:</em><a href="#monitoring.whizard.io/v1alpha1.S3">S3</a>)
</p>
<div>
<p>S3HTTPConfig stores the http.Transport configuration for the s3 minio client.</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>idleConnTimeout</code><br/>
<em>
github.com/prometheus/common/model.Duration
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>responseHeaderTimeout</code><br/>
<em>
github.com/prometheus/common/model.Duration
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>insecureSkipVerify</code><br/>
<em>
bool
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>tlsHandshakeTimeout</code><br/>
<em>
github.com/prometheus/common/model.Duration
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>expectContinueTimeout</code><br/>
<em>
github.com/prometheus/common/model.Duration
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>maxIdleConns</code><br/>
<em>
int
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>maxIdleConnsPerHost</code><br/>
<em>
int
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>maxConnsPerHost</code><br/>
<em>
int
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>tlsConfig</code><br/>
<em>
<a href="#monitoring.whizard.io/v1alpha1.TLSConfig">
TLSConfig
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>disableCompression</code><br/>
<em>
bool
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="monitoring.whizard.io/v1alpha1.S3SSEConfig">S3SSEConfig
</h3>
<p>
(<em>Appears on:</em><a href="#monitoring.whizard.io/v1alpha1.S3">S3</a>)
</p>
<div>
<p>S3SSEConfig deals with the configuration of SSE for Minio. The following options are valid:
kmsencryptioncontext == <a href="https://docs.aws.amazon.com/kms/latest/developerguide/services-s3.html#s3-encryption-context">https://docs.aws.amazon.com/kms/latest/developerguide/services-s3.html#s3-encryption-context</a></p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>type</code><br/>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>kmsKeyId</code><br/>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>kmsEncryptionContext</code><br/>
<em>
map[string]string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>encryptionKey</code><br/>
<em>
string
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="monitoring.whizard.io/v1alpha1.S3TraceConfig">S3TraceConfig
</h3>
<p>
(<em>Appears on:</em><a href="#monitoring.whizard.io/v1alpha1.S3">S3</a>)
</p>
<div>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>enable</code><br/>
<em>
bool
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="monitoring.whizard.io/v1alpha1.ServiceSpec">ServiceSpec
</h3>
<p>
(<em>Appears on:</em><a href="#monitoring.whizard.io/v1alpha1.Service">Service</a>)
</p>
<div>
<p>ServiceSpec defines the desired state of Service</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>tenantHeader</code><br/>
<em>
string
</em>
</td>
<td>
<p>HTTP header to determine tenant for remote write requests.</p>
</td>
</tr>
<tr>
<td>
<code>defaultTenantId</code><br/>
<em>
string
</em>
</td>
<td>
<p>Default tenant ID to use when none is provided via a header.</p>
</td>
</tr>
<tr>
<td>
<code>tenantLabelName</code><br/>
<em>
string
</em>
</td>
<td>
<p>Label name through which the tenant will be announced.</p>
</td>
</tr>
<tr>
<td>
<code>storage</code><br/>
<em>
<a href="#monitoring.whizard.io/v1alpha1.ObjectReference">
ObjectReference
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>remoteWrites</code><br/>
<em>
<a href="#monitoring.whizard.io/v1alpha1.RemoteWriteSpec">
[]RemoteWriteSpec
</a>
</em>
</td>
<td>
<p>RemoteWrites is the list of remote write configurations.
If it is configured, its targets will receive write requests from the Gateway and the Ruler.</p>
</td>
</tr>
<tr>
<td>
<code>remoteQuery</code><br/>
<em>
<a href="#monitoring.whizard.io/v1alpha1.RemoteQuerySpec">
RemoteQuerySpec
</a>
</em>
</td>
<td>
<p>RemoteQuery is the remote query configuration and the remote target should have prometheus-compatible Query APIs.
If not configured, the Gateway will proxy all read requests through the QueryFrontend to the Query,
If configured, the Gateway will proxy metrics read requests through the QueryFrontend to the remote target,
but proxy rules read requests directly to the Query.</p>
</td>
</tr>
<tr>
<td>
<code>gatewayTemplateSpec</code><br/>
<em>
<a href="#monitoring.whizard.io/v1alpha1.GatewaySpec">
GatewaySpec
</a>
</em>
</td>
<td>
<p>GatewayTemplateSpec defines the Gateway configuration template.</p>
</td>
</tr>
<tr>
<td>
<code>queryFrontendTemplateSpec</code><br/>
<em>
<a href="#monitoring.whizard.io/v1alpha1.QueryFrontendSpec">
QueryFrontendSpec
</a>
</em>
</td>
<td>
<p>QueryFrontendTemplateSpec defines the QueryFrontend configuration template.</p>
</td>
</tr>
<tr>
<td>
<code>queryTemplateSpec</code><br/>
<em>
<a href="#monitoring.whizard.io/v1alpha1.QuerySpec">
QuerySpec
</a>
</em>
</td>
<td>
<p>QueryTemplateSpec defines the Query configuration template.</p>
</td>
</tr>
<tr>
<td>
<code>rulerTemplateSpec</code><br/>
<em>
<a href="#monitoring.whizard.io/v1alpha1.RulerTemplateSpec">
RulerTemplateSpec
</a>
</em>
</td>
<td>
<p>RulerTemplateSpec defines the Ruler configuration template.</p>
</td>
</tr>
<tr>
<td>
<code>routerTemplateSpec</code><br/>
<em>
<a href="#monitoring.whizard.io/v1alpha1.RouterSpec">
RouterSpec
</a>
</em>
</td>
<td>
<p>RouterTemplateSpec defines the Router configuration template.</p>
</td>
</tr>
<tr>
<td>
<code>ingesterTemplateSpec</code><br/>
<em>
<a href="#monitoring.whizard.io/v1alpha1.IngesterTemplateSpec">
IngesterTemplateSpec
</a>
</em>
</td>
<td>
<p>IngesterTemplateSpec defines the Ingester configuration template.</p>
</td>
</tr>
<tr>
<td>
<code>storeTemplateSpec</code><br/>
<em>
<a href="#monitoring.whizard.io/v1alpha1.StoreSpec">
StoreSpec
</a>
</em>
</td>
<td>
<p>StoreTemplateSpec defines the Store configuration template.</p>
</td>
</tr>
<tr>
<td>
<code>compactorTemplateSpec</code><br/>
<em>
<a href="#monitoring.whizard.io/v1alpha1.CompactorTemplateSpec">
CompactorTemplateSpec
</a>
</em>
</td>
<td>
<p>CompactorTemplateSpec defines the Compactor configuration template.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="monitoring.whizard.io/v1alpha1.ServiceStatus">ServiceStatus
</h3>
<p>
(<em>Appears on:</em><a href="#monitoring.whizard.io/v1alpha1.Service">Service</a>)
</p>
<div>
<p>ServiceStatus defines the observed state of Service</p>
</div>
<h3 id="monitoring.whizard.io/v1alpha1.SidecarSpec">SidecarSpec
</h3>
<p>
(<em>Appears on:</em><a href="#monitoring.whizard.io/v1alpha1.IngesterSpec">IngesterSpec</a>, <a href="#monitoring.whizard.io/v1alpha1.QuerySpec">QuerySpec</a>, <a href="#monitoring.whizard.io/v1alpha1.RulerSpec">RulerSpec</a>)
</p>
<div>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>image</code><br/>
<em>
string
</em>
</td>
<td>
<p>Image is the envoy image with tag/version</p>
</td>
</tr>
<tr>
<td>
<code>resources</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#resourcerequirements-v1-core">
Kubernetes core/v1.ResourceRequirements
</a>
</em>
</td>
<td>
<p>Define resources requests and limits for sidecar container.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="monitoring.whizard.io/v1alpha1.StorageSpec">StorageSpec
</h3>
<p>
(<em>Appears on:</em><a href="#monitoring.whizard.io/v1alpha1.Storage">Storage</a>)
</p>
<div>
<p>StorageSpec defines the desired state of Storage</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>blockManager</code><br/>
<em>
<a href="#monitoring.whizard.io/v1alpha1.BlockManager">
BlockManager
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>S3</code><br/>
<em>
<a href="#monitoring.whizard.io/v1alpha1.S3">
S3
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="monitoring.whizard.io/v1alpha1.StorageStatus">StorageStatus
</h3>
<p>
(<em>Appears on:</em><a href="#monitoring.whizard.io/v1alpha1.Storage">Storage</a>)
</p>
<div>
<p>StorageStatus defines the observed state of Storage</p>
</div>
<h3 id="monitoring.whizard.io/v1alpha1.StoreSpec">StoreSpec
</h3>
<p>
(<em>Appears on:</em><a href="#monitoring.whizard.io/v1alpha1.Store">Store</a>, <a href="#monitoring.whizard.io/v1alpha1.ServiceSpec">ServiceSpec</a>)
</p>
<div>
<p>StoreSpec defines the desired state of Store</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>minTime</code><br/>
<em>
string
</em>
</td>
<td>
<p>MinTime specifies start of time range limit to serve</p>
</td>
</tr>
<tr>
<td>
<code>maxTime</code><br/>
<em>
string
</em>
</td>
<td>
<p>MaxTime specifies end of time range limit to serve</p>
</td>
</tr>
<tr>
<td>
<code>timeRanges</code><br/>
<em>
<a href="#monitoring.whizard.io/v1alpha1.TimeRange">
[]TimeRange
</a>
</em>
</td>
<td>
<p>TimeRanges is a list of TimeRange to partition Store.
If specified, the MinTime and MaxTime will be ignored.</p>
</td>
</tr>
<tr>
<td>
<code>indexCacheConfig</code><br/>
<em>
<a href="#monitoring.whizard.io/v1alpha1.IndexCacheConfig">
IndexCacheConfig
</a>
</em>
</td>
<td>
<p>IndexCacheConfig contains index cache configuration.</p>
</td>
</tr>
<tr>
<td>
<code>dataVolume</code><br/>
<em>
<a href="#monitoring.whizard.io/v1alpha1.KubernetesVolume">
KubernetesVolume
</a>
</em>
</td>
<td>
<p>DataVolume specifies how volume shall be used</p>
</td>
</tr>
<tr>
<td>
<code>replicas</code><br/>
<em>
int32
</em>
</td>
<td>
<p>Number of component instances to deploy.</p>
</td>
</tr>
<tr>
<td>
<code>image</code><br/>
<em>
string
</em>
</td>
<td>
<p>Component container image URL.</p>
</td>
</tr>
<tr>
<td>
<code>imagePullPolicy</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#pullpolicy-v1-core">
Kubernetes core/v1.PullPolicy
</a>
</em>
</td>
<td>
<p>Image pull policy.</p>
</td>
</tr>
<tr>
<td>
<code>resources</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#resourcerequirements-v1-core">
Kubernetes core/v1.ResourceRequirements
</a>
</em>
</td>
<td>
<p>Resources defines the resource requirements for single Pods.</p>
</td>
</tr>
<tr>
<td>
<code>logLevel</code><br/>
<em>
string
</em>
</td>
<td>
<p>Log level for component to be configured with.</p>
</td>
</tr>
<tr>
<td>
<code>logFormat</code><br/>
<em>
string
</em>
</td>
<td>
<p>Log format for component to be configured with.</p>
</td>
</tr>
<tr>
<td>
<code>flags</code><br/>
<em>
[]string
</em>
</td>
<td>
<p>Flags allows setting additional flags for the component container.</p>
</td>
</tr>
<tr>
<td>
<code>podMetadata</code><br/>
<em>
<a href="#monitoring.whizard.io/v1alpha1.EmbeddedObjectMetadata">
EmbeddedObjectMetadata
</a>
</em>
</td>
<td>
<p>PodMetadata configures labels and annotations which are propagated to the pods.</p>
</td>
</tr>
<tr>
<td>
<code>configMaps</code><br/>
<em>
[]string
</em>
</td>
<td>
<p>ConfigMaps is a list of ConfigMaps in the same namespace as the component
object, which shall be mounted into the default Pods.
Each ConfigMap is added to the StatefulSet/Deployment definition as a volume named <code>configmap-&lt;configmap-name&gt;</code>.
The ConfigMaps are mounted into /etc/whizard/configmaps/<configmap-name> in the default container.</p>
</td>
</tr>
<tr>
<td>
<code>secrets</code><br/>
<em>
[]string
</em>
</td>
<td>
<p>Secrets is a list of Secrets in the same namespace as the component
object, which shall be mounted into the Prometheus Pods.
Each Secret is added to the StatefulSet/Deployment definition as a volume named <code>secret-&lt;secret-name&gt;</code>.
The Secrets are mounted into /etc/whizard/secrets/<secret-name> in the default container.</p>
</td>
</tr>
<tr>
<td>
<code>containers</code><br/>
<em>
<a href="https://pkg.go.dev/k8s.io/apimachinery/pkg/runtime.#RawExtension">
k8s.io/apimachinery/pkg/runtime.RawExtension
</a>
</em>
</td>
<td>
<p>Containers allows injecting additional containers or modifying operator generated containers.
Containers described here modify an operator generated
container if they share the same name and modifications are done via a
strategic merge patch.</p>
</td>
</tr>
<tr>
<td>
<code>imagePullSecrets</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#localobjectreference-v1-core">
[]Kubernetes core/v1.LocalObjectReference
</a>
</em>
</td>
<td>
<p>An optional list of references to secrets in the same namespace
to use for pulling images from registries</p>
</td>
</tr>
<tr>
<td>
<code>securityContext</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#podsecuritycontext-v1-core">
Kubernetes core/v1.PodSecurityContext
</a>
</em>
</td>
<td>
<p>SecurityContext holds pod-level security attributes and common container settings.
This defaults to the default PodSecurityContext.</p>
</td>
</tr>
<tr>
<td>
<code>affinity</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#affinity-v1-core">
Kubernetes core/v1.Affinity
</a>
</em>
</td>
<td>
<p>If specified, the pod&rsquo;s scheduling constraints.</p>
</td>
</tr>
<tr>
<td>
<code>nodeSelector</code><br/>
<em>
map[string]string
</em>
</td>
<td>
<p>Define which Nodes the Pods are scheduled on.</p>
</td>
</tr>
<tr>
<td>
<code>tolerations</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#toleration-v1-core">
[]Kubernetes core/v1.Toleration
</a>
</em>
</td>
<td>
<p>If specified, the pod&rsquo;s tolerations.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="monitoring.whizard.io/v1alpha1.StoreStatus">StoreStatus
</h3>
<p>
(<em>Appears on:</em><a href="#monitoring.whizard.io/v1alpha1.Store">Store</a>)
</p>
<div>
<p>StoreStatus defines the observed state of Store</p>
</div>
<h3 id="monitoring.whizard.io/v1alpha1.TLSConfig">TLSConfig
</h3>
<p>
(<em>Appears on:</em><a href="#monitoring.whizard.io/v1alpha1.S3HTTPConfig">S3HTTPConfig</a>)
</p>
<div>
<p>TLSConfig configures the options for TLS connections.</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>ca</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#secretkeyselector-v1-core">
Kubernetes core/v1.SecretKeySelector
</a>
</em>
</td>
<td>
<p>The secret that including the CA cert.</p>
</td>
</tr>
<tr>
<td>
<code>cert</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#secretkeyselector-v1-core">
Kubernetes core/v1.SecretKeySelector
</a>
</em>
</td>
<td>
<p>The secret that including the client cert.</p>
</td>
</tr>
<tr>
<td>
<code>key</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#secretkeyselector-v1-core">
Kubernetes core/v1.SecretKeySelector
</a>
</em>
</td>
<td>
<p>The secret that including the client key.</p>
</td>
</tr>
<tr>
<td>
<code>serverName</code><br/>
<em>
string
</em>
</td>
<td>
<p>Used to verify the hostname for the targets.</p>
</td>
</tr>
<tr>
<td>
<code>insecureSkipVerify</code><br/>
<em>
bool
</em>
</td>
<td>
<p>Disable target certificate validation.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="monitoring.whizard.io/v1alpha1.TenantSpec">TenantSpec
</h3>
<p>
(<em>Appears on:</em><a href="#monitoring.whizard.io/v1alpha1.Tenant">Tenant</a>)
</p>
<div>
<p>TenantSpec defines the desired state of Tenant</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>tenant</code><br/>
<em>
string
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="monitoring.whizard.io/v1alpha1.TenantStatus">TenantStatus
</h3>
<p>
(<em>Appears on:</em><a href="#monitoring.whizard.io/v1alpha1.Tenant">Tenant</a>)
</p>
<div>
<p>TenantStatus defines the observed state of Tenant</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>ruler</code><br/>
<em>
<a href="#monitoring.whizard.io/v1alpha1.ObjectReference">
ObjectReference
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>compactor</code><br/>
<em>
<a href="#monitoring.whizard.io/v1alpha1.ObjectReference">
ObjectReference
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>ingester</code><br/>
<em>
<a href="#monitoring.whizard.io/v1alpha1.ObjectReference">
ObjectReference
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="monitoring.whizard.io/v1alpha1.TimeRange">TimeRange
</h3>
<p>
(<em>Appears on:</em><a href="#monitoring.whizard.io/v1alpha1.StoreSpec">StoreSpec</a>)
</p>
<div>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>minTime</code><br/>
<em>
string
</em>
</td>
<td>
<p>MinTime specifies start of time range limit to serve</p>
</td>
</tr>
<tr>
<td>
<code>maxTime</code><br/>
<em>
string
</em>
</td>
<td>
<p>MaxTime specifies end of time range limit to serve</p>
</td>
</tr>
</tbody>
</table>
<h3 id="monitoring.whizard.io/v1alpha1.WebConfig">WebConfig
</h3>
<p>
(<em>Appears on:</em><a href="#monitoring.whizard.io/v1alpha1.GatewaySpec">GatewaySpec</a>, <a href="#monitoring.whizard.io/v1alpha1.QueryFrontendSpec">QueryFrontendSpec</a>, <a href="#monitoring.whizard.io/v1alpha1.QuerySpec">QuerySpec</a>, <a href="#monitoring.whizard.io/v1alpha1.RouterSpec">RouterSpec</a>)
</p>
<div>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>httpServerTLSConfig</code><br/>
<em>
<a href="#monitoring.whizard.io/v1alpha1.HTTPServerTLSConfig">
HTTPServerTLSConfig
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>httpServerConfig</code><br/>
<em>
<a href="#monitoring.whizard.io/v1alpha1.HTTPServerConfig">
HTTPServerConfig
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>basicAuthUsers</code><br/>
<em>
<a href="#monitoring.whizard.io/v1alpha1.BasicAuth">
[]BasicAuth
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<hr/>

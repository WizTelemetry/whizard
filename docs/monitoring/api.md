# API Reference

## Packages
- [monitoring.whizard.io/v1alpha1](#monitoringwhizardiov1alpha1)


## monitoring.whizard.io/v1alpha1

Package v1alpha1 contains API Schema definitions for the monitoring v1alpha1 API group

Package v1alpha1 contains API Schema definitions for the monitoring v1alpha1 API group

### Resource Types
- [Compactor](#compactor)
- [Gateway](#gateway)
- [Ingester](#ingester)
- [Query](#query)
- [QueryFrontend](#queryfrontend)
- [Router](#router)
- [Ruler](#ruler)
- [Service](#service)
- [Storage](#storage)
- [Store](#store)
- [Tenant](#tenant)



#### AutoScaler







_Appears in:_
- [StoreSpec](#storespec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `minReplicas` _integer_ | minReplicas is the lower limit for the number of replicas to which the autoscaler<br />can scale down.  It defaults to 1 pod.  minReplicas is allowed to be 0 if the<br />alpha feature gate HPAScaleToZero is enabled and at least one Object or External<br />metric is configured.  Scaling is active as long as at least one metric value is<br />available. |  |  |
| `maxReplicas` _integer_ | maxReplicas is the upper limit for the number of replicas to which the autoscaler can scale up.<br />It cannot be less that minReplicas. |  |  |
| `metrics` _[MetricSpec](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26/#metricspec-v2beta2-autoscaling) array_ | metrics contains the specifications for which to use to calculate the<br />desired replica count (the maximum replica count across all metrics will<br />be used).  The desired replica count is calculated multiplying the<br />ratio between the target value and the current value by the current<br />number of pods.  Ergo, metrics used must decrease as the pod count is<br />increased, and vice-versa.  See the individual metric source types for<br />more information about how each type of metric must respond.<br />If not set, the default metric will be set to 80% average CPU utilization. |  |  |
| `behavior` _[HorizontalPodAutoscalerBehavior](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26/#horizontalpodautoscalerbehavior-v2beta2-autoscaling)_ | behavior configures the scaling behavior of the target<br />in both Up and Down directions (scaleUp and scaleDown fields respectively).<br />If not set, the default HPAScalingRules for scale up and scale down are used. |  |  |


#### BasicAuth



BasicAuth allow an endpoint to authenticate over basic authentication



_Appears in:_
- [HTTPClientConfig](#httpclientconfig)
- [RemoteQuerySpec](#remotequeryspec)
- [RemoteWriteSpec](#remotewritespec)
- [WebConfig](#webconfig)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `username` _[SecretKeySelector](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26/#secretkeyselector-v1-core)_ | The secret in the service monitor namespace that contains the username<br />for authentication. |  |  |
| `password` _[SecretKeySelector](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26/#secretkeyselector-v1-core)_ | The secret in the service monitor namespace that contains the password<br />for authentication. |  |  |


#### BlockGC







_Appears in:_
- [BlockManager](#blockmanager)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `enable` _boolean_ |  |  |  |
| `resources` _[ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26/#resourcerequirements-v1-core)_ | Define resources requests and limits for main container. |  |  |
| `image` _string_ | Image is the component image with tag/version. |  |  |
| `imagePullPolicy` _[PullPolicy](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26/#pullpolicy-v1-core)_ | Image pull policy.<br />One of Always, Never, IfNotPresent.<br />Defaults to Always if :latest tag is specified, or IfNotPresent otherwise.<br />Cannot be updated. |  |  |
| `gcInterval` _[Duration](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26/#duration-v1-meta)_ |  |  | Pattern: `^(0|(([0-9]+)y)?(([0-9]+)w)?(([0-9]+)d)?(([0-9]+)h)?(([0-9]+)m)?(([0-9]+)s)?(([0-9]+)ms)?)$` <br /> |
| `cleanupTimeout` _[Duration](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26/#duration-v1-meta)_ |  |  | Pattern: `^(0|(([0-9]+)y)?(([0-9]+)w)?(([0-9]+)d)?(([0-9]+)h)?(([0-9]+)m)?(([0-9]+)s)?(([0-9]+)ms)?)$` <br /> |
| `defaultTenantId` _string_ | Default tenant ID to use when none is provided via a header. |  |  |
| `tenantLabelName` _string_ | Label name through which the tenant will be announced. |  |  |


#### BlockManager







_Appears in:_
- [StorageSpec](#storagespec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `enable` _boolean_ |  |  |  |
| `containers` _[RawExtension](#rawextension)_ | Containers allows injecting additional containers or modifying operator generated containers.<br />Containers described here modify an operator generated<br />container if they share the same name and modifications are done via a<br />strategic merge patch. |  |  |
| `podMetadata` _[EmbeddedObjectMetadata](#embeddedobjectmetadata)_ | PodMetadata configures labels and annotations which are propagated to the pods.<br /><br />* "kubectl.kubernetes.io/default-container" annotation, set to main pod. |  |  |
| `secrets` _string array_ | Secrets is a list of Secrets in the same namespace as the component<br />object, which shall be mounted into the Prometheus Pods.<br />Each Secret is added to the StatefulSet/Deployment definition as a volume named `secret-<secret-name>`.<br />The Secrets are mounted into /etc/whizard/secrets/<secret-name> in the default container. |  |  |
| `configMaps` _string array_ | ConfigMaps is a list of ConfigMaps in the same namespace as the component<br />object, which shall be mounted into the default Pods.<br />Each ConfigMap is added to the StatefulSet/Deployment definition as a volume named `configmap-<configmap-name>`.<br />The ConfigMaps are mounted into /etc/whizard/configmaps/<configmap-name> in the default container. |  |  |
| `affinity` _[Affinity](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26/#affinity-v1-core)_ | If specified, the pod's scheduling constraints. |  |  |
| `nodeSelector` _object (keys:string, values:string)_ | Define which Nodes the Pods are scheduled on. |  |  |
| `tolerations` _[Toleration](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26/#toleration-v1-core) array_ | If specified, the pod's tolerations. |  |  |
| `resources` _[ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26/#resourcerequirements-v1-core)_ | Define resources requests and limits for main container. |  |  |
| `securityContext` _[PodSecurityContext](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26/#podsecuritycontext-v1-core)_ | SecurityContext holds pod-level security attributes and common container settings.<br />This defaults to the default PodSecurityContext. |  |  |
| `replicas` _integer_ | Number of replicas for a component. |  |  |
| `image` _string_ | Image is the component image with tag/version. |  |  |
| `imagePullPolicy` _[PullPolicy](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26/#pullpolicy-v1-core)_ | Image pull policy.<br />One of Always, Never, IfNotPresent.<br />Defaults to Always if :latest tag is specified, or IfNotPresent otherwise.<br />Cannot be updated. |  |  |
| `imagePullSecrets` _[LocalObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26/#localobjectreference-v1-core) array_ | An optional list of references to secrets in the same namespace<br />to use for pulling images from registries |  |  |
| `logLevel` _string_ | Log filtering level. Possible options: error, warn, info, debug. |  |  |
| `logFormat` _string_ | Log format to use. Possible options: logfmt or json. |  |  |
| `flags` _string array_ | Flags is the flags of component. |  |  |
| `serviceAccountName` _string_ | ServiceAccountName is the name of the ServiceAccount to use to run bucket Pods. |  |  |
| `nodePort` _integer_ | NodePort is the port used to expose the bucket service.<br />If this is a valid node port, the gateway service type will be set to NodePort accordingly. |  |  |
| `blockSyncInterval` _[Duration](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26/#duration-v1-meta)_ | Interval to sync block metadata from object storage |  | Pattern: `^(0|(([0-9]+)y)?(([0-9]+)w)?(([0-9]+)d)?(([0-9]+)h)?(([0-9]+)m)?(([0-9]+)s)?(([0-9]+)ms)?)$` <br /> |
| `gc` _[BlockGC](#blockgc)_ |  |  |  |


#### CacheProvider

_Underlying type:_ _string_





_Appears in:_
- [ResponseCacheProviderConfig](#responsecacheproviderconfig)



#### CommonSpec







_Appears in:_
- [BlockManager](#blockmanager)
- [CompactorSpec](#compactorspec)
- [CompactorTemplateSpec](#compactortemplatespec)
- [GatewaySpec](#gatewayspec)
- [IngesterSpec](#ingesterspec)
- [IngesterTemplateSpec](#ingestertemplatespec)
- [QueryFrontendSpec](#queryfrontendspec)
- [QuerySpec](#queryspec)
- [RouterSpec](#routerspec)
- [RulerSpec](#rulerspec)
- [RulerTemplateSpec](#rulertemplatespec)
- [StoreSpec](#storespec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `containers` _[RawExtension](#rawextension)_ | Containers allows injecting additional containers or modifying operator generated containers.<br />Containers described here modify an operator generated<br />container if they share the same name and modifications are done via a<br />strategic merge patch. |  |  |
| `podMetadata` _[EmbeddedObjectMetadata](#embeddedobjectmetadata)_ | PodMetadata configures labels and annotations which are propagated to the pods.<br /><br />* "kubectl.kubernetes.io/default-container" annotation, set to main pod. |  |  |
| `secrets` _string array_ | Secrets is a list of Secrets in the same namespace as the component<br />object, which shall be mounted into the Prometheus Pods.<br />Each Secret is added to the StatefulSet/Deployment definition as a volume named `secret-<secret-name>`.<br />The Secrets are mounted into /etc/whizard/secrets/<secret-name> in the default container. |  |  |
| `configMaps` _string array_ | ConfigMaps is a list of ConfigMaps in the same namespace as the component<br />object, which shall be mounted into the default Pods.<br />Each ConfigMap is added to the StatefulSet/Deployment definition as a volume named `configmap-<configmap-name>`.<br />The ConfigMaps are mounted into /etc/whizard/configmaps/<configmap-name> in the default container. |  |  |
| `affinity` _[Affinity](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26/#affinity-v1-core)_ | If specified, the pod's scheduling constraints. |  |  |
| `nodeSelector` _object (keys:string, values:string)_ | Define which Nodes the Pods are scheduled on. |  |  |
| `tolerations` _[Toleration](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26/#toleration-v1-core) array_ | If specified, the pod's tolerations. |  |  |
| `resources` _[ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26/#resourcerequirements-v1-core)_ | Define resources requests and limits for main container. |  |  |
| `securityContext` _[PodSecurityContext](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26/#podsecuritycontext-v1-core)_ | SecurityContext holds pod-level security attributes and common container settings.<br />This defaults to the default PodSecurityContext. |  |  |
| `replicas` _integer_ | Number of replicas for a component. |  |  |
| `image` _string_ | Image is the component image with tag/version. |  |  |
| `imagePullPolicy` _[PullPolicy](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26/#pullpolicy-v1-core)_ | Image pull policy.<br />One of Always, Never, IfNotPresent.<br />Defaults to Always if :latest tag is specified, or IfNotPresent otherwise.<br />Cannot be updated. |  |  |
| `imagePullSecrets` _[LocalObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26/#localobjectreference-v1-core) array_ | An optional list of references to secrets in the same namespace<br />to use for pulling images from registries |  |  |
| `logLevel` _string_ | Log filtering level. Possible options: error, warn, info, debug. |  |  |
| `logFormat` _string_ | Log format to use. Possible options: logfmt or json. |  |  |
| `flags` _string array_ | Flags is the flags of component. |  |  |


#### Compactor



Compactor is the Schema for the Compactor API





| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `monitoring.whizard.io/v1alpha1` | | |
| `kind` _string_ | `Compactor` | | |
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `spec` _[CompactorSpec](#compactorspec)_ |  |  |  |
| `status` _[CompactorStatus](#compactorstatus)_ |  |  |  |


#### CompactorSpec







_Appears in:_
- [Compactor](#compactor)
- [CompactorTemplateSpec](#compactortemplatespec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `containers` _[RawExtension](#rawextension)_ | Containers allows injecting additional containers or modifying operator generated containers.<br />Containers described here modify an operator generated<br />container if they share the same name and modifications are done via a<br />strategic merge patch. |  |  |
| `podMetadata` _[EmbeddedObjectMetadata](#embeddedobjectmetadata)_ | PodMetadata configures labels and annotations which are propagated to the pods.<br /><br />* "kubectl.kubernetes.io/default-container" annotation, set to main pod. |  |  |
| `secrets` _string array_ | Secrets is a list of Secrets in the same namespace as the component<br />object, which shall be mounted into the Prometheus Pods.<br />Each Secret is added to the StatefulSet/Deployment definition as a volume named `secret-<secret-name>`.<br />The Secrets are mounted into /etc/whizard/secrets/<secret-name> in the default container. |  |  |
| `configMaps` _string array_ | ConfigMaps is a list of ConfigMaps in the same namespace as the component<br />object, which shall be mounted into the default Pods.<br />Each ConfigMap is added to the StatefulSet/Deployment definition as a volume named `configmap-<configmap-name>`.<br />The ConfigMaps are mounted into /etc/whizard/configmaps/<configmap-name> in the default container. |  |  |
| `affinity` _[Affinity](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26/#affinity-v1-core)_ | If specified, the pod's scheduling constraints. |  |  |
| `nodeSelector` _object (keys:string, values:string)_ | Define which Nodes the Pods are scheduled on. |  |  |
| `tolerations` _[Toleration](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26/#toleration-v1-core) array_ | If specified, the pod's tolerations. |  |  |
| `resources` _[ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26/#resourcerequirements-v1-core)_ | Define resources requests and limits for main container. |  |  |
| `securityContext` _[PodSecurityContext](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26/#podsecuritycontext-v1-core)_ | SecurityContext holds pod-level security attributes and common container settings.<br />This defaults to the default PodSecurityContext. |  |  |
| `replicas` _integer_ | Number of replicas for a component. |  |  |
| `image` _string_ | Image is the component image with tag/version. |  |  |
| `imagePullPolicy` _[PullPolicy](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26/#pullpolicy-v1-core)_ | Image pull policy.<br />One of Always, Never, IfNotPresent.<br />Defaults to Always if :latest tag is specified, or IfNotPresent otherwise.<br />Cannot be updated. |  |  |
| `imagePullSecrets` _[LocalObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26/#localobjectreference-v1-core) array_ | An optional list of references to secrets in the same namespace<br />to use for pulling images from registries |  |  |
| `logLevel` _string_ | Log filtering level. Possible options: error, warn, info, debug. |  |  |
| `logFormat` _string_ | Log format to use. Possible options: logfmt or json. |  |  |
| `flags` _string array_ | Flags is the flags of component. |  |  |
| `disableDownsampling` _boolean_ | DisableDownsampling specifies whether to disable downsampling |  |  |
| `retention` _[Retention](#retention)_ | Retention configs how long to retain samples |  |  |
| `dataVolume` _[KubernetesVolume](#kubernetesvolume)_ | DataVolume specifies how volume shall be used |  |  |
| `tenants` _string array_ | Tenants if not empty indicates current config is for hard tenants; otherwise, it is for soft tenants. |  |  |


#### CompactorStatus



CompactorStatus defines the observed state of Compactor



_Appears in:_
- [Compactor](#compactor)



#### CompactorTemplateSpec







_Appears in:_
- [ServiceSpec](#servicespec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `containers` _[RawExtension](#rawextension)_ | Containers allows injecting additional containers or modifying operator generated containers.<br />Containers described here modify an operator generated<br />container if they share the same name and modifications are done via a<br />strategic merge patch. |  |  |
| `podMetadata` _[EmbeddedObjectMetadata](#embeddedobjectmetadata)_ | PodMetadata configures labels and annotations which are propagated to the pods.<br /><br />* "kubectl.kubernetes.io/default-container" annotation, set to main pod. |  |  |
| `secrets` _string array_ | Secrets is a list of Secrets in the same namespace as the component<br />object, which shall be mounted into the Prometheus Pods.<br />Each Secret is added to the StatefulSet/Deployment definition as a volume named `secret-<secret-name>`.<br />The Secrets are mounted into /etc/whizard/secrets/<secret-name> in the default container. |  |  |
| `configMaps` _string array_ | ConfigMaps is a list of ConfigMaps in the same namespace as the component<br />object, which shall be mounted into the default Pods.<br />Each ConfigMap is added to the StatefulSet/Deployment definition as a volume named `configmap-<configmap-name>`.<br />The ConfigMaps are mounted into /etc/whizard/configmaps/<configmap-name> in the default container. |  |  |
| `affinity` _[Affinity](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26/#affinity-v1-core)_ | If specified, the pod's scheduling constraints. |  |  |
| `nodeSelector` _object (keys:string, values:string)_ | Define which Nodes the Pods are scheduled on. |  |  |
| `tolerations` _[Toleration](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26/#toleration-v1-core) array_ | If specified, the pod's tolerations. |  |  |
| `resources` _[ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26/#resourcerequirements-v1-core)_ | Define resources requests and limits for main container. |  |  |
| `securityContext` _[PodSecurityContext](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26/#podsecuritycontext-v1-core)_ | SecurityContext holds pod-level security attributes and common container settings.<br />This defaults to the default PodSecurityContext. |  |  |
| `replicas` _integer_ | Number of replicas for a component. |  |  |
| `image` _string_ | Image is the component image with tag/version. |  |  |
| `imagePullPolicy` _[PullPolicy](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26/#pullpolicy-v1-core)_ | Image pull policy.<br />One of Always, Never, IfNotPresent.<br />Defaults to Always if :latest tag is specified, or IfNotPresent otherwise.<br />Cannot be updated. |  |  |
| `imagePullSecrets` _[LocalObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26/#localobjectreference-v1-core) array_ | An optional list of references to secrets in the same namespace<br />to use for pulling images from registries |  |  |
| `logLevel` _string_ | Log filtering level. Possible options: error, warn, info, debug. |  |  |
| `logFormat` _string_ | Log format to use. Possible options: logfmt or json. |  |  |
| `flags` _string array_ | Flags is the flags of component. |  |  |
| `disableDownsampling` _boolean_ | DisableDownsampling specifies whether to disable downsampling |  |  |
| `retention` _[Retention](#retention)_ | Retention configs how long to retain samples |  |  |
| `dataVolume` _[KubernetesVolume](#kubernetesvolume)_ | DataVolume specifies how volume shall be used |  |  |
| `tenants` _string array_ | Tenants if not empty indicates current config is for hard tenants; otherwise, it is for soft tenants. |  |  |
| `defaultTenantsPerCompactor` _integer_ | DefaultTenantsPerIngester Whizard default tenant count per ingester.<br />Default: 10 | 10 |  |


#### Duration

_Underlying type:_ _string_

Duration is a valid time unit
Supported units: y, w, d, h, m, s, ms Examples: `30s`, `1m`, `1h20m15s`

_Validation:_
- Pattern: `^(0|(([0-9]+)y)?(([0-9]+)w)?(([0-9]+)d)?(([0-9]+)h)?(([0-9]+)m)?(([0-9]+)s)?(([0-9]+)ms)?)$`

_Appears in:_
- [IngesterTemplateSpec](#ingestertemplatespec)
- [RemoteWriteSpec](#remotewritespec)
- [Retention](#retention)
- [RulerSpec](#rulerspec)
- [RulerTemplateSpec](#rulertemplatespec)



#### EmbeddedObjectMetadata



EmbeddedObjectMetadata contains a subset of the fields included in k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta
Only fields which are relevant to embedded resources are included.



_Appears in:_
- [BlockManager](#blockmanager)
- [CommonSpec](#commonspec)
- [CompactorSpec](#compactorspec)
- [CompactorTemplateSpec](#compactortemplatespec)
- [GatewaySpec](#gatewayspec)
- [IngesterSpec](#ingesterspec)
- [IngesterTemplateSpec](#ingestertemplatespec)
- [QueryFrontendSpec](#queryfrontendspec)
- [QuerySpec](#queryspec)
- [RouterSpec](#routerspec)
- [RulerSpec](#rulerspec)
- [RulerTemplateSpec](#rulertemplatespec)
- [StoreSpec](#storespec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `name` _string_ | Name must be unique within a namespace. Is required when creating resources, although<br />some resources may allow a client to request the generation of an appropriate name<br />automatically. Name is primarily intended for creation idempotence and configuration<br />definition.<br />Cannot be updated.<br />More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names#names |  |  |
| `labels` _object (keys:string, values:string)_ | Map of string keys and values that can be used to organize and categorize<br />(scope and select) objects. May match selectors of replication controllers<br />and services.<br />More info: http://kubernetes.io/docs/user-guide/labels |  |  |
| `annotations` _object (keys:string, values:string)_ | Annotations is an unstructured key value map stored with a resource that may be<br />set by external tools to store and retrieve arbitrary metadata. They are not<br />queryable and should be preserved when modifying objects.<br />More info: http://kubernetes.io/docs/user-guide/annotations |  |  |


#### Gateway



Gateway is the Schema for the monitoring gateway API





| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `monitoring.whizard.io/v1alpha1` | | |
| `kind` _string_ | `Gateway` | | |
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `spec` _[GatewaySpec](#gatewayspec)_ |  |  |  |
| `status` _[GatewayStatus](#gatewaystatus)_ |  |  |  |


#### GatewaySpec







_Appears in:_
- [Gateway](#gateway)
- [ServiceSpec](#servicespec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `containers` _[RawExtension](#rawextension)_ | Containers allows injecting additional containers or modifying operator generated containers.<br />Containers described here modify an operator generated<br />container if they share the same name and modifications are done via a<br />strategic merge patch. |  |  |
| `podMetadata` _[EmbeddedObjectMetadata](#embeddedobjectmetadata)_ | PodMetadata configures labels and annotations which are propagated to the pods.<br /><br />* "kubectl.kubernetes.io/default-container" annotation, set to main pod. |  |  |
| `secrets` _string array_ | Secrets is a list of Secrets in the same namespace as the component<br />object, which shall be mounted into the Prometheus Pods.<br />Each Secret is added to the StatefulSet/Deployment definition as a volume named `secret-<secret-name>`.<br />The Secrets are mounted into /etc/whizard/secrets/<secret-name> in the default container. |  |  |
| `configMaps` _string array_ | ConfigMaps is a list of ConfigMaps in the same namespace as the component<br />object, which shall be mounted into the default Pods.<br />Each ConfigMap is added to the StatefulSet/Deployment definition as a volume named `configmap-<configmap-name>`.<br />The ConfigMaps are mounted into /etc/whizard/configmaps/<configmap-name> in the default container. |  |  |
| `affinity` _[Affinity](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26/#affinity-v1-core)_ | If specified, the pod's scheduling constraints. |  |  |
| `nodeSelector` _object (keys:string, values:string)_ | Define which Nodes the Pods are scheduled on. |  |  |
| `tolerations` _[Toleration](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26/#toleration-v1-core) array_ | If specified, the pod's tolerations. |  |  |
| `resources` _[ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26/#resourcerequirements-v1-core)_ | Define resources requests and limits for main container. |  |  |
| `securityContext` _[PodSecurityContext](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26/#podsecuritycontext-v1-core)_ | SecurityContext holds pod-level security attributes and common container settings.<br />This defaults to the default PodSecurityContext. |  |  |
| `replicas` _integer_ | Number of replicas for a component. |  |  |
| `image` _string_ | Image is the component image with tag/version. |  |  |
| `imagePullPolicy` _[PullPolicy](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26/#pullpolicy-v1-core)_ | Image pull policy.<br />One of Always, Never, IfNotPresent.<br />Defaults to Always if :latest tag is specified, or IfNotPresent otherwise.<br />Cannot be updated. |  |  |
| `imagePullSecrets` _[LocalObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26/#localobjectreference-v1-core) array_ | An optional list of references to secrets in the same namespace<br />to use for pulling images from registries |  |  |
| `logLevel` _string_ | Log filtering level. Possible options: error, warn, info, debug. |  |  |
| `logFormat` _string_ | Log format to use. Possible options: logfmt or json. |  |  |
| `flags` _string array_ | Flags is the flags of component. |  |  |
| `webConfig` _[WebConfig](#webconfig)_ |  |  |  |
| `debug` _boolean_ | If debug mode is on, gateway will proxy Query UI |  |  |
| `enabledTenantsAdmission` _boolean_ | Deny unknown tenant data remote-write and query if enabled |  |  |
| `nodePort` _integer_ | NodePort is the port used to expose the gateway service.<br />If this is a valid node port, the gateway service type will be set to NodePort accordingly. |  |  |


#### GatewayStatus



GatewayStatus defines the observed state of Gateway



_Appears in:_
- [Gateway](#gateway)



#### HTTPClientConfig



HTTPClientConfig configures an HTTP client.



_Appears in:_
- [RemoteQuerySpec](#remotequeryspec)
- [RemoteWriteSpec](#remotewritespec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `basicAuth` _[BasicAuth](#basicauth)_ | The HTTP basic authentication credentials for the targets. |  |  |
| `bearerToken` _string_ | The bearer token for the targets. |  |  |


#### HTTPServerConfig







_Appears in:_
- [WebConfig](#webconfig)



#### HTTPServerTLSConfig







_Appears in:_
- [WebConfig](#webconfig)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `keySecret` _[SecretKeySelector](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26/#secretkeyselector-v1-core)_ | Secret containing the TLS key for the server. |  |  |
| `certSecret` _[SecretKeySelector](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26/#secretkeyselector-v1-core)_ | Contains the TLS certificate for the server. |  |  |
| `clientCASecret` _[SecretKeySelector](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26/#secretkeyselector-v1-core)_ | Contains the CA certificate for client certificate authentication to the server. |  |  |


#### InMemoryIndexCacheConfig







_Appears in:_
- [IndexCacheConfig](#indexcacheconfig)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `maxSize` _string_ | MaxSize represents overall maximum number of bytes cache can contain. |  |  |


#### InMemoryResponseCacheConfig



InMemoryResponseCacheConfig holds the configs for the in-memory cache provider.



_Appears in:_
- [ResponseCacheProviderConfig](#responsecacheproviderconfig)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `maxSize` _string_ | MaxSize represents overall maximum number of bytes cache can contain. |  |  |
| `maxSizeItems` _integer_ | MaxSizeItems represents the maximum number of entries in the cache. |  |  |
| `validity` _[Duration](#duration)_ | Validity represents the expiry duration for the cache. |  | Pattern: `^(0|(([0-9]+)y)?(([0-9]+)w)?(([0-9]+)d)?(([0-9]+)h)?(([0-9]+)m)?(([0-9]+)s)?(([0-9]+)ms)?)$` <br /> |


#### IndexCacheConfig



IndexCacheConfig specifies the index cache config.



_Appears in:_
- [StoreSpec](#storespec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `inMemory` _[InMemoryIndexCacheConfig](#inmemoryindexcacheconfig)_ |  |  |  |


#### Ingester



Ingester is the Schema for the Ingester API





| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `monitoring.whizard.io/v1alpha1` | | |
| `kind` _string_ | `Ingester` | | |
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `spec` _[IngesterSpec](#ingesterspec)_ |  |  |  |
| `status` _[IngesterStatus](#ingesterstatus)_ |  |  |  |


#### IngesterSpec



IngesterSpec defines the desired state of a Ingester



_Appears in:_
- [Ingester](#ingester)
- [IngesterTemplateSpec](#ingestertemplatespec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `containers` _[RawExtension](#rawextension)_ | Containers allows injecting additional containers or modifying operator generated containers.<br />Containers described here modify an operator generated<br />container if they share the same name and modifications are done via a<br />strategic merge patch. |  |  |
| `podMetadata` _[EmbeddedObjectMetadata](#embeddedobjectmetadata)_ | PodMetadata configures labels and annotations which are propagated to the pods.<br /><br />* "kubectl.kubernetes.io/default-container" annotation, set to main pod. |  |  |
| `secrets` _string array_ | Secrets is a list of Secrets in the same namespace as the component<br />object, which shall be mounted into the Prometheus Pods.<br />Each Secret is added to the StatefulSet/Deployment definition as a volume named `secret-<secret-name>`.<br />The Secrets are mounted into /etc/whizard/secrets/<secret-name> in the default container. |  |  |
| `configMaps` _string array_ | ConfigMaps is a list of ConfigMaps in the same namespace as the component<br />object, which shall be mounted into the default Pods.<br />Each ConfigMap is added to the StatefulSet/Deployment definition as a volume named `configmap-<configmap-name>`.<br />The ConfigMaps are mounted into /etc/whizard/configmaps/<configmap-name> in the default container. |  |  |
| `affinity` _[Affinity](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26/#affinity-v1-core)_ | If specified, the pod's scheduling constraints. |  |  |
| `nodeSelector` _object (keys:string, values:string)_ | Define which Nodes the Pods are scheduled on. |  |  |
| `tolerations` _[Toleration](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26/#toleration-v1-core) array_ | If specified, the pod's tolerations. |  |  |
| `resources` _[ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26/#resourcerequirements-v1-core)_ | Define resources requests and limits for main container. |  |  |
| `securityContext` _[PodSecurityContext](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26/#podsecuritycontext-v1-core)_ | SecurityContext holds pod-level security attributes and common container settings.<br />This defaults to the default PodSecurityContext. |  |  |
| `replicas` _integer_ | Number of replicas for a component. |  |  |
| `image` _string_ | Image is the component image with tag/version. |  |  |
| `imagePullPolicy` _[PullPolicy](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26/#pullpolicy-v1-core)_ | Image pull policy.<br />One of Always, Never, IfNotPresent.<br />Defaults to Always if :latest tag is specified, or IfNotPresent otherwise.<br />Cannot be updated. |  |  |
| `imagePullSecrets` _[LocalObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26/#localobjectreference-v1-core) array_ | An optional list of references to secrets in the same namespace<br />to use for pulling images from registries |  |  |
| `logLevel` _string_ | Log filtering level. Possible options: error, warn, info, debug. |  |  |
| `logFormat` _string_ | Log format to use. Possible options: logfmt or json. |  |  |
| `flags` _string array_ | Flags is the flags of component. |  |  |
| `ingesterTsdbCleanup` _[SidecarSpec](#sidecarspec)_ |  |  |  |
| `tenants` _string array_ | Tenants if not empty indicates current config is for hard tenants; otherwise, it is for soft tenants. |  |  |
| `localTsdbRetention` _string_ | LocalTsdbRetention configs how long to retain raw samples on local storage. |  |  |
| `dataVolume` _[KubernetesVolume](#kubernetesvolume)_ | DataVolume specifies how volume shall be used |  |  |


#### IngesterStatus



IngesterStatus defines the observed state of Ingester



_Appears in:_
- [Ingester](#ingester)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `tenants` _[IngesterTenantStatus](#ingestertenantstatus) array_ | Tenants contain all tenants that have been configured for this Ingester object,<br />except those Tenant objects that have been deleted. |  |  |


#### IngesterTemplateSpec







_Appears in:_
- [ServiceSpec](#servicespec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `containers` _[RawExtension](#rawextension)_ | Containers allows injecting additional containers or modifying operator generated containers.<br />Containers described here modify an operator generated<br />container if they share the same name and modifications are done via a<br />strategic merge patch. |  |  |
| `podMetadata` _[EmbeddedObjectMetadata](#embeddedobjectmetadata)_ | PodMetadata configures labels and annotations which are propagated to the pods.<br /><br />* "kubectl.kubernetes.io/default-container" annotation, set to main pod. |  |  |
| `secrets` _string array_ | Secrets is a list of Secrets in the same namespace as the component<br />object, which shall be mounted into the Prometheus Pods.<br />Each Secret is added to the StatefulSet/Deployment definition as a volume named `secret-<secret-name>`.<br />The Secrets are mounted into /etc/whizard/secrets/<secret-name> in the default container. |  |  |
| `configMaps` _string array_ | ConfigMaps is a list of ConfigMaps in the same namespace as the component<br />object, which shall be mounted into the default Pods.<br />Each ConfigMap is added to the StatefulSet/Deployment definition as a volume named `configmap-<configmap-name>`.<br />The ConfigMaps are mounted into /etc/whizard/configmaps/<configmap-name> in the default container. |  |  |
| `affinity` _[Affinity](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26/#affinity-v1-core)_ | If specified, the pod's scheduling constraints. |  |  |
| `nodeSelector` _object (keys:string, values:string)_ | Define which Nodes the Pods are scheduled on. |  |  |
| `tolerations` _[Toleration](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26/#toleration-v1-core) array_ | If specified, the pod's tolerations. |  |  |
| `resources` _[ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26/#resourcerequirements-v1-core)_ | Define resources requests and limits for main container. |  |  |
| `securityContext` _[PodSecurityContext](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26/#podsecuritycontext-v1-core)_ | SecurityContext holds pod-level security attributes and common container settings.<br />This defaults to the default PodSecurityContext. |  |  |
| `replicas` _integer_ | Number of replicas for a component. |  |  |
| `image` _string_ | Image is the component image with tag/version. |  |  |
| `imagePullPolicy` _[PullPolicy](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26/#pullpolicy-v1-core)_ | Image pull policy.<br />One of Always, Never, IfNotPresent.<br />Defaults to Always if :latest tag is specified, or IfNotPresent otherwise.<br />Cannot be updated. |  |  |
| `imagePullSecrets` _[LocalObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26/#localobjectreference-v1-core) array_ | An optional list of references to secrets in the same namespace<br />to use for pulling images from registries |  |  |
| `logLevel` _string_ | Log filtering level. Possible options: error, warn, info, debug. |  |  |
| `logFormat` _string_ | Log format to use. Possible options: logfmt or json. |  |  |
| `flags` _string array_ | Flags is the flags of component. |  |  |
| `ingesterTsdbCleanup` _[SidecarSpec](#sidecarspec)_ |  |  |  |
| `tenants` _string array_ | Tenants if not empty indicates current config is for hard tenants; otherwise, it is for soft tenants. |  |  |
| `localTsdbRetention` _string_ | LocalTsdbRetention configs how long to retain raw samples on local storage. |  |  |
| `dataVolume` _[KubernetesVolume](#kubernetesvolume)_ | DataVolume specifies how volume shall be used |  |  |
| `defaultTenantsPerIngester` _integer_ | DefaultTenantsPerIngester Whizard default tenant count per ingester.<br /><br />Default: 3 | 3 |  |
| `defaultIngesterRetentionPeriod` _[Duration](#duration)_ | DefaultIngesterRetentionPeriod Whizard default ingester retention period when it has no tenant.<br /><br />Default: "3h" | 3h | Pattern: `^(0|(([0-9]+)y)?(([0-9]+)w)?(([0-9]+)d)?(([0-9]+)h)?(([0-9]+)m)?(([0-9]+)s)?(([0-9]+)ms)?)$` <br /> |
| `disableTsdbCleanup` _boolean_ | DisableTSDBCleanup Disable the TSDB cleanup of ingester.<br />The cleanup will delete the blocks that belong to deleted tenants in the data directory of ingester TSDB.<br /><br />Default: true | true |  |


#### IngesterTenantStatus







_Appears in:_
- [IngesterStatus](#ingesterstatus)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `name` _string_ |  |  |  |
| `obsolete` _boolean_ | true represents that the tenant has been moved to other ingester but may left tsdb data in this ingester. |  |  |


#### KubernetesVolume



KubernetesVolume defines the configured volume for a instance.



_Appears in:_
- [CompactorSpec](#compactorspec)
- [CompactorTemplateSpec](#compactortemplatespec)
- [IngesterSpec](#ingesterspec)
- [IngesterTemplateSpec](#ingestertemplatespec)
- [RulerSpec](#rulerspec)
- [RulerTemplateSpec](#rulertemplatespec)
- [StoreSpec](#storespec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `emptyDir` _[EmptyDirVolumeSource](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26/#emptydirvolumesource-v1-core)_ |  |  |  |
| `persistentVolumeClaim` _[PersistentVolumeClaim](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26/#persistentvolumeclaim-v1-core)_ |  |  |  |
| `persistentVolumeClaimRetentionPolicy` _[StatefulSetPersistentVolumeClaimRetentionPolicy](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26/#statefulsetpersistentvolumeclaimretentionpolicy-v1-apps)_ | persistentVolumeClaimRetentionPolicy describes the lifecycle of persistent<br />volume claims created from persistentVolumeClaim.<br />This requires the kubernetes version >= 1.23 and its StatefulSetAutoDeletePVC feature gate to be enabled. |  |  |


#### ObjectReference







_Appears in:_
- [ServiceSpec](#servicespec)
- [TenantStatus](#tenantstatus)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `namespace` _string_ |  |  |  |
| `name` _string_ |  |  |  |


#### Query



Query is the Schema for the monitoring query API





| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `monitoring.whizard.io/v1alpha1` | | |
| `kind` _string_ | `Query` | | |
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `spec` _[QuerySpec](#queryspec)_ |  |  |  |
| `status` _[QueryStatus](#querystatus)_ |  |  |  |


#### QueryFrontend



QueryFrontend is the Schema for the monitoring queryfrontend API





| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `monitoring.whizard.io/v1alpha1` | | |
| `kind` _string_ | `QueryFrontend` | | |
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `spec` _[QueryFrontendSpec](#queryfrontendspec)_ |  |  |  |
| `status` _[QueryFrontendStatus](#queryfrontendstatus)_ |  |  |  |


#### QueryFrontendSpec







_Appears in:_
- [QueryFrontend](#queryfrontend)
- [ServiceSpec](#servicespec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `containers` _[RawExtension](#rawextension)_ | Containers allows injecting additional containers or modifying operator generated containers.<br />Containers described here modify an operator generated<br />container if they share the same name and modifications are done via a<br />strategic merge patch. |  |  |
| `podMetadata` _[EmbeddedObjectMetadata](#embeddedobjectmetadata)_ | PodMetadata configures labels and annotations which are propagated to the pods.<br /><br />* "kubectl.kubernetes.io/default-container" annotation, set to main pod. |  |  |
| `secrets` _string array_ | Secrets is a list of Secrets in the same namespace as the component<br />object, which shall be mounted into the Prometheus Pods.<br />Each Secret is added to the StatefulSet/Deployment definition as a volume named `secret-<secret-name>`.<br />The Secrets are mounted into /etc/whizard/secrets/<secret-name> in the default container. |  |  |
| `configMaps` _string array_ | ConfigMaps is a list of ConfigMaps in the same namespace as the component<br />object, which shall be mounted into the default Pods.<br />Each ConfigMap is added to the StatefulSet/Deployment definition as a volume named `configmap-<configmap-name>`.<br />The ConfigMaps are mounted into /etc/whizard/configmaps/<configmap-name> in the default container. |  |  |
| `affinity` _[Affinity](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26/#affinity-v1-core)_ | If specified, the pod's scheduling constraints. |  |  |
| `nodeSelector` _object (keys:string, values:string)_ | Define which Nodes the Pods are scheduled on. |  |  |
| `tolerations` _[Toleration](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26/#toleration-v1-core) array_ | If specified, the pod's tolerations. |  |  |
| `resources` _[ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26/#resourcerequirements-v1-core)_ | Define resources requests and limits for main container. |  |  |
| `securityContext` _[PodSecurityContext](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26/#podsecuritycontext-v1-core)_ | SecurityContext holds pod-level security attributes and common container settings.<br />This defaults to the default PodSecurityContext. |  |  |
| `replicas` _integer_ | Number of replicas for a component. |  |  |
| `image` _string_ | Image is the component image with tag/version. |  |  |
| `imagePullPolicy` _[PullPolicy](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26/#pullpolicy-v1-core)_ | Image pull policy.<br />One of Always, Never, IfNotPresent.<br />Defaults to Always if :latest tag is specified, or IfNotPresent otherwise.<br />Cannot be updated. |  |  |
| `imagePullSecrets` _[LocalObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26/#localobjectreference-v1-core) array_ | An optional list of references to secrets in the same namespace<br />to use for pulling images from registries |  |  |
| `logLevel` _string_ | Log filtering level. Possible options: error, warn, info, debug. |  |  |
| `logFormat` _string_ | Log format to use. Possible options: logfmt or json. |  |  |
| `flags` _string array_ | Flags is the flags of component. |  |  |
| `webConfig` _[WebConfig](#webconfig)_ |  |  |  |
| `cacheConfig` _[ResponseCacheProviderConfig](#responsecacheproviderconfig)_ | CacheProviderConfig ... |  |  |


#### QueryFrontendStatus



QueryFrontendStatus defines the observed state of QueryFrontend



_Appears in:_
- [QueryFrontend](#queryfrontend)



#### QuerySpec







_Appears in:_
- [Query](#query)
- [ServiceSpec](#servicespec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `containers` _[RawExtension](#rawextension)_ | Containers allows injecting additional containers or modifying operator generated containers.<br />Containers described here modify an operator generated<br />container if they share the same name and modifications are done via a<br />strategic merge patch. |  |  |
| `podMetadata` _[EmbeddedObjectMetadata](#embeddedobjectmetadata)_ | PodMetadata configures labels and annotations which are propagated to the pods.<br /><br />* "kubectl.kubernetes.io/default-container" annotation, set to main pod. |  |  |
| `secrets` _string array_ | Secrets is a list of Secrets in the same namespace as the component<br />object, which shall be mounted into the Prometheus Pods.<br />Each Secret is added to the StatefulSet/Deployment definition as a volume named `secret-<secret-name>`.<br />The Secrets are mounted into /etc/whizard/secrets/<secret-name> in the default container. |  |  |
| `configMaps` _string array_ | ConfigMaps is a list of ConfigMaps in the same namespace as the component<br />object, which shall be mounted into the default Pods.<br />Each ConfigMap is added to the StatefulSet/Deployment definition as a volume named `configmap-<configmap-name>`.<br />The ConfigMaps are mounted into /etc/whizard/configmaps/<configmap-name> in the default container. |  |  |
| `affinity` _[Affinity](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26/#affinity-v1-core)_ | If specified, the pod's scheduling constraints. |  |  |
| `nodeSelector` _object (keys:string, values:string)_ | Define which Nodes the Pods are scheduled on. |  |  |
| `tolerations` _[Toleration](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26/#toleration-v1-core) array_ | If specified, the pod's tolerations. |  |  |
| `resources` _[ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26/#resourcerequirements-v1-core)_ | Define resources requests and limits for main container. |  |  |
| `securityContext` _[PodSecurityContext](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26/#podsecuritycontext-v1-core)_ | SecurityContext holds pod-level security attributes and common container settings.<br />This defaults to the default PodSecurityContext. |  |  |
| `replicas` _integer_ | Number of replicas for a component. |  |  |
| `image` _string_ | Image is the component image with tag/version. |  |  |
| `imagePullPolicy` _[PullPolicy](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26/#pullpolicy-v1-core)_ | Image pull policy.<br />One of Always, Never, IfNotPresent.<br />Defaults to Always if :latest tag is specified, or IfNotPresent otherwise.<br />Cannot be updated. |  |  |
| `imagePullSecrets` _[LocalObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26/#localobjectreference-v1-core) array_ | An optional list of references to secrets in the same namespace<br />to use for pulling images from registries |  |  |
| `logLevel` _string_ | Log filtering level. Possible options: error, warn, info, debug. |  |  |
| `logFormat` _string_ | Log format to use. Possible options: logfmt or json. |  |  |
| `flags` _string array_ | Flags is the flags of component. |  |  |
| `webConfig` _[WebConfig](#webconfig)_ |  |  |  |
| `promqlEngine` _string_ |  |  |  |
| `stores` _[QueryStores](#querystores) array_ | Additional StoreApi servers from which Query component queries from |  |  |
| `selectorLabels` _object (keys:string, values:string)_ | Selector labels that will be exposed in info endpoint. |  |  |
| `replicaLabelNames` _string array_ | Labels to treat as a replica indicator along which data is deduplicated. |  |  |
| `envoy` _[SidecarSpec](#sidecarspec)_ | Envoy is used to config sidecar which proxies requests requiring auth to the secure stores |  |  |


#### QueryStatus



QueryStatus defines the observed state of Query



_Appears in:_
- [Query](#query)



#### QueryStores







_Appears in:_
- [QuerySpec](#queryspec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `addresses` _string array_ | Address is the addresses of StoreApi server, which may be prefixed with 'dns+' or 'dnssrv+' to detect StoreAPI servers through respective DNS lookups. |  |  |
| `caSecret` _[SecretKeySelector](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26/#secretkeyselector-v1-core)_ | Secret containing the CA cert to use for StoreApi connections |  |  |


#### RemoteQuerySpec



RemoteQuerySpec defines the configuration to query from remote service
which should have prometheus-compatible Query APIs.



_Appears in:_
- [ServiceSpec](#servicespec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `name` _string_ |  |  |  |
| `url` _string_ |  |  |  |
| `basicAuth` _[BasicAuth](#basicauth)_ | The HTTP basic authentication credentials for the targets. |  |  |
| `bearerToken` _string_ | The bearer token for the targets. |  |  |


#### RemoteWriteSpec



RemoteWriteSpec defines the remote write configuration.



_Appears in:_
- [ServiceSpec](#servicespec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `name` _string_ |  |  |  |
| `url` _string_ |  |  |  |
| `headers` _object (keys:string, values:string)_ | Custom HTTP headers to be sent along with each remote write request. |  |  |
| `remoteTimeout` _[Duration](#duration)_ | Timeout for requests to the remote write endpoint. |  | Pattern: `^(0|(([0-9]+)y)?(([0-9]+)w)?(([0-9]+)d)?(([0-9]+)h)?(([0-9]+)m)?(([0-9]+)s)?(([0-9]+)ms)?)$` <br /> |
| `basicAuth` _[BasicAuth](#basicauth)_ | The HTTP basic authentication credentials for the targets. |  |  |
| `bearerToken` _string_ | The bearer token for the targets. |  |  |


#### ResponseCacheProviderConfig



ResponseCacheProviderConfig is the initial ResponseCacheProviderConfig struct holder before parsing it into a specific cache provider.
Based on the config type the config is then parsed into a specific cache provider.



_Appears in:_
- [QueryFrontendSpec](#queryfrontendspec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `type` _[CacheProvider](#cacheprovider)_ |  |  |  |
| `inMemory` _[InMemoryResponseCacheConfig](#inmemoryresponsecacheconfig)_ |  |  |  |


#### Retention



Retention defines the config for retaining samples



_Appears in:_
- [CompactorSpec](#compactorspec)
- [CompactorTemplateSpec](#compactortemplatespec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `retentionRaw` _[Duration](#duration)_ | RetentionRaw specifies how long to retain raw samples in bucket |  | Pattern: `^(0|(([0-9]+)y)?(([0-9]+)w)?(([0-9]+)d)?(([0-9]+)h)?(([0-9]+)m)?(([0-9]+)s)?(([0-9]+)ms)?)$` <br /> |
| `retention5m` _[Duration](#duration)_ | Retention5m specifies how long to retain samples of 5m resolution in bucket |  | Pattern: `^(0|(([0-9]+)y)?(([0-9]+)w)?(([0-9]+)d)?(([0-9]+)h)?(([0-9]+)m)?(([0-9]+)s)?(([0-9]+)ms)?)$` <br /> |
| `retention1h` _[Duration](#duration)_ | Retention1h specifies how long to retain samples of 1h resolution in bucket |  | Pattern: `^(0|(([0-9]+)y)?(([0-9]+)w)?(([0-9]+)d)?(([0-9]+)h)?(([0-9]+)m)?(([0-9]+)s)?(([0-9]+)ms)?)$` <br /> |


#### Router



Router is the Schema for the monitoring router API





| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `monitoring.whizard.io/v1alpha1` | | |
| `kind` _string_ | `Router` | | |
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `spec` _[RouterSpec](#routerspec)_ |  |  |  |
| `status` _[RouterStatus](#routerstatus)_ |  |  |  |


#### RouterSpec







_Appears in:_
- [Router](#router)
- [ServiceSpec](#servicespec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `containers` _[RawExtension](#rawextension)_ | Containers allows injecting additional containers or modifying operator generated containers.<br />Containers described here modify an operator generated<br />container if they share the same name and modifications are done via a<br />strategic merge patch. |  |  |
| `podMetadata` _[EmbeddedObjectMetadata](#embeddedobjectmetadata)_ | PodMetadata configures labels and annotations which are propagated to the pods.<br /><br />* "kubectl.kubernetes.io/default-container" annotation, set to main pod. |  |  |
| `secrets` _string array_ | Secrets is a list of Secrets in the same namespace as the component<br />object, which shall be mounted into the Prometheus Pods.<br />Each Secret is added to the StatefulSet/Deployment definition as a volume named `secret-<secret-name>`.<br />The Secrets are mounted into /etc/whizard/secrets/<secret-name> in the default container. |  |  |
| `configMaps` _string array_ | ConfigMaps is a list of ConfigMaps in the same namespace as the component<br />object, which shall be mounted into the default Pods.<br />Each ConfigMap is added to the StatefulSet/Deployment definition as a volume named `configmap-<configmap-name>`.<br />The ConfigMaps are mounted into /etc/whizard/configmaps/<configmap-name> in the default container. |  |  |
| `affinity` _[Affinity](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26/#affinity-v1-core)_ | If specified, the pod's scheduling constraints. |  |  |
| `nodeSelector` _object (keys:string, values:string)_ | Define which Nodes the Pods are scheduled on. |  |  |
| `tolerations` _[Toleration](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26/#toleration-v1-core) array_ | If specified, the pod's tolerations. |  |  |
| `resources` _[ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26/#resourcerequirements-v1-core)_ | Define resources requests and limits for main container. |  |  |
| `securityContext` _[PodSecurityContext](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26/#podsecuritycontext-v1-core)_ | SecurityContext holds pod-level security attributes and common container settings.<br />This defaults to the default PodSecurityContext. |  |  |
| `replicas` _integer_ | Number of replicas for a component. |  |  |
| `image` _string_ | Image is the component image with tag/version. |  |  |
| `imagePullPolicy` _[PullPolicy](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26/#pullpolicy-v1-core)_ | Image pull policy.<br />One of Always, Never, IfNotPresent.<br />Defaults to Always if :latest tag is specified, or IfNotPresent otherwise.<br />Cannot be updated. |  |  |
| `imagePullSecrets` _[LocalObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26/#localobjectreference-v1-core) array_ | An optional list of references to secrets in the same namespace<br />to use for pulling images from registries |  |  |
| `logLevel` _string_ | Log filtering level. Possible options: error, warn, info, debug. |  |  |
| `logFormat` _string_ | Log format to use. Possible options: logfmt or json. |  |  |
| `flags` _string array_ | Flags is the flags of component. |  |  |
| `webConfig` _[WebConfig](#webconfig)_ |  |  |  |
| `replicationFactor` _integer_ | How many times to replicate incoming write requests |  |  |


#### RouterStatus



RouterStatus defines the observed state of Query



_Appears in:_
- [Router](#router)



#### Ruler



Ruler is the Schema for the Ruler API





| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `monitoring.whizard.io/v1alpha1` | | |
| `kind` _string_ | `Ruler` | | |
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `spec` _[RulerSpec](#rulerspec)_ |  |  |  |
| `status` _[RulerStatus](#rulerstatus)_ |  |  |  |


#### RulerSpec



RulerSpec defines the desired state of a Ruler



_Appears in:_
- [Ruler](#ruler)
- [RulerTemplateSpec](#rulertemplatespec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `containers` _[RawExtension](#rawextension)_ | Containers allows injecting additional containers or modifying operator generated containers.<br />Containers described here modify an operator generated<br />container if they share the same name and modifications are done via a<br />strategic merge patch. |  |  |
| `podMetadata` _[EmbeddedObjectMetadata](#embeddedobjectmetadata)_ | PodMetadata configures labels and annotations which are propagated to the pods.<br /><br />* "kubectl.kubernetes.io/default-container" annotation, set to main pod. |  |  |
| `secrets` _string array_ | Secrets is a list of Secrets in the same namespace as the component<br />object, which shall be mounted into the Prometheus Pods.<br />Each Secret is added to the StatefulSet/Deployment definition as a volume named `secret-<secret-name>`.<br />The Secrets are mounted into /etc/whizard/secrets/<secret-name> in the default container. |  |  |
| `configMaps` _string array_ | ConfigMaps is a list of ConfigMaps in the same namespace as the component<br />object, which shall be mounted into the default Pods.<br />Each ConfigMap is added to the StatefulSet/Deployment definition as a volume named `configmap-<configmap-name>`.<br />The ConfigMaps are mounted into /etc/whizard/configmaps/<configmap-name> in the default container. |  |  |
| `affinity` _[Affinity](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26/#affinity-v1-core)_ | If specified, the pod's scheduling constraints. |  |  |
| `nodeSelector` _object (keys:string, values:string)_ | Define which Nodes the Pods are scheduled on. |  |  |
| `tolerations` _[Toleration](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26/#toleration-v1-core) array_ | If specified, the pod's tolerations. |  |  |
| `resources` _[ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26/#resourcerequirements-v1-core)_ | Define resources requests and limits for main container. |  |  |
| `securityContext` _[PodSecurityContext](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26/#podsecuritycontext-v1-core)_ | SecurityContext holds pod-level security attributes and common container settings.<br />This defaults to the default PodSecurityContext. |  |  |
| `replicas` _integer_ | Number of replicas for a component. |  |  |
| `image` _string_ | Image is the component image with tag/version. |  |  |
| `imagePullPolicy` _[PullPolicy](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26/#pullpolicy-v1-core)_ | Image pull policy.<br />One of Always, Never, IfNotPresent.<br />Defaults to Always if :latest tag is specified, or IfNotPresent otherwise.<br />Cannot be updated. |  |  |
| `imagePullSecrets` _[LocalObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26/#localobjectreference-v1-core) array_ | An optional list of references to secrets in the same namespace<br />to use for pulling images from registries |  |  |
| `logLevel` _string_ | Log filtering level. Possible options: error, warn, info, debug. |  |  |
| `logFormat` _string_ | Log format to use. Possible options: logfmt or json. |  |  |
| `flags` _string array_ | Flags is the flags of component. |  |  |
| `rulerQueryProxy` _[SidecarSpec](#sidecarspec)_ |  |  |  |
| `rulerWriteProxy` _[SidecarSpec](#sidecarspec)_ |  |  |  |
| `prometheusConfigReloader` _[SidecarSpec](#sidecarspec)_ |  |  |  |
| `ruleSelectors` _[LabelSelector](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26/#labelselector-v1-meta) array_ | Label selectors to select which PrometheusRules to mount for alerting and recording.<br />The result of multiple selectors are ORed. |  |  |
| `ruleNamespaceSelector` _[LabelSelector](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26/#labelselector-v1-meta)_ | Namespaces to be selected for PrometheusRules discovery. If unspecified, only<br />the same namespace as the Ruler object is in is used. |  |  |
| `shards` _integer_ | Number of shards to take the hash of fully qualified name of the rule group in order to split rules.<br />Each shard of rules will be bound to one separate statefulset.<br />Default: 1 | 1 |  |
| `tenant` _string_ | Tenant if not empty indicates which tenant's data is evaluated for the selected rules;<br />otherwise, it is for all tenants. |  |  |
| `queryConfig` _[SecretKeySelector](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26/#secretkeyselector-v1-core)_ |  |  |  |
| `remoteWriteConfig` _[SecretKeySelector](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26/#secretkeyselector-v1-core)_ |  |  |  |
| `labels` _object (keys:string, values:string)_ | Labels configure the external label pairs to Ruler. A default replica label<br />`ruler_replica` will be always added  as a label with the value of the pod's name and it will be dropped in the alerts. |  |  |
| `alertDropLabels` _string array_ | AlertDropLabels configure the label names which should be dropped in Ruler alerts.<br />The replica label `ruler_replica` will always be dropped in alerts. |  |  |
| `alertmanagersUrl` _string array_ | Define URLs to send alerts to Alertmanager.<br />Note: this field will be ignored if AlertmanagersConfig is specified.<br />Maps to the `alertmanagers.url` arg. |  |  |
| `alertmanagersConfig` _[SecretKeySelector](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26/#secretkeyselector-v1-core)_ | Define configuration for connecting to alertmanager. Maps to the `alertmanagers.config` arg. |  |  |
| `evaluationInterval` _[Duration](#duration)_ | Interval between consecutive evaluations.<br /><br />Default: "1m" | 1m | Pattern: `^(0|(([0-9]+)y)?(([0-9]+)w)?(([0-9]+)d)?(([0-9]+)h)?(([0-9]+)m)?(([0-9]+)s)?(([0-9]+)ms)?)$` <br /> |
| `dataVolume` _[KubernetesVolume](#kubernetesvolume)_ | DataVolume specifies how volume shall be used |  |  |


#### RulerStatus



RulerStatus defines the observed state of Ruler



_Appears in:_
- [Ruler](#ruler)



#### RulerTemplateSpec







_Appears in:_
- [ServiceSpec](#servicespec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `containers` _[RawExtension](#rawextension)_ | Containers allows injecting additional containers or modifying operator generated containers.<br />Containers described here modify an operator generated<br />container if they share the same name and modifications are done via a<br />strategic merge patch. |  |  |
| `podMetadata` _[EmbeddedObjectMetadata](#embeddedobjectmetadata)_ | PodMetadata configures labels and annotations which are propagated to the pods.<br /><br />* "kubectl.kubernetes.io/default-container" annotation, set to main pod. |  |  |
| `secrets` _string array_ | Secrets is a list of Secrets in the same namespace as the component<br />object, which shall be mounted into the Prometheus Pods.<br />Each Secret is added to the StatefulSet/Deployment definition as a volume named `secret-<secret-name>`.<br />The Secrets are mounted into /etc/whizard/secrets/<secret-name> in the default container. |  |  |
| `configMaps` _string array_ | ConfigMaps is a list of ConfigMaps in the same namespace as the component<br />object, which shall be mounted into the default Pods.<br />Each ConfigMap is added to the StatefulSet/Deployment definition as a volume named `configmap-<configmap-name>`.<br />The ConfigMaps are mounted into /etc/whizard/configmaps/<configmap-name> in the default container. |  |  |
| `affinity` _[Affinity](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26/#affinity-v1-core)_ | If specified, the pod's scheduling constraints. |  |  |
| `nodeSelector` _object (keys:string, values:string)_ | Define which Nodes the Pods are scheduled on. |  |  |
| `tolerations` _[Toleration](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26/#toleration-v1-core) array_ | If specified, the pod's tolerations. |  |  |
| `resources` _[ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26/#resourcerequirements-v1-core)_ | Define resources requests and limits for main container. |  |  |
| `securityContext` _[PodSecurityContext](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26/#podsecuritycontext-v1-core)_ | SecurityContext holds pod-level security attributes and common container settings.<br />This defaults to the default PodSecurityContext. |  |  |
| `replicas` _integer_ | Number of replicas for a component. |  |  |
| `image` _string_ | Image is the component image with tag/version. |  |  |
| `imagePullPolicy` _[PullPolicy](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26/#pullpolicy-v1-core)_ | Image pull policy.<br />One of Always, Never, IfNotPresent.<br />Defaults to Always if :latest tag is specified, or IfNotPresent otherwise.<br />Cannot be updated. |  |  |
| `imagePullSecrets` _[LocalObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26/#localobjectreference-v1-core) array_ | An optional list of references to secrets in the same namespace<br />to use for pulling images from registries |  |  |
| `logLevel` _string_ | Log filtering level. Possible options: error, warn, info, debug. |  |  |
| `logFormat` _string_ | Log format to use. Possible options: logfmt or json. |  |  |
| `flags` _string array_ | Flags is the flags of component. |  |  |
| `rulerQueryProxy` _[SidecarSpec](#sidecarspec)_ |  |  |  |
| `rulerWriteProxy` _[SidecarSpec](#sidecarspec)_ |  |  |  |
| `prometheusConfigReloader` _[SidecarSpec](#sidecarspec)_ |  |  |  |
| `ruleSelectors` _[LabelSelector](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26/#labelselector-v1-meta) array_ | Label selectors to select which PrometheusRules to mount for alerting and recording.<br />The result of multiple selectors are ORed. |  |  |
| `ruleNamespaceSelector` _[LabelSelector](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26/#labelselector-v1-meta)_ | Namespaces to be selected for PrometheusRules discovery. If unspecified, only<br />the same namespace as the Ruler object is in is used. |  |  |
| `shards` _integer_ | Number of shards to take the hash of fully qualified name of the rule group in order to split rules.<br />Each shard of rules will be bound to one separate statefulset.<br />Default: 1 | 1 |  |
| `tenant` _string_ | Tenant if not empty indicates which tenant's data is evaluated for the selected rules;<br />otherwise, it is for all tenants. |  |  |
| `queryConfig` _[SecretKeySelector](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26/#secretkeyselector-v1-core)_ |  |  |  |
| `remoteWriteConfig` _[SecretKeySelector](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26/#secretkeyselector-v1-core)_ |  |  |  |
| `labels` _object (keys:string, values:string)_ | Labels configure the external label pairs to Ruler. A default replica label<br />`ruler_replica` will be always added  as a label with the value of the pod's name and it will be dropped in the alerts. |  |  |
| `alertDropLabels` _string array_ | AlertDropLabels configure the label names which should be dropped in Ruler alerts.<br />The replica label `ruler_replica` will always be dropped in alerts. |  |  |
| `alertmanagersUrl` _string array_ | Define URLs to send alerts to Alertmanager.<br />Note: this field will be ignored if AlertmanagersConfig is specified.<br />Maps to the `alertmanagers.url` arg. |  |  |
| `alertmanagersConfig` _[SecretKeySelector](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26/#secretkeyselector-v1-core)_ | Define configuration for connecting to alertmanager. Maps to the `alertmanagers.config` arg. |  |  |
| `evaluationInterval` _[Duration](#duration)_ | Interval between consecutive evaluations.<br /><br />Default: "1m" | 1m | Pattern: `^(0|(([0-9]+)y)?(([0-9]+)w)?(([0-9]+)d)?(([0-9]+)h)?(([0-9]+)m)?(([0-9]+)s)?(([0-9]+)ms)?)$` <br /> |
| `dataVolume` _[KubernetesVolume](#kubernetesvolume)_ | DataVolume specifies how volume shall be used |  |  |
| `disableAlertingRulesAutoSelection` _boolean_ | DisableAlertingRulesAutoSelection disable auto select alerting rules in tenant ruler<br /><br />Default: true | true |  |


#### S3



Config stores the configuration for s3 bucket.



_Appears in:_
- [StorageSpec](#storagespec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `bucket` _string_ |  |  |  |
| `endpoint` _string_ |  |  |  |
| `region` _string_ |  |  |  |
| `awsSdkAuth` _boolean_ |  |  |  |
| `accessKey` _[SecretKeySelector](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26/#secretkeyselector-v1-core)_ |  |  |  |
| `insecure` _boolean_ |  |  |  |
| `signatureVersion2` _boolean_ |  |  |  |
| `secretKey` _[SecretKeySelector](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26/#secretkeyselector-v1-core)_ |  |  |  |
| `putUserMetadata` _object (keys:string, values:string)_ |  |  |  |
| `httpConfig` _[S3HTTPConfig](#s3httpconfig)_ |  |  |  |
| `trace` _[S3TraceConfig](#s3traceconfig)_ |  |  |  |
| `listObjectsVersion` _string_ |  |  |  |
| `partSize` _integer_ | PartSize used for multipart upload. Only used if uploaded object size is known and larger than configured PartSize.<br />NOTE we need to make sure this number does not produce more parts than 10 000. |  |  |
| `sseConfig` _[S3SSEConfig](#s3sseconfig)_ |  |  |  |
| `stsEndpoint` _string_ |  |  |  |


#### S3HTTPConfig



S3HTTPConfig stores the http.Transport configuration for the s3 minio client.



_Appears in:_
- [S3](#s3)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `idleConnTimeout` _[Duration](#duration)_ |  |  | Pattern: `^(0|(([0-9]+)y)?(([0-9]+)w)?(([0-9]+)d)?(([0-9]+)h)?(([0-9]+)m)?(([0-9]+)s)?(([0-9]+)ms)?)$` <br /> |
| `responseHeaderTimeout` _[Duration](#duration)_ |  |  | Pattern: `^(0|(([0-9]+)y)?(([0-9]+)w)?(([0-9]+)d)?(([0-9]+)h)?(([0-9]+)m)?(([0-9]+)s)?(([0-9]+)ms)?)$` <br /> |
| `insecureSkipVerify` _boolean_ |  |  |  |
| `tlsHandshakeTimeout` _[Duration](#duration)_ |  |  | Pattern: `^(0|(([0-9]+)y)?(([0-9]+)w)?(([0-9]+)d)?(([0-9]+)h)?(([0-9]+)m)?(([0-9]+)s)?(([0-9]+)ms)?)$` <br /> |
| `expectContinueTimeout` _[Duration](#duration)_ |  |  | Pattern: `^(0|(([0-9]+)y)?(([0-9]+)w)?(([0-9]+)d)?(([0-9]+)h)?(([0-9]+)m)?(([0-9]+)s)?(([0-9]+)ms)?)$` <br /> |
| `maxIdleConns` _integer_ |  |  |  |
| `maxIdleConnsPerHost` _integer_ |  |  |  |
| `maxConnsPerHost` _integer_ |  |  |  |
| `tlsConfig` _[TLSConfig](#tlsconfig)_ |  |  |  |


#### S3SSEConfig



S3SSEConfig deals with the configuration of SSE for Minio. The following options are valid:
kmsencryptioncontext == https://docs.aws.amazon.com/kms/latest/developerguide/services-s3.html#s3-encryption-context



_Appears in:_
- [S3](#s3)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `type` _string_ |  |  |  |
| `kmsKeyId` _string_ |  |  |  |
| `kmsEncryptionContext` _object (keys:string, values:string)_ |  |  |  |
| `encryptionKey` _string_ |  |  |  |


#### S3TraceConfig







_Appears in:_
- [S3](#s3)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `enable` _boolean_ |  |  |  |


#### Service



Service is the Schema for the monitoring service API





| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `monitoring.whizard.io/v1alpha1` | | |
| `kind` _string_ | `Service` | | |
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `spec` _[ServiceSpec](#servicespec)_ |  |  |  |
| `status` _[ServiceStatus](#servicestatus)_ |  |  |  |


#### ServiceSpec



ServiceSpec defines the desired state of a Service



_Appears in:_
- [Service](#service)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `tenantHeader` _string_ | HTTP header to determine tenant for remote write requests. |  |  |
| `defaultTenantId` _string_ | Default tenant ID to use when none is provided via a header. |  |  |
| `tenantLabelName` _string_ | Label name through which the tenant will be announced. |  |  |
| `storage` _[ObjectReference](#objectreference)_ |  |  |  |
| `remoteWrites` _[RemoteWriteSpec](#remotewritespec) array_ | RemoteWrites is the list of remote write configurations.<br />If it is configured, its targets will receive write requests from the Gateway and the Ruler. |  |  |
| `remoteQuery` _[RemoteQuerySpec](#remotequeryspec)_ | RemoteQuery is the remote query configuration and the remote target should have prometheus-compatible Query APIs.<br />If not configured, the Gateway will proxy all read requests through the QueryFrontend to the Query,<br />If configured, the Gateway will proxy metrics read requests through the QueryFrontend to the remote target,<br />but proxy rules read requests directly to the Query. |  |  |
| `gatewayTemplateSpec` _[GatewaySpec](#gatewayspec)_ |  |  |  |
| `queryFrontendTemplateSpec` _[QueryFrontendSpec](#queryfrontendspec)_ |  |  |  |
| `queryTemplateSpec` _[QuerySpec](#queryspec)_ |  |  |  |
| `rulerTemplateSpec` _[RulerTemplateSpec](#rulertemplatespec)_ |  |  |  |
| `routerTemplateSpec` _[RouterSpec](#routerspec)_ |  |  |  |
| `ingesterTemplateSpec` _[IngesterTemplateSpec](#ingestertemplatespec)_ |  |  |  |
| `storeTemplateSpec` _[StoreSpec](#storespec)_ |  |  |  |
| `compactorTemplateSpec` _[CompactorTemplateSpec](#compactortemplatespec)_ |  |  |  |


#### ServiceStatus



ServiceStatus defines the observed state of Service



_Appears in:_
- [Service](#service)



#### SidecarSpec







_Appears in:_
- [IngesterSpec](#ingesterspec)
- [IngesterTemplateSpec](#ingestertemplatespec)
- [QuerySpec](#queryspec)
- [RulerSpec](#rulerspec)
- [RulerTemplateSpec](#rulertemplatespec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `image` _string_ | Image is the envoy image with tag/version |  |  |
| `resources` _[ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26/#resourcerequirements-v1-core)_ | Define resources requests and limits for sidecar container. |  |  |


#### Storage









| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `monitoring.whizard.io/v1alpha1` | | |
| `kind` _string_ | `Storage` | | |
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `spec` _[StorageSpec](#storagespec)_ |  |  |  |
| `status` _[StorageStatus](#storagestatus)_ |  |  |  |


#### StorageSpec







_Appears in:_
- [Storage](#storage)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `blockManager` _[BlockManager](#blockmanager)_ |  |  |  |
| `S3` _[S3](#s3)_ |  |  |  |


#### StorageStatus







_Appears in:_
- [Storage](#storage)



#### Store



Store is the Schema for the Store API





| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `monitoring.whizard.io/v1alpha1` | | |
| `kind` _string_ | `Store` | | |
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `spec` _[StoreSpec](#storespec)_ |  |  |  |
| `status` _[StoreStatus](#storestatus)_ |  |  |  |


#### StoreSpec







_Appears in:_
- [ServiceSpec](#servicespec)
- [Store](#store)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `containers` _[RawExtension](#rawextension)_ | Containers allows injecting additional containers or modifying operator generated containers.<br />Containers described here modify an operator generated<br />container if they share the same name and modifications are done via a<br />strategic merge patch. |  |  |
| `podMetadata` _[EmbeddedObjectMetadata](#embeddedobjectmetadata)_ | PodMetadata configures labels and annotations which are propagated to the pods.<br /><br />* "kubectl.kubernetes.io/default-container" annotation, set to main pod. |  |  |
| `secrets` _string array_ | Secrets is a list of Secrets in the same namespace as the component<br />object, which shall be mounted into the Prometheus Pods.<br />Each Secret is added to the StatefulSet/Deployment definition as a volume named `secret-<secret-name>`.<br />The Secrets are mounted into /etc/whizard/secrets/<secret-name> in the default container. |  |  |
| `configMaps` _string array_ | ConfigMaps is a list of ConfigMaps in the same namespace as the component<br />object, which shall be mounted into the default Pods.<br />Each ConfigMap is added to the StatefulSet/Deployment definition as a volume named `configmap-<configmap-name>`.<br />The ConfigMaps are mounted into /etc/whizard/configmaps/<configmap-name> in the default container. |  |  |
| `affinity` _[Affinity](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26/#affinity-v1-core)_ | If specified, the pod's scheduling constraints. |  |  |
| `nodeSelector` _object (keys:string, values:string)_ | Define which Nodes the Pods are scheduled on. |  |  |
| `tolerations` _[Toleration](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26/#toleration-v1-core) array_ | If specified, the pod's tolerations. |  |  |
| `resources` _[ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26/#resourcerequirements-v1-core)_ | Define resources requests and limits for main container. |  |  |
| `securityContext` _[PodSecurityContext](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26/#podsecuritycontext-v1-core)_ | SecurityContext holds pod-level security attributes and common container settings.<br />This defaults to the default PodSecurityContext. |  |  |
| `replicas` _integer_ | Number of replicas for a component. |  |  |
| `image` _string_ | Image is the component image with tag/version. |  |  |
| `imagePullPolicy` _[PullPolicy](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26/#pullpolicy-v1-core)_ | Image pull policy.<br />One of Always, Never, IfNotPresent.<br />Defaults to Always if :latest tag is specified, or IfNotPresent otherwise.<br />Cannot be updated. |  |  |
| `imagePullSecrets` _[LocalObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26/#localobjectreference-v1-core) array_ | An optional list of references to secrets in the same namespace<br />to use for pulling images from registries |  |  |
| `logLevel` _string_ | Log filtering level. Possible options: error, warn, info, debug. |  |  |
| `logFormat` _string_ | Log format to use. Possible options: logfmt or json. |  |  |
| `flags` _string array_ | Flags is the flags of component. |  |  |
| `minTime` _string_ | MinTime specifies start of time range limit to serve |  |  |
| `maxTime` _string_ | MaxTime specifies end of time range limit to serve |  |  |
| `timeRanges` _[TimeRange](#timerange) array_ | TimeRanges is a list of TimeRange to partition Store.<br />If specified, the MinTime and MaxTime will be ignored. |  |  |
| `indexCacheConfig` _[IndexCacheConfig](#indexcacheconfig)_ | IndexCacheConfig contains index cache configuration. |  |  |
| `dataVolume` _[KubernetesVolume](#kubernetesvolume)_ | DataVolume specifies how volume shall be used |  |  |
| `scaler` _[AutoScaler](#autoscaler)_ |  |  |  |


#### StoreStatus



StoreStatus defines the observed state of Store



_Appears in:_
- [Store](#store)



#### TLSConfig



TLSConfig configures the options for TLS connections.



_Appears in:_
- [S3HTTPConfig](#s3httpconfig)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `ca` _[SecretKeySelector](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26/#secretkeyselector-v1-core)_ | The secret that including the CA cert. |  |  |
| `cert` _[SecretKeySelector](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26/#secretkeyselector-v1-core)_ | The secret that including the client cert. |  |  |
| `key` _[SecretKeySelector](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26/#secretkeyselector-v1-core)_ | The secret that including the client key. |  |  |
| `serverName` _string_ | Used to verify the hostname for the targets. |  |  |
| `insecureSkipVerify` _boolean_ | Disable target certificate validation. |  |  |


#### Tenant



Tenant is the Schema for the monitoring Tenant API





| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `monitoring.whizard.io/v1alpha1` | | |
| `kind` _string_ | `Tenant` | | |
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `spec` _[TenantSpec](#tenantspec)_ |  |  |  |
| `status` _[TenantStatus](#tenantstatus)_ |  |  |  |


#### TenantSpec







_Appears in:_
- [Tenant](#tenant)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `tenant` _string_ |  |  |  |


#### TenantStatus







_Appears in:_
- [Tenant](#tenant)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `ruler` _[ObjectReference](#objectreference)_ |  |  |  |
| `compactor` _[ObjectReference](#objectreference)_ |  |  |  |
| `ingester` _[ObjectReference](#objectreference)_ |  |  |  |


#### TimeRange







_Appears in:_
- [StoreSpec](#storespec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `minTime` _string_ | MinTime specifies start of time range limit to serve |  |  |
| `maxTime` _string_ | MaxTime specifies end of time range limit to serve |  |  |


#### WebConfig



WebConfig defines the configuration for the HTTP server.
More info: https://prometheus.io/docs/prometheus/latest/configuration/https/



_Appears in:_
- [GatewaySpec](#gatewayspec)
- [QueryFrontendSpec](#queryfrontendspec)
- [QuerySpec](#queryspec)
- [RouterSpec](#routerspec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `httpServerTLSConfig` _[HTTPServerTLSConfig](#httpservertlsconfig)_ |  |  |  |
| `httpServerConfig` _[HTTPServerConfig](#httpserverconfig)_ |  |  |  |
| `basicAuthUsers` _[BasicAuth](#basicauth) array_ |  |  |  |



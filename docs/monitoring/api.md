

# API Docs

This Document documents the types introduced by the whizard to be consumed by users.

> Note this document is generated from code comments. When contributing a change to this document please do so by changing the code comments.

## Table of Contents
* [AutoScaler](#autoscaler)
* [Compactor](#compactor)
* [CompactorList](#compactorlist)
* [CompactorSpec](#compactorspec)
* [EnvoySpec](#envoyspec)
* [Gateway](#gateway)
* [InMemoryIndexCacheConfig](#inmemoryindexcacheconfig)
* [InMemoryResponseCacheConfig](#inmemoryresponsecacheconfig)
* [IndexCacheConfig](#indexcacheconfig)
* [Ingester](#ingester)
* [IngesterList](#ingesterlist)
* [IngesterSpec](#ingesterspec)
* [KubernetesVolume](#kubernetesvolume)
* [ObjectReference](#objectreference)
* [Query](#query)
* [QueryFrontend](#queryfrontend)
* [QueryStores](#querystores)
* [ResponseCacheProviderConfig](#responsecacheproviderconfig)
* [Retention](#retention)
* [Router](#router)
* [Ruler](#ruler)
* [RulerList](#rulerlist)
* [RulerSpec](#rulerspec)
* [Service](#service)
* [ServiceList](#servicelist)
* [ServiceSpec](#servicespec)
* [Store](#store)
* [StoreList](#storelist)
* [StoreSpec](#storespec)
* [S3](#s3)
* [S3HTTPConfig](#s3httpconfig)
* [S3SSEConfig](#s3sseconfig)
* [S3TraceConfig](#s3traceconfig)
* [Storage](#storage)
* [StorageList](#storagelist)
* [StorageSpec](#storagespec)
* [TLSConfig](#tlsconfig)
* [Tenant](#tenant)
* [TenantList](#tenantlist)
* [TenantSpec](#tenantspec)
* [TenantStatus](#tenantstatus)

## AutoScaler



| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| minReplicas | minReplicas is the lower limit for the number of replicas to which the autoscaler can scale down.  It defaults to 1 pod.  minReplicas is allowed to be 0 if the alpha feature gate HPAScaleToZero is enabled and at least one Object or External metric is configured.  Scaling is active as long as at least one metric value is available. | *int32 | false |
| maxReplicas | maxReplicas is the upper limit for the number of replicas to which the autoscaler can scale up. It cannot be less that minReplicas. | int32 | true |
| metrics | metrics contains the specifications for which to use to calculate the desired replica count (the maximum replica count across all metrics will be used).  The desired replica count is calculated multiplying the ratio between the target value and the current value by the current number of pods.  Ergo, metrics used must decrease as the pod count is increased, and vice-versa.  See the individual metric source types for more information about how each type of metric must respond. If not set, the default metric will be set to 80% average CPU utilization. | []v2beta2.MetricSpec | false |
| behavior | behavior configures the scaling behavior of the target in both Up and Down directions (scaleUp and scaleDown fields respectively). If not set, the default HPAScalingRules for scale up and scale down are used. | *v2beta2.HorizontalPodAutoscalerBehavior | false |

[Back to TOC](#table-of-contents)

## Compactor

Compactor is the Schema for the Compactor API

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| metadata |  | [metav1.ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.21/#objectmeta-v1-meta) | false |
| spec |  | [CompactorSpec](#compactorspec) | false |
| status |  | [CompactorStatus](#compactorstatus) | false |

[Back to TOC](#table-of-contents)

## CompactorList

CompactorList contains a list of Compactor

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| metadata |  | [metav1.ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.21/#listmeta-v1-meta) | false |
| items |  | [][Compactor](#compactor) | true |

[Back to TOC](#table-of-contents)

## CompactorSpec



| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| affinity | If specified, the pod's scheduling constraints. | *corev1.Affinity | false |
| nodeSelector | Define which Nodes the Pods are scheduled on. | map[string]string | false |
| tolerations | If specified, the pod's tolerations. | []corev1.Toleration | false |
| resources | Define resources requests and limits for main container. | [corev1.ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.21/#resourcerequirements-v1-core) | false |
| replicas | Number of replicas for a component | *int32 | false |
| image | Image is the image with tag/version | string | false |
| imagePullPolicy | Image pull policy. One of Always, Never, IfNotPresent. Defaults to Always if :latest tag is specified, or IfNotPresent otherwise. Cannot be updated. | [corev1.PullPolicy](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.21/#container-v1-core) | false |
| logLevel | Log filtering level. Possible options: error, warn, info, debug | string | false |
| logFormat | Log format to use. Possible options: logfmt or json | string | false |
| downsamplingDisable | DownsamplingDisable specifies whether to disable downsampling | *bool | false |
| retention | Retention configs how long to retain samples | *[Retention](#retention) | false |
| storage |  | *[ObjectReference](#objectreference) | false |
| flags | Flags is the flags of compactor. | []string | false |
| dataVolume | DataVolume specifies how volume shall be used | *[KubernetesVolume](#kubernetesvolume) | false |
| tenants | Tenants if not empty indicates current config is for hard tenants; otherwise, it is for soft tenants. | []string | false |

[Back to TOC](#table-of-contents)

## EnvoySpec

EnvoySpec defines the desired state of envoy proxy sidecar which delegates requests to the secure stores

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| image | Image is the envoy image with tag/version | string | false |
| resources | Define resources requests and limits for envoy container. | [corev1.ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.21/#resourcerequirements-v1-core) | false |

[Back to TOC](#table-of-contents)

## Gateway



| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| affinity | If specified, the pod's scheduling constraints. | *corev1.Affinity | false |
| nodeSelector | Define which Nodes the Pods are scheduled on. | map[string]string | false |
| tolerations | If specified, the pod's tolerations. | []corev1.Toleration | false |
| resources | Define resources requests and limits for main container. | [corev1.ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.21/#resourcerequirements-v1-core) | false |
| replicas | Number of replicas for a component | *int32 | false |
| image | Image is the gateway image with tag/version. | string | false |
| logLevel | Log filtering level. Possible options: error, warn, info, debug. | string | false |
| logFormat | Log format to use. Possible options: logfmt or json. | string | false |
| serverCertificate | Secret name for HTTP Server certificate (Kubernetes TLS secret type) | string | false |
| clientCaCertificate | Secret name for HTTP Client CA certificate (Kubernetes TLS secret type) | string | false |

[Back to TOC](#table-of-contents)

## InMemoryIndexCacheConfig



| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| maxSize | MaxSize represents overall maximum number of bytes cache can contain. | string | true |

[Back to TOC](#table-of-contents)

## InMemoryResponseCacheConfig

InMemoryResponseCacheConfig holds the configs for the in-memory cache provider.

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| maxSize | MaxSize represents overall maximum number of bytes cache can contain. | string | true |
| maxSizeItems | MaxSizeItems represents the maximum number of entries in the cache. | int | true |
| validity | Validity represents the expiry duration for the cache. | time.Duration | true |

[Back to TOC](#table-of-contents)

## IndexCacheConfig

IndexCacheConfig specifies the index cache config.

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| inMemory |  | *[InMemoryIndexCacheConfig](#inmemoryindexcacheconfig) | false |

[Back to TOC](#table-of-contents)

## Ingester

Ingester is the Schema for the Ingester API

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| metadata |  | [metav1.ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.21/#objectmeta-v1-meta) | false |
| spec |  | [IngesterSpec](#ingesterspec) | false |
| status |  | [IngesterStatus](#ingesterstatus) | false |

[Back to TOC](#table-of-contents)

## IngesterList

IngesterList contains a list of Ingester

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| metadata |  | [metav1.ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.21/#listmeta-v1-meta) | false |
| items |  | [][Ingester](#ingester) | true |

[Back to TOC](#table-of-contents)

## IngesterSpec

IngesterSpec defines the desired state of a Ingester

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| tenants | Tenants if not empty indicates current config is for hard tenants; otherwise, it is for soft tenants. | []string | false |
| affinity | If specified, the pod's scheduling constraints. | *corev1.Affinity | false |
| nodeSelector | Define which Nodes the Pods are scheduled on. | map[string]string | false |
| tolerations | If specified, the pod's tolerations. | []corev1.Toleration | false |
| resources | Define resources requests and limits for main container. | [corev1.ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.21/#resourcerequirements-v1-core) | false |
| replicas | Number of replicas for a component. | *int32 | false |
| image | Image is the image with tag/version | string | false |
| logLevel | Log filtering level. Possible options: error, warn, info, debug | string | false |
| logFormat | Log format to use. Possible options: logfmt or json | string | false |
| localTsdbRetention | LocalTsdbRetention configs how long to retain raw samples on local storage. | string | false |
| flags | Flags is the flags of ingester. | []string | false |
| storage | If specified, the object key of Storage for long term storage. | *[ObjectReference](#objectreference) | false |
| dataVolume | DataVolume specifies how volume shall be used | *[KubernetesVolume](#kubernetesvolume) | false |

[Back to TOC](#table-of-contents)

## KubernetesVolume

KubernetesVolume defines the configured volume for a instance.

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| emptyDir |  | *corev1.EmptyDirVolumeSource | false |
| pvc |  | *corev1.PersistentVolumeClaim | false |

[Back to TOC](#table-of-contents)

## ObjectReference



| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| namespace |  | string | false |
| name |  | string | false |

[Back to TOC](#table-of-contents)

## Query



| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| affinity | If specified, the pod's scheduling constraints. | *corev1.Affinity | false |
| nodeSelector | Define which Nodes the Pods are scheduled on. | map[string]string | false |
| tolerations | If specified, the pod's tolerations. | []corev1.Toleration | false |
| resources | Define resources requests and limits for main container. | [corev1.ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.21/#resourcerequirements-v1-core) | false |
| replicas | Number of replicas for a component | *int32 | false |
| image | Image is the image with tag/version | string | false |
| logLevel | Log filtering level. Possible options: error, warn, info, debug | string | false |
| logFormat | Log format to use. Possible options: logfmt or json | string | false |
| stores | Additional StoreApi servers from which Query component queries from | [][QueryStores](#querystores) | false |
| selectorLabels | Selector labels that will be exposed in info endpoint. | map[string]string | false |
| replicaLabelNames | Labels to treat as a replica indicator along which data is deduplicated. | []string | false |
| flags | Flags is the flags of query. | []string | false |
| envoy | Envoy is used to config sidecar which proxies requests requiring auth to the secure stores | [EnvoySpec](#envoyspec) | false |

[Back to TOC](#table-of-contents)

## QueryFrontend



| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| affinity | If specified, the pod's scheduling constraints. | *corev1.Affinity | false |
| nodeSelector | Define which Nodes the Pods are scheduled on. | map[string]string | false |
| tolerations | If specified, the pod's tolerations. | []corev1.Toleration | false |
| resources | Define resources requests and limits for main container. | [corev1.ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.21/#resourcerequirements-v1-core) | false |
| replicas | Number of replicas for a component | *int32 | false |
| image | Image is the image with tag/version | string | false |
| logLevel | Log filtering level. Possible options: error, warn, info, debug | string | false |
| logFormat | Log format to use. Possible options: logfmt or json | string | false |
| flags | Flags is the flags of query frontend. | []string | false |
| cacheConfig | CacheProviderConfig ... | *[ResponseCacheProviderConfig](#responsecacheproviderconfig) | false |

[Back to TOC](#table-of-contents)

## QueryStores



| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| addresses | Address is the addresses of StoreApi server, which may be prefixed with 'dns+' or 'dnssrv+' to detect StoreAPI servers through respective DNS lookups. | []string | false |
| caSecret | Secret containing the CA cert to use for StoreApi connections | *corev1.SecretKeySelector | false |

[Back to TOC](#table-of-contents)

## ResponseCacheProviderConfig

ResponseCacheProviderConfig is the initial ResponseCacheProviderConfig struct holder before parsing it into a specific cache provider. Based on the config type the config is then parsed into a specific cache provider.

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| type |  | CacheProvider | true |
| inMemory |  | *[InMemoryResponseCacheConfig](#inmemoryresponsecacheconfig) | false |

[Back to TOC](#table-of-contents)

## Retention

Retention defines the config for retaining samples

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| retentionRaw | RetentionRaw specifies how long to retain raw samples in bucket | Duration | false |
| retention5m | Retention5m specifies how long to retain samples of 5m resolution in bucket | Duration | false |
| retention1h | Retention1h specifies how long to retain samples of 1h resolution in bucket | Duration | false |

[Back to TOC](#table-of-contents)

## Router



| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| affinity | If specified, the pod's scheduling constraints. | *corev1.Affinity | false |
| nodeSelector | Define which Nodes the Pods are scheduled on. | map[string]string | false |
| tolerations | If specified, the pod's tolerations. | []corev1.Toleration | false |
| resources | Define resources requests and limits for main container. | [corev1.ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.21/#resourcerequirements-v1-core) | false |
| replicas | Number of replicas for a component. | *int32 | false |
| image | Image is the image with tag/version | string | false |
| logLevel | Log filtering level. Possible options: error, warn, info, debug | string | false |
| logFormat | Log format to use. Possible options: logfmt or json | string | false |
| replicationFactor | How many times to replicate incoming write requests | *uint64 | false |
| flags | Flags is the flags of router. | []string | false |

[Back to TOC](#table-of-contents)

## Ruler

Ruler is the Schema for the Ruler API

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| metadata |  | [metav1.ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.21/#objectmeta-v1-meta) | false |
| spec |  | [RulerSpec](#rulerspec) | false |
| status |  | [RulerStatus](#rulerstatus) | false |

[Back to TOC](#table-of-contents)

## RulerList

RulerList contains a list of Ruler

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| metadata |  | [metav1.ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.21/#listmeta-v1-meta) | false |
| items |  | [][Ruler](#ruler) | true |

[Back to TOC](#table-of-contents)

## RulerSpec

RulerSpec defines the desired state of a Ruler

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| affinity | If specified, the pod's scheduling constraints. | *corev1.Affinity | false |
| nodeSelector | Define which Nodes the Pods are scheduled on. | map[string]string | false |
| tolerations | If specified, the pod's tolerations. | []corev1.Toleration | false |
| resources | Define resources requests and limits for main container. | [corev1.ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.21/#resourcerequirements-v1-core) | false |
| replicas | Number of replicas for a component. | *int32 | false |
| image | Image is the image with tag/version | string | false |
| logLevel | Log filtering level. Possible options: error, warn, info, debug | string | false |
| logFormat | Log format to use. Possible options: logfmt or json | string | false |
| ruleSelector | A label selector to select which PrometheusRules to mount for alerting and recording. | *[metav1.LabelSelector](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.21/#labelselector-v1-meta) | false |
| ruleNamespaceSelector | Namespaces to be selected for PrometheusRules discovery. If unspecified, only the same namespace as the Ruler object is in is used. | *[metav1.LabelSelector](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.21/#labelselector-v1-meta) | false |
| shards | Number of shards to take the hash of fully qualified name of the rule group in order to split rules. Each shard of rules will be bound to one separate statefulset. Default: `1` | *int32 | false |
| tenant | Tenant if not empty indicates which tenant's data is evaluated for the selected rules; otherwise, it is for all tenants. | string | false |
| labels | Labels configure the external label pairs to Ruler. A default replica label `ruler_replica` will be always added  as a label with the value of the pod's name and it will be dropped in the alerts. | map[string]string | false |
| alertDropLabels | AlertDropLabels configure the label names which should be dropped in Ruler alerts. The replica label `ruler_replica` will always be dropped in alerts. | []string | false |
| alertmanagersConfig | Define configuration for connecting to alertmanager. Maps to the `alertmanagers.config` arg. | *corev1.SecretKeySelector | false |
| evaluationInterval | Interval between consecutive evaluations. Default: `30s` | Duration | false |
| flags | Flags is the flags of ruler. | []string | false |
| dataVolume | DataVolume specifies how volume shall be used | *[KubernetesVolume](#kubernetesvolume) | false |

[Back to TOC](#table-of-contents)

## Service

Service is the Schema for the monitoring service API

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| metadata |  | [metav1.ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.21/#objectmeta-v1-meta) | false |
| spec |  | [ServiceSpec](#servicespec) | false |
| status |  | [ServiceStatus](#servicestatus) | false |

[Back to TOC](#table-of-contents)

## ServiceList

ServiceList contains a list of Service

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| metadata |  | [metav1.ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.21/#listmeta-v1-meta) | false |
| items |  | [][Service](#service) | true |

[Back to TOC](#table-of-contents)

## ServiceSpec

ServiceSpec defines the desired state of a Service

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| tenantHeader | HTTP header to determine tenant for remote write requests. | string | false |
| defaultTenantId | Default tenant ID to use when none is provided via a header. | string | false |
| tenantLabelName | Label name through which the tenant will be announced. | string | false |
| storage |  | *[ObjectReference](#objectreference) | false |
| gateway | Gateway to proxy and auth requests to Query and Router. | *[Gateway](#gateway) | false |
| query | Query component querys from the backends such as Ingester and Store by automated discovery. | *[Query](#query) | false |
| router | Receive Router component routes to the backends such as Ingester by automated discovery. | *[Router](#router) | false |
| queryFrontend | QueryFrontend component implements a service deployed in front of queriers to improve query parallelization and caching. | *[QueryFrontend](#queryfrontend) | false |

[Back to TOC](#table-of-contents)

## Store

Store is the Schema for the Store API

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| metadata |  | [metav1.ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.21/#objectmeta-v1-meta) | false |
| spec |  | [StoreSpec](#storespec) | false |
| status |  | [StoreStatus](#storestatus) | false |

[Back to TOC](#table-of-contents)

## StoreList

StoreList contains a list of Store

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| metadata |  | [metav1.ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.21/#listmeta-v1-meta) | false |
| items |  | [][Store](#store) | true |

[Back to TOC](#table-of-contents)

## StoreSpec



| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| affinity | If specified, the pod's scheduling constraints. | *corev1.Affinity | false |
| nodeSelector | Define which Nodes the Pods are scheduled on. | map[string]string | false |
| tolerations | If specified, the pod's tolerations. | []corev1.Toleration | false |
| resources | Define resources requests and limits for main container. | [corev1.ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.21/#resourcerequirements-v1-core) | false |
| replicas | Number of replicas for a component | *int32 | false |
| image | Image is the image with tag/version | string | false |
| imagePullPolicy | Image pull policy. One of Always, Never, IfNotPresent. Defaults to Always if :latest tag is specified, or IfNotPresent otherwise. Cannot be updated. | [corev1.PullPolicy](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.21/#container-v1-core) | false |
| logLevel | Log filtering level. Possible options: error, warn, info, debug | string | false |
| logFormat | Log format to use. Possible options: logfmt or json | string | false |
| storage |  | *[ObjectReference](#objectreference) | false |
| minTime | MinTime specifies start of time range limit to serve | string | false |
| maxTime | MaxTime specifies end of time range limit to serve | string | false |
| indexCacheConfig | IndexCacheConfig contains index cache configuration. | *[IndexCacheConfig](#indexcacheconfig) | false |
| flags | Flags is the flag of store. | []string | false |
| dataVolume | DataVolume specifies how volume shall be used | *[KubernetesVolume](#kubernetesvolume) | false |
| scaler |  | *[AutoScaler](#autoscaler) | false |

[Back to TOC](#table-of-contents)

## S3

Config stores the configuration for s3 bucket.

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| bucket |  | string | true |
| endpoint |  | string | true |
| region |  | string | false |
| awsSdkAuth |  | bool | false |
| accessKey |  | *corev1.SecretKeySelector | true |
| insecure |  | bool | false |
| signatureVersion2 |  | bool | false |
| secretKey |  | *corev1.SecretKeySelector | true |
| putUserMetadata |  | map[string]string | false |
| httpConfig |  | [S3HTTPConfig](#s3httpconfig) | false |
| trace |  | [S3TraceConfig](#s3traceconfig) | false |
| listObjectsVersion |  | string | false |
| partSize | PartSize used for multipart upload. Only used if uploaded object size is known and larger than configured PartSize. NOTE we need to make sure this number does not produce more parts than 10 000. | uint64 | false |
| sseConfig |  | [S3SSEConfig](#s3sseconfig) | false |
| stsEndpoint |  | string | false |

[Back to TOC](#table-of-contents)

## S3HTTPConfig

S3HTTPConfig stores the http.Transport configuration for the s3 minio client.

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| idleConnTimeout |  | model.Duration | false |
| responseHeaderTimeout |  | model.Duration | false |
| insecureSkipVerify |  | bool | false |
| tlsHandshakeTimeout |  | model.Duration | false |
| expectContinueTimeout |  | model.Duration | false |
| maxIdleConns |  | int | false |
| maxIdleConnsPerHost |  | int | false |
| maxConnsPerHost |  | int | false |
| tlsConfig |  | [TLSConfig](#tlsconfig) | false |

[Back to TOC](#table-of-contents)

## S3SSEConfig

S3SSEConfig deals with the configuration of SSE for Minio. The following options are valid: kmsencryptioncontext == https://docs.aws.amazon.com/kms/latest/developerguide/services-s3.html#s3-encryption-context

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| type |  | string | false |
| kmsKeyId |  | string | false |
| kmsEncryptionContext |  | map[string]string | false |
| encryptionKey |  | string | false |

[Back to TOC](#table-of-contents)

## S3TraceConfig



| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| enable |  | bool | false |

[Back to TOC](#table-of-contents)

## Storage



| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| metadata |  | [metav1.ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.21/#objectmeta-v1-meta) | false |
| spec |  | [StorageSpec](#storagespec) | false |
| status |  | [StorageStatus](#storagestatus) | false |

[Back to TOC](#table-of-contents)

## StorageList



| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| metadata |  | [metav1.ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.21/#listmeta-v1-meta) | false |
| items |  | [][Storage](#storage) | true |

[Back to TOC](#table-of-contents)

## StorageSpec



| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| S3 |  | *[S3](#s3) | false |

[Back to TOC](#table-of-contents)

## TLSConfig

TLSConfig configures the options for TLS connections.

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| ca | The secret that including the CA cert. | *corev1.SecretKeySelector | false |
| cert | The secret that including the client cert. | *corev1.SecretKeySelector | false |
| key | The secret that including the client key. | *corev1.SecretKeySelector | false |
| serverName | Used to verify the hostname for the targets. | string | false |
| insecureSkipVerify | Disable target certificate validation. | bool | false |

[Back to TOC](#table-of-contents)

## Tenant

Tenant is the Schema for the monitoring Tenant API

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| metadata |  | [metav1.ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.21/#objectmeta-v1-meta) | false |
| spec |  | [TenantSpec](#tenantspec) | false |
| status |  | [TenantStatus](#tenantstatus) | false |

[Back to TOC](#table-of-contents)

## TenantList



| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| metadata |  | [metav1.ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.21/#listmeta-v1-meta) | false |
| items |  | [][Tenant](#tenant) | true |

[Back to TOC](#table-of-contents)

## TenantSpec



| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| tenant |  | string | false |
| storage |  | *[ObjectReference](#objectreference) | false |

[Back to TOC](#table-of-contents)

## TenantStatus



| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| ruler |  | *[ObjectReference](#objectreference) | false |
| compactor |  | *[ObjectReference](#objectreference) | false |
| ingester |  | *[ObjectReference](#objectreference) | false |

[Back to TOC](#table-of-contents)

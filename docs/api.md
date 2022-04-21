<br>
# API Docs

This Document documents the types introduced by the paodin-monitoring to be consumed by users.

> Note this document is generated from code comments. When contributing a change to this document please do so by changing the code comments.

## Table of Contents
* [EnvoySpec](#envoyspec)
* [Gateway](#gateway)
* [KubernetesVolume](#kubernetesvolume)
* [ObjectReference](#objectreference)
* [Query](#query)
* [QueryStores](#querystores)
* [Retention](#retention)
* [Service](#service)
* [ServiceList](#servicelist)
* [ServiceSpec](#servicespec)
* [Store](#store)
* [StoreList](#storelist)
* [StoreSpec](#storespec)
* [Thanos](#thanos)
* [ThanosCompact](#thanoscompact)
* [ThanosReceiveIngestor](#thanosreceiveingestor)
* [ThanosReceiveIngestorList](#thanosreceiveingestorlist)
* [ThanosReceiveIngestorSpec](#thanosreceiveingestorspec)
* [ThanosReceiveRouter](#thanosreceiverouter)
* [ThanosStore](#thanosstore)
* [ThanosStoreGateway](#thanosstoregateway)

## EnvoySpec

EnvoySpec defines the desired state of envoy proxy sidecar which delegates requests to the secure thanos stores

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| image | Image is the thanos image with tag/version | string | false |
| resources | Define resources requests and limits for envoy container. | [corev1.ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.21/#resourcerequirements-v1-core) | false |

[Back to TOC](#table-of-contents)

## Gateway



| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| affinity | If specified, the pod's scheduling constraints. | *corev1.Affinity | false |
| nodeSelector | Define which Nodes the Pods are scheduled on. | map[string]string | false |
| tolerations | If specified, the pod's tolerations. | []corev1.Toleration | false |
| resources | Define resources requests and limits for main container. | [corev1.ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.21/#resourcerequirements-v1-core) | false |
| replicas | Number of replicas for a thanos component | *int32 | false |
| image | Image is the gateway image with tag/version. | string | false |
| logLevel | Log filtering level. Possible options: error, warn, info, debug. | string | false |
| logFormat | Log format to use. Possible options: logfmt or json. | string | false |
| serverCertificate | Secret name for HTTP Server certificate (Kubernetes TLS secret type) | string | false |
| clientCaCertificate | Secret name for HTTP Client CA certificate (Kubernetes TLS secret type) | string | false |

[Back to TOC](#table-of-contents)

## KubernetesVolume

KubernetesVolume defines the configured volume for a thanos instance.

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
| replicas | Number of replicas for a thanos component | *int32 | false |
| image | Image is the thanos image with tag/version | string | false |
| logLevel | Log filtering level. Possible options: error, warn, info, debug | string | false |
| logFormat | Log format to use. Possible options: logfmt or json | string | false |
| stores | Additional StoreApi servers from which Thanos Query component queries from | [][QueryStores](#querystores) | false |
| selectorLabels | Selector labels that will be exposed in info endpoint. | map[string]string | false |
| replicaLabelNames | Labels to treat as a replica indicator along which data is deduplicated. | []string | false |
| envoy | Envoy is used to config sidecar which proxies requests requiring auth to the secure stores | [EnvoySpec](#envoyspec) | false |

[Back to TOC](#table-of-contents)

## QueryStores



| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| addresses | Address is the addresses of StoreApi server, which may be prefixed with 'dns+' or 'dnssrv+' to detect StoreAPI servers through respective DNS lookups. For more info, see https://thanos.io/tip/thanos/service-discovery.md/#dns-service-discovery | []string | false |
| caSecret | Secret containing the CA cert to use for StoreApi connections | *corev1.SecretKeySelector | false |

[Back to TOC](#table-of-contents)

## Retention

Retention defines the config for retaining samples

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| retentionRaw | RetentionRaw specifies how long to retain raw samples in bucket | string | false |
| retention5m | Retention5m specifies how long to retain samples of 5m resolution in bucket | string | false |
| retention1h | Retention1h specifies how long to retain samples of 1h resolution in bucket | string | false |

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
| gateway | Gateway to proxy and auth requests to Thanos Query and Thanos Receive Router defined in Thanos. | *[Gateway](#gateway) | false |
| thanos | Thanos cluster contains explicit Thanos Query and Thanos Receive Router, and implicit Thanos Receive Ingestor and Thanos Store Gateway and Thanos Compact which are selected by label selector `monitoring.paodin.io/service=<service_namespace>.<service_name>`. | *[Thanos](#thanos) | false |

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

StoreSpec defines the desired state of a Store

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| objectStorageConfig | ObjectStorageConfig allows specifying a key of a Secret containing object store configuration | *corev1.SecretKeySelector | false |
| thanos | Thanos contains Thanos Store Gateway and Thanos Compact. | *[ThanosStore](#thanosstore) | false |

[Back to TOC](#table-of-contents)

## Thanos



| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| query | Thanos Query component querys from the backends such as Thanos Receive Ingestor and Thanos Store Gateway by automated discovery. | *[Query](#query) | false |
| receiveRouter | Thanos Receive Router component routes to the backends such as Thanos Receive Ingestor by automated discovery. | *[ThanosReceiveRouter](#thanosreceiverouter) | false |

[Back to TOC](#table-of-contents)

## ThanosCompact



| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| affinity | If specified, the pod's scheduling constraints. | *corev1.Affinity | false |
| nodeSelector | Define which Nodes the Pods are scheduled on. | map[string]string | false |
| tolerations | If specified, the pod's tolerations. | []corev1.Toleration | false |
| resources | Define resources requests and limits for main container. | [corev1.ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.21/#resourcerequirements-v1-core) | false |
| replicas | Number of replicas for a thanos component | *int32 | false |
| image | Image is the thanos image with tag/version | string | false |
| logLevel | Log filtering level. Possible options: error, warn, info, debug | string | false |
| logFormat | Log format to use. Possible options: logfmt or json | string | false |
| downsamplingDisable | DownsamplingDisable specifies whether to disable downsampling | *bool | false |
| retention | Retention configs how long to retain samples | *[Retention](#retention) | false |
| dataVolume | DataVolume specifies how volume shall be used | *[KubernetesVolume](#kubernetesvolume) | false |

[Back to TOC](#table-of-contents)

## ThanosReceiveIngestor

ThanosReceiveIngestor is the Schema for the ThanosReceiveIngestor API

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| metadata |  | [metav1.ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.21/#objectmeta-v1-meta) | false |
| spec |  | [ThanosReceiveIngestorSpec](#thanosreceiveingestorspec) | false |
| status |  | [ThanosReceiveIngestorStatus](#thanosreceiveingestorstatus) | false |

[Back to TOC](#table-of-contents)

## ThanosReceiveIngestorList

ThanosReceiveIngestorList contains a list of Store

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| metadata |  | [metav1.ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.21/#listmeta-v1-meta) | false |
| items |  | [][ThanosReceiveIngestor](#thanosreceiveingestor) | true |

[Back to TOC](#table-of-contents)

## ThanosReceiveIngestorSpec

ThanosReceiveIngestorSpec defines the desired state of a ThanosReceiveIngestor

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| tenants | Tenants if not empty indicates current config is for hard tenants; otherwise, it is for soft tenants. | []string | false |
| affinity | If specified, the pod's scheduling constraints. | *corev1.Affinity | false |
| nodeSelector | Define which Nodes the Pods are scheduled on. | map[string]string | false |
| tolerations | If specified, the pod's tolerations. | []corev1.Toleration | false |
| resources | Define resources requests and limits for main container. | [corev1.ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.21/#resourcerequirements-v1-core) | false |
| replicas | Number of replicas for a thanos component. | *int32 | false |
| image | Image is the thanos image with tag/version | string | false |
| logLevel | Log filtering level. Possible options: error, warn, info, debug | string | false |
| logFormat | Log format to use. Possible options: logfmt or json | string | false |
| localTsdbRetention | LocalTsdbRetention configs how long to retain raw samples on local storage. | string | false |
| longTermStore | If specified, the object key of Store for long term storage. | *[ObjectReference](#objectreference) | false |
| dataVolume | DataVolume specifies how volume shall be used | *[KubernetesVolume](#kubernetesvolume) | false |

[Back to TOC](#table-of-contents)

## ThanosReceiveRouter



| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| affinity | If specified, the pod's scheduling constraints. | *corev1.Affinity | false |
| nodeSelector | Define which Nodes the Pods are scheduled on. | map[string]string | false |
| tolerations | If specified, the pod's tolerations. | []corev1.Toleration | false |
| resources | Define resources requests and limits for main container. | [corev1.ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.21/#resourcerequirements-v1-core) | false |
| replicas | Number of replicas for a thanos component. | *int32 | false |
| image | Image is the thanos image with tag/version | string | false |
| logLevel | Log filtering level. Possible options: error, warn, info, debug | string | false |
| logFormat | Log format to use. Possible options: logfmt or json | string | false |
| replicationFactor | How many times to replicate incoming write requests | *uint64 | false |

[Back to TOC](#table-of-contents)

## ThanosStore



| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| storeGateway | Thanos Store Gateway will be selected as query backends by Service. | *[ThanosStoreGateway](#thanosstoregateway) | false |
| compact | Thanos Compact as object storage data compactor and lifecycle manager. | *[ThanosCompact](#thanoscompact) | false |

[Back to TOC](#table-of-contents)

## ThanosStoreGateway



| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| affinity | If specified, the pod's scheduling constraints. | *corev1.Affinity | false |
| nodeSelector | Define which Nodes the Pods are scheduled on. | map[string]string | false |
| tolerations | If specified, the pod's tolerations. | []corev1.Toleration | false |
| resources | Define resources requests and limits for main container. | [corev1.ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.21/#resourcerequirements-v1-core) | false |
| replicas | Number of replicas for a thanos component | *int32 | false |
| image | Image is the thanos image with tag/version | string | false |
| logLevel | Log filtering level. Possible options: error, warn, info, debug | string | false |
| logFormat | Log format to use. Possible options: logfmt or json | string | false |
| minTime | MinTime specifies start of time range limit to serve | string | false |
| maxTime | MaxTime specifies end of time range limit to serve | string | false |
| dataVolume | DataVolume specifies how volume shall be used | *[KubernetesVolume](#kubernetesvolume) | false |

[Back to TOC](#table-of-contents)

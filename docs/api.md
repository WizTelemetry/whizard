<br>
# API Docs

This Document documents the types introduced by the paodin-monitoring to be consumed by users.

> Note this document is generated from code comments. When contributing a change to this document please do so by changing the code comments.

## Table of Contents
* [CommonThanosFields](#commonthanosfields)
* [Compact](#compact)
* [EnvoySpec](#envoyspec)
* [KubernetesVolume](#kubernetesvolume)
* [Query](#query)
* [QueryStores](#querystores)
* [Receive](#receive)
* [ReceiveIngestor](#receiveingestor)
* [ReceiveRouter](#receiverouter)
* [Retention](#retention)
* [StoreGateway](#storegateway)
* [Thanos](#thanos)
* [ThanosList](#thanoslist)
* [ThanosSpec](#thanosspec)

## CommonThanosFields



| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| image | Image is the thanos image with tag/version | string | false |
| logLevel | Log filtering level. Possible options: error, warn, info, debug | string | false |
| logFormat | Log format to use. Possible options: logfmt or json | string | false |

[Back to TOC](#table-of-contents)

## Compact



| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| affinity | If specified, the pod's scheduling constraints. | *corev1.Affinity | false |
| nodeSelector | Define which Nodes the Pods are scheduled on. | map[string]string | false |
| tolerations | If specified, the pod's tolerations. | []corev1.Toleration | false |
| resources | Define resources requests and limits for main container. | [corev1.ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.21/#resourcerequirements-v1-core) | false |
| replicas | Number of replicas for a thanos component | *int32 | false |
| downsamplingDisable | DownsamplingDisable specifies whether to disable downsampling | *bool | false |
| retention | Retention configs how long to retain samples | *[Retention](#retention) | false |
| objectStorageConfig | ObjectStorageConfig allows specifying a key of a Secret containing object store configuration | *corev1.SecretKeySelector | false |
| dataVolume | DataVolume specifies how volume shall be used | *[KubernetesVolume](#kubernetesvolume) | false |

[Back to TOC](#table-of-contents)

## EnvoySpec

EnvoySpec defines the desired state of envoy proxy sidecar which delegates requests to the secure thanos stores

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| image | Image is the thanos image with tag/version | string | false |
| resources | Define resources requests and limits for envoy container. | [corev1.ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.21/#resourcerequirements-v1-core) | false |

[Back to TOC](#table-of-contents)

## KubernetesVolume

KubernetesVolume defines the configured volume for a thanos instance.

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| emptyDir |  | *corev1.EmptyDirVolumeSource | false |
| pvc |  | *corev1.PersistentVolumeClaim | false |

[Back to TOC](#table-of-contents)

## Query



| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| affinity | If specified, the pod's scheduling constraints. | *corev1.Affinity | false |
| nodeSelector | Define which Nodes the Pods are scheduled on. | map[string]string | false |
| tolerations | If specified, the pod's tolerations. | []corev1.Toleration | false |
| resources | Define resources requests and limits for main container. | [corev1.ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.21/#resourcerequirements-v1-core) | false |
| replicas | Number of replicas for a thanos component | *int32 | false |
| stores | Additional StoreApi servers from which Thanos Query component queries from | [][QueryStores](#querystores) | false |
| selectorLabels | Selector labels that will be exposed in info endpoint. | map[string]string | false |
| envoy | Envoy is used to config sidecar which proxies requests requiring auth to the secure stores | [EnvoySpec](#envoyspec) | false |

[Back to TOC](#table-of-contents)

## QueryStores



| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| addresses | Address is the addresses of StoreApi server, which may be prefixed with 'dns+' or 'dnssrv+' to detect StoreAPI servers through respective DNS lookups. For more info, see https://thanos.io/tip/thanos/service-discovery.md/#dns-service-discovery | []string | false |
| caSecret | Secret containing the CA cert to use for StoreApi connections | *corev1.SecretKeySelector | false |

[Back to TOC](#table-of-contents)

## Receive



| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| router |  | [ReceiveRouter](#receiverouter) | false |
| ingestors |  | [][ReceiveIngestor](#receiveingestor) | false |

[Back to TOC](#table-of-contents)

## ReceiveIngestor



| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| name | Ingestor name must be unique within current thanos cluster, which follows the regulation for k8s resource name. | string | false |
| tenants | Tenants if not empty indicates current config is for hard tenants; otherwise, it is for soft tenants. | []string | false |
| affinity | If specified, the pod's scheduling constraints. | *corev1.Affinity | false |
| nodeSelector | Define which Nodes the Pods are scheduled on. | map[string]string | false |
| tolerations | If specified, the pod's tolerations. | []corev1.Toleration | false |
| resources | Define resources requests and limits for main container. | [corev1.ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.21/#resourcerequirements-v1-core) | false |
| replicas | Number of replicas for a thanos component | *int32 | false |
| localTsdbRetention | LocalTsdbRetention configs how long to retain raw samples on local storage | string | false |
| objectStorageConfig | ObjectStorageConfig allows specifying a key of a Secret containing object store configuration | *corev1.SecretKeySelector | false |
| dataVolume | DataVolume specifies how volume shall be used | *[KubernetesVolume](#kubernetesvolume) | false |

[Back to TOC](#table-of-contents)

## ReceiveRouter



| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| affinity | If specified, the pod's scheduling constraints. | *corev1.Affinity | false |
| nodeSelector | Define which Nodes the Pods are scheduled on. | map[string]string | false |
| tolerations | If specified, the pod's tolerations. | []corev1.Toleration | false |
| resources | Define resources requests and limits for main container. | [corev1.ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.21/#resourcerequirements-v1-core) | false |
| replicas | Number of replicas for a thanos component | *int32 | false |
| replicationFactor |  | *uint64 | false |
| tenantHeader | HTTP header to determine tenant for remote write requests | string | false |
| defaultTenantId | Default tenant ID to use when none is provided via a header | string | false |
| tenantLabelName | Label name through which the tenant will be announced. | string | false |

[Back to TOC](#table-of-contents)

## Retention

Retention defines the config for retaining samples

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| retentionRaw | RetentionRaw specifies how long to retain raw samples in bucket | string | false |
| retention5m | Retention5m specifies how long to retain samples of 5m resolution in bucket | string | false |
| retention1h | Retention1h specifies how long to retain samples of 1h resolution in bucket | string | false |

[Back to TOC](#table-of-contents)

## StoreGateway



| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| affinity | If specified, the pod's scheduling constraints. | *corev1.Affinity | false |
| nodeSelector | Define which Nodes the Pods are scheduled on. | map[string]string | false |
| tolerations | If specified, the pod's tolerations. | []corev1.Toleration | false |
| resources | Define resources requests and limits for main container. | [corev1.ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.21/#resourcerequirements-v1-core) | false |
| replicas | Number of replicas for a thanos component | *int32 | false |
| minTime | MinTime specifies start of time range limit to serve | string | false |
| maxTime | MaxTime specifies end of time range limit to serve | string | false |
| objectStorageConfig | ObjectStorageConfig allows specifying a key of a Secret containing object store configuration | *corev1.SecretKeySelector | false |
| dataVolume | DataVolume specifies how volume shall be used | *[KubernetesVolume](#kubernetesvolume) | false |

[Back to TOC](#table-of-contents)

## Thanos

Thanos is the Schema for the thanos API

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| metadata |  | [metav1.ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.21/#objectmeta-v1-meta) | false |
| spec |  | [ThanosSpec](#thanosspec) | false |
| status |  | [ThanosStatus](#thanosstatus) | false |

[Back to TOC](#table-of-contents)

## ThanosList

ThanosList contains a list of Thanos

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| metadata |  | [metav1.ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.21/#listmeta-v1-meta) | false |
| items |  | [][Thanos](#thanos) | true |

[Back to TOC](#table-of-contents)

## ThanosSpec

ThanosSpec defines the desired state of a Thanos cluster

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| defaultFields |  | [CommonThanosFields](#commonthanosfields) | false |
| query |  | *[Query](#query) | false |
| receive |  | *[Receive](#receive) | false |
| storeGateway |  | *[StoreGateway](#storegateway) | false |
| compact |  | *[Compact](#compact) | false |

[Back to TOC](#table-of-contents)

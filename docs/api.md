<br>
# API Docs

This Document documents the types introduced by the paodin-monitoring to be consumed by users.

> Note this document is generated from code comments. When contributing a change to this document please do so by changing the code comments.

## Table of Contents
* [EnvoySpec](#envoyspec)
* [IngressSpec](#ingressspec)
* [QueryStore](#querystore)
* [ThanosQuery](#thanosquery)
* [ThanosQueryList](#thanosquerylist)
* [ThanosQuerySpec](#thanosqueryspec)
* [KubernetesVolume](#kubernetesvolume)
* [ReceiveIngestorSpec](#receiveingestorspec)
* [ReceiveRouterSpec](#receiverouterspec)
* [RouterHashringConfig](#routerhashringconfig)
* [ThanosReceive](#thanosreceive)
* [ThanosReceiveList](#thanosreceivelist)
* [ThanosReceiveSpec](#thanosreceivespec)
* [ThanosCompactSpec](#thanoscompactspec)
* [ThanosStorage](#thanosstorage)
* [ThanosStorageList](#thanosstoragelist)
* [ThanosStorageRetention](#thanosstorageretention)
* [ThanosStorageSpec](#thanosstoragespec)
* [ThanosStoreGatewaySpec](#thanosstoregatewayspec)

## EnvoySpec

EnvoySpec defines the desired state of envoy proxy sidecar which delegates requests to the secure thanos stores

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| image | Image is the thanos image with tag/version | string | false |
| imagePullPolicy |  | [corev1.PullPolicy](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.21/#container-v1-core) | false |
| resources | Define resources requests and limits for envoy container. | [corev1.ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.21/#resourcerequirements-v1-core) | false |

[Back to TOC](#table-of-contents)

## IngressSpec



| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| host |  | string | false |
| path |  | string | false |
| secretName |  | string | false |

[Back to TOC](#table-of-contents)

## QueryStore



| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| address | Address is the address of a store api server, which may be prefixed with 'dns+' or 'dnssrv+' to detect store API servers through respective DNS lookups. For more info, see https://thanos.io/tip/thanos/service-discovery.md/#dns-service-discovery | string | false |
| secretName |  | string | false |

[Back to TOC](#table-of-contents)

## ThanosQuery

ThanosQuery is the Schema for the thanosqueries API

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| metadata |  | [metav1.ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.21/#objectmeta-v1-meta) | false |
| spec |  | [ThanosQuerySpec](#thanosqueryspec) | false |
| status |  | [ThanosQueryStatus](#thanosquerystatus) | false |

[Back to TOC](#table-of-contents)

## ThanosQueryList

ThanosQueryList contains a list of ThanosQuery

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| metadata |  | [metav1.ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.21/#listmeta-v1-meta) | false |
| items |  | [][ThanosQuery](#thanosquery) | true |

[Back to TOC](#table-of-contents)

## ThanosQuerySpec

ThanosQuerySpec defines the desired state of ThanosQuery

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| affinity | If specified, the pod's scheduling constraints. | *corev1.Affinity | false |
| nodeSelector | Define which Nodes the Pods are scheduled on. | map[string]string | false |
| tolerations | If specified, the pod's tolerations. | []corev1.Toleration | false |
| resources | Define resources requests and limits for main container. | [corev1.ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.21/#resourcerequirements-v1-core) | false |
| image | Image is the thanos image with tag/version | string | false |
| imagePullPolicy |  | [corev1.PullPolicy](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.21/#container-v1-core) | false |
| replicas | Number of replicas for a thanos query component | *int32 | false |
| stores | Stores config store api servers from where series are queried | [][QueryStore](#querystore) | false |
| selectorLabels | SelectorLabels config query selector labels that will be exposed in info endpoint | map[string]string | false |
| httpIngress | HttpIngress configs http request entry from services outside the cluster | *[IngressSpec](#ingressspec) | false |
| grpcIngress | HttpIngress configs grpc request entry from services outside the cluster | *[IngressSpec](#ingressspec) | false |
| level | LogLevel configs log filtering level. Possible options: error, warn, info, debug | string | false |
| format | LogFormat configs log format to use. Possible options: logfmt or json | string | false |
| envoy | Envoy is used to config envoy sidecar which delegates requests to the secure stores | *[EnvoySpec](#envoyspec) | false |

[Back to TOC](#table-of-contents)

## KubernetesVolume

KubernetesVolume defines the configured volume for a thanos receiver.

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| emptyDir |  | *corev1.EmptyDirVolumeSource | false |
| pvc |  | *corev1.PersistentVolumeClaim | false |

[Back to TOC](#table-of-contents)

## ReceiveIngestorSpec

ReceiveIngestorSpec defines the configs that thanos receive running in the ingestor mode requires

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| localTsdbRetention | LocalTSDBRetention configs how long to retain raw samples on local storage | string | false |
| objectStorageConfig | ObjectStorageConfig allows specifying a key of a Secret containing object store configuration | *corev1.SecretKeySelector | false |
| dataVolume | DataVolume specifies how volume shall be used | *[KubernetesVolume](#kubernetesvolume) | false |

[Back to TOC](#table-of-contents)

## ReceiveRouterSpec

ReceiveRouterSpec defines the configs that thanos receive running in the router mode requires

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| hashringsRefreshInterval | HashringsRefreshInterval configs refresh interval to re-read the hashring configuration file | string | false |
| hardTenantHashrings | HardTenantHashrings are hashrings with non-empty tenants which match the tenant in the request | []*[RouterHashringConfig](#routerhashringconfig) | false |
| softTenantHashring | SoftTenantHashring is a hashring with empty tenants which is used when the tenant in the request cannot be found in HardTenantHashrings. | *[RouterHashringConfig](#routerhashringconfig) | false |
| replicationFactor |  | *uint64 | false |
| remoteWriteIngress | RemoteWriteIngress configs remote write request entry from services outside the cluster | *[IngressSpec](#ingressspec) | false |

[Back to TOC](#table-of-contents)

## RouterHashringConfig

RouterHashringConfig defines the hashring config for a team of tenants

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| name |  | string | false |
| tenants | Tenants if not empty indicates current config is for hard tenants; otherwise, it is for soft tenants. | []string | false |
| endpoints | Endpoints are statically configured endpoints which receive requests for the specified tenants | []string | false |
| endpointsNamespaceSelector | EndpointsNamespaceSelector and EndpointsSelector select endpoints which receive requests for the specified tenants They only work when Endpoints is empty | *[metav1.LabelSelector](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.21/#labelselector-v1-meta) | false |
| endpointsSelector |  | *[metav1.LabelSelector](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.21/#labelselector-v1-meta) | false |

[Back to TOC](#table-of-contents)

## ThanosReceive

ThanosReceive is the Schema for the thanosreceives API

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| metadata |  | [metav1.ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.21/#objectmeta-v1-meta) | false |
| spec |  | [ThanosReceiveSpec](#thanosreceivespec) | false |
| status |  | [ThanosReceiveStatus](#thanosreceivestatus) | false |

[Back to TOC](#table-of-contents)

## ThanosReceiveList

ThanosReceiveList contains a list of ThanosReceive

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| metadata |  | [metav1.ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.21/#listmeta-v1-meta) | false |
| items |  | [][ThanosReceive](#thanosreceive) | true |

[Back to TOC](#table-of-contents)

## ThanosReceiveSpec

ThanosReceiveSpec defines the desired state of ThanosReceive

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| affinity | If specified, the pod's scheduling constraints. | *corev1.Affinity | false |
| nodeSelector | Define which Nodes the Pods are scheduled on. | map[string]string | false |
| tolerations | If specified, the pod's tolerations. | []corev1.Toleration | false |
| resources | Define resources requests and limits for single Pods. | [corev1.ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.21/#resourcerequirements-v1-core) | false |
| image | Image is the thanos image with tag/version | string | false |
| imagePullPolicy |  | [corev1.PullPolicy](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.21/#container-v1-core) | false |
| replicas | Number of replicas for a thanos receive component | *int32 | false |
| router | Router specifies the configs that thanos receive running in the router mode requires | *[ReceiveRouterSpec](#receiverouterspec) | false |
| ingestor | Ingestor specifies the configs that thanos receive running in the ingestor mode requires | *[ReceiveIngestorSpec](#receiveingestorspec) | false |
| tenantHeader | TenantHeader configs the HTTP header specifying the replica number of a write request to thanos receive | string | false |
| defaultTenantId | DefaultTenantId configs the default tenant ID to use when none is provided via a header | string | false |
| tenantLabelName | TenantLabelName configs the label name through which the tenant will be announced. | string | false |
| level | LogLevel configs log filtering level. Possible options: error, warn, info, debug | string | false |
| format | LogFormat configs log format to use. Possible options: logfmt or json | string | false |

[Back to TOC](#table-of-contents)

## ThanosCompactSpec

ThanosCompactSpec defines the desired state of thanos compact

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| resources | Define resources requests and limits for single Pods. | [corev1.ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.21/#resourcerequirements-v1-core) | false |
| downsamplingDisable | DownsamplingDisable specifies whether to disable downsampling | *bool | false |
| retention | Retention configs how long to retain samples | *[ThanosStorageRetention](#thanosstorageretention) | false |
| dataVolume | DataVolume specifies how volume shall be used | *[KubernetesVolume](#kubernetesvolume) | false |

[Back to TOC](#table-of-contents)

## ThanosStorage

ThanosStorage is the Schema for the thanosstorages API

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| metadata |  | [metav1.ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.21/#objectmeta-v1-meta) | false |
| spec |  | [ThanosStorageSpec](#thanosstoragespec) | false |
| status |  | [ThanosStorageStatus](#thanosstoragestatus) | false |

[Back to TOC](#table-of-contents)

## ThanosStorageList

ThanosStorageList contains a list of ThanosStorage

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| metadata |  | [metav1.ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.21/#listmeta-v1-meta) | false |
| items |  | [][ThanosStorage](#thanosstorage) | true |

[Back to TOC](#table-of-contents)

## ThanosStorageRetention

ThanosStorageRetention defines the config for retaining samples

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| retentionRaw | RetentionRaw specifies how long to retain raw samples in bucket | string | false |
| retention5m | Retention5m specifies how long to retain samples of 5m resolution in bucket | string | false |
| retention1h | Retention1h specifies how long to retain samples of 1h resolution in bucket | string | false |

[Back to TOC](#table-of-contents)

## ThanosStorageSpec

ThanosStorageSpec defines the desired state of ThanosStorage

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| affinity | If specified, the pod's scheduling constraints. | *corev1.Affinity | false |
| nodeSelector | Define which Nodes the Pods are scheduled on. | map[string]string | false |
| tolerations | If specified, the pod's tolerations. | []corev1.Toleration | false |
| image | Image is the thanos image with tag/version | string | false |
| imagePullPolicy |  | [corev1.PullPolicy](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.21/#container-v1-core) | false |
| objectStorageConfig | ObjectStorageConfig allows specifying a key of a Secret containing object store configuration | *corev1.SecretKeySelector | false |
| level | LogLevel configs log filtering level. Possible options: error, warn, info, debug | string | false |
| format | LogFormat configs log format to use. Possible options: logfmt or json | string | false |
| gateway | Gateway specifies the configs for thanos store gateway | *[ThanosStoreGatewaySpec](#thanosstoregatewayspec) | false |
| compact | Compact specifies the configs for thanos compact | *[ThanosCompactSpec](#thanoscompactspec) | false |

[Back to TOC](#table-of-contents)

## ThanosStoreGatewaySpec

ThanosStoreGatewaySpec defines the desired state of thanos store gateway

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| resources | Define resources requests and limits for single Pods. | [corev1.ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.21/#resourcerequirements-v1-core) | false |
| replicas |  | *int32 | false |
| dataVolume | DataVolume specifies how volume shall be used | *[KubernetesVolume](#kubernetesvolume) | false |
| minTime | MinTime specifies start of time range limit to serve | string | false |
| maxTime | MaxTime specifies end of time range limit to serve | string | false |

[Back to TOC](#table-of-contents)

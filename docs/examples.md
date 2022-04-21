# Host/Members Mode

Refer to the following example for multi-cluster monitoring of KubeSphere platform.

> Please firstly refer to [here](../README.md#quickstart) to install paodin-controller-manager.

## Scattered

1. On all clusters, configure Thanos Sidecar and external labels for Prometheus:

  ```shell
  kubectl -n kubesphere-monitoring-system patch prometheus k8s --patch='{"spec":{"externalLabels":{"cluster":"<cluster_name>"},"thanos":{}}}' --type=merge
  ```

2. On host cluster, deploy Thanos Query to proxy all Thanos Sidecar stores: 

  ```shell
  cat <<EOF | kubectl apply -f -
  apiVersion: monitoring.paodin.io/v1alpha1
  kind: Service
  metadata:
    name: scattered
    namespace: kubesphere-monitoring-system
  spec:
    stores:
    - addresses: 
      - prometheus-operated:10901
      - <member-prometheus-svc>:10901
  EOF
  ```


## Central

> Use object storage for long term storage.

1. On host cluster, create secret for object storage config: 

  ```shell
  cat <<EOF | kubectl apply -f -
  apiVersion: v1
  kind: Secret
  metadata:
    name: objectstorage
    namespace: kubesphere-monitoring-system
  type: Opaque
  data:
    thanos.yaml: |-
      type: s3
      config:
        bucket: thanos-storage
        region: sh1a
        endpoint: s3.sh1a.qingstor.com
        access_key: <access_key>
        secret_key: <secret_key>
  EOF
  ```

2. On host cluster, create a service instance: 

  ```shell
  cat <<EOF | kubectl apply -f -
  apiVersion: monitoring.paodin.io/v1alpha1
  kind: Service
  metadata:
    name: central
    namespace: kubesphere-monitoring-system
  spec:
    tenantHeader: cluster
    defaultTenantId: unknown
    tenantLabelName: cluster
    gateway: {}
    thanos: 
      query:
        replicaLabelNames:
        - prometheus_replica
        - receive_replica
        - thanos_ruler_replica
      receiveRouter: 
        replicationFactor: 2
  EOF
  ```

3. On host cluster, create a store instance for long term storage:

  ```shell
  cat <<EOF | kubectl apply -f -
  apiVersion: monitoring.paodin.io/v1alpha1
  kind: Store
  metadata:
    name: longterm
    namespace: kubesphere-monitoring-system
    labels: 
      monitoring.paodin.io/service: kubesphere-monitoring-system.central
  spec:
    objectStorageConfig: 
      name: objectstorage
      key: thanos.yaml
    thanos: 
      storeGateway: {}
      compact: {}
  EOF
  ```

4. On host cluster, create an ingestor instance. The soft tenants instance with empty tenants as follows can receive requests with all tenants:

  ```shell
  cat <<EOF | kubectl apply -f -
  apiVersion: monitoring.paodin.io/v1alpha1
  kind: ThanosReceiveIngestor
  metadata:
    name: softs
    namespace: kubesphere-monitoring-system
    labels: 
      monitoring.paodin.io/service: kubesphere-monitoring-system.central
  spec:
    tenants: []
    replicas: 2
    longTermStore: 
      namespace: kubesphere-monitoring-system
      name: longterm
  EOF
  ```

5. On all clusters, configure Prometheus to write to gateway:  

  ```shell
  kubectl -n kubesphere-monitoring-system patch prometheus k8s --patch='{"spec":{"remoteWrite":[{"url":"http://<gateway_address>:9090/<cluster_name>/api/v1/receive"}]}}' --type=merge
  ```

6. On all clusters, configure ks-apiserver to read from gateway:  

  update monitoring endpoint as follows by `kubectl -n kubesphere-system edit cm kubesphere-config`:   

  ```yaml
  ...
  data:
    kubesphere.yaml: |
      ...
      monitoring:
        endpoint: http://<gateway_address>:9090/<cluster_name>
        ...
      ...
  ...
  ```



The following briefly describes how to integrate PaodinMonitoring and configure related components in several multi-cluster monitoring modes. 

## host cluster performs cross-cluster query when data is scattered in each cluster

Each cluster stores data on its own cluster, and single-cluster bussiness still queries data from itself, and multi-cluster business may request all clusters.

Here are the pros and cons of this mode:  

- pros:  
    1. single-cluster bussiness query logic per cluster remains unchanged.  
    2. cross-cluster query only involves multi-cluster business data.  
    3. fewer new components, just a sidecar for Prometheus on each cluster, a Thanos Query and a Paodin Monitoring operator on host cluster.  

- cons:  
    1. require network connectivity from host cluster to member clusters.  
    2. each multi-cluster business query may require data transfer accross the clusters, unless a Query Frontend is introduced as a cache layer.    
    3. member clusters are still heavy. 

Follow the steps below to deploy and configure host cluster to query monitoring data accross all clusters:  

1. On each cluster, config a Thanos Sidecar for Prometheus to expose data query, as well as an external label to mark the cluster to which the data belongs.  

    ```shell
    kubectl -n kubesphere-monitoring-system patch prometheus k8s --patch='{"spec":{"externalLabels":{"cluster":"host"},"thanos":{}}}' --type=merge

    kubectl -n kubesphere-monitoring-system patch prometheus k8s --patch='{"spec":{"externalLabels":{"cluster":"member"},"thanos":{}}}' --type=merge
    ```
    > A service or ingress may be configured to make the Thanos Sidecar accessible to host cluster.  

2. On host cluster, deploy the operator of PaodinMonitoring, and then create a Thanos cluster only with Query to connect the Thanos Sidecar accross all clusters.  
    ```shell
    kubectl apply -f https://raw.githubusercontent.com/kubesphere/paodin-monitoring/master/config/bundle.yaml

    cat <<EOF | kubectl apply -f -
    apiVersion: monitoring.paodin.io/v1alpha1
    kind: Service
    metadata:
      name: t1
      namespace: kubesphere-monitoring-system
    spec:
      thanos
        query:
          stores:
          - addresses:
            - prometheus-operated:10901 # for host cluster
            - <member-thanos-sidecar-adress>:10901
    EOF
    ```

3. On host cluster, query data from multi clusters.

    ```shell
    curl --data "query=up{cluster='host|member'}" t1-query-operated.kubesphere-monitoring-system.svc:10902/api/v1/query
    ```

## host cluster queries multi-cluster data pushed by member clusters

Each cluster stores data on its own cluster, and single-cluster bussiness still queries data from itself. Member clusters push data involving multi-cluster business to host cluster through remote-write, and host cluster query these data while processing multi-cluster business.  

Here are the pros and cons of this mode:  

- pros:  
    1. single-cluster bussiness query logic per cluster remains unchanged.  
    2. multi-cluster business may only request local cluster data.  
    3. support object storage for long-term data.  

- cons:  
    1. require network connectivity from member clusters to host cluster.  
    2. member clusters are still heavy. 


1. On host cluster, deploy the operator of PaodinMonitoring, and then create a Thanos cluster with a receive router to receive data from member clusters, a receive ingestor to land data.    
    ```shell
    kubectl apply -f https://raw.githubusercontent.com/kubesphere/paodin-monitoring/master/config/bundle.yaml

    cat <<EOF | kubectl apply -f -
    apiVersion: monitoring.paodin.io/v1alpha1
    kind: Service
    metadata:
      name: t2
    spec:
      thanos: 
        query: 
          replicas: 2
        receive:
          router:
            replicas: 2
          ingestors:
          - name: softs
            replicas: 2
            dataVolume:
              pvc:
                spec:
                  resources:
                    requests:
                      storage: 5Gi
    EOF
    ```
    > A service or ingress may be configured to make the Thanos Receive router accessible to member clusters.  

2. On member clusters, configure Prometheus to push data through remote-write to Thanos Receive router in host cluster.
    ```shell
    kubectl -n kubesphere-monitoring-system patch prometheus k8s --patch='{"spec":{"remoteWrite":{"url":"<member-thanos-sidecar-adress>:19291","writeRelabelConfigs": []}}}' --type=merge
    ```
    > configure `writeRelabelConfigs` to only push data related to multi-cluster business.

## member clusters query its own data from host cluster

Member clusters push all data to host cluster and do not store data on its own cluster. Both single-cluster bussiness and multi-cluster bussiness have to query data from host cluster.  

Here are the pros and cons of this mode:  

- pros:  
    1. member cluster can be light weight.  
    2. multi-cluster business may only request local cluster data.  
    3. support object storage for long-term data.  

- cons:  
    1. require network connectivity from member clusters to host cluster. 
    2. each single-cluster business query on member clusters requires data transfer, unless a Query Frontend is introduced as a cache layer.  
    3. require additional process to restrict member clusters to only have access to its own data.  
    4. require host cluster highly avaiable and the data has to be replicated across multiple host clusters.  
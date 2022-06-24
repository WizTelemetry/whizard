# Paodin

Paodin is an observability platform for Kubernetes infrastructure and applications, which integrates your metrics, logs, events and so on accross kubernetes clusters.

## Status

| Component | Function | Status | Comment
|--------|-----------|--------|--------|
| *paodin-controller-manager*  | `Service` (CRD) |  | Define one Service instance, which contains *monitoring-gateway*, *Thanos Query* and *Thanos Receive Router* to handle and route metric read and write requests.
|| `ThanosReceiveIngestor` (CRD) |  | Define one *Thanos Receive Ingestor* instance which lands metric data.
|| `Store` (CRD) |  | Define one Store instace, which contains *Thanos Store Gateway* and *Thanos Compact* for long term storage.
|| `ThanosRuler` (CRD) | | Define one *Thanos Ruler* instance.
|| `AlertingRule` (CRD) | | Define one single alerting rule.
|| `RuleGroup` (CRD) | | Define one rule group.
| *monitoring-gateway* | Auth/Proxy for monitoring service |  |
| *monitoring-agent-proxy* | Proxy for prometheus agent
| *paodin-apiserver* |

## Install

- Install paodin-controller-manager:

    ```shell
    kubectl apply -f https://raw.githubusercontent.com/kubesphere/paodin/master/config/bundle.yaml
    ```

- See [here](./docs/monitoring/examples.md) to deploy your monitoring service accross clusters.
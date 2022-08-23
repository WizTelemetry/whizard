# Whizard

Whizard is an observability platform for Kubernetes infrastructure and applications, which integrates your metrics, logs, events and so on accross kubernetes clusters.

## Status

| Component | Function | Status | Comment
|--------|-----------|--------|--------|
| *whizard-controller-manager*  | `Service` (CRD) |  | Define one Service instance, which contains *monitoring-gateway*, *Query* and *Router* to handle and route metric read and write requests.
|| `Ingester` (CRD) |  | Define one *Receive Ingester* instance which lands metric data.
|| `Store` (CRD) |  | Define one Store instace for long term storage.
|| `Ruler` (CRD) | | Define one *Ruler* instance.
|| `AlertingRule` (CRD) | | Define one single alerting rule.
|| `RuleGroup` (CRD) | | Define one rule group.
| *monitoring-gateway* | Auth/Proxy for monitoring service |  |
| *monitoring-agent-proxy* | Proxy for prometheus agent
| *whizard-apiserver* |

## Install

- Install whizard-controller-manager:

    ```shell
    kubectl apply -f https://raw.githubusercontent.com/kubesphere/whizard/master/config/bundle.yaml
    ```

- See [here](./docs/monitoring/examples.md) to deploy your monitoring service accross clusters.
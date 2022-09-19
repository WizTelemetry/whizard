# Whizard

Whizard is a Cloud Native observability platform for Cloud Native and traditional infrastructure and applications.
Currently, Whizard supports managing massive metrics data from multiple tenants, and the support for managing logs and tracing data will be added in the future.

## Architecture

<div align=center><img src=docs/images/whizard.svg></div>

## Status

| Component | Function | Status | Comment
|--------|-----------|--------|--------|
| *whizard-controller-manager*  | `Service` (CRD) |  | Define one Service instance, which contains *monitoring-gateway*, *Query* and *Router* to handle and route metric read and write requests.
|| `Ingester` (CRD) |  | Define one *Receive Ingester* instance which lands metric data.
|| `Store` (CRD) |  | Define one Store instace for long term storage.
|| `Ruler` (CRD) | | Define one *Ruler* instance.
| *monitoring-gateway* | Auth/Proxy for monitoring service |  |
| *monitoring-agent-proxy* | Proxy for prometheus agent
| *whizard-apiserver* |

## Install

- Install whizard-controller-manager:

    ```shell
    kubectl apply -f https://raw.githubusercontent.com/kubesphere/whizard/master/config/bundle.yaml
    ```

- See [here](./docs/monitoring/examples.md) to deploy your monitoring service accross clusters.
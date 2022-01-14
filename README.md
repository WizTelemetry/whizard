# PaodinMonitoring

PaodinMonitoring provides an optimized multi-cluster monitoring solution based on prometheus(server/agent) and thanos.

## Design

![design](./docs/images/design.png "Multi-Cluster Monitoring Architecture")

The architecture diagram above covers a variety of scenarios for which PaodinMonitoring applies.

## CRDs

PaodinMonitoring contains an operator that acts on the following CRDs to deploy some key components: 

- `ThanosQuery`, which defines a desired thanos query deployment. The operator will inject an envoy sidecar container to it, which proxies query requests requiring auth to secure stores.
- `ThanosReceive`, which defines a desired thanos receive deployment.
- `ThanosStorage`, which defines two desired deployments: thanos store gateway, thanos compact. An object storage usually requires a thanos store gateway component to provide metrics query apis, and a thanos compact component to compact metrics block and do metrics lifecycle management, so define them in one CRD.

To learn more about them have a look at the [api doc](docs/api.md).

## Install

Install the CRDs and the PaodinMonitoring:

```shell
kubectl apply -f https://raw.githubusercontent.com/kubesphere/paodin-monitoring/master/config/bundle.yaml
```

## Usage

See [here](./docs/usage.md) to learn how to use PaodinMonitoring for multi-cluster monitoring.

## Roadmap

- [x] Define CRDs to configure and deploy thanos components: query, receive, store gateway, compact.
- [x] Add ingress configurations for http server and grpc server of thanos query, remote-write server of thanos receive
- [x] Add configuration of envoy sidecar which is used to proxy query requests to the secure stores.
- [x] Optional running mode configuration to the CRD for thanos receive: router, ingestor.
- [x] Automatically discovers thanos receive ingestor endpoints (by endpoint selector) for thanos receive router.
- [x] Configure soft and hard tenants separately and clearly for thanos receive router hashring.
- [ ] Add more configurable items to CRDs for configuring individual components flexibly.  
- [ ] Support configuration automatic reloading, mainly for secrets.  

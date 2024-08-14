# Whizard

## Overview

The Whizard provides [Kubernetes](https://kubernetes.io/) native deployment and management of [Thanos](https://thanos.io/) and related monitoring components. 

The Whizard includes, but is not limited to, the following features:

- **Cloud-Native Deployment and Operation**: All components support definition and maintenance via Custom Resource Definitions (CRDs), simplifying configuration and operation. This includes Thanos components (Router, Ingester, Compactor, Store, Query, QueryFrontend, Ruler) and Whizard-specific components (Service, Tenant, Storage).
- **Tenant-Based Automatic Horizontal Scaling**: Recognizing that Horizontal Pod Autoscalers (HPA) based on CPU and Memory may not meet the stability requirements of enterprise-level stateful workloads, Whizard introduces a tenant-based workload scaling mechanism. Components such as Ingester, Compactor, and Ruler scale horizontally with tenant creation and deletion, ensuring stable operation and providing tenant-level horizontal scaling and resource reclamation.
- **Support for Multi-Cluster Management in K8s**: To enhance monitoring and alerting for multi-cluster K8s environments, Whizard's maintainers developed the whizard-adapter. This tool automatically creates or deletes Whizard tenants based on the creation or deletion of K8s or KubeSphere clusters, triggering the automatic scaling of Thanos stateful workloads.

For an introduction to the Whizard, see the [getting started](https://whizardtelemetry.github.io/docs/whizard-docs/intro) guide.

## Architecture

<div align=center><img src=docs/images/whizard.svg></div>

## CustomResourceDefinitions

A core feature of the Whizard is to monitor the Kubernetes API server for changes
to specific objects and ensure that the current component deployments match these objects.
The Operator acts on the following [Custom Resource Definitions (CRDs)](https://kubernetes.io/docs/tasks/access-kubernetes-api/extend-api-custom-resource-definitions/):

- **Compactor**: Defines the Compactor component, which does the block compaction and life cycle management for the object storages.
- **Gateway**: Defines the Gateway component, which provides a unified entry point for metrics read and write requests.
- **Ingester**: Defines the Ingester component, which receives metrics data from routers, caches data in memory, flushes data to disk, and finally upload metrics blocks to object storage.
- **Query**: Defines the Query component, which fetches data from the ingesters and/or stores and then evaluates the query.
- **QueryFrontend**: Defines the Query Frontend component, which improves the query performance by request splitting and result cache.
- **Router**: Defines the Router component, which routes and replicates the metrics to the ingesters.
- **Ruler**: Defines the Ruler component, which evaluates recording and alerting rules.
- **Service**: Defines a Whizard Service, which connects different whizard components together to provide a complete monitoring service. It also contains shared configurations of different components.
- **Storage**: Defines the Storage instance, which contains the object storage configuration, and a block manager component for the block inspection and GC.
- **Store**: Defines the Store instance, which facilitates metrics reads from the object storage.
- **Tenant**: Defines a tenant which is the basic unit of resource isolation and auto-scaling.

The Whizard automatically detects changes in the Kubernetes API server to any of the above objects, and ensures that matching deployments and configurations are kept in sync.

## Quickstart

### Prerequisites

Kubernetes v1.19.0+ is necessary for running Whizard.

### Deploy Whizard with YAML

To quickly try out *just* the Whizard inside a cluster, **choose a release** and run the following command:

```sh
kubectl create --server-side -f config/bundle.yaml
```

### Deploy Whizard with Helm

> Note: For the helm based install, Helm v3.2.1 or higher is needed

```sh
helm install whizard --create-namespace -n kubesphere-monitoring-system charts/whizard/
```

## Contributing

We welcome contributions to this repository! If you have ideas for additional tests or would like to contribute code, please submit a pull request.

## Security

If you find a security vulnerability related to the Prometheus Operator, please do not report it by opening a GitHub issue, but instead please send an e-mail to the maintainers of the project found in the [MAINTAINERS.md](MAINTAINERS.md) file.

## Issues

If you encounter any issues while using Whizard or running the tests in this repository, please submit an issue on this repository.
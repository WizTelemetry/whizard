# v0.9.0-rc.0 / 2023-09-07

### FEATURES

* Gateway supports server-side basic_auth authentication configuration([#332](https://github.com/WhizardTelemetry/whizard/pull/332))
* Add persistentVolumeClaimRetentionPolicy config([#330](https://github.com/WhizardTelemetry/whizard/pull/330))
* Add securityContext ([#329](https://github.com/WhizardTelemetry/whizard/pull/329))
* Add genSignedCert in charts([#316](https://github.com/WhizardTelemetry/whizard/pull/316))
* Add time partition support([#311](https://github.com/WhizardTelemetry/whizard/pull/311))
* Gateway supports tenant access control([#310](https://github.com/WhizardTelemetry/whizard/pull/310))
* Gateway enhancement and refactor([#310](https://github.com/WhizardTelemetry/whizard/pull/310))

### ENHANCEMENT

* Upgrade thanos version to v0.32.2([#335](https://github.com/WhizardTelemetry/whizard/pull/335))
* Fix dependabot security alerts([#329](https://github.com/WhizardTelemetry/whizard/pull/329))
* Upgrade go version to 1.21([#306](https://github.com/WhizardTelemetry/whizard/pull/329))

# v0.8.0 / 2023-08-07

### FEATURES

* Support external remote write and query ([#290](https://github.com/WhizardTelemetry/whizard/pull/290))

### ENHANCEMENT

* Gateway: Replace tenant matcher of query param for query api ([#297](https://github.com/WhizardTelemetry/whizard/pull/297))
* Optimize ingester data cleanup when a tenant is deleted ([#291](https://github.com/WhizardTelemetry/whizard/pull/291))
* Optimize transport parameters for agent proxy component ([#286](https://github.com/WhizardTelemetry/whizard/pull/286))

### BUGFIX

* Fix values config in the chart ([#292](https://github.com/WhizardTelemetry/whizard/pull/292))

# v0.8.0-rc.0 / 2023-07-27

### FEATURES

* Support external remote write and query ([#290](https://github.com/WhizardTelemetry/whizard/pull/290))

### ENHANCEMENT

* Optimize ingester data cleanup when a tenant is deleted ([#291](https://github.com/WhizardTelemetry/whizard/pull/291))
* Optimize transport parameters for agent proxy component ([#286](https://github.com/WhizardTelemetry/whizard/pull/286))

### BUGFIX

* Fix values config in the chart ([#292](https://github.com/WhizardTelemetry/whizard/pull/292))

# v0.7.0 / 2023-06-30

### FEATURES

* Allow to override --query flag for global ruler to query external data sources([#277](https://github.com/WhizardTelemetry/whizard/pull/277))

### BUGFIX

* ruler watches router([#276](https://github.com/WhizardTelemetry/whizard/pull/276))
* remove alpn_protocols in envoy config([#275](https://github.com/WhizardTelemetry/whizard/pull/275))


# v0.7.0-rc.0 / 2023-06-21

### FEATURES

* Components support https configuration([#264](https://github.com/WhizardTelemetry/whizard/pull/264))
* Gateway supports tls configuration for downstream services([#263](https://github.com/WhizardTelemetry/whizard/pull/263))

### BUGFIX

* Fix ruler name conflicts ([#265](https://github.com/WhizardTelemetry/whizard/pull/265))
* Fixed ruler name being too long([#250](https://github.com/WhizardTelemetry/whizard/pull/250))
* Optionally deploy Store HPA([#246](https://github.com/WhizardTelemetry/whizard/pull/246))
* Degrade Thanos Query to v0.30.2([#252](https://github.com/WhizardTelemetry/whizard/pull/252))

### ENHANCEMENT

* Add cherry pick action in CI([#247](https://github.com/WhizardTelemetry/whizard/pull/247))
* Add `--tsdb.out-of-order.time-window=10m` flag to ingester([#252](https://github.com/WhizardTelemetry/whizard/pull/252))

# v0.6.2 / 2023-05-12
### Features

* Support to configure imagePullSecrets for private registry([#241](https://github.com/WhizardTelemetry/whizard/pull/241))


# v0.6.1 / 2023-04-21
### BUGFIX

* Donot copy all labels of custom resources to managed workloads to fix that managed workloads cannot be upgraded([#230](https://github.com/WhizardTelemetry/whizard/pull/230))
* Fix object storage config in chart([#231](https://github.com/WhizardTelemetry/whizard/pull/231))

# v0.6.0 / 2023-04-14

### Features

* Optimize tenant data cleaning on ingester([#182](https://github.com/WhizardTelemetry/whizard/pull/182))
* Add tenant selector in store([#171](https://github.com/WhizardTelemetry/whizard/pull/171))
* Allow tenants to monopolize resources([#170](https://github.com/WhizardTelemetry/whizard/pull/170))

### ENHANCEMENT

* Allows global configuration to update compactor.retention([#186](https://github.com/WhizardTelemetry/whizard/pull/186))
* Adjust ingester retention period([#185](https://github.com/WhizardTelemetry/whizard/pull/185))
* Upgrade dependencies([#188](https://github.com/WhizardTelemetry/whizard/pull/188), [#157](https://github.com/WhizardTelemetry/whizard/pull/157), [#208](https://github.com/WhizardTelemetry/whizard/pull/208))
* Update charts([#187](https://github.com/WhizardTelemetry/whizard/pull/187), [#162](https://github.com/WhizardTelemetry/whizard/pull/162))
* Support mutil-arch image build([#136](https://github.com/WhizardTelemetry/whizard/pull/136))
* Upgrade Thanos to v0.31.0([#208](https://github.com/WhizardTelemetry/whizard/pull/208))
* Update Ruler to query from QueryFrontend with a higher performance than Query([#211](https://github.com/WhizardTelemetry/whizard/pull/211))

### BUGFIX

* Fix tls secret volume mount in storage component([#183](https://github.com/WhizardTelemetry/whizard/pull/183))
* Fix mapstructure decode bugs and add config unit test([#137](https://github.com/WhizardTelemetry/whizard/pull/137))
* Fix override method bugs and add options unit test([#160](https://github.com/WhizardTelemetry/whizard/pull/160))
* Fix resources config parse and override error ([#208](https://github.com/WhizardTelemetry/whizard/pull/208))

# 0.6.0-rc.2 / 2023-04-03

### ENHANCEMENT

* Update the ruler to query from the QueryFrontend with a higher performance than the Query([#211](https://github.com/WhizardTelemetry/whizard/pull/211))

### BUGFIX

* Update go4.org/unsafe/assume-no-moving-gc to fix the build failed to run with go1.20([#211](https://github.com/WhizardTelemetry/whizard/pull/211))

# 0.6.0-rc.1 / 2023-03-24

### ENHANCEMENT

* Upgrade Thanos to v0.31.0([#208](https://github.com/WhizardTelemetry/whizard/pull/208))
* Upgrade dependencies([#208](https://github.com/WhizardTelemetry/whizard/pull/208))

### BUGFIX

* Fix some bugs([#201](https://github.com/WhizardTelemetry/whizard/pull/201))



# 0.6.0-rc.0 / 2023-03-08

### Features

* Optimize tenant data cleaning on ingester([#182](https://github.com/WhizardTelemetry/whizard/pull/182))
* Add tenant selector in store([#171](https://github.com/WhizardTelemetry/whizard/pull/171))
* Allow tenants to monopolize resources([#170](https://github.com/WhizardTelemetry/whizard/pull/170))

### ENHANCEMENT

* Allows global configuration to update compactor.retention([#186](https://github.com/WhizardTelemetry/whizard/pull/186))
* Adjust ingester retention period([#185](https://github.com/WhizardTelemetry/whizard/pull/185))
* Upgrade dependencies([#188](https://github.com/WhizardTelemetry/whizard/pull/188), [#157](https://github.com/WhizardTelemetry/whizard/pull/157))
* Update charts([#187](https://github.com/WhizardTelemetry/whizard/pull/187), [#162](https://github.com/WhizardTelemetry/whizard/pull/162))
* Support mutil-arch image build([#136](https://github.com/WhizardTelemetry/whizard/pull/136))

### BUGFIX

* Fix tls secret volume mount in storage component([#183](https://github.com/WhizardTelemetry/whizard/pull/183))
* Fix mapstructure decode bugs and add config unit test([#137](https://github.com/WhizardTelemetry/whizard/pull/137))
* Fix override method bugs and add options unit test([#160](https://github.com/WhizardTelemetry/whizard/pull/160))



# 0.5.0-rc.0 / 2022-09-29

This is the first release of whizard.

## What's new

Whizard is a distributed cloud observability platform that provides unified observability (currently monitoring and alerting) for Multi-Cloud, On-Premise and Edge infrastructures and applications. 

### Kubernetes native deployment and management

The Whizard Controller Manager simplifies and automates the configuration and deployment of the whizard components by the following CRDs:  

- `Compactor`: Defines the Compactor component, which does the block compaction and life cycle management for the object storages.
- `Gateway`: Defines the Gateway component, which provides an unified entrypoints for metrics read and write requests. 
- `Ingester`: Defines the Ingester component, which receives metrics data from routers, caches data in memory, flushs data to disk, and finally uploads metrics blocks to object storage.
- `Query`: Defines the Query component, which fetches data from the ingesters and/or stores and then evaluates the query.
- `QueryFrontend`: Defines the Query Frontend component, which improves the query performance by request spliting and result caching.
- `Router`: Defines the Router component, which routes and replicates the metrics to the ingesters. 
- `Ruler`: Defines the Ruler component, which evaluates recording and alerting rules.
- `Service`: Defines a Whizard Service, which connects different whizard components together to provide a complete monitoring service. It also contains shared configurations of different components. 
- `Storage`: Defines the Storage instance, which contains the object storage configuration, and a block manager component for the block inspection and GC.
- `Store`: Defines the Store instance, which facilitates metrics reads from the object storage.
- `Tenant`: Defines a tenant which is the basic unit of resource isolation and auto scaling.  

### Multi-tenancy and Auto scaling

- Whizard components support multi-tenancy and are able to auto scale. 
- The store component supports to auto scale based on its actual load. 
- The ruler component also can also scale based on rule group sharding for a single tenant with too many rules.  

### Data management and GC

- Whizard provides metrics data life cycle management for data on disk or in object storage. If a tenant is deleted, Whizard can automatically cleanup this tenant's blocks in the object storage or on local disk.

### Others

- Whizard also has an agent proxy component that implements the Prometheus HTTP v1 API (reads/writes), which can be used as a data collection agent and a query proxy.


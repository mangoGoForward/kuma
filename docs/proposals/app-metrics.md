# Application metrics
 
## Context
 
Currently, when an application exposes metrics as regular service it's protected and traffic is secured (mTLS enabled), it's a bit tricky
to [expose those metrics without mTLS](https://kuma.io/docs/dev/policies/traffic-metrics/#expose-metrics-from-applications). Also, by configuring paths user's can by mistake expose more than it's necessary. To avoid it we should come up with a solution that fulfills all the requirements: scraping metrics with/without mTLS and safer exposing of metrics.
 
## Requirements
* support different application ports/paths
* support more than one container/application
* support no mTLS communication with Prometheus
 
## Idea
 
Instead of exposing application metrics directly, we propose creating a metrics aggregator in `kuma-dp`. The aggregator will be responsible for the scraping of metrics from applications proxied by this instance of `kuma-dp`.
`kuma-dp` can then expose this port with or without mTLS. It's worth mentioning that this idea can support only Prometheus endpoint scraping and exposing metrics as Prometheus.

### Mesh level configuration

We should allow configuring also metrics from which endpoints should be aggregated at the mesh level. Most of the environments might have similar applications exposing metrics at the same port/path. That would speed up configuration and allow operators to expose metrics easier.
```
type: Mesh
name: default
metrics:
 enabledBackend: prometheus-1
 backends:
 - name: prometheus-1
   type: prometheus
   conf:
    skipMTLS: false
    port: 5670 # default exposed port with metrics 
    path: /metrics
    tags:   
      kuma.io/service: dataplane-metrics
    aggregate:
     - name: app
       port: 1236
       path: /metrics
       default: true
     - name: opa
       port: 12345
       path: /stats
       default: false
```

If the default will be `true` then we take the mesh configuration and add the one that was defined by the user during deployment. If the value will be `false`, then we take the one provided by the user one override it. Name is going to be a key to distinguishing applications configuration. In case user provide different name then we are going to append list.
This might be a good solution for environments that have common services running within the cluster and the operator doesn't want users to override this configuration.

 
### Universal
At the mesh level, it will be possible to define the default exposed path and port at which metrics are going to be accessible.
```
type: Mesh
name: default
metrics:
 enabledBackend: prometheus-1
 backends:
 - name: prometheus-1
   type: prometheus
   conf:
    skipMTLS: false
    port: 5670 # default exposed port with metrics 
    path: /metrics
    tags:   
      kuma.io/service: dataplane-metrics
    aggregate:
     - name: app
       port: 1236
       path: /metrics
       default: false
     - name: opa
       port: 12345
       path: /stats
       default: true
```
`Dataplane` configuration enables at which port/path `kuma-dp` is going to expose aggregated metrics from Envoy and other applications running within one workflow. Also, it's possible to override the `Mesh` configuration of endpoints to scrap by the `Dataplane` configuration. They are going to be identified by `name` and if `Mesh` and `Dataplane` configurations have definitions with the same name, then the `Dataplane` configuration has precedence before `Mesh` and is merged.
```
type: Dataplane
mesh: default
name: redis-1
networking:
 address: 192.168.0.1
 inbound:
 - port: 6379
   servicePort: 1234
   tags:
     kuma.io/service: redis
 metrics:
   type: prometheus
   conf:
     skipMTLS: true
     port: 1234
     path: /non-standard-path
     aggregate:
     - name: app
       port: 1237
       path: /metrics/prom
       # at dataplane level we don't to define `default` <- only mesh level, so we are going to validate it
```

Finally `kuma-dp` is going to scrape metrics from `app`, defined at dataplane level, and `opa` defined at `Mesh` level.
If user will try to override `opa` configuration we are going to throw an error, so below example won't work.
```
type: Dataplane
mesh: default
name: redis-1
networking:
 address: 192.168.0.1
 inbound:
 - port: 6379
   servicePort: 1234
   tags:
     kuma.io/service: redis
 metrics:
   type: prometheus
   conf:
     skipMTLS: true
     port: 1234
     path: /non-standard-path
     aggregate:
     - name: opa # it's default so cannot override
       port: 1238
       path: /metrics/prom
```
### Kubernetes
 
The same `Mesh` configuration works for K8s and Universal.
```
apiVersion: kuma.io/v1alpha1
kind: Mesh
metadata:
 name: default
spec:
 metrics:
   enabledBackend: prometheus-1
   backends:
   - name: prometheus-1
     type: prometheus
     conf:
       skipMTLS: false
       port: 5670
       path: /metrics
       tags: # tags that can be referred in Traffic Permission when metrics are secured by mTLS 
         kuma.io/service: dataplane-metrics
       aggregate:
       - name: app
         port: 1236
         path: /metrics
         default: false
       - name: opa
         port: 12345
         path: /stats
         default: true
```
In the case of Kubernetes configuration of endpoints to scrape is done by label `prometheus.metrics.kuma.io/aggregate-{name}-port` and `prometheus.metrics.kuma.io/aggregate-{name}-path`, where `{name}` is used for distinguish between config for different containers/apps running in the pod. The merge of configuration at the `Mesh` level and `Dataplane` level, like in universal, is going to append or override depends on configuration of the Mesh . 
 
```
apiVersion: apps/v1
kind: Deployment
metadata:
 name: example-app
 namespace: kuma-example
spec:
 ...
 template:
   metadata:
     ...
     labels:
       # indicate to Kuma that this Pod doesn't need a sidecar
       kuma.io/sidecar-injection: enabled
     annotations:
       prometheus.metrics.kuma.io/aggregate-app-port: 1234 # format is aggregate-{name}-port
       prometheus.metrics.kuma.io/aggregate-app-path: "/metrics"
       prometheus.metrics.kuma.io/aggregate-opa-port: 12345 # because opa is default we cannot override it
       prometheus.metrics.kuma.io/aggregate-opa-path: "/stats"
   spec:
     containers:
       ...
```
 
## Use existing prometheus tags

Prometheus on K8s uses annotation `prometheus.io/{path/port/scrapping}`, we might use it to map them for kuma specifc 
annotations that are going to be used to scrap endpoints. At the begging we might not implement this feature but in future we might get this.

## Labels

We are not going to change any metric's label so stats from applications are going to be unchanged.

### How to get information about path/port for `kuma-dp` to expose metrics?
 
 * Bootstrap config is going to return information about path/part. This configuration is going to be static and won't be possible to change it during application runtime.

## What about applications that don't have common applications?

We should allow users to exclude applications from scraping some endpoints. In this case, we might provide the `Mesh` level configuration that allows excluding some applications from appending metrics configuration for those services.
```
       aggregate:
       - name: app
         port: 1236
         path: /metrics
         default: false
       - name: opa
         port: 12345
         path: /stats
         default: true
         excludedServices: [service-1, service-2]
```
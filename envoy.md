# A Short Introduction to Envoy Proxy

[Envoy Proxy](https://www.envoyproxy.io/) is an open-source proxy server (it fulfils the same role as HAProxy, Nginx, and similar tools). Envoy is highly flexible and configurable via its plugin architecture.

Envoy was developed [at Lyft](https://mattklein123.dev/2021/09/14/5-years-envoy-oss/) and open-sourced in 2016.
From the outset it was intended to have a role in both observability and in mitigating and avoiding load related issues.

## Config

Envoy always has a static configuration file, which can be YAML or JSON (in this workshop we will be using 
`envoy.yaml` as the configuration file). Envoy supports having most of its configuration supplied dynamically via 
a control plane - almost everything except the configuration for the control plane itself - so things can be 
changed centrally without needing to restart Envoy. However, we won't be using
dynamic configuration here, for simplicity's sake and to make the configurations clear.

## Concepts and Terminology

### Downstream and Upstream

Some terminology you should know is `downstream` and `upstream`.
Downstreams are basically clients, which initiate connections to servers (upstreams).
However, the use of the terms downstream and upstream is preferred, as many upstreams are also downstreams of something else: the concept should always be seen as specific to the context of a given Envoy configuration.

In this workshop we will only have one layer of downstreams and upstreams in each demo, for simplicity.

### Listeners, filters, clusters

`Listeners` listen for connections on some port. When a connection is received, it will be processed through
a chain of `filters`, which can modify requests in various ways.

Most Envoy configuration relates to filter configuration, and most of the demos in this workshop are primarily
about configuring filters.

`Clusters` manage sets of upstreams, plus loadbalancing configurations.

### Envoy plugins

Envoy is [designed to be extensible](https://www.envoyproxy.io/docs/envoy/latest/extending/extending).
Envoy is written in C++, as are many of its plugins. 
However, extensions (specifically, HTTP filters) exist to permit plugins to be written in other languages, such as [Go](https://www.envoyproxy.io/docs/envoy/latest/configuration/http/http_filters/golang_filter) and [Lua](https://www.envoyproxy.io/docs/envoy/latest/start/sandboxes/lua).


## Envoy admin interface

When running any demo in this workshop you will find Envoy's admin interface and status page at
[http://localhost:9901](http://localhost:9901).

See documentation on the admin interface [here](https://www.envoyproxy.io/docs/envoy/latest/operations/admin)

You can try this out by running the [plain Envoy](./plain-envoy/README.md) non-demo.

## Envoy statistics

When running any demo in this workshop you will find Envoy's exported metrics at
[http://localhost:9901/stats/prometheus](http://localhost:9901/stats/prometheus).
It can be useful to explore them here so you know what's available for graphing in Prometheus or Grafana.

You can try this out by running the [plain Envoy](./plain-envoy/README.md) non-demo.

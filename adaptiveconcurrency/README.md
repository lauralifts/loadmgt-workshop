# Envoy Adaptive Concurrency

## What it is and how it works

Adaptive concurrency is similar to circuitbreaking, but rather than requiring the operator to determine
how many requests may be outstanding to the upstreams, it dynamically calculates this based
on observed latency. This allows a system to adapt to changes in performance and capacity.
It also reduces the burden on operators to maintain the correct configuration settings as workloads change.

Adaptive concurrency is calculated on a per-host basis, so it is sensitive to variations in capacity between 
different hosts (unlike circuitreaking). This can be useful if you have a heterogenous fleet.

### The minRTT measurement process

Envoy uses a component called the Gradient Controller to determine the number of outstanding requests which should
be permitted.

Operators can configure the maximum permissible concurrency (per upstream) and how often the concurrency should be recalculated.

What Envoy periodically does is to limit the concurrency to a subset of the upstreams in a cluster, and measure the 
latency of the upstreams under these conditions. This measurement should represent the performance of the upstreams
when they are not under excessive load. The minRTT then is based on a configurable percentile value of that 
set of measurements. The number of requests sampled depends on the request rate. 

The [Envoy adaptive concurrency](https://www.envoyproxy.io/docs/envoy/latest/configuration/http/http_filters/adaptive_concurrency_filter)
documentation describes the calculation of the per-upstream concurrent request limit. 
The calculation is designed to iteratively increase the concurrency limit when the current performance of an upstream is close to 
the ideal minRTT value, and to reduce it when performance is slower.

The gradient controller configuration includes a jitter parameter. This is useful to make sure that the minRTT calculation -
and thus the severe limiting of concurrency - does not happen for all upstreams at the same time. Along with this, it is recommended
to enable the [previous hosts retry predicate](https://www.envoyproxy.io/docs/envoy/latest/intro/arch_overview/http/http_connection_management#arch-overview-http-retry-plugins) when using retries, so as to avoid requests repeatedly being 
allocated to upstreams which are recalculating minRTT.

### Avoiding including healthchecks in minRTT calculation

Healthchecks should be excluded from the minRTT calculation - they are often much faster than other kinds of requests.
This should be achieved by using a different listener and filter chain for the healthcheck endpoint. 


## Limitations and Gotchas

 * The adaptive concurrency filter is a HTTP filter, so it can't be used with non-HTTP servers
 * There is no connection-oriented version of adaptive concurrency (as it needs to operate on request latency)
 * Unlike circuitbreaking, it does not support multiple priorities (i.e. circuitbreaking's default and high priorities)
 * Healthchecks should be excluded from the minRTT calculation - they are often much faster than other kinds of requests
 * The adaptive concurrency filter needs to be able to limit the concurrency for a cluster - this means that the cluster should not be receiving requests from other sources
 * As with most other Envoy load management mechanisms, there is no coordination between Envoy instances - so, combined with the above point, this implies that adaptive concurrency may not work very well when the fleet of Envoys is large

## Configuration

https://www.envoyproxy.io/docs/envoy/latest/configuration/http/http_filters/adaptive_concurrency_filter

Adaptive concurrency is not enabled by default: you must enable the filter for it to become active.
However, when adaptive concurrency is in use, the `max_concurrency_limit` defaults to 1000. This might 
be too low for some large-scale systems.

## Demos TODO

 * [Adaptive concurrency demo](./demo-ac/README.md)

## Links


https://www.alibabacloud.com/blog/brief-analysis-of-envoy-adaptive-concurrency-filter_600658
https://klaviyo.tech/adaptive-concurrency-control-for-mixed-analytical-workloads-51350439aeec


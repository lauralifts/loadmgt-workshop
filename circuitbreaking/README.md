# Circuitbreaking with Envoy

## What it is and how it works

The canonical [circuit breaker](https://en.wikipedia.org/wiki/Circuit_breaker_design_pattern) design pattern is 
intended to prevent cascading failures occurring in software, by breaking the vicious cycle where client services 
of an overloaded upstream service continue to send requests (often more requests, due to retries) 
and upstreams are thus prevented from working through their load backlog and becoming responsive again.
Traditionally, circuit breakers are either open or closed, based on the downstream service's perception of the 
upstream service's recent performance (closed means traffic flows).

[Envoy circuit breaking](https://www.envoyproxy.io/docs/envoy/latest/intro/arch_overview/upstream/circuit_breaking) allows you to 
do circuit breaking based on any of the following parameters:
 * Max concurrent outstanding requests
 * Max concurrent connections
 * Maximum number of outstanding automatic retries

## Limitations and Gotchas

Envoy circuit breaking is fully-distributed; i.e. each Envoy instance makes its own local decisions without
coordination with other Envoys. It is therefore  most suitable for situations where there is a small homogeneous 
pool of proxies fronting a service, rather than a large heterogeneous deployment of sidecars.

Envoy's implementation of circuitbreaking uses eventually-consistent state (shared between worker threads).
As a result of this, connection limits can be exceeded temporarily.

Circuitbreakers are enabled by default: see [disabling circuit breaking](https://www.envoyproxy.io/docs/envoy/latest/faq/load_balancing/disable_circuit_breaking#faq-disable-circuit-breaking) for information including default values (which may be too small for large-scale installation). 

Envoy's circuitbreaking doesn't attempt to be fair or have any QoS mechanism. 
If two clients send load to an upstream, overloading it and triggering circuitbreaking, both clients 
will have the same proportion of their requests throttled, even if one client is sending the vast majority of the load - a noisy neighbour can seriously degrade the service for all.

Envoy circuitbreaking does always not work well for custom protocols which multiplex multiple requests 
onto one connection, because Envoy has no visibility into the individual requests. In some cases 
connection-oriented circuitbreaking will be sufficient, but it depends on the request profile.


## Configuration

See 
https://www.envoyproxy.io/docs/envoy/latest/api-v3/config/cluster/v3/circuit_breaker.proto#envoy-v3-api-msg-config-cluster-v3-circuitbreakers


## Demos 
 * [Basic request-based and connection based circuitbreaking](./demo-basic/README.md)
 * [Circuitbreaking on retries](./demo-retries/README.md)
 * [Prioritising requests](./demo-prios/README.md)

TODO: demo with multiple upstreams, to demonstrate clarify what is per cluster and what is per host

## Useful Links

https://blog.turbinelabs.io/circuit-breaking-da855a96a61d
https://www.envoyproxy.io/docs/envoy/latest/intro/arch_overview/upstream/circuit_breaking

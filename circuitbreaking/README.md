# Circuitbreaking with Envoy

## What it is and how it works

TODO rephrase this more elegantly
The canonical [circuit breaker](https://en.wikipedia.org/wiki/Circuit_breaker_design_pattern) design pattern is 
intended to prevent cascading failures occurring in software, by breaking the cycle where client services of an
overloaded upstream service continue to send requests (often more requests, due to retries) and upstreams are
thus prevented from working through their load backlog and becoming responsive again.
Traditionally, circuit breakers are either open or closed, based on the downstream service's perception of the 
upstream service's recent performance. 

[Envoy circuit breaking](https://www.envoyproxy.io/docs/envoy/latest/intro/arch_overview/upstream/circuit_breaking) is a little different to the archetypal circuit breaking described above. 
Envoy's implementation allows you to control any of the following parameters:
 * Max concurrent requests
 * Max concurrent connections
 * Maximum number of automatic retries


todo routes and priorities

## Limitations and Gotchas

Envoy circuit breaking is fully-distributed; i.e. each Envoy instance makes its own local decisions without
coordination with other Envoys. It is probably most suitable for situations where there is a small homogeneous 
pool of proxies fronting a service, rather than a large heterogeneous deployment of sidecars.
TODO a diagram

Ofc can use dedicated Envoy pools per-client-pool as a way of doing QoS.


Envoy's implementation of circuitbreaking uses eventually-consistent state (shared between worker threads).
As a result of this, connection limits can be exceeded temporarily.

Circuitbreakers are enabled by default. 
The default values are TODO
https://www.envoyproxy.io/docs/envoy/latest/faq/load_balancing/disable_circuit_breaking#faq-disable-circuit-breaking

Envoy's circuitbreaking doesn't attempt to be fair or have any QoS mechanism. 
If two clients send load to an upstream, overloading it and triggering circuitbreaking, both clients 
will have the same proportion of their requests throttled, even if one client is sending the vast majority of the load - a noisy neighbour can break the service for all.

Finally, Envoy circuitbreaking does not work well for custom protocols which multiplex multiple requests onto one
connection. todo explain why

## Configuration

https://www.envoyproxy.io/docs/envoy/latest/api-v3/config/cluster/v3/circuit_breaker.proto#envoy-v3-api-msg-config-cluster-v3-circuitbreakers

priorities, retry budgets

## Demos TODO

 * [Basic request-based and connection based circuitbreaking](./demo-basic/README.md)
 * [Circuitbreaking on retries](./demo-retries/README.md)

 * priorities
 * noisy neighbours

## Useful Links
todo prune 

https://blog.turbinelabs.io/circuit-breaking-da855a96a61d
https://github.com/envoyproxy/envoy/issues/14550
https://stackoverflow.com/questions/60162863/envoy-circuit-breaking-non-deterministic-behaviour
https://www.envoyproxy.io/docs/envoy/latest/intro/arch_overview/upstream/circuit_breaking

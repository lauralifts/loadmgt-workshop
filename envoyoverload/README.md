
# Envoy Overload Manager

## What it is and how it works

The other sections in this workshop focus on how Envoy can protect upstream services from overload.
However, Envoy itself can also become overloaded. Envoy's [Overload Manager](https://www.envoyproxy.io/docs/envoy/latest/intro/arch_overview/operations/overload_manager) is designed to help
protect Envoy from overload. 

Currently, extensions exist to trigger protection measures based on utilisation of CPU (linux only), downstream connections, 
heap, as well as a synthetic resource (which reads from a file) which may be used for testing.

There are a few key concepts involved in Envoy overload management:
 * Resource Monitors are Envoy extensions which monitor utilisation of a system resource. They always result in a utilisation value which is a float between 0 and 1.
 * Overload Actions are things which Envoy can do to limit resource usage, such as to stop accepting connections or requests, to disable HTTP keepalives, to reduce various timeouts, and reset expensive streams
 * Loadshed points are similar to Overload Actions, but allow load to be shed at very specific points in the connection or request lifecycle
 * Triggers are configuration which connect Resource Monitors to Overload Actions or Loadshed Points. Triggers can be either scaled or threshold based, so load shedding may either be applied gradually as utilisation increases above a threshold, or all at once when a threshold is reached.

## Limitations and Gotchas

Several parts of the overload manager bear warnings such as `This extension is functional but has not had substantial production burn time, use only with this caveat`.

It is a good idea to exempt your healthchecks and probes from overload manager loadshedding - see the second demo below.
Envoy [Listener configs](https://www.envoyproxy.io/docs/envoy/latest/api-v3/config/listener/v3/listener.proto) now include a 
`bypass_overload_manager` flag for this purpose.

The Envoy Overload Manager doesn't do any prioritisation of requests or connections to be throttled (beyond anything you do using listener configs and bypass).

## Configuration

https://www.envoyproxy.io/docs/envoy/latest/configuration/operations/overload_manager/overload_manager

https://www.envoyproxy.io/docs/envoy/latest/api-v3/config/overload/v3/overload.proto#extension-category-envoy-resource-monitors

By default the Envoy overload manager is not active - it must be specifically enabled in configuration.


## Demos 

 * [Max Downstream Connections](./demo-maxconn/README.md)
 * [Synthetic Resource based request shedding and bypassing the overload manager](./demo-ir/README.md)

## Links

https://www.envoyproxy.io/docs/envoy/latest/intro/arch_overview/operations/overload_manager


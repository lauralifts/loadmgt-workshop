
# Envoy Overload Manager

## What it is and how it works

The other sections in this workshop focus on how Envoy can protect upstream services from overload.
However, Envoy itself can also become overloaded. Envoy's Overload Manager is designed to help
protect Envoy from overload. 

Currently, extensions exist to trigger protection measures based on utilisation of CPU, downstream connections, 
heap, as well as a synthetic resource (which reads from a file) which may be used for testing, or to implement 
protection based on other kinds of resources for which no Envoy monitor exists. 

There are a few key concepts involved in Envoy overload management:
 * Resource Monitors are Envoy extensions which monitor utilisation of a system resource. They always result in a utilisation value which is a float between 0 and 1.
 * Overload Actions are things which Envoy can do to limit resource usage, such as to stop accepting connections or requests, to disable HTTP keepalives, to reduce various timeouts, and reset expensive streams
 * Loadshed points are similar to Overload Actions, but allow load to be shed at very specific points in the connection or request lifecycle
 * Triggers are configuration which connect Resource Monitors to Overload Actions or Loadshed Points. Triggers can be either scaled or threshold based, so load shedding may either be applied gradually as utilisation increases, or all at once when a threshold is reached.

## Limitations and Gotchas

https://github.com/envoyproxy/envoy/issues/23843
Healthchecks - newish bypass overload manager flag - https://www.envoyproxy.io/docs/envoy/latest/api-v3/config/listener/v3/listener.proto

Several parts of the overload manager bear warnings such as `This extension is functional but has not had substantial production burn time, use only with this caveat.`

No prioritisation.

## Configuration

https://www.envoyproxy.io/docs/envoy/latest/configuration/operations/overload_manager/overload_manager

https://www.envoyproxy.io/docs/envoy/latest/api-v3/config/overload/v3/overload.proto#extension-category-envoy-resource-monitors

What is the default behaviour 

## Demos TODO

inc stats and/or the header sent
 x-envoy-local-overloaded

Demo bypassing overload manager

## Links

https://www.envoyproxy.io/docs/envoy/latest/intro/arch_overview/operations/overload_manager


# Envoy admission control 

## What it is and how it works

Envoy admission control is a form of client-side throttling, as described in the [Handling Overload](https://sre.google/sre-book/handling-overload/) section of the SRE Book.

Admission control is a good fit when you have a big fleet of Envoys, particularly ones running colocated with downstreams, as in a service mesh. It is also useful when you want to be sensitive to increases in request failure rate, rather than latency increases.

Unlike the adaptive concurrency filter, the observation of requests is totally passive - the use of this filter
does not cause Envoy to modify the flow of traffic to the upstreams in any way.

## Limitations and Gotchas

The filter is a HTTP filter, so it can't be used with non-HTTP protocols.

As with circuit breaking, there's no concept of request weight, which can matter when costs to service requests may vary.

Unlike circuit breaking, there isn't a concept of different priority requests - it isn't a natural fit with the error-rate driven logic of admission control, anyway. 

## Configuration

[Admission Control Config](https://www.envoyproxy.io/docs/envoy/latest/api-v3/extensions/filters/http/admission_control/v3/admission_control.proto#envoy-v3-api-msg-extensions-filters-http-admission-control-v3-admissioncontrol)

## Demos

See [demo](demo/README.md).

## Links

* [Admission Control Filter](https://www.envoyproxy.io/docs/envoy/latest/configuration/http/http_filters/admission_control_filter)
* [Handling Overload](https://sre.google/sre-book/handling-overload/)


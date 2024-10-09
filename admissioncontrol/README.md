# Envoy admission control 

## What it is and how it works

Envoy admission control is a form of client-side throttling, as described in the [Handling Overload](https://sre.google/sre-book/handling-overload/) section of the SRE Book.

## Limitations and Gotchas

The filter is a HTTP filter, so it can't be used with non-HTTP protocols.

As with circuit breaking, there's no concept of request weight, which can matter when costs to service requests may vary.

Unlike circuit breaking, there isn't a concept of priority (todo: can this be combined with CB? does that make sense?)

This one is good when you have a big fleet of client sidecar envoys, and when you want to be sensitive to increases in request failure rate.

TODO think this one does failure only, can't react to increased latenc

## Configuration

https://www.envoyproxy.io/docs/envoy/latest/api-v3/extensions/filters/http/admission_control/v3/admission_control.proto#envoy-v3-api-msg-extensions-filters-http-admission-control-v3-admissioncontrol

What is the default behaviour - don't think this one is active by default, you have to add the filter

## Demos TODO

* Plain client circuit breaking
* Different levels of aggression
* Combined with circuitbreaking? 
* enabled flag - operating as passthru (presumably still with metrics)

## Links

[Handling Overload](https://sre.google/sre-book/handling-overload/)

https://www.envoyproxy.io/docs/envoy/latest/configuration/http/http_filters/admission_control_filter


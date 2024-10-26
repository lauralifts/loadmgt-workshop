# Envoy Isolation and Ratelimiting

In all of the sections we have seen to date, load management has been something of a blunt tool.
Traffic can be differentiated only by endpoint (for example, avoiding throttling HTTP traffic by using
a separate filter chain) or by two priority levels in the case of circuit breaking. 

In many situations, it is useful to be able to limit traffic in other ways.
In particular, it is very helpful to be able to apply per-tenant rate limits, as a form of isolation.
This means that a bad actor, or a tenant with runaway automation, cannot adversely impact other tenants.

## What it is and how it works

Where every other example in this workshop is completely decentralised - i.e. each Envoy instance
is making local decisions about how to route traffic - isolation and ratelimiting across a service
does require some centralised coordination.  

[Envoy Ratelimit](https://github.com/envoyproxy/ratelimit) is a service, written in Golang and backed by 
Redis, which provides centralised ratelimiting services which can be used with Envoy.
It makes the decision to ratelimit - or not ratelimit - on a per-request or per-connection basis.
It is well integrated into the Envoy ecosystem - for instance, it uses similar configuration mechanisms.
And, of course, it supports the ratelimit APIs which Envoy Proxy uses.

## Auth API 

We often want to ratelimit based on application-level concepts (such as tenant ID, or subscription levels, 
for example). However, this data is probably not on your incoming requests. A common pattern is to use the 
[auth](https://www.envoyproxy.io/docs/envoy/latest/intro/arch_overview/security/ext_authz_filter.html) filter
to add that data to the requests, prior to useing the ratelimit service. 

For example, it is often necessary to look up metadata about tenant IDs based on bearer tokens.
In general, any data that we do not trust the client to send is a good candidate for this pattern.

## Local ratelimiting 

Envoy also supports [local ratelimiting](https://www.envoyproxy.io/docs/envoy/latest/intro/arch_overview/other_features/local_rate_limiting#arch-overview-local-rate-limit). This can be used independently, 
or can be used with centralised ratelimiting (to limit the load on the centralised ratelimit service).

## Limitations and Gotchas

Envoy also has the ability to perform decentralised quota-based ratelimiting, where each Envoy accepts
an assigned quota configuration from a rate limiting service - see the [Envoy docs](https://www.envoyproxy.io/docs/envoy/latest/intro/arch_overview/other_features/global_rate_limiting#quota-based-rate-limiting).
However, there is currently no open-source implementation of the rate limiting control plane.

## Configuration

https://www.envoyproxy.io/docs/envoy/v1.5.0/configuration/http_filters/rate_limit_filter#config-http-filters-rate-limit

https://www.envoyproxy.io/docs/envoy/v1.5.0/configuration/network_filters/rate_limit_filter#config-network-filters-rate-limit



## Demos
 * [Local ratelimiting demo](./demo-local/README.md)
 * [Basic centralised ratelimiting demo](./demo-rq/README.md) 
 * [Centralised ratelimiting plus auth](./demo-rlimit/README.md) 


 * ratelimit by ip
 https://www.envoyproxy.io/docs/envoy/latest/api-v3/config/route/v3/route_components.proto#envoy-v3-api-msg-config-route-v3-ratelimit-action

If you get to this point and have time available, try creating your own demo combining both local ratelimiting and centralised ratelimiting.

## Links

https://www.envoyproxy.io/docs/envoy/v1.5.0/intro/arch_overview/global_rate_limiting
https://serialized.net/2019/05/envoy-ratelimits/
https://www.aboutwayfair.com/tech-innovation/understanding-envoy-rate-limits



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

### Auth API for ratelimiting

One way that rate limiting can be achieved in Envoy is using the [auth](https://www.envoyproxy.io/docs/envoy/
latest/intro/arch_overview/security/ext_authz_filter.html) filter, and implementing the external authorization 
API with some rate limiting logic. This can be helpful if you already have some sort of ratelimiting 
implementation and you want to do the simplest possible thing to make it compatible with Envoy.

Another useful thing you can do with an auth API which is helpful for ratelimiting is to use the auth process to 
set headers which can then be used for ratelimiting. So the auth filter comes first, decorating the request
with useful metadata, and then ratelimiting occurs based on that. 
For example, it is often necessary to look up metadata about tenant IDs based on bearer tokens.
In general, any data that we do not trust the client to send is a good candidate for this pattern.


### Envoy ratelimiting API

The other way to do ratelimiting is to use the purpose-built APIs and filters.
todo

## Envoy Ratelimit service

[Envoy Ratelimit](https://github.com/envoyproxy/ratelimit) is a service, written in Golang and backed by 
Redis, which provides centralised ratelimiting services which can be used with Envoy.
It is well integrated into the Envoy ecosystem - for instance, it uses similar configuration mechanisms.
And, of course, it supports the ratelimit APIs which Envoy Proxy uses.

## Local ratelimiting 

https://www.envoyproxy.io/docs/envoy/latest/intro/arch_overview/other_features/local_rate_limiting#arch-overview-local-rate-limit

todo

## Bandwidth limiting


https://www.envoyproxy.io/docs/envoy/latest/configuration/http/http_filters/bandwidth_limit_filter#config-http-filters-bandwidth-limit
todo

## todo loadstats
if you're doing something v fancy with global traffic shaping
maybe this goes elsewhere todo

## Limitations and Gotchas

The main one is that requests can vary in cost, for several reasons - often very significantly. 
Ratelimiting does not take this into account - all connections or requests are treated similarly.

## Configuration

https://www.envoyproxy.io/docs/envoy/v1.5.0/configuration/http_filters/rate_limit_filter#config-http-filters-rate-limit

https://www.envoyproxy.io/docs/envoy/v1.5.0/configuration/network_filters/rate_limit_filter#config-network-filters-rate-limit



## Demos

 * [Ratelimiting connection demo](./demo-basic/README.md) todo finish this
    - enabled, enforcing - https://www.envoyproxy.io/docs/envoy/v1.5.0/configuration/http_filters/rate_limit_filter#config-http-filters-rate-limit

 * ratelimit request demo with auth
    complex ratelimit - tenant id x global operation limit
 * local plus global todo
 * ratelimit actions todo


## Links

https://www.envoyproxy.io/docs/envoy/v1.5.0/intro/arch_overview/global_rate_limiting
https://serialized.net/2019/05/envoy-ratelimits/
https://www.aboutwayfair.com/tech-innovation/understanding-envoy-rate-limits



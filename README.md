# Load Management in Envoy - Workshop

Effective load management is a core aspect of the SRE role. In this workshop, participants will be introduced to a number of Envoy proxy features that are used for loadshedding and isolation, such as circuit breaking, adaptive concurrency, and ratelimiting. Participants will also use custom Go plugins to perform loadshedding. As part of the practical element of the workshop, participants will interact with Envoy configurations and status/control pages and endpoints, as well as Envoy’s telemetry.  

## Takeaways

* Familiarity with a variety of Envoy load management features
* Understanding when to use each form of load management, and the limitations of each
* Practical experience with Envoy configuration, controls, status, and metrics

## Prerequisites

Please bring a laptop to the workshop. Your laptop should have a working Docker and Docker Compose installation - see [Docker Docs](https://docs.docker.com/compose/install/) for installation instructions.

## Envoy Proxy

If you are not familiar with Envoy Proxy (or if you'd like a refresher), please spend a few minutes on
[the Envoy intro](./envoy.md).

## Other tools

As well as Envoy, the demos in this workshop also use [Prometheus](https://prometheus.io/) for collecting 
metrics, and [Grafana](https://grafana.com/docs/grafana/latest/introduction/) for dashboards.

These are provided for you, but you may want to extend them. 
Grafana's UI is reasonably intuitive (ask for help if you need it).
Some familiarity with Prometheus' query language, [PromQL](https://prometheus.io/docs/prometheus/latest/querying/basics/) may help - the docs have a number of examples - again, ask if you need help.

## Provided upstream/downstream code samples

The demos use some simple clients and servers - you'll find the Go source code for these in the [code](./code) 
directory here. The containers for the demos are built when running them, so you can change the source code and
have the changes take effect in your demo (if you want to try something the demo doesn't do).
The same upstreams/downstreams are used throughout the sequence of demos, so if you change something you could break other demos - if that's the case, you can just `git stash` your local changes and revert to the published code.

When running any demo, Prometheus will be available on [http://localhost:9090](http://localhost:9090) and
Grafana on [http://localhost:3000](http://localhost:3000).

## Sections

It is best to follow these in order, if possible.

[Envoy Circuit Breaking](/circuitbreaking/README.md): this section demonstrates how Envoy's circuitbreaking can be used to avoid cascading failures, based on connection count, concurrent requests, or concurrent retries.

[Envoy Adaptive Concurrency](/adaptiveconcurrency/README.md) is similar to circuitbreaking, but uses dynamic 
measurements of upstream performance to determine when to shed load, rather than preconfigured limits.
TODO finish demo.

[Envoy Admission Control](/admissioncontrol/README.md) is Envoy's version of client-side throttling, which
is particularly suitable in a service mesh scenario where there are a large number of Envoys running as sidecars with downstreams.

[Envoy Isolation and Ratelimiting](/isolation/README.md) demonstrates Envoy's mechanisms for ratelimiting, both
local and centralised.
TODO finish demo.

[How Envoy Adapts to Overload](/envoyoverload/README.md) shows some of the mechanisms that can be used to protect
Envoy itself from performance degradation in the face of excessive load.

# loadmgt-workshop

Effective load management is a core aspect of the SRE role. In this workshop, participants will be introduced to a number of Envoy proxy features that are used for loadshedding and isolation, such as circuit breaking, adaptive concurrency, and ratelimiting. Participants will also use custom Go plugins to perform loadshedding. As part of the practical element of the workshop, participants will interact with Envoy configurations and status/control pages and endpoints, as well as Envoy’s telemetry.  

## Takeaways

* Familiarity with Envoy load management features
* Understanding when to use each form of load management, and the limitations of each
* Practical experience with Envoy configuration, controls, status, and metrics


## Prerequisites

Please bring a laptop to the workshop. Your laptop should have a working Docker and Docker Compose installation - see [Docker Docs](https://docs.docker.com/compose/install/) for installation instructions.

Most of this workshop can be completed without writing code. Source code for the demo programs that drive load and which serve requests is available, and you might want to experiment with modifying that code. All of it is in Golang.

The final exercise involves using an Envoy custom ratelimiting plugin (again, in Golang). We provide both an incomplete outline version that you can complete yourself, and a fully-completed version, so it is possible for you to choose whether you want to write code or not. 

## Sections

Follow these in order.

 * [Envoy Circuit Breaking](/circuitbreaking/README.md)
 * [Envoy Admission Control](/admissioncontrol/README.md)
 * [Envoy Adaptive Concurrency](/adaptiveconcurrency/README.md)
 * [Envoy Isolation and Ratelimiting](/isolation/README.md)
 * [How Envoy Adapts to Overload](/envoyoverload/README.md)
 * [Envoy Custom Plugins for Ratelimiting](/plugins/README.md) (in Golang)
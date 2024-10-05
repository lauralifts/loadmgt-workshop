# Envoy Adaptive Concurrency

## What it is and how it works

Lower maintenance way of circuitbreaking with sane defaults
Takes account of latency 
minrtt measurement process

## Limitations and Gotchas

HTTP only
No priorities
How minRTT measurement is done and using jitter to avoid, previous hosts retry predicate
Take steps to avoid healthchecks being in the minRTT measurement
Again like cb it works best when the envoy fleet is smaller than the number of clients
And again it doesn't do fairness/QoS, isn't coordinated, doesn't do request weight

## Configuration

https://www.envoyproxy.io/docs/envoy/latest/configuration/http/http_filters/adaptive_concurrency_filter

What is the default behaviour - don't think this one is active by default, you have to add the filter

## Demos TODO

Multiple upstreams, observe the impact of the minRTT window
Healthchecks distorting minRTT
Stats should be in the demos see config
enabled on/off

## Links


https://www.alibabacloud.com/blog/brief-analysis-of-envoy-adaptive-concurrency-filter_600658
https://klaviyo.tech/adaptive-concurrency-control-for-mixed-analytical-workloads-51350439aeec


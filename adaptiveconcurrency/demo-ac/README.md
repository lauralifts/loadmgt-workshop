# Envoy Adaptive Concurrency

Bring down any other docker-compose environment running as part of this workshop (to avoid port clashes).

Bring this environment up by running 

```
 docker-compose up --build -d
```

## Adaptive concurrency config

The adaptive concurrency configuration is here:

```
          http_filters:
          - name: envoy.filters.http.adaptive_concurrency
            typed_config:
              "@type": type.googleapis.com/envoy.extensions.filters.http.adaptive_concurrency.v3.AdaptiveConcurrency
              gradient_controller_config:
                sample_aggregate_percentile:
                  value: 90
                concurrency_limit_params:
                  concurrency_update_interval: 0.1s
                min_rtt_calc_params:
                  jitter:
                    value: 10
                  interval: 60s
                  request_count: 50
              enabled:
                default_value: false
                runtime_key: "adaptive_concurrency.enabled"
```

Let's start by sending [100 qps](http://localhost:9094/config?http_rate=100&http_max_parallelism=100) to the upstreams.
There are two upstreams in this setup. 
One is configured to respond immediately to requests, and the other is configured to increase its latency as the request rate increases.

Use the [Grafana dashboard](http://localhost:3000/d/workshop/load-management-workshop?orgId=1&refresh=5s) to see the traffic coming to the two upstreams.

Adaptive concurrency is not enabled yet, so you should see no throttling.

You should see:
 * TODO

## Enabling adaptive concurrency


todo observe effect by varying req rate - increasing delay kicks in after 20 qps
todo observe that healthchecks aren't affected
todo observe the minrtt calc... 
todo observe different upstream characteristics

    - multiple upstreams
    - healthcheck endpoint excluded
    - jitter and retries with retry predicate
    - enabled on/off via runtime flag (curl post)

Multiple upstreams, observe the impact of the minRTT window, observe impact of different capacities
 - long minRTT

Healthchecks distorting minRTT
Stats should be in the demos see config

impact of different percentiles

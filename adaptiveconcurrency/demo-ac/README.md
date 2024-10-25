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

Let's start by sending [10 qps](http://localhost:9094/config?http_rate=3&http_max_parallelism=100) to the upstream.
The upstream is configured to respond to requests at 100ms latency, and to increase its latency as a square of the number of current inflight requests.
So at  qps, there should be little to no latency increase.

Use the [Grafana dashboard](http://localhost:3000/d/workshop/load-management-workshop?orgId=1&refresh=5s) to see the traffic.

Adaptive concurrency is not enabled yet, so you should see no throttling.

You should see:
 * No throttling - same number of downstream requests as upstream, and all status 2xx
 * You should see p50 latency is around 100 milliseconds

 Now send [50 qps](http://localhost:9094/config?http_rate=50&http_max_parallelism=100) to the upstream.

You should see:
 * Massive increase in latency
 * Requests still succeeding, but very slow


## Enabling adaptive concurrency

Now enable adaptive concurrency by changing your `envoy.yaml`

```
              enabled:
                default_value: true
                runtime_key: "adaptive_concurrency.enabled"
```

Run `docker-compose restart envoy` so that Envoy will pick up these changes.

You will immediately see adaptive concurrency kick in: Envoy will start sending a lot of 503s downstream. 
Less traffic makes it through to the upstreams, but some does. 
Latency for successful requests drops.

TODO explain MinRTT and Gradient here.


todo observe effect by varying req rate - increasing delay kicks in after 20 qps
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

Adding outlier detection
https://www.envoyproxy.io/docs/envoy/latest/intro/arch_overview/upstream/outlier#arch-overview-outlier-detection

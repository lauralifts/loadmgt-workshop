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

Let's start by sending 3 requests per second to the upstream: http://localhost:9094/config?http_rate=3&http_max_parallelism=100
The upstream is configured to respond to requests at 100ms latency, and to increase its latency as a square of the number of current inflight requests.
So at 3 qps, there should be little to no latency increase.

Use the [Grafana dashboard](http://localhost:3000/d/workshop/load-management-workshop?orgId=1&refresh=5s) to see the traffic.

Adaptive concurrency is not enabled yet, so you should see no throttling.

You should see:
 * No throttling - same number of downstream requests as upstream, and all status 2xx
 * You should see p50 latency is around 100 milliseconds

 Now send 50 requests per second: http://localhost:9094/config?http_rate=50&http_max_parallelism=100

You should see:
 * A massive increase in latency
 * Some requests still succeeding, but very slowly

## Enabling adaptive concurrency

Go back to sending [3 qps](http://localhost:9094/config?http_rate=3&http_max_parallelism=100) to the upstream.

Now enable adaptive concurrency by changing your `envoy.yaml`

```
              enabled:
                default_value: true
                runtime_key: "adaptive_concurrency.enabled"
```

Run `docker-compose restart envoy` so that Envoy will pick up these changes.

You will immediately see adaptive concurrency kick in: 
 * Envoy will start sending a lot of 503s downstream. 
 * Less traffic makes it through to the upstreams, but some does. 
 * Latency for successful requests drops.

You will see several things happening to the Adaptive Concurrency metrics shown on the dashboard.
 * The adaptive concurrency limit, minRTT, gradient and sampleRTT metrics now have nonzero values

However, you should not see any throttling occurring.

## Increase traffic again

Send 50 requests per second: http://localhost:9094/config?http_rate=50&http_max_parallelism=100
As before, latency will go through the roof.

You should also see:
 * The adaptive concurrency limit drop sharply to about 3
 * The minRTT and gradient metrics will change
 * Envoy will start to throttle a lot of traffic with 503s - see the downstream requests by code 

You will see that the sample RTT msecs changes every minute, as Envoy re-samples the latency. 
You should see that the gradient and the adaptive concurrency limit increase when the sample RTT goes down,
and vice-versa.

## Reducing the request rate and latency

Reduce the request rate again: http://localhost:9094/config?http_rate=3&http_max_parallelism=100
You should see that Envoy stops throttling as much. 

## Healthchecks

This example includes a separate HTTP listener for healthchecks, and a simple program that
pings the healthcheck listener. 
These healthchecks are routed to a separate listener which does not include an adaptive concurrency filter.
This is to make sure that healthcheck traffic - which is often very lightweight - does not skew the minRTT calculation. 
Look at `envoy.yaml` and understand how this has been done.

## Bring the demo down

Run 

```
 docker-compose down
```

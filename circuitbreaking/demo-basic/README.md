# Envoy Request-Based Circuitbreaking 

Bring down any other docker-compose environment running as part of this workshop (to avoid port clashes).

Bring this environment up by running 

```
 docker-compose up --build
```

## Basic request-based circuitbreaking config

Envoy here is configured with two sets of listeners and two backend clusters.
One set serves HTTP requests and one set gRPC.

There is a [Grafana dashboard](http://localhost:3000/d/workshop/load-management-workshop?orgId=1&refresh=5s) with visualisations showing requests upstream and downstream, and the state of its circuitbreakers.

The HTTP cluster has the following circuit breaker config:

```
    circuit_breakers:
      thresholds:
      - priority: DEFAULT
        max_connections: 1024
        max_pending_requests: 1024
        max_requests: 1024
        max_retries: 10
        track_remaining: true
```

### Normal traffic with no circuitbreaking

Let's send some HTTP traffic through Envoy.
Use the config endpoint of our downstream load-generator program: [100 qps and high parallelism](http://localhost:9094/config?http_rate=100&http_max_parallelism=2000)

You should see in Grafana that the request-based CB is closed, and that there are many requests and connections left before the CB would close. 

### Degraded upstream performance anc circuitbreaking

Now, let's change the performance of the upstream: [100msec latency and parallelism 1](http://localhost:9092/config?latency=100&parallelism=1)

Now the server can only handle 10 qps, and we are sending 100 qps. 
The connections and requests quickly pile up, and we quickly see the CBs close. 

The downstreams continue to make requests, but Envoy will send 504s.

### Restoring normal upstream performance, circuit breakers close

If we restore the upstream to its normal performance, then the circruit breakers close and the requests flow again. 
http://localhost:9092/config?latency=1&parallelism=1000

Note - if the downstreams can't connect it may be due to socket exhaustion (the code is intentionally written as a connection hog to demo some of Envoy's capabilities).
Restart the downstream container in case of issues.
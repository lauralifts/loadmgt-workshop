# Envoy Retry-Based Circuitbreaking 

Bring down any other docker-compose environment running as part of this workshop (to avoid port clashes).

Bring this environment up by running 

```
 docker-compose up --build -d
```

## Basic retry-based circuitbreaking config

Envoy is set up in a similar way to the previous demo.

There is a [Grafana dashboard](http://localhost:3000/d/workshop/load-management-workshop?orgId=1&refresh=5s) with visualisations showing requests upstream and downstream, and the state of its circuitbreakers.

### Enabling retries

The only difference here is that retries are enabled on the route used by the HTTP listener:

```
                  retry_policy:
                    retry_on: "5xx"
                    num_retries: 3
```

As before, there is a circuit breaker configured to open when there are 10 concurrent outstanding retries.
This gives a mechanism for reducing traffic to upstreams when they are experiencing high error rates, which
can often be the result of some transient problem. Avoiding a retry storm prevents a short-lived issue from 
becoming a [metastable failure](https://charap.co/metastable-failures-in-the-wild/).

### Sending HTTP traffic

Let's send some HTTP traffic through Envoy.
Use the config endpoint of our downstream load-generator program: [config - 100 qps](http://localhost:9094/config?http_rate=100&http_max_parallelism=100)

You should see in the [Grafana dashboard](http://localhost:3000/d/workshop/load-management-workshop?orgId=1&refresh=5s) that all CBs are closed,
that all requests are making it to the upstream and are succeeding.

### Upstreams have a high error rate and CBs open

Now, let's change the behaviour of the upstream: [20% error rate and 1000ms latency](http://localhost:9092/config?latency=1000&error_rate=0.2)
With 100 QPS sent this should consume the retry budget and cause the retry-based circuitbreaker to open, cutting off traffic to the upstream.

### Error rate returns to normal and CBs close

Restore the error rate to 0: http://localhost:9092/config?latency=1000&error_rate=0

You should see the rate of 5xx responses drop on the graph, and the circuit breaker will close. 

## Bring the demo down

Run 

```
 docker-compose down
```

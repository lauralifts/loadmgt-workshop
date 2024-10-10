# Envoy Prioritisation and Circuit Breaking

Bring down any other docker-compose environment running as part of this workshop (to avoid port clashes).

Bring this environment up by running 

```
 docker-compose up --build
```

## Prioritisation

Envoy is set up in a similar way to the previous demos, with the addition of a new route and a new priority-based
circuitbreaking configuration.



The config below allows more connections, pending requests, and outstanding retries to traffic that is 
on a high priority route.

```
      - priority: DEFAULT
        max_connections: 512
        max_pending_requests: 512
        max_requests: 512
        max_retries: 10
        track_remaining: true
      - priority: HIGH
        max_connections: 1024
        max_pending_requests: 1024
        max_requests: 1024
        max_retries: 20
        track_remaining: true
```

Traffic to the HTTP service now has two prioritised routes:

```
            routes:
              - match: { prefix: "/hipri" }
                route:  
                  cluster: some_http_service 
                  priority: HIGH
                  retry_policy:
                    retry_on: "5xx"
                    num_retries: 3
              - match: { prefix: "/" }
                route:  
                  cluster: some_http_service 
                  retry_policy:
                    retry_on: "5xx"
                    num_retries: 3
```



As usual, there is a [Grafana dashboard](http://localhost:3000/d/workshop/load-management-workshop?orgId=1&refresh=5s) with visualisations showing requests upstream and downstream, and the state of its circuitbreakers.

In this example, we have two downstreams, one on port 9094 and one at 9095.

Let's make the upstream flaky and slow: http://localhost:9092/config?latency=1000&error_rate=0.1&parallelism=1000

Let's send some default priority HTTP traffic through Envoy.
Use the config endpoint of our downstream load-generator program: [config - 100 qps](http://localhost:9094/config?http_rate=1000&http_max_parallelism=1000)

That should be sufficient to trip the default priority request circuit breaker: see this happen in the Grafana dash.

Now let's try to make some high priority requests. http://localhost:9095/config?hipri=true&http_rate=10&http_max_parallelism=10

We should see the hipri requests largely succeeding and more of the default priority requests being loadshed, as the default pririty circuit breaker opens.

TODO make this clearer in the stats

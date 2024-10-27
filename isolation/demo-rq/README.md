# Envoy centralised ratelimiting

In this demo we will explore Envoy ratelimiting with the centralised [ratelimit](https://github.com/envoyproxy/ratelimit) service.

The demo uses the following Envoy config:

```
      - name: envoy.rate_limit
        typed_config: 
          "@type": type.googleapis.com/envoy.extensions.filters.network.ratelimit.v3.RateLimit
          stat_prefix: ingress_grpc_rlimit
          domain: backend
          failure_mode_deny: false
          descriptors:
          - entries:
            - key: client
              value: foo
          rate_limit_service:
            grpc_service:
              envoy_grpc:
                cluster_name: some_ratelimit_service
              timeout: 0.25s 
```

This will call the ratelimit service for each gRPC request, with the descriptor pair `client=foo`.

The ratelimit config service is very simple, permitting 10 requests per second for that pair:

```
domain: backend
descriptors:
  - key: client
    value: foo
    rate_limit:
      requests_per_unit: 10
      unit: second
```

Bring this environment up by running 

```
 docker-compose up --buil -d
```

## Generating gRPC traffic

You can generate some traffic: http://localhost:9094/config?grpc_rate=100&grpc_max_parallelism=2000
This sends 100 requests per second to the upstream via Envoy.

As usual, open [Grafana](http://localhost:3000/d/workshop/load-management-workshop?orgId=1&refresh=5s).

You will see 100 requests of gRPC traffic being served. If you look at the gRPC status codes graphed from the perspective of the downstreams you will see 
that most requests are geting a 'Unavailable' status code.
The `Envoy ratelimit over limit` graph will also show that many requests are over limit.

The ratelimit service has a status page at http://localhost:6070/
Open it and check:
 * stats at http://localhost:6070/stats/
 * config at http://localhost:6070/rlconfig


## Failure mode deny

In `envoy.yaml` you will see the line

```
          failure_mode_deny: false
```

This means that if the ratelimit service is unavailable, Envoy will not deny requests.
Stop the rate limit service by running

```
docker-compose stop ratelimit
```

Look at Grafana. You should see:
 * Requests to upstreams increasing
 * Downstreams seeing 100% gRPC status code OK

You'll see some Envoy ratelimit errors in Grafana. 

If you open [Envoy clusters](http://localhost:9901/clusters) you'll see no hosts.
You can also see the Envoy ratelimit healthy upstreams go to zero in Grafana.

Run `docker-compose start ratelimit` to bring the ratelimiter back.

## Modifying the rate limit

You can see the current rate limit config by opening http://localhost:6070/rlconfig

Modify it by changing `ratelimit/config/config.yaml` in this directory.

Change it to 

```
domain: backend
descriptors:
  - key: client
    value: foo
    rate_limit:
      requests_per_unit: 100
      unit: second
```

Run

 ```
docker-compose restart ratelimit
```

Reload http://localhost:6070/rlconfig to see it take effect.
You should see that 100 requests per second are now permitted.

## Shutting the environment down

There are a variety of ways to do this: use Docker Desktop if you have it installed, or run 
```
docker-compose down
```
in this directory.

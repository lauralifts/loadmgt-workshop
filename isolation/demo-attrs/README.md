# Envoy centralised ratelimiting

In this demo we will further explore Envoy ratelimiting with the centralised [ratelimit](https://github.com/envoyproxy/ratelimit) service.

In this scenario, we use a header - `x-level` sent by the downstreams to do ratelimiting (in a real production scenario this might be set by
an external auth service). 60% of requests sent will have level `free`, and 20% each `paid` and `enterprise`.

The Envoy configuration is as follows:

```
          route_config:
            name: local_route
            virtual_hosts:
            - name: local_service
              domains: ["*"]
              routes:
              - match: { prefix: "/" }
                route: { cluster: some_http_service }
              rate_limits:
                - actions:
                  - request_headers:
                      descriptor_key: level
                      header_name: x-level
          http_filters:
          - name: envoy.rate_limit
            typed_config: 
              "@type": type.googleapis.com/envoy.extensions.filters.http.ratelimit.v3.RateLimit
              stat_prefix: http_rlimit
              domain: backend
              failure_mode_deny: false
              rate_limit_service:
                grpc_service:
                  envoy_grpc:
                    cluster_name: some_ratelimit_service
                  timeout: 0.25s  
```

The ratelimit configuration is

```
domain: backend
descriptors:
  - key: level
    value: free
    rate_limit:
      requests_per_unit: 10
      unit: second

  - key: level
    value: paid
    rate_limit:
      requests_per_unit: 100
      unit: second

  - key: level
    value: enterprise
    rate_limit:
      requests_per_unit: 1000
      unit: second
```


Bring this environment up by running 

```
 docker-compose up --build -d
```

## Generating HTTP traffic

You can generate some traffic: http://localhost:9094/config?http_rate=100&http_max_parallelism=2000
This sends 100 requests per second to the upstream via Envoy.

As usual, open [Grafana](http://localhost:3000/d/workshop/load-management-workshop?orgId=1&refresh=5s).

You should observe that lots of `free` tier traffic is throttled, but no paid or enterprise tier.

## Experiment with changing the ratelimits, or adding extra descriptors

You can add more descriptors under `actions` in the `envoy.yaml` configuration - see the [rate limit filter](https://www.envoyproxy.io/docs/envoy/latest/configuration/http/http_filters/rate_limit_filter)
docs.

For example

```
              rate_limits:
                - actions:
                  - request_headers:
                      descriptor_key: level
                      header_name: x-level
                  - request_headers:
                      descriptor_key: agent
                      header_name: User-Agent
                  - {remote_address: {}}
```

Restart envoy as usual - `docker-compose restart envoy` and observe the logs for the ratelimit service.

Either use the Docker Desktop UI or run `docker-compose logs ratelimit`
to see what values the ratelimit service receives.

Experiment with changing the the ratelimit configuration in `ratelimit/config/config.yaml`.
You will need to restart the ratelimit service each time.

It's useful to watch the debug logs from the ratelimit service as you do this.
You can also see the current ratelimit config: http://localhost:6070/rlconfig

Some examples of ratelimit service config can be found [here](https://github.com/envoyproxy/ratelimit?tab=readme-ov-file#examples).

## Shutting the environment down

There are a variety of ways to do this: use Docker Desktop if you have it installed, or run 
```
docker-compose down
```
in this directory.

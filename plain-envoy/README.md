# Basic Envoy setup with upstream, downstream, and monitoring

This isn't an exercise section - it's just a 'plain' Envoy setup with:
  * Envoy 
  * Envoy config that gets copied to Envoy container
  * A 'load driver' that generates load (HTTP or gRPC)
  * A server that serves RPCs (with configurable max capacity, latency)
  * Prometheus, which collects metrics from envoy, the load driver, and the server
  * Grafana, with some pre-configured dashboards

Copying this is the jumping-off point for each of the exercise sections.

Bring this environment up by running 

```
 docker-compose up --build
```

## Envoy admin interface

With the docker-compose cluster running, access [http://localhost:9901/help](http://localhost:9901/help).
See documentation on the admin inteeface [here](https://www.envoyproxy.io/docs/envoy/latest/operations/admin).

## Envoy listeners
HTTP on 9902, gRPC on 9903.

## Upstream

Hardcoded to port 9092 for HTTP, 9093 for gRPC.
Serves HTTP on /.

Prometheus metrics: http://localhost:9092/metrics

Has /config endpoint - http://localhost:9092/config

Config endpoint params can be used to update the config, e.g.

```
http://localhost:9092/config?latency=10&parallelism=40
```

## Downstream

Takes address of HTTP and gRPC test servers in env.
Has /config endpoint - http://localhost:9094/config

Config endpoint params can be used to update the config, e.g.

```
http://localhost:9094/config?http_rate=1&http_max_parallelism=10
http://localhost:9094/config?grpc_rate=100&grpc_max_parallelism=2
```

Prometheus metrics: http://localhost:9094/metrics

## Prometheus

TODO

## Grafana 

TODO
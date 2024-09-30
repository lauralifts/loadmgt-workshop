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
See documentation on the admin inteface [here](https://www.envoyproxy.io/docs/envoy/latest/operations/admin).

## Upstream

Hardcoded to port 9092. 
Serves HTTP on /.
Exports counter for HTTP requests seen/served.

configurable latency? configurable max parallelism? configurable port?
todo add and env file

todo all same for grpc
todo add a status page with configs
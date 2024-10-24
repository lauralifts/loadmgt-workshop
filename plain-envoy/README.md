# Basic Envoy setup with upstream, downstream, and monitoring

This isn't an exercise section - it's just a 'plain' Envoy setup with:
  * Envoy 
  * Envoy config that gets copied to Envoy container
  * A 'load driver' that generates load (HTTP or gRPC)
  * A server that serves RPCs (with configurable max capacity, latency)
  * Prometheus, which collects metrics from envoy, the load driver, and the server
  * Grafana, with some pre-configured dashboards

Copying this is the jumping-off point for each of the exercise sections.
You can also use it to kick the tires on Envoy, or clone it to try something else.

Bring this environment up by running 

```
 docker-compose up --build -d
```

You can generate some HTTP traffic using the config endpoint of our downstream load-generator program: [config - 100 qps](http://localhost:9094/config?http_rate=100&http_max_parallelism=2000). This sends 100 requests per second to the upstream via Envoy.

Prometheus will be available on [http://localhost:9090](http://localhost:9090) and
Grafana on [http://localhost:3000](http://localhost:3000/d/workshop/load-management-workshop?orgId=1&refresh=5s).

You can browse [the Envoy admin interface](http://localhost:9901) and the [Envoy statistics endpoint](http://localhost:9901/stats/prometheus).

## Shutting the environment down

There are a variety of ways to do this: use Docker Desktop if you have it installed, or run 
```
docker-compose down
```
in this directory.

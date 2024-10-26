# Envoy ratelimiting by connection

In this demo we will explore Envoy ratelimiting by TCP connection.
This can be useful for non-HTTP services.

Bring this environment up by running 

```
 docker-compose up --build --remove-orphans -d
```



You can generate some HTTP traffic using the config endpoint of our downstream load-generator program: [config - 100 qps](http://localhost:9094/config?http_rate=100&http_max_parallelism=2000). This sends 100 requests per second to the upstream via Envoy.

Prometheus will be available on [http://localhost:9090](http://localhost:9090) and
Grafana on [http://localhost:3000](http://localhost:3000).

You can browse [the Envoy admin interface](http://localhost:9901) and the [Envoy statistics endpoint](http://localhost:9901/stats/prometheus).

todo demo shadow mode

http://localhost:6070/rlconfig



## Shutting the environment down

There are a variety of ways to do this: use Docker Desktop if you have it installed, or run 
```
docker-compose down
```
in this directory.
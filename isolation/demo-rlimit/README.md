# Envoy centralised ratelimiting

In this demo we will explore Envoy ratelimiting with the centralised [ratelimit](https://github.com/envoyproxy/ratelimit) service, 
plus an authentication service.

https://www.funnel-labs.io/2022/10/10/envoyproxy-3-sophisticated-rate-limiting/
https://serialized.net/2019/05/envoy-ratelimits/



todo discuss cofnig

Bring this environment up by running 

```
 docker-compose up --build --remove-orphans -d
```

## Generating gRPC traffic

You can generate some traffic : [100 qps](http://localhost:9094/config?grpc_rate=100&grpc_max_parallelism=2000). This sends 100 requests per second to the upstream via Envoy.

As usual, open [Grafana](http://localhost:3000/d/workshop/load-management-workshop?orgId=1&refresh=5s).

## Shutting the environment down

There are a variety of ways to do this: use Docker Desktop if you have it installed, or run 
```
docker-compose down
```
in this directory.
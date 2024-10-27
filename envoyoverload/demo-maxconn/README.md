# Envoy Overload Manager Max Connections Demo

This demo shows how we can use Envoy to limit the number of downstream connections.

You will see that the `envoy.yaml` in this folder contains a snippet which configures the Overload Manager:

```
overload_manager:
  refresh_interval:
    seconds: 0
    nanos: 250000000
  resource_monitors:
  - name: "envoy.resource_monitors.global_downstream_max_connections"
    typed_config:
      "@type": type.googleapis.com/envoy.extensions.resource_monitors.downstream_connections.v3.DownstreamConnectionsConfig
      max_active_downstream_connections: 1000
```

When using the downstream connections resource monitor, we don't need to specify any triggers, overload actions, or loadshed points,
because the resource monitor configuration by itself serves to limit incoming downstream connections.

## Running the demo

Make sure any previous demos from this workshop have been stopped.

Bring this environment up by running 
```
 docker-compose up --build -d
```

Open [Grafana](http://localhost:3000/d/workshop/load-management-workshop?orgId=1&refresh=5s) and expand the Envoy row.

Start sending gRPC requests:
[http://localhost:9094/config?grpc_rate=1000&grpc_max_parallelism=1000]([http://localhost:9094/config?grpc_rate=1000&grpc_max_parallelism=1000](http://localhost:9094/config?grpc_rate=10&grpc_max_parallelism=10)
)

You should see the graph of Envoy connections increase and level off at 1000 connections. You should see 
an equal number of incoming requests to the upstreams as outgoing requests at the downstreams. 
There is no loadshedding (or only minor loadshedding)

Now increase the parallelism of your requests to 1100:
http://localhost:9094/config?grpc_rate=2000&grpc_max_parallelism=1100

Now you should begin to see connections being rejected. You won't see requests being rejected; the downstreams will
just send requests over the connections that are available to them.

However, you will see the graph of rejected downstream connections increase above zero, indicating that
Envoy is performing throttling here.

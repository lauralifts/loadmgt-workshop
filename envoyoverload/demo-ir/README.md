# Envoy Injected Resource Demo

This demo shows how Envoy loadshedding based on an injected resource works.

You will see that the `envoy.yaml` in this folder contains a snippet which configures the Overload Manager:

```
overload_manager:
  refresh_interval:
    seconds: 0
    nanos: 250000000
  resource_monitors:
  - name: "envoy.resource_monitors.injected_resource"
    typed_config:
      "@type": type.googleapis.com/envoy.extensions.resource_monitors.injected_resource.v3.InjectedResourceConfig
      filename: "/res.txt"
  actions:
  - name: "envoy.overload_actions.stop_accepting_requests"
    triggers:
    - name: "envoy.resource_monitors.injected_resource"
      scaled:
        scaling_threshold: 0.70
        saturation_threshold: 0.85
```

This means that when the `res.txt` file (which is mapped to res.txt in this folder) is above 0.7, Envoy will
begin to reject requests in a scaled manner. Above 0.85 all requests will be rejected.

## Running the demo

Make sure any previous demos from this workshop have been stopped.

Bring this environment up by running 
```
 docker-compose up --build -d
```

Open [Grafana](http://localhost:3000/d/workshop/load-management-workshop?orgId=1&refresh=5s) and expand the Envoy row.

You will see graphs for injected resource pressure, and whether the stop accepting requests action is active.
You should see that injected resource pressure is 50 percent, and that stop accepting requests is not active.

Start sending gRPC requests:
[http://localhost:9094/config?grpc_rate=1000&grpc_max_parallelism=1000](http://localhost:9094/config?grpc_rate=1000&grpc_max_parallelism=1000)

You should see all requests proceeding.

Now, change the contents of `res.txt` to 0.9.
You will have to also restart the Envoy container to pick that up.

```
docker-compose restart envoy
```

In Grafana you should now see 
 * Injected resource pressure at 90
 * The stop accepting requests action is active
 * Requests downstream but none upstream

Now, change the contents of `res.txt` to 0.75 - between the scaling threshold and the saturation threshold.
Restart Envoy

```
docker-compose restart envoy
```

Now you should see that stop accepting requests is at around 33%, and Envoy does shed some requests, and you can see this because the downstream requests are higher than upstream.

Note that loadshed gRPC requests will have HTTP 200 responses in the Envoy graphs - but the downstream graphs will show you ther gRPC status codes.

## Bypassing the overload manager

Send HTTP requests.
[http://localhost:9094/config?http_rate=100&http_max_parallelism=10](http://localhost:9094/config?http_rate=100&http_max_parallelism=100)

Wait a few seconds and these will appear on the Envoy graphs.
Unlike the gRPC requests, you should notice that the HTTP requests are not throttled - the upstream request count is the same as the downstream.

Take a look in `envoy.yaml` at the HTTP listener - you will notice that `bypass_overload_manager: true` is set on 
the HTTP listener, but not the gRPC listener.

This can be useful for exempting high priority traffic from overload manager loadshedding.

This is a particularly good idea for healthcheck traffic - in an overload state, if you shed your healthcheck 
traffic you may end up in a state where your orchestration service constantly restarts your proxies, which
further reduces your capacity and can contribute to a cascading failure scenario or metastable failure state.

## Bring the demo down

Run 

```
 docker-compose down
```


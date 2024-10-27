# Envoy local ratelimiting

In this demo we will demonstrate Envoy's local ratelimiting functionality.
This can be useful in context where you don't need coordination - for instance, where Envoys are
colocated with clients and you want to limit each client to a certain level of traffic.

Additionally, local ratelimiting can be used in conjunction with other forms of ratelimiting. 

Bring this environment up by running 

```
 docker-compose up --build -d
```

Envoy has the following local ratelimit configuration for HTTP traffic:

```
  http_filters:
          - name: envoy.filters.http.local_ratelimit
            typed_config:
              "@type": type.googleapis.com/envoy.extensions.filters.http.local_ratelimit.v3.LocalRateLimit
              stat_prefix: http_local_rate_limiter
              token_bucket:
                max_tokens: 10
                tokens_per_fill: 2
                fill_interval: 1s
              filter_enabled:
                runtime_key: local_rate_limit_enabled
                default_value:
                  numerator: 100
                  denominator: HUNDRED
              filter_enforced:
                runtime_key: local_rate_limit_enforced
                default_value:
                  numerator: 0
                  denominator: HUNDRED
              response_headers_to_add:
              - append_action: OVERWRITE_IF_EXISTS_OR_ADD
                header:
                  key: x-local-rate-limit
                  value: 'true'
```

This limits HTTP traffic to 2 requests per second (with the ability to 'store' up to 10 tokens to smooth peaks and troughs).
However, the limit is not actually enforced: it is only advisory. This is useful as a way to test the impact of ratelimits before enforcing them.
As Twitter demonstrated in the recent past, this is a good way to have an outage! Gradual enforcement of ratelimiting is always recommended.

## Observing the ratelimit without enforcement

You can generate some HTTP traffic using the config endpoint of our downstream load-generator program: http://localhost:9094/config?http_rate=5&http_max_parallelism=50

As usual, open the [Grafana dashboard](http://localhost:3000/d/workshop/load-management-workshop?orgId=1&refresh=5s).
You should see successful HTTP requests with no throttling applied yet. In the Envoy section you will see graphs showing
for how many requests the rate limit is enabled, and for how many it is enforced.

Rate limiting is enabled for all requests, but enforced for none.
We can see from the `Rate Limit OK` and `Rate Limited Requests` graphs that if rate limiting were enabled, we would reject 3 requests per second.

Increase the load to 50 qps: http://localhost:9094/config?http_rate=50&http_max_parallelism=100

You'll see the `Rate Limited Requests` graph increase to around 48 qps, but all requests will still succeed.

## Enforcing 50%

Change the filter configuration in `envoy.yaml` as follows:

```
  filter_enforced:
                runtime_key: local_rate_limit_enforced
                default_value:
                  numerator: 50
                  denominator: HUNDRED
```

Then run `docker-compose restart envoy` to pick that change up.

In Grafana, you should start to see some rate limting. 
If we were 100% enforcing rate limiting, we would block 48 requests per second.
At 50%, we should see around 24 requests per second blocked.
You will see downstreams receiving HTTP 4xx status codes. 

## Rate limit headers and actions

Run 

```
curl -I -v localhost:9902
```

and look at the response. Do this a few times.

You will see some requests succeed, but some should fail with HTTP 429 Too Many Requests status code.
You will also see that Envoy adds the `x-local-rate-limit: true` header - this is because we added an `append_action` in the configuration (see above).

Remove this section from the configuration:

```
              response_headers_to_add:
              - append_action: OVERWRITE_IF_EXISTS_OR_ADD
                header:
                  key: x-local-rate-limit
                  value: 'true'
```

Then run `docker-compose restart envoy` to pick that change up.

Run 

```
curl -I -v localhost:9902
```

again. You will see that you no longer get the `x-local-rate-limit` header.

## Enforcing 100%

Change the filter configuration in `envoy.yaml` as follows:

```
  filter_enforced:
                runtime_key: local_rate_limit_enforced
                default_value:
                  numerator: 100
                  denominator: HUNDRED
```

Then run `docker-compose restart envoy` to pick that change up.

You will see that Envoy now limits the request rate to 2 requests per second and the rest are loadshed.

## Increasing the rate limit

Change the token bucket configuration in `envoy.yaml` as follows:

```
            token_bucket:
                max_tokens: 50
                tokens_per_fill: 20
                fill_interval: 1s
```

Then run `docker-compose restart envoy` to pick that change up.

You will see that Envoy now allows 20 qps to be served.

## Shutting the environment down

There are a variety of ways to do this: use Docker Desktop if you have it installed, or run 
```
docker-compose down
```
in this directory.

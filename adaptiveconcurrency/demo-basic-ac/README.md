# Envoy Adaptive Concurrency

Bring down any other docker-compose environment running as part of this workshop (to avoid port clashes).

Bring this environment up by running 

```
 docker-compose up --build
```

## Basic adaptive concurrency config

todo describe configs
todo observe effect by varying req rate - increasing delay kicks in after 20 qps
todo observe that healthchecks aren't affected
todo observe the minrtt calc... 
todo observe different upstream characteristics
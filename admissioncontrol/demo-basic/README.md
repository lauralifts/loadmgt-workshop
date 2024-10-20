# Basic Admission Control demo

In this demo we have 10 downstream clients and 3 Envoy instances.
There is one upstream. 

Envoy admission control config is todos

We start off with each downstream sending 50 qps, which the upstream can serve.
TODO look at dash at admission control stats

Now we set the upstream error rate to 20%: todo config link

Observe the Grafana dash and see todo

Now let's disable admission control by modifying the Envoy runtime variable todo.
Observe that admission control no longer applies todo.

todo turn it back on again.
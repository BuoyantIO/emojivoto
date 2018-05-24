# Conduit lifecycle test configuration

Production testing the proxy's discovery & caching.

The goal of this test suite is to run an outbound proxy for a prolonged amount
of time in a dynamically-scheduled environment in order to exercise:
- Route resource lifecyle (i.e. routes are properly evicted)
- Telemetry resource lifecycle (i.e. prometheus can run steadily for a long
  time, proxy doesn't leak memory in exporter).
- Service discovery lifecycle (i.e. updates are honored correctly, doesn't get
  out sync).

## First time setup

[`lifecycle.yml`](lifecycle.yml) creates a `ClusterRole`, which requires your user to have this
ability.

```bash
kubectl create clusterrolebinding cluster-admin-binding-$USER \
  --clusterrole=cluster-admin --user=$(gcloud config get-value account)
```

## Deploy

Install Conduit service mesh:

```bash
conduit install --conduit-namespace conduit-lifecycle | kubectl apply -f -
conduit dashboard --conduit-namespace conduit-lifecycle
```

Deploy test framework to `lifecycle` namespace:

```bash
cat lifecycle.yml | conduit inject --conduit-namespace conduit-lifecycle - | kubectl apply -f -
```

## Observe

Browse to Grafana:

```bash
conduit dashboard --conduit-namespace conduit-lifecycle --show grafana
```

Tail slow-cooker logs:

```bash
kubectl -n lifecycle logs -f $(
  kubectl -n lifecycle get po --selector=job-name=slow-cooker -o jsonpath='{.items[*].metadata.name}'
) slow-cooker
```

Relevant Grafana dashboards to observe
- `Conduit Deployment`, for route lifecycle and service discovery lifecycle
- `Prometheus 2.0 Stats`, for telemetry resource lifecycle


## Teardown

```bash
kubectl delete ns lifecycle
kubectl delete ns conduit-lifecycle
```

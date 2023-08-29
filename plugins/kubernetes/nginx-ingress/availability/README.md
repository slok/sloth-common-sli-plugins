# Kubernetes nginx-ingress availability

Availability plugin for services using the nginx ingress controller for exposure.

Uses the HTTP response status codes to measure the events as good or bad. It counts an error event when HTTP response status code is >=500 or 429.

In other words, it counts as good events the <500 and !429 HTTP response status codes.

## Options

- `filter`: (**Optional**) A prometheus filter string using concatenated labels (e.g: `job="k8sapiserver",env="production",cluster="k8s-42"`)

## Metric requirements

- `nginx_ingress_controller_requests`.

## Usage examples

### Without filter

```yaml
sli:
  plugin:
    id: "sloth-common/kubernetes/nginx-ingress/availability"
```

### With custom filter

```yaml
sli:
  plugin:
    id: "sloth-common/kubernetes/nginx-ingress/availability"
    options:
      filter: ingress="k8sapiserver",host="my.example.com",exported_namespace="my-production-app"
```

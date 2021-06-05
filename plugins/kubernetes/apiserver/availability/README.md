# Kubernetes apiserver availability

Availability plugin for the Kubernetes apiserver.

Uses the API HTTP response status codes to measure the events as good or bad. It counts an error event when HTTP response status code is >=500 or 429.

In other words, it counts as good events the <500 and !429 HTTP response status codes.

## Options

- `filter`: (**Optional**) A prometheus filter string using concatenated labels (e.g: `job="k8sapiserver",env="production",cluster="k8s-42"`)

## Metric requirements

- `apiserver_request_total`.

## Usage examples

### Without filter

```yaml
sli:
  plugin:
    id: "sloth-common/kubernetes/apiserver/availability"
```

### With custom filter

```yaml
sli:
  plugin:
    id: "sloth-common/kubernetes/apiserver/availability"
    options:
      filter: job="k8sapiserver",env="production",cluster="k8s-42"
```

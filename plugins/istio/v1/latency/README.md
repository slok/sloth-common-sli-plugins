# Istio V1 latency

Latency plugin for [Istio V1][istio] services.

Uses Istio v1 request metrics to get the latency on a service.

## Options

- `namespace`: (**required**) Kubernetes namespace of the service.
- `service`: (**required**) Service name.
- `bucket`: (**Required**) The max latency allowed histogram bucket.
- `exclude_errors`: (**Optional**) Boolean that will exclude errored requests from valid events
- `filter`: (**Optional**) A prometheus filter string using concatenated labels

## Metric requirements

- `istio_request_duration_milliseconds_bucket`: From [istio].
- `istio_request_duration_milliseconds_count`: From [istio].

## Usage examples

### Without filter

```yaml
sli:
  plugin:
    id: "sloth-common/istio/v1/latency"
    options:
      namespace: "default"
      service: "test"
      bucket: "300"
```

### With filters

```yaml
sli:
  plugin:
    id: "sloth-common/istio/v1/latency"
    options:
      namespace: "default"
      service: "test"
      bucket: "300"
      filter: request_protocol="http"
```

### Excluding errors (5xx)

```yaml
sli:
  plugin:
    id: "sloth-common/istio/v1/latency"
    options:
      namespace: "default"
      service: "test"
      bucket: "300"
      filter: request_protocol="http"
      exclude_errors: true
```

[Istio]: https://istio.io/v1.10/docs/
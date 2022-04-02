# Istio V1 availability

Availability plugin for [Istio V1][istio] services.

Uses Istio v1 request metrics to get the availability on a service.

## Options

- `namespace`: (**required**) Kubernetes namespace of the service.
- `service`: (**required**) Service name.
- `filter`: (**Optional**) A prometheus filter string using concatenated labels

## Metric requirements

- `istio_requests_total`: From [istio].

## Usage examples

### Without filter

```yaml
sli:
  plugin:
    id: "sloth-common/istio/v1/availability"
    options:
      namespace: "default"
      service: "test"
```

### With filters

```yaml
sli:
  plugin:
    id: "sloth-common/istio/v1/availability"
    options:
      namespace: "default"
      service: "test"
      filter: request_protocol="http"
```

[Istio]: https://istio.io/v1.10/docs/
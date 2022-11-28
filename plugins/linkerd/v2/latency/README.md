# Linkerd V2 latency

Availability plugin for [Linkerd V2](https://linkerd.io/) services.

Uses Linkerd v2 service metrics to get the correct and invalid availability on the serving services.

## Options

- `filter`: (**Optional**) A prometheus filter string using concatenated labels

## Metric requirements

- `request_total`: From [linkerd].

## Usage examples

### With filters

```yaml
sli:
  plugin:
    id: "sloth-common/linkerd/v2/latency"
    options:
      filter: `deployment="test-deploy", namespace="test-ns"`
```

[linkerd]: https://linkerd.io/2.12/overview/

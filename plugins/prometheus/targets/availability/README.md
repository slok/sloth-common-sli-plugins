# Prometheus Target availability

Availability plugin for the [Prometheus] targets.

A target will be counted as an error when is not up to be scraped by Prometheus.

## Options

- `filter`: (**Optional**) Prometheus extra label filter.

## Metric requirements

- `up`.

## Usage examples

### Without filter

```yaml
sli:
  plugin:
    id: "sloth-common/prometheus/targets/availability"
```

### With filter

```yaml
sli:
  plugin:
    id: "sloth-common/prometheus/targets/availability"
    options:
      filter: job="svc1",env="prod"
```

[prometheus]: https://prometheus.io

# Prometheus rules evaluation availability

Availability plugin for the [Prometheus] rules evaluation.

## Options

- `filter`: (**Optional**) Prometheus extra label filter.

## Metric requirements

- `prometheus_rule_evaluation_failures_total`.
- `prometheus_rule_evaluations_total`

## Usage examples

### Without filter

```yaml
sli:
  plugin:
    id: "sloth-common/prometheus/rules/eval-availability"
```

### With filter

```yaml
sli:
  plugin:
    id: "sloth-common/prometheus/rules/eval-availability"
    options:
      filter: job="svc1",env="prod"
```

[prometheus]: https://prometheus.io

# Kubernetes Kooper latency

Latency plugin for the Kubernetes [Kooper] controller/operator library.

This SLI will measure the latency of handling events/k8s objects on the Kooper controller handlers.

## Options

- `controller`: (**Required**) The controller being measured.
- `bucket`: (**Required**) The max latency allowed histogram bucket.
- `filter`: (**Optional**) Prometheus extra label filter.

## Metric requirements

- `kooper_controller_processed_event_duration_seconds_bucket`.
- `kooper_controller_processed_event_duration_seconds_count`.

## Usage examples

### Without filter

```yaml
sli:
  plugin:
    id: "sloth-common/kubernetes/kooper/latency"
    options:
      controller: "sloth"
      bucket: "0.25"
```

### With filter

```yaml
sli:
  plugin:
    id: "sloth-common/kubernetes/kooper/latency"
    options:
      controller: "sloth"
      bucket: "0.25"
      filter: job="svc1",env="prod"
```

[kooper]: https://github.com/spotahome/kooper

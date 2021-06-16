# Kubernetes Kooper Queue congestion

Queue congestion plugin for the Kubernetes [Kooper] controller/operator library.

This SLI will measure if there is a congestion on the event queue based on the duration an event is hold in the queue before being processed.

## Options

- `controller`: (**Required**) The controller being measured.
- `bucket`: (**Required**) The max latency allowed hitogram bucket.
- `filter`: (**Optional**) Prometheus extra label filter.

## Metric requirements

- `kooper_controller_event_in_queue_duration_seconds_count`.
- `kooper_controller_event_in_queue_duration_seconds_bucket`.

## Usage examples

### Without filter

```yaml
sli:
  plugin:
    id: "sloth-common/kubernetes/kooper/queue-congestion"
    options:
      controller: "sloth"
      bucket: "0.25"
```

### With filter

```yaml
sli:
  plugin:
    id: "sloth-common/kubernetes/kooper/queue-congestion"
    options:
      controller: "sloth"
      bucket: "0.25"
      filter: job="svc1",env="prod"
```

[kooper]: https://github.com/spotahome/kooper

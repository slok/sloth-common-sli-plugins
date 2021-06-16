# Kubernetes Kooper availability

Availability plugin for the Kubernetes [Kooper] controller/operator library.

Kooper library tracks the retries. These retries will be subtracted from the total events, so the number of retries doesn't affect the error ratio.

Kooper retries are processed as a correct metrics.

## Options

- `controller`: (**Required**) The controller being measured.
- `filter`: (**Optional**) Prometheus extra label filter.

## Metric requirements

- `kooper_controller_processed_event_duration_seconds_count`.
- `kooper_controller_queued_events_total`.

## Usage examples

### Without filter

```yaml
sli:
  plugin:
    id: "sloth-common/kubernetes/kooper/availability"
    options:
      controller: "sloth"
```

### With filter

```yaml
sli:
  plugin:
    id: "sloth-common/kubernetes/kooper/availability"
    options:
      controller: "sloth"
      filter: job="svc1",env="prod"
```

[kooper]: https://github.com/spotahome/kooper

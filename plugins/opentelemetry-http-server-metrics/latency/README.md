# opentelemetry-http-server-metrics latency

Latency plugin for [opentelemetry/http-server-metrics].

Uses the HTTP response histogram to measure the events as good or bad. It counts as error events that don't fall in the provided bucket.

For example if the bucket is `0.25`, We will measure as an error, the requests that fall in the buckets greater than `0.25`, in other words, that took longer than `250ms` .

WARNING: While the specification says the bucket units should be in seconds, some implementations (e.g. Backstage) use milliseconds. You may have to adapt this as needed for the service you are measuring.

## Options

- `bucket`: (**Required**) The max latency allowed histogram bucket.
- `filter`: (**Required**) A prometheus filter string using concatenated labels (e.g: `exported_service="alertmanager-api",service="alertgram"`)
- `exclude_errors`: (**Optional**) Boolean that will exclude errored requests from valid events when measuring latency requests.

## Metric requirements

- `http_server_duration_bucket`: From [example]
- `http_server_duration_count`: From [example]

## Usage examples

### Default

```yaml
sli:
  plugin:
    id: "sloth-common/opentelemetry-http-server-metrics/latency"
    options:
      filter: exported_service="alertmanager-api",service="alertgram"
      bucket: "1"
```

### Excluding errors (5xx)

```yaml
sli:
  plugin:
    id: "sloth-common/opentelemetry-http-server-metrics/latency"
    options:
      filter: exported_service="alertmanager-api",service="alertgram"
      bucket: "0.25"
      exclude_errors: "true"
```

[opentelemetry/http-server-metrics]: https://opentelemetry.io/docs/specs/semconv/http/http-metrics/#metric-httpserverrequestduration
[example]: https://github.com/open-telemetry/opentelemetry.io/blob/8411a8ca60efc531f0378e6835d3a311b4d0518e/content/en/docs/collector/internal-telemetry.md?plain=1#L305C3-L305C24

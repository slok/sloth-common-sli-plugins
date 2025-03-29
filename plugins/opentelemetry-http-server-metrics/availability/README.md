# opentelemetry-http-server-metrics availability

Availability plugin for [opentelemetry/http-server-metrics].

Uses the API HTTP response status codes to measure the events as good or bad. It counts an error event when HTTP response status code is >=500 or 429.

In other words, it counts as good events the <500 and !429 HTTP response status codes.

## Options

- `filter`: (**Required**) A prometheus filter string using concatenated labels (e.g: `exported_service="alertmanager-api",service="alertgram"`)

## Metric requirements

- `http_server_duration_count`: From [example].

## Usage examples

```yaml
sli:
  plugin:
    id: "sloth-common/opentelemetry-http-server-metrics/availability"
    options:
      filter: exported_service="alertmanager-api",service="alertgram"
```

[opentelemetry/http-server-metrics]: https://opentelemetry.io/docs/specs/semconv/http/http-metrics/#metric-httpserverrequestduration
[example]: https://github.com/open-telemetry/opentelemetry.io/blob/8411a8ca60efc531f0378e6835d3a311b4d0518e/content/en/docs/collector/internal-telemetry.md?plain=1#L305C3-L305C24

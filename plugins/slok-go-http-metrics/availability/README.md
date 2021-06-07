# slok/go-http-metrics availability

Availability plugin for [slok/go-http-metrics].

Uses the API HTTP response status codes to measure the events as good or bad. It counts an error event when HTTP response status code is >=500 or 429.

In other words, it counts as good events the <500 and !429 HTTP response status codes.

## Options

- `filter`: (**Required**) A prometheus filter string using concatenated labels (e.g: `exported_service="alertmanager-api",service="alertgram"`)

## Metric requirements

- `http_request_duration_seconds_count`: From [slok/go-http-metrics].

## Usage examples

```yaml
sli:
  plugin:
    id: "sloth-common/slok-go-http-metrics/availability"
    options:
      filter: exported_service="alertmanager-api",service="alertgram"
```

[slok/go-http-metrics]: https://github.com/slok/go-http-metrics

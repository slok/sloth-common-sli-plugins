# slok/go-http-metrics latency

Latency plugin for [slok/go-http-metrics].

Uses the HTTP response histogram to measure the events as good or bad. It counts as error events that don't fall in the provided bucket.

For example if the bucket is `0.25`, We will measure as an error, the requests that fall in the buckets greater than `0.25`, in other words, that took longer than `250ms` .

## Options

- `bucket`: (**Required**) The max latency allowed histogram bucket.
- `filter`: (**Required**) A prometheus filter string using concatenated labels (e.g: `exported_service="alertmanager-api",service="alertgram"`)
- `exclude_errors`: (**Optional**) Boolean that will exclude errored requests from valid events when measuring latency requests.

## Metric requirements

- `http_request_duration_seconds_bucket`: From [slok/go-http-metrics].
- `http_request_duration_seconds_count`: From [slok/go-http-metrics].

## Usage examples

### Default

```yaml
sli:
  plugin:
    id: "sloth-common/slok-go-http-metrics/latency"
    options:
      filter: exported_service="alertmanager-api",service="alertgram"
      bucket: "1"
```

### Excluding errors (5xx)

```yaml
sli:
  plugin:
    id: "sloth-common/slok-go-http-metrics/latency"
    options:
      filter: exported_service="alertmanager-api",service="alertgram"
      bucket: "0.25"
      exclude_errors: "true"
```

[slok/go-http-metrics]: https://github.com/slok/go-http-metrics

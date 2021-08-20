# Traefik V1 latency

Latency plugin for [Traefik V1][traefik] backends.

Uses Traefik v1 backend metrics to get the latency on the serving backends.

## Options

- `bucket`: (**Required**) The max latency allowed histogram bucket.
- `backend_regex`: (**required**) Regex to match the traefik backends.
- `filter`: (**Optional**) A prometheus filter string using concatenated labels
- `exclude_errors`: (**Optional**) Boolean that will exclude errored requests from valid events when measuring latency requests.

## Metric requirements

- `traefik_backend_request_duration_seconds_bucket`: From [traefik].
- `traefik_backend_request_duration_seconds_count`: From [traefik].

## Usage examples

### Without filter

```yaml
sli:
  plugin:
    id: "sloth-common/traefik/v1/latency"
    options:
      backend_regex: "^github.com/slok/sloth/?$"
      bucket: "0.3"
```

### With filters

```yaml
sli:
  plugin:
    id: "sloth-common/traefik/v1/latency"
    options:
      backend_regex: "^github.com/slok/sloth/?$"
      bucket: "0.3"
      filter: method="GET"
```

### Excluding errors (5xx)

```yaml
sli:
  plugin:
    id: "sloth-common/traefik/v1/latency"
    options:
      backend_regex: "^github.com/slok/sloth/?$"
      bucket: "0.3"
      filter: method="GET"
      exclude_errors: "true"
```

[traefik]: https://doc.traefik.io/traefik/v1.7/

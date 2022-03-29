# Traefik V1 availability

Availability plugin for [Traefik V1][traefik] backends.

Uses Traefik v1 backend metrics to get the correct and invalid availability on the serving backends.

## Options

- `filter`: (**Optional**) A prometheus filter string using concatenated labels
- `backend_regex`: (**required**) Regex to match the traefik backends.

## Metric requirements

- `traefik_backend_requests_total`: From [traefik].

## Usage examples

### Without filter

```yaml
sli:
  plugin:
    id: "sloth-common/traefik/v1/availability"
    options:
      backend_regex: "^github.com/slok/sloth/?$"
```

### With filters

```yaml
sli:
  plugin:
    id: "sloth-common/traefik/v1/availability"
    options:
      backend_regex: "^github.com/slok/sloth/?$"
      filter: method="GET"
```

[traefik]: https://doc.traefik.io/traefik/v1.7/

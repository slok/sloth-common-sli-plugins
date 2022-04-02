# Traefik V2 availability

Availability plugin for [Traefik V2][traefik] services.

Uses Traefik v2 service metrics to get the correct and invalid availability on the serving services.

## Options

- `filter`: (**Optional**) A prometheus filter string using concatenated labels
- `service_name_regex`: (**required**) Regex to match the traefik services.

## Metric requirements

- `traefik_service_requests_total`: From [traefik].

## Usage examples

### Without filter

```yaml
sli:
  plugin:
    id: "sloth-common/traefik/v2/availability"
    options:
      service_name_regex: "^default-slok-sloth$"
```

### With filters

```yaml
sli:
  plugin:
    id: "sloth-common/traefik/v2/availability"
    options:
      service_name_regex: "^default-slok-sloth$"
      filter: method="GET"
```

[traefik]: https://doc.traefik.io/traefik/v2.6/

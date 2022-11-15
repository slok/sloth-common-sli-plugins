# HAProxy availability

Availability plugin for the haproxy 2.x prometheus metrics.

Uses the standard prometheus metrics from a haproxy 2.x release to create an sli
against the frontend metrics based on http response codes.

In other words, it counts as good events the <500  HTTP response status codes.

## Options

- `filter`: (**Optional**) A prometheus filter string using concatenated labels (e.g: `instance=~"hostname.pattern.*"`)

## Metric requirements

- `haproxy_frontend_http_responses_total`.

## Usage examples

### Without filter

```yaml
sli:
  plugin:
    id: "sloth-common/haproxy/frontend/availability"
```

### With custom filter

```yaml
sli:
  plugin:
    id: "sloth-common/haproxy/frontend/availability"
    options:
      filter: instance=~"web-front-load-balancers.*"
```

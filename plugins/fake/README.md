# Fake

A plugin that can fake a burning rate based on speed/rate factor options. A good plugin for exploring/learning how SLOs, SLIs, error budget and SLO alerting works.

You can get more information on the [SRE workbook `Alert on burn rate`][examples-sre-book] section.

## Options

- `burn_rate`: (**Required**) A number that tells the burn rate factor (e.g: `1`, `2`, `10`...).
- `jitter_percent`: (**Optional**) A percent number that will add/remove jitter on the burned rate.

## Metric requirements

Doesn't have any metric requirements.

## Usage examples

### 0.5x speed `30d` window, consumed in `60d`

```yaml
sli:
  plugin:
    id: "sloth-common/fake"
    options:
      burn_rate: "0.5"
```

### 1x speed `30d` window, consumed in `30d`

```yaml
sli:
  plugin:
    id: "sloth-common/fake"
    options:
      burn_rate: "1"
```

### 2x speed `30d` window, consumed in `15d`

```yaml
sli:
  plugin:
    id: "sloth-common/fake"
    options:
      burn_rate: "2"
```

### 10x speed `30d` window, consumed in `3d`

```yaml
sli:
  plugin:
    id: "sloth-common/fake"
    options:
      burn_rate: "10"
```

### 10x speed `30d` window, consumed in `43m`

```yaml
sli:
  plugin:
    id: "sloth-common/fake"
    options:
      burn_rate: "1000"
```

### 1x speed `30d` window, consumed in `30d` using jitter

```yaml
sli:
  plugin:
    id: "sloth-common/fake"
    options:
      burn_rate: "1"
      jitter_percent: "10"
```

[examples-sre-book]: https://sre.google/workbook/alerting-on-slos/#burn_rates_and_time_to_complete_budget_ex

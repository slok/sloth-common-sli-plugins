# Kubernetes apiserver latency

Latency plugin for the Kubernetes apiserver.

Uses the API HTTP response histogram to measure the events as good or bad. It counts as error events that don't fall in the provided bucket.

For example if the bucket is `0.2`, We will measure as an error, the requests that fall in the buckets greater than `0.2`, in other words, that took longer than `200ms` .

## Options

- `bucket`: (**Required**) The max latency allowed hitogram bucket.
- `filter`: (**Optional**) A prometheus filter string using concatenated labels (e.g: `job="k8sapiserver",env="production",cluster="k8s-42"`)

## Metric requirements

- `apiserver_request_duration_seconds_count`.
- `apiserver_request_duration_seconds_bucket`.

## Usage examples

### Don't allow requests >50ms

```yaml
sli:
  plugin:
    id: "sloth-common/kubernetes/apiserver/latency"
    options:
      bucket: "0.05"
```

### Don't allow requests >200ms

```yaml
sli:
  plugin:
    id: "sloth-common/kubernetes/apiserver/latency"
    options:
      bucket: "0.2"
```

### Don't allow requests >1s

```yaml
sli:
  plugin:
    id: "sloth-common/kubernetes/apiserver/latency"
    options:
      bucket: "1"
```

### With custom filter

```yaml
sli:
  plugin:
    id: "sloth-common/kubernetes/apiserver/latency"
    options:
      bucket: "0.2"
      filter: job="k8sapiserver",env="production",cluster="k8s-42"
```

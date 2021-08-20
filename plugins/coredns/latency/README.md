# CoreDNS latency

Request/response latency plugin for [CoreDNS].

## Options

- `bucket`: (**Required**) The max latency allowed histogram bucket.
- `filter`: (**Optional**) A prometheus filter string using concatenated labels (e.g: `exported_service="alertmanager-api",service="alertgram"`)

## Metric requirements

- `coredns_dns_request_duration_seconds_bucket`: From [coreDNS].
- `coredns_dns_request_duration_seconds_count`: From [coreDNS].

## Usage examples

### Without filter

```yaml
sli:
  plugin:
    id: "sloth-common/coredns/latency"
    options:
      bucket: "0.25"
```

### With filter

```yaml
sli:
  plugin:
    id: "sloth-common/coredns/latency"
    options:
      bucket: "0.032"
      filter: job="kube-dns",server="dns://:53",service="kube-dns", zone="."
```

[coredns]: https://coredns.io/

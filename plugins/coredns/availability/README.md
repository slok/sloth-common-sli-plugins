# CoreDNS availability

Availability plugin for [CoreDNS].

Uses the coreDNS dns response metrics and error [rcodes] to get the correct and invalid availability.

By default the plugin will measure as errors the `SERVFAIL` rcodes, however this can be overwrite with a regex using `custom_rcode_regex` option.

## Options

- `filter`: (**Optional**) A prometheus filter string using concatenated labels (e.g: `exported_service="alertmanager-api",service="alertgram"`)
- `custom_rcode_regex`: (**Optional**) Custom regex to match error rcodes (e.g: `(NXDOMAIN|SERVFAIL|FORMERR)`). More codes [here][rcodes-list]. If not set it will be used `SERVFAIL` by default.

## Metric requirements

- `coredns_dns_responses_total`: From [coreDNS].

## Usage examples

### Without filter nor custom rcode regex

```yaml
sli:
  plugin:
    id: "sloth-common/coredns/availability"
```

### Default rcode and custom filter

```yaml
sli:
  plugin:
    id: "sloth-common/coredns/availability"
    options:
      filter: job="kube-dns",server="dns://:53",service="kube-dns", zone="."
```

### Custom rcode without filters

```yaml
sli:
  plugin:
    id: "sloth-common/coredns/availability"
    options:
      custom_rcode_regex: (NXDOMAIN|SERVFAIL|FORMERR)
```

### Custom rcode with filters

```yaml
sli:
  plugin:
    id: "sloth-common/coredns/availability"
    options:
      filter: job="kube-dns",server="dns://:53",service="kube-dns", zone="."
      custom_rcode_regex: (NXDOMAIN|SERVFAIL|FORMERR)
```

[coredns]: https://coredns.io/
[rcodes]: https://www.iana.org/assignments/dns-parameters/dns-parameters.xhtml#dns-parameters-6
[rcodes-list]: https://github.com/miekg/dns/blob/ab67aa64230094bdd0167ee5360e00e0a250a3ac/msg.go#L137-L159

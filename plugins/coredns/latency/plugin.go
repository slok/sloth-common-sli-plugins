package latency

import (
	"bytes"
	"context"
	"fmt"
	"strconv"
	"strings"
	"text/template"
)

const (
	// SLIPluginVersion is the version of the plugin spec.
	SLIPluginVersion = "prometheus/v1"
	// SLIPluginID is the registering ID of the plugin.
	SLIPluginID = "sloth-common/coredns/latency"
)

// We use a combination of `>0` on the division denominator (so we don't divide by 0 to avoid getting `NaN`)
// and then `OR on vector(1)` when the `>0` took effect to return a `0` error ratio.
var queryTpl = template.Must(template.New("").Option("missingkey=error").Parse(`
1 - ((
  sum(rate(coredns_dns_request_duration_seconds_bucket{ {{ .filterError }}le="{{ .bucket }}" }[{{"{{ .window }}"}}]))
  /
  (sum(rate(coredns_dns_request_duration_seconds_count{ {{ .filterTotal }} }[{{"{{ .window }}"}}])) > 0)
) OR on() vector(1))
`))

// SLIPlugin will return a query that will calculate the the request handling latency in coreDNS.
func SLIPlugin(ctx context.Context, meta, labels, options map[string]string) (string, error) {
	bucket, err := getBucket(options)
	if err != nil {
		return "", fmt.Errorf(`could not get bucket: %w`, err)
	}

	filter := getFilter(options)
	filterTotal := filter
	filterError := filter
	if filterError != "" {
		filterError = filter + ","
	}

	// Create query.
	var b bytes.Buffer
	data := map[string]string{
		"bucket":      bucket,
		"filterTotal": filterTotal,
		"filterError": filterError,
	}
	err = queryTpl.Execute(&b, data)
	if err != nil {
		return "", fmt.Errorf("could not render query template: %w", err)
	}

	return b.String(), nil
}

func getBucket(options map[string]string) (string, error) {
	bucket := options["bucket"]
	if bucket == "" {
		return "", fmt.Errorf(`"bucket" option is required`)
	}

	_, err := strconv.ParseFloat(bucket, 64)
	if err != nil {
		return "", fmt.Errorf("not a valid bucket, can't parse to float64: %w", err)
	}

	return bucket, nil
}

func getFilter(options map[string]string) string {
	filter := options["filter"]
	filter = strings.Trim(filter, "{},")

	return filter
}

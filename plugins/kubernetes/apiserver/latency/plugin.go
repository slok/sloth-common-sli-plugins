package latency

import (
	"bytes"
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"text/template"
)

const (
	// SLIPluginVersion is the version of the plugin spec.
	SLIPluginVersion = "prometheus/v1"
	// SLIPluginID is the registering ID of the plugin.
	SLIPluginID = "sloth-common/kubernetes/apiserver/latency"
)

var queryTpl = template.Must(template.New("").Option("missingkey=error").Parse(`
(
  sum(rate(apiserver_request_duration_seconds_count{ {{.filter}}verb!="WATCH" }[{{"{{.window}}"}}]))
  -
  sum(rate(apiserver_request_duration_seconds_bucket{ {{.filter}}le="{{.bucket}}",verb!="WATCH" }[{{"{{.window}}"}}]))
)
/
sum(rate(apiserver_request_duration_seconds_count{ {{.filter}}verb!="WATCH" }[{{"{{.window}}"}}]))
`))

// SLIPlugin will return a query that will return the latency by using the http requests histogram buckets.
// We will count as errors the requests that are not on the expected bucket.
func SLIPlugin(ctx context.Context, meta, labels, options map[string]string) (string, error) {
	bucket, err := getBucket(options)
	if err != nil {
		return "", fmt.Errorf(`could not get bucket: %w`, err)
	}

	filter, err := getFilter(options)
	if err != nil {
		return "", fmt.Errorf("could not get filter: %w", err)
	}

	if filter != "" {
		filter += ","
	}

	// Create query.
	var b bytes.Buffer
	data := map[string]string{
		"filter": filter,
		"bucket": bucket,
	}
	err = queryTpl.Execute(&b, data)
	if err != nil {
		return "", fmt.Errorf("could not render query template: %w", err)
	}

	return b.String(), nil
}

var filterRegex = regexp.MustCompile(`(?m)^{?([^=]+="[^=,"]+",)*([^=]+="[^=,"]+")$`)

func getFilter(options map[string]string) (string, error) {
	filter, ok := options["filter"]
	if !ok || (ok && filter == "") {
		return "", nil
	}

	// Sanitize and check filter.
	filter = strings.Trim(filter, "{},")
	match := filterRegex.MatchString(filter)
	if !match {
		return "", fmt.Errorf("invalid prometheus filter: %s", filter)
	}

	return filter, nil
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

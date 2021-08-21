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
	SLIPluginID = "sloth-common/slok-go-http-metrics/latency"
)

var queryTpl = template.Must(template.New("").Option("missingkey=error").Parse(`
1 - (
  sum(rate(http_request_duration_seconds_bucket{ {{.filterBucket}}le="{{.bucket}}" }[{{"{{.window}}"}}]))
  /
  sum(rate(http_request_duration_seconds_count{ {{.filterCount}} }[{{"{{.window}}"}}]))
)
`))

// SLIPlugin will return a query that will return the latency error based on https://github.com/slok/go-http-metrics
// response request latency buckets.
// Counts as an error event the requests that are not part of the required latency bucket.
// Accepts "exclude_errors" bool option so we don't count the errors as valid events.
func SLIPlugin(ctx context.Context, meta, labels, options map[string]string) (string, error) {
	bucket, err := getBucket(options)
	if err != nil {
		return "", fmt.Errorf(`could not get bucket: %w`, err)
	}

	filter, err := getFilter(options)
	if err != nil {
		return "", fmt.Errorf("could not get filter: %w", err)
	}

	// Add error code exclusion filter if required.
	excludeErrors, err := getExcludeErrors(options)
	if err != nil {
		return "", fmt.Errorf(`could not get exclude_errors: %w`, err)
	}

	if excludeErrors {
		filter = strings.Join([]string{filter, `code!~"5.."`}, ",")
	}

	filterCount := filter
	filterBucket := filter
	if filterBucket != "" {
		filterBucket = filter + ","
	}

	// Create query.
	var b bytes.Buffer
	data := map[string]string{
		"filterBucket": filterBucket,
		"filterCount":  filterCount,
		"bucket":       bucket,
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
		return "", fmt.Errorf("filter is required")
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

func getExcludeErrors(options map[string]string) (bool, error) {
	excludeErrorsS, ok := options["exclude_errors"]
	if !ok || (ok && excludeErrorsS == "") {
		return false, nil
	}

	excludeErrors, err := strconv.ParseBool(excludeErrorsS)
	if err != nil {
		return false, fmt.Errorf("not a valid exclude_errors, can't parse to bool: %w", err)
	}

	return excludeErrors, nil
}

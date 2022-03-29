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
	SLIPluginID = "sloth-common/istio/v1/latency"
)

var queryTpl = template.Must(template.New("").Option("missingkey=error").Parse(`
1 - ((
  sum(rate(istio_request_duration_milliseconds_bucket{ {{.filter}}destination_service_name="{{.service}}",destination_service_namespace="{{.namespace}}",le="{{.bucket}}" }[{{"{{.window}}"}}]))
  /
  (sum(rate(istio_request_duration_milliseconds_count{ {{.filter}}destination_service_name="{{.service}}",destination_service_namespace="{{.namespace}}" }[{{"{{.window}}"}}])) > 0)
) OR on() vector(1))
`))

// SLIPlugin will return a query that will return the latency based on istio V1 request metrics.
// Counts as an error event the requests that are not part of the required latency bucket.
// Accepts "exclude_errors" bool option so we don't count the errors as valid events.
func SLIPlugin(ctx context.Context, meta, labels, options map[string]string) (string, error) {
	bucket, err := getBucket(options)
	if err != nil {
		return "", fmt.Errorf(`could not get bucket: %w`, err)
	}

	service, err := getRequiredStringValue("service", options)
	if err != nil {
		return "", fmt.Errorf("could not get service: %w", err)
	}

	namespace, err := getRequiredStringValue("namespace", options)
	if err != nil {
		return "", fmt.Errorf("could not get namespace: %w", err)
	}

	filter, err := getFilter(options)
	if err != nil {
		return "", fmt.Errorf("could not get filter: %w", err)
	}

	// Create query.
	var b bytes.Buffer
	data := map[string]string{
		"filter":    filter,
		"service":   service,
		"namespace": namespace,
		"bucket":    bucket,
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

const errorFilter = `code!~"5.."`

func getFilter(options map[string]string) (string, error) {
	filter := options["filter"]
	filter = strings.Trim(filter, "{},")

	// Add error code exclusion filter if required.
	excludeErrors, err := getExcludeErrors(options)
	if err != nil {
		return "", fmt.Errorf(`could not get exclude_errors: %w`, err)
	}

	switch {
	case !excludeErrors && filter == "":
		return "", nil
	case !excludeErrors && filter != "":
		return filter + ",", nil
	case excludeErrors && filter == "":
		return errorFilter + `,`, nil
	case excludeErrors && filter != "":
		return filter + `,` + errorFilter + `,`, nil
	}

	return "", fmt.Errorf("invalid case")
}

func getExcludeErrors(options map[string]string) (bool, error) {
	excludeErrorsS := options["exclude_errors"]
	if excludeErrorsS == "" {
		return false, nil
	}

	excludeErrors, err := strconv.ParseBool(excludeErrorsS)
	if err != nil {
		return false, fmt.Errorf("not a valid exclude_errors, can't parse to bool: %w", err)
	}

	return excludeErrors, nil
}

func getRequiredStringValue(key string, options map[string]string) (string, error) {
	value, exists := options[key]
	value = strings.TrimSpace(value)

	if !exists || value == "" {
		return "", fmt.Errorf("%s is required", key)
	}

	return value, nil
}

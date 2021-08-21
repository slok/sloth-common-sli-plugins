package availability

import (
	"bytes"
	"context"
	"fmt"
	"regexp"
	"strings"
	"text/template"
)

const (
	// SLIPluginVersion is the version of the plugin spec.
	SLIPluginVersion = "prometheus/v1"
	// SLIPluginID is the registering ID of the plugin.
	SLIPluginID = "sloth-common/slok-go-http-metrics/availability"
)

var queryTpl = template.Must(template.New("").Option("missingkey=error").Parse(`
sum(rate(http_request_duration_seconds_count{ {{.filterError}}code=~"(5..|429)" }[{{"{{.window}}"}}]))
/
sum(rate(http_request_duration_seconds_count{ {{.filterTotal}} }[{{"{{.window}}"}}]))
`))

// SLIPlugin will return a query that will return the availability error based on https://github.com/slok/go-http-metrics
// status response codes.
// Counts as an error event the requests that have >=500 and 429 status codes.
func SLIPlugin(ctx context.Context, meta, labels, options map[string]string) (string, error) {
	filter, err := getFilter(options)
	if err != nil {
		return "", fmt.Errorf("could not get filter: %w", err)
	}

	filterTotal := filter
	filterError := filter
	if filterError != "" {
		filterError = filter + ","
	}

	// Create query.
	var b bytes.Buffer
	data := map[string]string{
		"filterError": filterError,
		"filterTotal": filterTotal,
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

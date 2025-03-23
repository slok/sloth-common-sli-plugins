package availability

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"text/template"
)

const (
	// SLIPluginVersion is the version of the plugin spec.
	SLIPluginVersion = "prometheus/v1"
	// SLIPluginID is the registering ID of the plugin.
	SLIPluginID = "sloth-common/haproxy/frontend/availability"
)

var queryTpl = template.Must(template.New("").Option("missingkey=error").Parse(`
sum(rate(haproxy_frontend_http_responses_total{ {{.filterError}}code="5xx" }[{{"{{.window}}"}}]))
/
sum(rate(haproxy_frontend_http_responses_total{ {{.filterTotal}} }[{{"{{.window}}"}}]))
`))

// SLIPlugin will return a query that will return the availability error based on frontend responses
// Counts as an error event the requests that have 5xx status code.
func SLIPlugin(ctx context.Context, meta, labels, options map[string]string) (string, error) {
	filter, err := getFilter(options)
	if err != nil {
		return "", fmt.Errorf("could not get filter: %w", err)
	}

	filterTotal := filter
	filterError := filter

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

func getFilter(options map[string]string) (string, error) {
	filter, ok := options["filter"]
	if !ok || (ok && filter == "") {
		return "", nil
	}
	filter = strings.Trim(filter, "{},")
	if filter != "" {
		filter += ","
	}

	return filter, nil
}

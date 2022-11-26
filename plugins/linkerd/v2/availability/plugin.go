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
	SLIPluginID = "sloth-common/linkerd/v2/availability"
)

var queryTpl = template.Must(template.New("").Option("missingkey=error").Parse(`
(
  sum(increase(response_total{ {{.filter}}classification="failure", direction="inbound" }[{{"{{.window}}"}}]))
  /
  sum(increase(response_total{ {{.filter}}direction="inbound" }[{{"{{.window}}"}}]))
) OR on() vector(0)
`))

// SLIPlugin will return a query that always will be 0.
func SLIPlugin(ctx context.Context, meta, labels, options map[string]string) (string, error) {
	filter, err := getFilter(options)
	if err != nil {
		return filter, err
	}
	if filter != "" {
		filter += ","
	}
	// Create query.
	var b bytes.Buffer
	data := map[string]string{
		"filter": filter,
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

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
	SLIPluginID = "sloth-common/traefik/v1/availability"
)

var queryTpl = template.Must(template.New("").Option("missingkey=error").Parse(`
(
  sum(rate(traefik_backend_requests_total{ {{.filter}}backend=~"{{.backendRegex}}",code=~"(5..|429)" }[{{"{{.window}}"}}]))
  /
  (sum(rate(traefik_backend_requests_total{ {{.filter}}backend=~"{{.backendRegex}}" }[{{"{{.window}}"}}])) > 0)
) OR on() vector(0)
`))

// SLIPlugin will return a query that will return the availability error based on traefik V1 backend metrics.
func SLIPlugin(ctx context.Context, meta, labels, options map[string]string) (string, error) {
	backendRegex, err := getBackendRegex(options)
	if err != nil {
		return "", fmt.Errorf("could not get backend regex: %w", err)
	}

	// Create query.
	var b bytes.Buffer
	data := map[string]string{
		"filter":       getFilter(options),
		"backendRegex": backendRegex,
	}
	err = queryTpl.Execute(&b, data)
	if err != nil {
		return "", fmt.Errorf("could not render query template: %w", err)
	}

	return b.String(), nil
}

func getFilter(options map[string]string) string {
	filter := options["filter"]
	filter = strings.Trim(filter, "{},")
	if filter != "" {
		filter += ","
	}

	return filter
}

func getBackendRegex(options map[string]string) (string, error) {
	backend := options["backend_regex"]
	backend = strings.TrimSpace(backend)

	if backend == "" {
		return "", fmt.Errorf("backend is required")
	}

	_, err := regexp.Compile(backend)
	if err != nil {
		return "", fmt.Errorf("invalid regex: %w", err)
	}

	return backend, nil
}

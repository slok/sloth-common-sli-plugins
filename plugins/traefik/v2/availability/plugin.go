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
	SLIPluginID = "sloth-common/traefik/v2/availability"
)

var queryTpl = template.Must(template.New("").Option("missingkey=error").Parse(`
(
  sum(rate(traefik_service_requests_total{ {{.filter}}service=~"{{.serviceName}}",code=~"(5..)" }[{{"{{.window}}"}}]))
  /
  (sum(rate(traefik_service_requests_total{ {{.filter}}service=~"{{.serviceName}}" }[{{"{{.window}}"}}])) > 0)
) OR on() vector(0)
`))

// SLIPlugin will return a query that will return the availability error based on traefik V1 service metrics.
func SLIPlugin(ctx context.Context, meta, labels, options map[string]string) (string, error) {
	service_name, err := getServiceName(options)
	if err != nil {
		return "", fmt.Errorf("could not get service name: %w", err)
	}

	// Create query.
	var b bytes.Buffer
	data := map[string]string{
		"filter":      getFilter(options),
		"serviceName": service_name,
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

func getServiceName(options map[string]string) (string, error) {
	service_name := options["service_name_regex"]
	service_name = strings.TrimSpace(service_name)

	if service_name == "" {
		return "", fmt.Errorf("service name is required")
	}

	_, err := regexp.Compile(service_name)
	if err != nil {
		return "", fmt.Errorf("invalid regex: %w", err)
	}

	return service_name, nil
}

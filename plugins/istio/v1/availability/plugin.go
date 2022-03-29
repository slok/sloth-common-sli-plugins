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
	SLIPluginID = "sloth-common/istio/v1/availability"
)

var queryTpl = template.Must(template.New("").Option("missingkey=error").Parse(`
(
  sum(rate(istio_requests_total{ {{.filter}}destination_service_name="{{.service}}",destination_service_namespace="{{.namespace}}",response_code=~"(5..|429)" }[{{"{{.window}}"}}])) 
  /          
  (sum(rate(istio_requests_total{ {{.filter}}destination_service_name="{{.service}}",destination_service_namespace="{{.namespace}}" }[{{"{{.window}}"}}])) > 0)
) OR on() vector(0)
`))

// SLIPlugin will return a query that will return the availability error based on istio V1 request metrics.
func SLIPlugin(ctx context.Context, meta, labels, options map[string]string) (string, error) {

	service, err := getRequiredStringValue("service", options)
	if err != nil {
		return "", fmt.Errorf("could not get service: %w", err)
	}

	namespace, err := getRequiredStringValue("namespace", options)
	if err != nil {
		return "", fmt.Errorf("could not get namespace: %w", err)
	}

	// Create query.
	var b bytes.Buffer
	data := map[string]string{
		"filter":    getFilter(options),
		"service":   service,
		"namespace": namespace,
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

func getRequiredStringValue(key string, options map[string]string) (string, error) {
	value := options[key]
	value = strings.TrimSpace(value)

	if value == "" {
		return "", fmt.Errorf("%s is required", key)
	}

	return value, nil
}

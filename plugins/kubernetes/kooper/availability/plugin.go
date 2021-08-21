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
	SLIPluginID = "sloth-common/kubernetes/kooper/availability"
)

// We use a combination of "> 0" on the division denominator (so we don't divide by 0 to avoid getting "NaN")
// and then "OR on vector(0)" when the ">0" took effect to return a "0" error ratio.
var queryTpl = template.Must(template.New("").Option("missingkey=error").Parse(`
(
  sum(rate(kooper_controller_processed_event_duration_seconds_count{ {{ .filter }}controller="{{ .controller }}",success="false" }[{{"{{ .window }}"}}]))
  /
  ((
    sum(rate(kooper_controller_processed_event_duration_seconds_count{ {{ .filter }}controller="{{ .controller }}" }[{{"{{ .window }}"}}]))
    -
    (sum(rate(kooper_controller_queued_events_total{ {{ .filter }}controller="{{ .controller }}",requeue="true" }[{{"{{ .window }}"}}])) OR on() vector(0))
  ) > 0)
) OR on() vector(0)
`))

// SLIPlugin will return a query that will return the availability by using the kooper queued and handled metrics.
// Requeued events will be subtracted total so we get a real ratio of errors per object independently of the retries.
func SLIPlugin(ctx context.Context, meta, labels, options map[string]string) (string, error) {
	controller, err := getController(options)
	if err != nil {
		return "", fmt.Errorf(`could not get controller: %w`, err)
	}

	// Create query.
	var b bytes.Buffer
	data := map[string]string{
		"controller": controller,
		"filter":     getFilter(options),
	}
	err = queryTpl.Execute(&b, data)
	if err != nil {
		return "", fmt.Errorf("could not render query template: %w", err)
	}

	return b.String(), nil
}

func getController(options map[string]string) (string, error) {
	controller, ok := options["controller"]
	if !ok || (ok && controller == "") {
		return "", fmt.Errorf(`"controller" option is required`)
	}

	return controller, nil
}

func getFilter(options map[string]string) string {
	filter := options["filter"]
	filter = strings.Trim(filter, "{},")
	if filter != "" {
		filter += ","
	}

	return filter
}

package queuecongestion

import (
	"bytes"
	"context"
	"fmt"
	"strconv"
	"strings"
	"text/template"
)

const (
	SLIPluginVersion = "prometheus/v1"
	SLIPluginID      = "sloth-common/kubernetes/kooper/queue-congestion"
)

var queryTpl = template.Must(template.New("").Option("missingkey=error").Parse(`
1 - (
  sum(rate(kooper_controller_event_in_queue_duration_seconds_bucket{ {{ .filter }}controller="{{ .controller }}",le="{{ .bucket }}" }[{{"{{ .window }}"}}]))
  /
  sum(rate(kooper_controller_event_in_queue_duration_seconds_count{ {{ .filter }}controller="{{ .controller }}" }[{{"{{ .window }}"}}]))
)
`))

// SLIPlugin will return a query that will return the congestion on the event queue by using the kooper's
// event queued time measurements.
func SLIPlugin(ctx context.Context, meta, labels, options map[string]string) (string, error) {
	bucket, err := getBucket(options)
	if err != nil {
		return "", fmt.Errorf(`could not get bucket: %w`, err)
	}

	controller, err := getController(options)
	if err != nil {
		return "", fmt.Errorf(`could not get controller: %w`, err)
	}

	// Create query.
	var b bytes.Buffer
	data := map[string]string{
		"controller": controller,
		"bucket":     bucket,
		"filter":     getFilter(options),
	}
	err = queryTpl.Execute(&b, data)
	if err != nil {
		return "", fmt.Errorf("could not render query template: %w", err)
	}

	return b.String(), nil
}

func getBucket(options map[string]string) (string, error) {
	bucket, ok := options["bucket"]
	if !ok || (ok && bucket == "") {
		return "", fmt.Errorf(`"bucket" option is required`)
	}

	_, err := strconv.ParseFloat(bucket, 64)
	if err != nil {
		return "", fmt.Errorf("not a valid bucket, can't parse to float64: %w", err)
	}

	return bucket, nil
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

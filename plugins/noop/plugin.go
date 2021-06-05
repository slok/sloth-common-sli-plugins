package noop

import (
	"context"
)

const (
	SLIPluginVersion = "prometheus/v1"
	SLIPluginID      = "sloth_common/noop"
)

// SLIPlugin will return a query that always will be 0.
func SLIPlugin(ctx context.Context, meta, labels, options map[string]string) (string, error) {
	return `max_over_time(vector(0)[{{.window}}:])`, nil
}

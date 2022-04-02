package availability_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	latency "github.com/slok/sloth-common-sli-plugins/plugins/traefik/v2/latency"
)

func TestSLIPlugin(t *testing.T) {
	tests := map[string]struct {
		meta     map[string]string
		labels   map[string]string
		options  map[string]string
		expQuery string
		expErr   bool
	}{
		"Without service name, should fail.": {
			options: map[string]string{"bucket": "0.5"},
			expErr:  true,
		},

		"An invalid service name query, should fail.": {
			options: map[string]string{"bucket": "0.5", "service_name_regex": "([xyz"},
			expErr:  true,
		},

		"An empty service name query, should fail.": {
			options: map[string]string{"bucket": "0.5", "service_name_regex": ""},
			expErr:  true,
		},

		"Not having a bucket, should fail.": {
			options: map[string]string{"service_name_regex": ".*"},
			expErr:  true,
		},

		"Having an invalid bucket, should fail.": {
			options: map[string]string{"service_name_regex": ".*", "bucket": "something"},
			expErr:  true,
		},

		"Having invalid exclude_errors, should fail.": {
			options: map[string]string{"service_name_regex": ".*", "bucket": "0.5", "exclude_errors": "10"},
			expErr:  true,
		},

		"Default use case, should return a valid query.": {
			options: map[string]string{"service_name_regex": "github.com/slok/sloth/?", "bucket": "0.5"},
			expQuery: `
1 - ((
  sum(rate(traefik_service_request_duration_seconds_bucket{ service=~"github.com/slok/sloth/?",le="0.5" }[{{.window}}]))
  /
  (sum(rate(traefik_service_request_duration_seconds_count{ service=~"github.com/slok/sloth/?" }[{{.window}}])) > 0)
) OR on() vector(1))
`,
		},

		"Using exclude errors, should return a valid query.": {
			options: map[string]string{"service_name_regex": "github.com/slok/sloth/?", "bucket": "0.5", "exclude_errors": "true"},
			expQuery: `
1 - ((
  sum(rate(traefik_service_request_duration_seconds_bucket{ code!~"5..",service=~"github.com/slok/sloth/?",le="0.5" }[{{.window}}]))
  /
  (sum(rate(traefik_service_request_duration_seconds_count{ code!~"5..",service=~"github.com/slok/sloth/?" }[{{.window}}])) > 0)
) OR on() vector(1))
`,
		},

		"Using filter, should return a valid query.": {
			options: map[string]string{
				"service_name_regex": "github.com/slok/sloth/?",
				"bucket":             "0.5",
				"filter":             `k1="v2",k2="v2"`,
			},
			expQuery: `
1 - ((
  sum(rate(traefik_service_request_duration_seconds_bucket{ k1="v2",k2="v2",service=~"github.com/slok/sloth/?",le="0.5" }[{{.window}}]))
  /
  (sum(rate(traefik_service_request_duration_seconds_count{ k1="v2",k2="v2",service=~"github.com/slok/sloth/?" }[{{.window}}])) > 0)
) OR on() vector(1))
`,
		},

		"Using filter and excluding errors, should return a valid query.": {
			options: map[string]string{
				"service_name_regex": "github.com/slok/sloth/?",
				"bucket":             "0.5",
				"filter":             `k1="v2",k2="v2"`,
				"exclude_errors":     "true",
			},
			expQuery: `
1 - ((
  sum(rate(traefik_service_request_duration_seconds_bucket{ k1="v2",k2="v2",code!~"5..",service=~"github.com/slok/sloth/?",le="0.5" }[{{.window}}]))
  /
  (sum(rate(traefik_service_request_duration_seconds_count{ k1="v2",k2="v2",code!~"5..",service=~"github.com/slok/sloth/?" }[{{.window}}])) > 0)
) OR on() vector(1))
`,
		},

		"Filter should be sanitized with ','.": {
			options: map[string]string{
				"service_name_regex": "github.com/slok/sloth/?",
				"bucket":             "0.5",
				"filter":             `k1="v2",k2="v2",`,
			},
			expQuery: `
1 - ((
  sum(rate(traefik_service_request_duration_seconds_bucket{ k1="v2",k2="v2",service=~"github.com/slok/sloth/?",le="0.5" }[{{.window}}]))
  /
  (sum(rate(traefik_service_request_duration_seconds_count{ k1="v2",k2="v2",service=~"github.com/slok/sloth/?" }[{{.window}}])) > 0)
) OR on() vector(1))
`,
		},

		"Filter should be sanitized with '{'.": {
			options: map[string]string{
				"service_name_regex": "github.com/slok/sloth/?",
				"bucket":             "0.5",
				"filter":             `{k1="v2",k2="v2",},`,
			},
			expQuery: `
1 - ((
  sum(rate(traefik_service_request_duration_seconds_bucket{ k1="v2",k2="v2",service=~"github.com/slok/sloth/?",le="0.5" }[{{.window}}]))
  /
  (sum(rate(traefik_service_request_duration_seconds_count{ k1="v2",k2="v2",service=~"github.com/slok/sloth/?" }[{{.window}}])) > 0)
) OR on() vector(1))
`,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)

			gotQuery, err := latency.SLIPlugin(context.TODO(), test.meta, test.labels, test.options)

			if test.expErr {
				assert.Error(err)
			} else if assert.NoError(err) {
				assert.Equal(test.expQuery, gotQuery)
			}
		})
	}
}

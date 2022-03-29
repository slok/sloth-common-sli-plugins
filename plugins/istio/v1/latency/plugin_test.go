package latency_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/slok/sloth-common-sli-plugins/plugins/istio/v1/latency"
)

func TestSLIPlugin(t *testing.T) {
	tests := map[string]struct {
		meta     map[string]string
		labels   map[string]string
		options  map[string]string
		expQuery string
		expErr   bool
	}{
		"Without namespace, should fail.": {
			options: map[string]string{"bucket": "500", "service": "test"},
			expErr:  true,
		},

		"An empty namespace query, should fail.": {
			options: map[string]string{"bucket": "500", "namespace": "", "service": "test"},
			expErr:  true,
		},

		"Without service, should fail.": {
			options: map[string]string{"bucket": "500", "namespace": "default"},
			expErr:  true,
		},

		"An empty service query, should fail.": {
			options: map[string]string{"bucket": "500", "namespace": "default", "service": ""},
			expErr:  true,
		},

		"Not having a bucket, should fail.": {
			options: map[string]string{"namespace": "default", "service": "test"},
			expErr:  true,
		},

		"Having an invalid bucket, should fail.": {
			options: map[string]string{"namespace": "default", "service": "test", "bucket": "something"},
			expErr:  true,
		},

		"Having invalid exclude_errors, should fail.": {
			options: map[string]string{"namespace": "default", "service": "test", "bucket": "500", "exclude_errors": "10"},
			expErr:  true,
		},

		"Default use case, should return a valid query.": {
			options: map[string]string{"namespace": "default", "service": "test", "bucket": "500"},
			expQuery: `
1 - ((
  sum(rate(istio_request_duration_milliseconds_bucket{ destination_service_name="test",destination_service_namespace="default",le="500" }[{{.window}}]))
  /
  (sum(rate(istio_request_duration_milliseconds_count{ destination_service_name="test",destination_service_namespace="default" }[{{.window}}])) > 0)
) OR on() vector(1))
`,
		},

		"Using exclude errors, should return a valid query.": {
			options: map[string]string{"namespace": "default", "service": "test", "bucket": "500", "exclude_errors": "true"},
			expQuery: `
1 - ((
  sum(rate(istio_request_duration_milliseconds_bucket{ code!~"5..",destination_service_name="test",destination_service_namespace="default",le="500" }[{{.window}}]))
  /
  (sum(rate(istio_request_duration_milliseconds_count{ code!~"5..",destination_service_name="test",destination_service_namespace="default" }[{{.window}}])) > 0)
) OR on() vector(1))
`,
		},

		"Using filter, should return a valid query.": {
			options: map[string]string{
				"namespace": "default",
				"service":   "test",
				"bucket":    "500",
				"filter":    `k1="v2",k2="v2"`,
			},
			expQuery: `
1 - ((
  sum(rate(istio_request_duration_milliseconds_bucket{ k1="v2",k2="v2",destination_service_name="test",destination_service_namespace="default",le="500" }[{{.window}}]))
  /
  (sum(rate(istio_request_duration_milliseconds_count{ k1="v2",k2="v2",destination_service_name="test",destination_service_namespace="default" }[{{.window}}])) > 0)
) OR on() vector(1))
`,
		},

		"Using filter and excluding errors, should return a valid query.": {
			options: map[string]string{
				"namespace":      "default",
				"service":        "test",
				"bucket":         "500",
				"filter":         `k1="v2",k2="v2"`,
				"exclude_errors": "true",
			},
			expQuery: `
1 - ((
  sum(rate(istio_request_duration_milliseconds_bucket{ k1="v2",k2="v2",code!~"5..",destination_service_name="test",destination_service_namespace="default",le="500" }[{{.window}}]))
  /
  (sum(rate(istio_request_duration_milliseconds_count{ k1="v2",k2="v2",code!~"5..",destination_service_name="test",destination_service_namespace="default" }[{{.window}}])) > 0)
) OR on() vector(1))
`,
		},

		"Filter should be sanitized with ','.": {
			options: map[string]string{
				"namespace": "default",
				"service":   "test",
				"bucket":    "500",
				"filter":    `k1="v2",k2="v2",`,
			},
			expQuery: `
1 - ((
  sum(rate(istio_request_duration_milliseconds_bucket{ k1="v2",k2="v2",destination_service_name="test",destination_service_namespace="default",le="500" }[{{.window}}]))
  /
  (sum(rate(istio_request_duration_milliseconds_count{ k1="v2",k2="v2",destination_service_name="test",destination_service_namespace="default" }[{{.window}}])) > 0)
) OR on() vector(1))
`,
		},

		"Filter should be sanitized with '{'.": {
			options: map[string]string{
				"namespace": "default",
				"service":   "test",
				"bucket":    "500",
				"filter":    `{k1="v2",k2="v2",},`,
			},
			expQuery: `
1 - ((
  sum(rate(istio_request_duration_milliseconds_bucket{ k1="v2",k2="v2",destination_service_name="test",destination_service_namespace="default",le="500" }[{{.window}}]))
  /
  (sum(rate(istio_request_duration_milliseconds_count{ k1="v2",k2="v2",destination_service_name="test",destination_service_namespace="default" }[{{.window}}])) > 0)
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

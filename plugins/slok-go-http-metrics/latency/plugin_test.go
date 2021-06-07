package latency_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/slok/sloth-common-sli-plugins/plugins/slok-go-http-metrics/latency"
)

func TestSLIPlugin(t *testing.T) {
	tests := map[string]struct {
		meta     map[string]string
		labels   map[string]string
		options  map[string]string
		expQuery string
		expErr   bool
	}{
		"Having a wrong filter, it should fail.": {
			options: map[string]string{"bucket": "0.2", "filter": "something="},
			expErr:  true,
		},

		"Not having a filter should fail.": {
			options: map[string]string{"bucket": "0.2"},
			expErr:  true,
		},

		"Not having a bucket should fail.": {
			options: map[string]string{"filter": `something="something"`},
			expErr:  true,
		},

		"Having a filter should return the query with the filters.": {
			options: map[string]string{"bucket": "0.2", "filter": `k1="v1",k2="v2"`},
			expQuery: `
1 - (
  sum(rate(http_request_duration_seconds_bucket{ k1="v1",k2="v2",le="0.2" }[{{.window}}]))
  /
  sum(rate(http_request_duration_seconds_count{ k1="v1",k2="v2" }[{{.window}}]))
)
`,
		},

		"Having exclude_errors options should return the query with the error exclusion.": {
			options: map[string]string{"bucket": "0.2", "filter": `k1="v1",k2="v2"`, "exclude_errors": "true"},
			expQuery: `
1 - (
  sum(rate(http_request_duration_seconds_bucket{ k1="v1",k2="v2",code!~"5..",le="0.2" }[{{.window}}]))
  /
  sum(rate(http_request_duration_seconds_count{ k1="v1",k2="v2",code!~"5.." }[{{.window}}]))
)
`,
		},

		"Having a filter with `{}` should return the query with the filters.": {
			options: map[string]string{"bucket": "0.2", "filter": `{k1="v1",k2="v2"}`},
			expQuery: `
1 - (
  sum(rate(http_request_duration_seconds_bucket{ k1="v1",k2="v2",le="0.2" }[{{.window}}]))
  /
  sum(rate(http_request_duration_seconds_count{ k1="v1",k2="v2" }[{{.window}}]))
)
`,
		},

		"Having a filter with `,` should return the query with the filters.": {
			options: map[string]string{"bucket": "0.2", "filter": `k1="v1",k2="v2",`},
			expQuery: `
1 - (
  sum(rate(http_request_duration_seconds_bucket{ k1="v1",k2="v2",le="0.2" }[{{.window}}]))
  /
  sum(rate(http_request_duration_seconds_count{ k1="v1",k2="v2" }[{{.window}}]))
)
`,
		},

		"Having a filter with `{}` and `,` should return the query with the filters.": {
			options: map[string]string{"bucket": "0.2", "filter": `{k1="v1",k2="v2",}`},
			expQuery: `
1 - (
  sum(rate(http_request_duration_seconds_bucket{ k1="v1",k2="v2",le="0.2" }[{{.window}}]))
  /
  sum(rate(http_request_duration_seconds_count{ k1="v1",k2="v2" }[{{.window}}]))
)
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

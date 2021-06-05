package latency_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/slok/sloth-common-sli-plugins/plugins/kubernetes/apiserver/latency"
)

func TestSLIPlugin(t *testing.T) {
	tests := map[string]struct {
		meta     map[string]string
		labels   map[string]string
		options  map[string]string
		expQuery string
		expErr   bool
	}{
		"Not having bucket option, it should fail.": {
			options: map[string]string{},
			expErr:  true,
		},

		"Having a wrong bucket, it should fail.": {
			options: map[string]string{
				"bucket": "zero point two",
			},
			expErr: true,
		},

		"Having a wrong filter, it should fail.": {
			options: map[string]string{
				"bucket": "0.2",
				"filter": "something=",
			},
			expErr: true,
		},

		"Not having a filter shouldn't fail and return the correct query.": {
			options: map[string]string{"bucket": "0.2"},
			expQuery: `
(
  sum(rate(apiserver_request_duration_seconds_count{ verb!="WATCH" }[{{.window}}]))
  -
  sum(rate(apiserver_request_duration_seconds_bucket{ le="0.2",verb!="WATCH" }[{{.window}}]))
)
/
sum(rate(apiserver_request_duration_seconds_count{ verb!="WATCH" }[{{.window}}]))
`,
		},

		"Having a filter shouldn't fail and return the correct query.": {
			options: map[string]string{
				"bucket": "0.2",
				"filter": `k1="v2",k2="v2"`,
			},
			expQuery: `
(
  sum(rate(apiserver_request_duration_seconds_count{ k1="v2",k2="v2",verb!="WATCH" }[{{.window}}]))
  -
  sum(rate(apiserver_request_duration_seconds_bucket{ k1="v2",k2="v2",le="0.2",verb!="WATCH" }[{{.window}}]))
)
/
sum(rate(apiserver_request_duration_seconds_count{ k1="v2",k2="v2",verb!="WATCH" }[{{.window}}]))
`,
		},

		"Having a filter with `{}` should return the query with the filters.": {
			options: map[string]string{
				"bucket": "0.2",
				"filter": `{k1="v2",k2="v2"}`,
			},
			expQuery: `
(
  sum(rate(apiserver_request_duration_seconds_count{ k1="v2",k2="v2",verb!="WATCH" }[{{.window}}]))
  -
  sum(rate(apiserver_request_duration_seconds_bucket{ k1="v2",k2="v2",le="0.2",verb!="WATCH" }[{{.window}}]))
)
/
sum(rate(apiserver_request_duration_seconds_count{ k1="v2",k2="v2",verb!="WATCH" }[{{.window}}]))
`,
		},

		"Having a filter with `,` should return the query with the filters.": {
			options: map[string]string{
				"bucket": "0.2",
				"filter": `k1="v2",k2="v2",`,
			},
			expQuery: `
(
  sum(rate(apiserver_request_duration_seconds_count{ k1="v2",k2="v2",verb!="WATCH" }[{{.window}}]))
  -
  sum(rate(apiserver_request_duration_seconds_bucket{ k1="v2",k2="v2",le="0.2",verb!="WATCH" }[{{.window}}]))
)
/
sum(rate(apiserver_request_duration_seconds_count{ k1="v2",k2="v2",verb!="WATCH" }[{{.window}}]))
`,
		},

		"Having a filter with `{}` and `,` should return the query with the filters.": {
			options: map[string]string{
				"filter": `{k1="v2",k2="v2",}`,
				"bucket": "0.2",
			},
			expQuery: `
(
  sum(rate(apiserver_request_duration_seconds_count{ k1="v2",k2="v2",verb!="WATCH" }[{{.window}}]))
  -
  sum(rate(apiserver_request_duration_seconds_bucket{ k1="v2",k2="v2",le="0.2",verb!="WATCH" }[{{.window}}]))
)
/
sum(rate(apiserver_request_duration_seconds_count{ k1="v2",k2="v2",verb!="WATCH" }[{{.window}}]))
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

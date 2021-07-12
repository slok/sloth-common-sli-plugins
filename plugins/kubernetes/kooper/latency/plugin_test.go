package latency_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/slok/sloth-common-sli-plugins/plugins/kubernetes/kooper/latency"
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
			options: map[string]string{"controller": "sloth"},
			expErr:  true,
		},

		"Not having controller option, it should fail.": {
			options: map[string]string{"bucket": "0.25"},
			expErr:  true,
		},

		"Having no filter should return the correct query.": {
			options: map[string]string{"bucket": "0.25", "controller": "sloth"},
			expQuery: `
1 - ((
  sum(rate(kooper_controller_processed_event_duration_seconds_bucket{ controller="sloth",le="0.25" }[{{ .window }}]))
  /
  (sum(rate(kooper_controller_processed_event_duration_seconds_count{ controller="sloth" }[{{ .window }}])) > 0)
) OR on() vector(1))
`,
		},

		"Having filter should return the correct query.": {
			options: map[string]string{"bucket": "0.25", "controller": "sloth", "filter": `k1="v1",k2="v2"`},
			expQuery: `
1 - ((
  sum(rate(kooper_controller_processed_event_duration_seconds_bucket{ k1="v1",k2="v2",controller="sloth",le="0.25" }[{{ .window }}]))
  /
  (sum(rate(kooper_controller_processed_event_duration_seconds_count{ k1="v1",k2="v2",controller="sloth" }[{{ .window }}])) > 0)
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

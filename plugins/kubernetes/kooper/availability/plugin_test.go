package availability_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/slok/sloth-common-sli-plugins/plugins/kubernetes/kooper/availability"
)

func TestSLIPlugin(t *testing.T) {
	tests := map[string]struct {
		meta     map[string]string
		labels   map[string]string
		options  map[string]string
		expQuery string
		expErr   bool
	}{
		"Not having controller option, it should fail.": {
			options: map[string]string{},
			expErr:  true,
		},

		"Having no filter option should return the correct query.": {
			options: map[string]string{"controller": "sloth"},
			expQuery: `
(
  sum(rate(kooper_controller_processed_event_duration_seconds_count{ controller="sloth",success="false" }[{{ .window }}]))
  /
  ((
    sum(rate(kooper_controller_processed_event_duration_seconds_count{ controller="sloth" }[{{ .window }}]))
    -
    (sum(rate(kooper_controller_queued_events_total{ controller="sloth",requeue="true" }[{{ .window }}])) OR on() vector(0))
  ) > 0)
) OR on() vector(0)
`,
		},

		"Having filter option should return the correct query.": {
			options: map[string]string{"controller": "sloth", "filter": `k1="v1",k2="v2"`},
			expQuery: `
(
  sum(rate(kooper_controller_processed_event_duration_seconds_count{ k1="v1",k2="v2",controller="sloth",success="false" }[{{ .window }}]))
  /
  ((
    sum(rate(kooper_controller_processed_event_duration_seconds_count{ k1="v1",k2="v2",controller="sloth" }[{{ .window }}]))
    -
    (sum(rate(kooper_controller_queued_events_total{ k1="v1",k2="v2",controller="sloth",requeue="true" }[{{ .window }}])) OR on() vector(0))
  ) > 0)
) OR on() vector(0)
`,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)

			gotQuery, err := availability.SLIPlugin(context.TODO(), test.meta, test.labels, test.options)

			if test.expErr {
				assert.Error(err)
			} else if assert.NoError(err) {
				assert.Equal(test.expQuery, gotQuery)
			}
		})
	}
}

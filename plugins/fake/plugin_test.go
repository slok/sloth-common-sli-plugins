package fake_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/slok/sloth-common-sli-plugins/plugins/fake"
)

func TestSLIPlugin(t *testing.T) {
	tests := map[string]struct {
		meta     map[string]string
		labels   map[string]string
		options  map[string]string
		expQuery string
		expErr   bool
	}{
		"Not having objective in metadata, it should fail.": {
			meta:    map[string]string{},
			options: map[string]string{"burn_rate": "1"},
			expErr:  true,
		},

		"Having a wrong objective in metadata, it should fail.": {
			meta:    map[string]string{"objective": "ninety nine"},
			options: map[string]string{"burn_rate": "1"},
			expErr:  true,
		},

		"Not having burn rate in options, it should fail.": {
			meta:    map[string]string{"objective": "99"},
			options: map[string]string{},
			expErr:  true,
		},

		"Having a wrong burn rate in options, it should fail.": {
			meta:    map[string]string{"objective": "99"},
			options: map[string]string{"burn_rate": "one"},
			expErr:  true,
		},

		"Having objective and burn rate, should return a correct query (speed 1).": {
			meta:     map[string]string{"objective": "99"},
			options:  map[string]string{"burn_rate": "1"},
			expQuery: `max_over_time(vector(0.010000)[{{.window}}:])`,
		},

		"Having objective and burn rate, should return a correct query (speed 2).": {
			meta:     map[string]string{"objective": "99"},
			options:  map[string]string{"burn_rate": "2"},
			expQuery: `max_over_time(vector(0.020000)[{{.window}}:])`,
		},

		"Having objective and burn rate, should return a correct query (speed 10).": {
			meta:     map[string]string{"objective": "99.9"},
			options:  map[string]string{"burn_rate": "10"},
			expQuery: `max_over_time(vector(0.010000)[{{.window}}:])`,
		},

		"Having objective and burn rate, should return a correct query (speed 0.5).": {
			meta:     map[string]string{"objective": "98"},
			options:  map[string]string{"burn_rate": "0.5"},
			expQuery: `max_over_time(vector(0.010000)[{{.window}}:])`,
		},

		"Having objective, burn rate and jitter, should return a correct query (speed 1, +-10%).": {
			meta:     map[string]string{"objective": "99"},
			options:  map[string]string{"burn_rate": "1", "jitter_percent": "10"},
			expQuery: `(max_over_time(vector(0.010000)[{{.window}}:])) - ((0.010000 * (10.000000 - ((time() * minute() * hour() * day_of_week() * month()) % 20.000000))) / 100)`,
		},

		"Having objective, burn rate and jitter, should return a correct query (speed ~2 +-50%).": {
			meta:     map[string]string{"objective": "99.9"},
			options:  map[string]string{"burn_rate": "2", "jitter_percent": "50"},
			expQuery: `(max_over_time(vector(0.002000)[{{.window}}:])) - ((0.002000 * (50.000000 - ((time() * minute() * hour() * day_of_week() * month()) % 100.000000))) / 100)`,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)

			gotQuery, err := fake.SLIPlugin(context.TODO(), test.meta, test.labels, test.options)

			if test.expErr {
				assert.Error(err)
			} else if assert.NoError(err) {
				assert.Equal(test.expQuery, gotQuery)
			}
		})
	}
}

package availability_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/slok/sloth-common-sli-plugins/plugins/traefik/v1/availability"
)

func TestSLIPlugin(t *testing.T) {
	tests := map[string]struct {
		meta     map[string]string
		labels   map[string]string
		options  map[string]string
		expQuery string
		expErr   bool
	}{
		"Without backend, should fail.": {
			options: map[string]string{},
			expErr:  true,
		},

		"An invalid backend query, should fail.": {
			options: map[string]string{"backend_regex": "([xyz"},
			expErr:  true,
		},

		"An empty backend query, should fail.": {
			options: map[string]string{"backend_regex": ""},
			expErr:  true,
		},

		"Not having a filter and with backend should return a valid query.": {
			options: map[string]string{"backend_regex": "github.com/slok/sloth/?"},
			expQuery: `
(
  sum(rate(traefik_backend_requests_total{ backend=~"github.com/slok/sloth/?",code=~"(5..|429)" }[{{.window}}]))
  /
  (sum(rate(traefik_backend_requests_total{ backend=~"github.com/slok/sloth/?" }[{{.window}}])) > 0)
) OR on() vector(0)
`,
		},

		"Having a filter and with backend should return a valid query.": {
			options: map[string]string{
				"filter":        `k1="v2",k2="v2"`,
				"backend_regex": "github.com/slok/sloth/?",
			},
			expQuery: `
(
  sum(rate(traefik_backend_requests_total{ k1="v2",k2="v2",backend=~"github.com/slok/sloth/?",code=~"(5..|429)" }[{{.window}}]))
  /
  (sum(rate(traefik_backend_requests_total{ k1="v2",k2="v2",backend=~"github.com/slok/sloth/?" }[{{.window}}])) > 0)
) OR on() vector(0)
`,
		},

		"Filter should be sanitized with ','.": {
			options: map[string]string{
				"filter":        `k1="v2",k2="v2",`,
				"backend_regex": "github.com/slok/sloth/?",
			},
			expQuery: `
(
  sum(rate(traefik_backend_requests_total{ k1="v2",k2="v2",backend=~"github.com/slok/sloth/?",code=~"(5..|429)" }[{{.window}}]))
  /
  (sum(rate(traefik_backend_requests_total{ k1="v2",k2="v2",backend=~"github.com/slok/sloth/?" }[{{.window}}])) > 0)
) OR on() vector(0)
`,
		},

		"Filter should be sanitized with '{'.": {
			options: map[string]string{
				"filter":        `{k1="v2",k2="v2",},`,
				"backend_regex": "github.com/slok/sloth/?",
			},
			expQuery: `
(
  sum(rate(traefik_backend_requests_total{ k1="v2",k2="v2",backend=~"github.com/slok/sloth/?",code=~"(5..|429)" }[{{.window}}]))
  /
  (sum(rate(traefik_backend_requests_total{ k1="v2",k2="v2",backend=~"github.com/slok/sloth/?" }[{{.window}}])) > 0)
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

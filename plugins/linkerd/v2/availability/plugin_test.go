package availability_test

import (
	"context"
	"testing"

	availability "github.com/slok/sloth-common-sli-plugins/plugins/linkerd/v2/availability"
	"github.com/stretchr/testify/assert"
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
			options: map[string]string{"filter": "something="},
			expErr:  true,
		},

		"Not having a filter should fail.": {
			expErr: true,
		},

		"Having a filter should return the query with the filters.": {
			options: map[string]string{"filter": `k1="v2",k2="v2"`},
			expQuery: `
(
  sum(increase(response_total{ k1="v2",k2="v2",classification="failure", direction="inbound" }[{{.window}}]))
  /
  sum(increase(response_total{ k1="v2",k2="v2",direction="inbound" }[{{.window}}]))
) OR on() vector(0)
`,
		},

		"Having a filter with `{}` should return the query with the filters.": {
			options: map[string]string{"filter": `{k1="v2",k2="v2"}`},
			expQuery: `
(
  sum(increase(response_total{ k1="v2",k2="v2",classification="failure", direction="inbound" }[{{.window}}]))
  /
  sum(increase(response_total{ k1="v2",k2="v2",direction="inbound" }[{{.window}}]))
) OR on() vector(0)
`,
		},

		"Having a filter with `,` should return the query with the filters.": {
			options: map[string]string{"filter": `k1="v2",k2="v2",`},
			expQuery: `
(
  sum(increase(response_total{ k1="v2",k2="v2",classification="failure", direction="inbound" }[{{.window}}]))
  /
  sum(increase(response_total{ k1="v2",k2="v2",direction="inbound" }[{{.window}}]))
) OR on() vector(0)
`,
		},

		"Having a filter with `{}` and `,` should return the query with the filters.": {
			options: map[string]string{"filter": `{k1="v2",k2="v2",}`},
			expQuery: `
(
  sum(increase(response_total{ k1="v2",k2="v2",classification="failure", direction="inbound" }[{{.window}}]))
  /
  sum(increase(response_total{ k1="v2",k2="v2",direction="inbound" }[{{.window}}]))
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

package availability_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/slok/sloth-common-sli-plugins/plugins/prometheus/targets/availability"
)

func TestSLIPlugin(t *testing.T) {
	tests := map[string]struct {
		meta     map[string]string
		labels   map[string]string
		options  map[string]string
		expQuery string
		expErr   bool
	}{
		"Having no filter option should return the correct query.": {
			options: map[string]string{},
			expQuery: `
sum(count_over_time((up{  } == 0)[{{ .window }}:])) or vector(0)
/
sum(count_over_time((up{  })[{{ .window }}:]))
`,
		},

		"Having filter option should return the correct query.": {
			options: map[string]string{"filter": `k1="v1",k2="v2"`},
			expQuery: `
sum(count_over_time((up{ k1="v1",k2="v2" } == 0)[{{ .window }}:])) or vector(0)
/
sum(count_over_time((up{ k1="v1",k2="v2" })[{{ .window }}:]))
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

package availability_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/slok/sloth-common-sli-plugins/plugins/kubernetes/nginx-ingress/availability"
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

		"Not having a filter shouldn't fail and return the correct query.": {
			expQuery: "\nsum(rate(nginx_ingress_controller_requests{ status=~\"(5..|429)\" }[{{.window}}]))\n/\nsum(rate(nginx_ingress_controller_requests{  }[{{.window}}]))\nor on() vector(0)\n",
		},

		"Having a filter should return the query with the filters.": {
			options:  map[string]string{"filter": `k1="v2",k2="v2"`},
			expQuery: "\nsum(rate(nginx_ingress_controller_requests{ k1=\"v2\",k2=\"v2\",status=~\"(5..|429)\" }[{{.window}}]))\n/\nsum(rate(nginx_ingress_controller_requests{ k1=\"v2\",k2=\"v2\" }[{{.window}}]))\nor on() vector(0)\n",
		},

		"Having a filter with `{}` should return the query with the filters.": {
			options:  map[string]string{"filter": `{k1="v2",k2="v2"}`},
			expQuery: "\nsum(rate(nginx_ingress_controller_requests{ k1=\"v2\",k2=\"v2\",status=~\"(5..|429)\" }[{{.window}}]))\n/\nsum(rate(nginx_ingress_controller_requests{ k1=\"v2\",k2=\"v2\" }[{{.window}}]))\nor on() vector(0)\n",
		},

		"Having a filter with `,` should return the query with the filters.": {
			options:  map[string]string{"filter": `k1="v2",k2="v2",`},
			expQuery: "\nsum(rate(nginx_ingress_controller_requests{ k1=\"v2\",k2=\"v2\",status=~\"(5..|429)\" }[{{.window}}]))\n/\nsum(rate(nginx_ingress_controller_requests{ k1=\"v2\",k2=\"v2\" }[{{.window}}]))\nor on() vector(0)\n",
		},

		"Having a filter with `{}` and `,` should return the query with the filters.": {
			options:  map[string]string{"filter": `{k1="v2",k2="v2",}`},
			expQuery: "\nsum(rate(nginx_ingress_controller_requests{ k1=\"v2\",k2=\"v2\",status=~\"(5..|429)\" }[{{.window}}]))\n/\nsum(rate(nginx_ingress_controller_requests{ k1=\"v2\",k2=\"v2\" }[{{.window}}]))\nor on() vector(0)\n",
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

package availability_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/slok/sloth-common-sli-plugins/plugins/istio/v1/availability"
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
			options: map[string]string{"service": "test"},
			expErr:  true,
		},

		"An empty namespace query, should fail.": {
			options: map[string]string{"namespace": "", "service": "test"},
			expErr:  true,
		},

		"Without service, should fail.": {
			options: map[string]string{"namespace": "default"},
			expErr:  true,
		},

		"An empty service query, should fail.": {
			options: map[string]string{"namespace": "default", "service": ""},
			expErr:  true,
		},

		"Not having a filter and with namespace + service should return a valid query.": {
			options: map[string]string{"namespace": "default", "service": "test"},
			expQuery: `
(
  sum(rate(istio_requests_total{ destination_service_name="test",destination_service_namespace="default",response_code=~"(5..|429)" }[{{.window}}])) 
  /          
  (sum(rate(istio_requests_total{ destination_service_name="test",destination_service_namespace="default" }[{{.window}}])) > 0)
) OR on() vector(0)
`,
		},

		"Having a filter, with namespace and service should return a valid query.": {
			options: map[string]string{
				"filter":    `k1="v2",k2="v2"`,
				"namespace": "default",
				"service":   "test",
			},
			expQuery: `
(
  sum(rate(istio_requests_total{ k1="v2",k2="v2",destination_service_name="test",destination_service_namespace="default",response_code=~"(5..|429)" }[{{.window}}])) 
  /          
  (sum(rate(istio_requests_total{ k1="v2",k2="v2",destination_service_name="test",destination_service_namespace="default" }[{{.window}}])) > 0)
) OR on() vector(0)
`,
		},

		"Filter should be sanitized with ','.": {
			options: map[string]string{
				"filter":    `k1="v2",k2="v2",`,
				"namespace": "default",
				"service":   "test",
			},
			expQuery: `
(
  sum(rate(istio_requests_total{ k1="v2",k2="v2",destination_service_name="test",destination_service_namespace="default",response_code=~"(5..|429)" }[{{.window}}])) 
  /          
  (sum(rate(istio_requests_total{ k1="v2",k2="v2",destination_service_name="test",destination_service_namespace="default" }[{{.window}}])) > 0)
) OR on() vector(0)
`,
		},

		"Filter should be sanitized with '{'.": {
			options: map[string]string{
				"filter":    `{k1="v2",k2="v2",},`,
				"namespace": "default",
				"service":   "test",
			},
			expQuery: `
(
  sum(rate(istio_requests_total{ k1="v2",k2="v2",destination_service_name="test",destination_service_namespace="default",response_code=~"(5..|429)" }[{{.window}}])) 
  /          
  (sum(rate(istio_requests_total{ k1="v2",k2="v2",destination_service_name="test",destination_service_namespace="default" }[{{.window}}])) > 0)
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

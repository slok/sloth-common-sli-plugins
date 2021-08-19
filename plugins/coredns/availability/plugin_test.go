package availability_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/slok/sloth-common-sli-plugins/plugins/coredns/availability"
)

func TestSLIPlugin(t *testing.T) {
	tests := map[string]struct {
		meta     map[string]string
		labels   map[string]string
		options  map[string]string
		expQuery string
		expErr   bool
	}{
		"Not having a filter and without custom rcode should return a valid query.": {
			options:  map[string]string{},
			expQuery: "\nsum(rate(coredns_dns_responses_total{ rcode=~\"SERVFAIL\" }[{{.window}}]))\n/\n(sum(rate(coredns_dns_responses_total{  }[{{.window}}])) > 0)\n",
		},

		"Not having a filter and with custom rcode should return a valid query.": {
			options:  map[string]string{"custom_rcode_regex": "(SERVFAIL|FORMERR)"},
			expQuery: "\nsum(rate(coredns_dns_responses_total{ rcode=~\"(SERVFAIL|FORMERR)\" }[{{.window}}]))\n/\n(sum(rate(coredns_dns_responses_total{  }[{{.window}}])) > 0)\n",
		},

		"custom rcode regex should fail.": {
			options: map[string]string{"custom_rcode_regex": "(SERVFAIL|"},
			expErr:  true,
		},

		"Having a filter and without custom rcode should return a valid query.": {
			options:  map[string]string{"filter": `k1="v2",k2="v2"`},
			expQuery: "\nsum(rate(coredns_dns_responses_total{ k1=\"v2\",k2=\"v2\",rcode=~\"SERVFAIL\" }[{{.window}}]))\n/\n(sum(rate(coredns_dns_responses_total{ k1=\"v2\",k2=\"v2\" }[{{.window}}])) > 0)\n",
		},

		"Filter should be sanitized with ','.": {
			options:  map[string]string{"filter": `k1="v2",k2="v2",`},
			expQuery: "\nsum(rate(coredns_dns_responses_total{ k1=\"v2\",k2=\"v2\",rcode=~\"SERVFAIL\" }[{{.window}}]))\n/\n(sum(rate(coredns_dns_responses_total{ k1=\"v2\",k2=\"v2\" }[{{.window}}])) > 0)\n",
		},

		"Filter should be sanitized with '{'.": {
			options:  map[string]string{"filter": `{k1="v2",k2="v2",},`},
			expQuery: "\nsum(rate(coredns_dns_responses_total{ k1=\"v2\",k2=\"v2\",rcode=~\"SERVFAIL\" }[{{.window}}]))\n/\n(sum(rate(coredns_dns_responses_total{ k1=\"v2\",k2=\"v2\" }[{{.window}}])) > 0)\n",
		},

		"Having a filter and with custom rcode should return a valid query.": {
			options: map[string]string{
				"filter":             `k1="v2",k2="v2"`,
				"custom_rcode_regex": "(SERVFAIL|FORMERR)",
			},
			expQuery: "\nsum(rate(coredns_dns_responses_total{ k1=\"v2\",k2=\"v2\",rcode=~\"(SERVFAIL|FORMERR)\" }[{{.window}}]))\n/\n(sum(rate(coredns_dns_responses_total{ k1=\"v2\",k2=\"v2\" }[{{.window}}])) > 0)\n",
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

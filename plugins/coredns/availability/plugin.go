package availability

import (
	"bytes"
	"context"
	"fmt"
	"regexp"
	"strings"
	"text/template"
)

const (
	// SLIPluginVersion is the version of the plugin spec.
	SLIPluginVersion = "prometheus/v1"
	// SLIPluginID is the registering ID of the plugin.
	SLIPluginID = "sloth-common/coredns/availability"
)

var queryTpl = template.Must(template.New("").Option("missingkey=error").Parse(`
sum(rate(coredns_dns_responses_total{ {{.filterError}}rcode=~"{{.rcodeRegex}}" }[{{"{{.window}}"}}]))
/
(sum(rate(coredns_dns_responses_total{ {{.filterTotal}} }[{{"{{.window}}"}}])) > 0)
`))

// SLIPlugin will return a query that will return the availability error based on coreDNS rcodes.
func SLIPlugin(ctx context.Context, meta, labels, options map[string]string) (string, error) {
	rcodeRegex, err := getRCodeRegex(options)
	if err != nil {
		return "", fmt.Errorf("could not get custom rcodes regex: %w", err)
	}

	filter := getFilter(options)
	filterTotal := filter
	filterError := filter
	if filterError != "" {
		filterError = filter + ","
	}

	// Create query.
	var b bytes.Buffer
	data := map[string]string{
		"filterError": filterError,
		"filterTotal": filterTotal,
		"rcodeRegex":  rcodeRegex,
	}
	err = queryTpl.Execute(&b, data)
	if err != nil {
		return "", fmt.Errorf("could not render query template: %w", err)
	}

	return b.String(), nil
}

func getFilter(options map[string]string) string {
	filter := options["filter"]
	filter = strings.Trim(filter, "{},")

	return filter
}

// For the full rcodes list, Check these:
// - https://github.com/miekg/dns/blob/ab67aa64230094bdd0167ee5360e00e0a250a3ac/types.go#L124-#L144
// - https://github.com/miekg/dns/blob/ab67aa64230094bdd0167ee5360e00e0a250a3ac/msg.go#L137-L159
func getRCodeRegex(options map[string]string) (string, error) {
	customRCode := options["custom_rcode_regex"]
	customRCode = strings.TrimSpace(customRCode)

	// Set a safe default for errors.
	if customRCode == "" {
		return "SERVFAIL", nil
	}

	_, err := regexp.Compile(customRCode)
	if err != nil {
		return "", fmt.Errorf("invalid regex: %w", err)
	}

	return customRCode, nil
}

package fake

import (
	"context"
	"fmt"
	"strconv"
)

const (
	// SLIPluginVersion is the version of the plugin spec.
	SLIPluginVersion = "prometheus/v1"
	// SLIPluginID is the registering ID of the plugin.
	SLIPluginID = "sloth-common/fake"
)

// SLIPlugin will return a query that will fake a burning error budget at the desired speed using
// `burn_rate` option.
// The plugins also accepts a `jitter_percent` that will add/remove a jitter in the range of that jitter percent.
func SLIPlugin(ctx context.Context, meta, labels, options map[string]string) (string, error) {
	objective, err := getSLOObjective(meta)
	if err != nil {
		return "", fmt.Errorf("could not get objective: %w", err)
	}

	burnRate, err := getBurnRate(options)
	if err != nil {
		return "", fmt.Errorf("could not get burn rate: %w", err)
	}

	jitter, err := getJitterPercent(options)
	if err != nil {
		return "", fmt.Errorf("could not get jitter percent: %w", err)
	}

	// Get the error budget.
	errorBudgetPerc := 100 - objective
	errorBudgetRatio := errorBudgetPerc / 100

	// Apply factor (burn rate) to error budget.
	expectedSLIError := errorBudgetRatio * burnRate

	// Create regular query and if no jitter required then regular query
	query := fmt.Sprintf(`max_over_time(vector(%f)[{{.window}}:])`, expectedSLIError)
	if jitter == 0 {
		return query, nil
	}

	// Create jitter (this is the best random I could came up with) and rest to the regular query.
	jitterQuery := fmt.Sprintf(`(%f * (%f - ((time() * minute() * hour() * day_of_week() * month()) %% %f))) / 100`, expectedSLIError, jitter, jitter*2)

	return fmt.Sprintf("(%s) - (%s)", query, jitterQuery), nil
}

func getBurnRate(options map[string]string) (float64, error) {
	burnRateS, ok := options["burn_rate"]
	if !ok {
		return 0, fmt.Errorf("'burn_rate' option is required")
	}

	burnRate, err := strconv.ParseFloat(burnRateS, 64)
	if err != nil {
		return 0, fmt.Errorf("not a valid burn_rate, can't parse to float64: %w", err)
	}

	return burnRate, nil
}

func getJitterPercent(options map[string]string) (float64, error) {
	jitterPercentS, ok := options["jitter_percent"]
	if !ok {
		return 0, nil
	}

	jitterPercent, err := strconv.ParseFloat(jitterPercentS, 64)
	if err != nil {
		return 0, fmt.Errorf("not a valid jitter_percent, can't parse to float64: %w", err)
	}

	return jitterPercent, nil
}

func getSLOObjective(meta map[string]string) (float64, error) {
	objectiveS, ok := meta["objective"]
	if !ok {
		return 0, fmt.Errorf("'objective' metadata missing")
	}

	objective, err := strconv.ParseFloat(objectiveS, 64)
	if err != nil {
		return 0, fmt.Errorf("not a valid objective, can't parse to float64: %w", err)
	}

	return objective, nil
}

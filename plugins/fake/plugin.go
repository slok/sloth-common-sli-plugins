package fake

import (
	"context"
	"fmt"
	"strconv"
)

const (
	SLIPluginVersion = "prometheus/v1"
	SLIPluginID      = "sloth-common/fake"
)

// SLIPlugin will return a query that will fake a burning error budget at the desired speed using
// `burn_rate` option.
func SLIPlugin(ctx context.Context, meta, labels, options map[string]string) (string, error) {
	objective, err := getSLOObjective(meta)
	if err != nil {
		return "", fmt.Errorf("could not get objective: %w", err)
	}

	burnRate, err := getBurnRate(options)
	if err != nil {
		return "", fmt.Errorf("could not get burn rate: %w", err)
	}

	// Get the error budget.
	errorBudgetPerc := 100 - objective
	errorBudgetRatio := errorBudgetPerc / 100

	// Apply factor (burn rate) to error budget.
	expectedSLIError := errorBudgetRatio * burnRate

	query := fmt.Sprintf(`max_over_time(vector(%f)[{{.window}}:])`, expectedSLIError)

	return query, nil
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

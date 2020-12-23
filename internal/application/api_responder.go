package application

import (
	"encoding/json"

	"go.dev.pztrn.name/metricator/internal/models"
)

// Responds with needed data. First parameter is a type of data needed (like metric name),
// second parameter is actual metric name. Second parameter also can be empty.
func (a *Application) respond(rInfo *models.RequestInfo) string {
	// If metric was requested - return only it.
	if rInfo.Metric != "" {
		metric, err := a.storage.Get(rInfo.Metric)
		if err != nil {
			return ""
		}

		return metric.GetValue()
	}

	// Otherwise we should get all metrics as slice and return them as string.
	// This is needed for metrics autodiscovery.
	metrics := a.storage.GetAsSlice()

	metricsBytes, err := json.Marshal(metrics)
	if err != nil {
		// ToDo: log error
		return ""
	}

	return string(metricsBytes)
}

package application

import (
	"go.dev.pztrn.name/metricator/internal/models"
)

// Responds with needed data. First parameter is a type of data needed (like metric name),
// second parameter is actual metric name. Second parameter also can be empty.
func (a *Application) respond(rInfo *models.RequestInfo) string {
	metric, err := a.storage.Get(rInfo.Metric)
	if err != nil {
		return ""
	}

	return metric.GetValue()
}

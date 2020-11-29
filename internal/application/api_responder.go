package application

import (
	"context"

	"go.dev.pztrn.name/metricator/internal/common"
)

func (a *Application) respond(ctx context.Context) string {
	metricName := ctx.Value(common.ContextKeyMetric).(string)

	return a.storage.Get(metricName)
}

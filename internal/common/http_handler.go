package common

import "go.dev.pztrn.name/metricator/internal/models"

// HTTPHandlerFunc describes signature of HTTP requests handling function.
type HTTPHandlerFunc func(*models.RequestInfo) string

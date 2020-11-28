package common

// ContextKey is a type of context.Context keys.
type ContextKey string

const (
	// ContextKeyApplication specifies that returned value is a name of application.
	ContextKeyApplication ContextKey = "applicationName"
	// ContextKeyMetric specifies that returned value is a name of metric of application.
	ContextKeyMetric ContextKey = "metricName"
)

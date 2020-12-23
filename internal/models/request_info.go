package models

// RequestInfo is a parsed request information to throw into application's handler.
type RequestInfo struct {
	// Application is a name of application. We should ask it's handler for metrics.
	Application string
	// Metric is a metric name with parameters (e.g. requests{path='/',code=200} will
	// be "requests/path:\//code:200").
	Metric string
	// Request type is a type of request. Currently known: "apps_list", "info", and "metrics".
	// All other request types will produce HTTP 400 error.
	RequestType string
	// APIVersion is a version of API requested.
	APIVersion int
}

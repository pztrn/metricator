package models

// RequestInfo is a parsed request information to throw into application's handler.
type RequestInfo struct {
	ApiVersion  int
	Application string
	Metric      string
	RequestType string
}

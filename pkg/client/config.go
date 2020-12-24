package client

// Config is a Metricator client configuration.
type Config struct {
	// Host is a host where Metricator is available for requests.
	Host string
	// Timeout specifies HTTP client timeout.
	Timeout int
}

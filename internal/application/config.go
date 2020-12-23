package application

import "time"

// Config is a generic application configuration.
type Config struct {
	// Headers is a list of headers that should be added to metrics request.
	Headers map[string]string `yaml:"headers"`
	// Endpoint is a remote application endpoint which should give us metrics
	// in Prometheus format.
	Endpoint string `yaml:"endpoint"`
	// TimeBetweenRequests is a minimal amount of time which should pass
	// between requests.
	TimeBetweenRequests time.Duration `yaml:"time_between_requests"`
}

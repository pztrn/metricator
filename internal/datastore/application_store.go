package datastore

import "sync"

// This is application-specific data storage.
type applicationStorage struct {
	metrics      map[string]string
	metricsMutex sync.RWMutex
}

// Creates new application-specific storage.
func newApplicationStorage() *applicationStorage {
	as := &applicationStorage{}
	as.initialize()

	return as
}

// Initializes internal things.
func (as *applicationStorage) initialize() {
	as.metrics = make(map[string]string)
}

// Starts application-specific things, like goroutine for HTTP requests.
func (as *applicationStorage) start() {}

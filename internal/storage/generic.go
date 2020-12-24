package storage

import "go.dev.pztrn.name/metricator/pkg/schema"

// GenericStorage describes interface every other storage should embed
// and conform to as it contains essential things like context handling.
type GenericStorage interface {
	// Get returns data from storage by key.
	Get(string) (schema.Metric, error)
	// GetAsSlice returns all data from storage as slice.
	GetAsSlice() []schema.Metric
	// GetDoneChan returns a channel which should be used to block execution
	// until storage's routines are completed.
	GetDoneChan() chan struct{}
	// Put puts passed data into storage.
	Put(map[string]schema.Metric)
	// Start starts asynchronous things if needed.
	Start()
}

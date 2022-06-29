package memory

import (
	"context"
	"errors"
	"sync"

	"go.dev.pztrn.name/metricator/internal/logger"
	"go.dev.pztrn.name/metricator/pkg/schema"
)

// ErrMetricNotFound appears if requested metric wasn't found in storage.
var ErrMetricNotFound = errors.New("metric not found")

// Storage is an in-memory storage.
type Storage struct {
	dataMutex sync.RWMutex
	ctx       context.Context
	doneChan  chan struct{}
	logger    *logger.Logger
	data      map[string]schema.Metric
	name      string
}

// NewStorage creates new in-memory storage to use.
func NewStorage(ctx context.Context, name string, logger *logger.Logger) (*Storage, chan struct{}) {
	// nolint:exhaustruct
	storage := &Storage{
		ctx:      ctx,
		doneChan: make(chan struct{}),
		logger:   logger,
		name:     name,
		data:     make(map[string]schema.Metric),
	}
	storage.initialize()

	return storage, storage.doneChan
}

// Get returns data from storage by key.
func (s *Storage) Get(key string) (schema.Metric, error) {
	s.logger.Debugln("Retrieving data for", key, "key from storage...")

	s.dataMutex.RLock()
	defer s.dataMutex.RUnlock()

	data, found := s.data[key]
	if !found {
		s.logger.Infoln("Key", key, "not found in storage!")

		return schema.NewMetric("", "", "", nil), ErrMetricNotFound
	}

	s.logger.Debugf("Key %s found: %+v\n", key, data)

	return data, nil
}

// GetAsSlice returns all data from storage as slice.
func (s *Storage) GetAsSlice() []schema.Metric {
	s.logger.Debugln("Returning all stored metrics as slice...")

	metrics := make([]schema.Metric, 0, len(s.data))

	for _, metric := range s.data {
		metrics = append(metrics, metric)
	}

	return metrics
}

// GetDoneChan returns a channel which should be used to block execution
// until storage's routines are completed.
func (s *Storage) GetDoneChan() chan struct{} {
	return s.doneChan
}

// Initializes internal things.
func (s *Storage) initialize() {
	s.data = make(map[string]schema.Metric)
}

// Put puts passed data into storage.
func (s *Storage) Put(data map[string]schema.Metric) {
	s.dataMutex.Lock()
	defer s.dataMutex.Unlock()

	for k, v := range data {
		// We should not put valueless metrics.
		if v.Value != "" {
			s.data[k] = v
		}
	}

	s.logger.Debugln("Put", len(data), "items in", s.name)
}

// Start starts asynchronous things if needed.
func (s *Storage) Start() {
	// The Context Listening Goroutine.
	go func() {
		<-s.ctx.Done()

		s.logger.Infoln("In-memory storage", s.name, "done")

		s.doneChan <- struct{}{}
	}()
}

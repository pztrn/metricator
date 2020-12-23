package memory

import (
	"context"
	"errors"
	"log"
	"sync"

	"go.dev.pztrn.name/metricator/internal/models"
)

var ErrMetricNotFound = errors.New("metric not found")

// Storage is an in-memory storage.
type Storage struct {
	ctx      context.Context
	doneChan chan struct{}
	name     string

	data      map[string]models.Metric
	dataMutex sync.RWMutex
}

// NewStorage creates new in-memory storage to use.
func NewStorage(ctx context.Context, name string) (*Storage, chan struct{}) {
	s := &Storage{
		ctx:      ctx,
		doneChan: make(chan struct{}),
		name:     name,
	}
	s.initialize()

	return s, s.doneChan
}

// Get returns data from storage by key.
func (s *Storage) Get(key string) (models.Metric, error) {
	s.dataMutex.RLock()
	defer s.dataMutex.RUnlock()

	data, found := s.data[key]
	if !found {
		return models.NewMetric("", "", "", nil), ErrMetricNotFound
	}

	return data, nil
}

// GetAsSlice returns all data from storage as slice.
func (s *Storage) GetAsSlice() []models.Metric {
	metrics := make([]models.Metric, 0, len(s.data))

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
	s.data = make(map[string]models.Metric)
}

// Put puts passed data into storage.
func (s *Storage) Put(data map[string]models.Metric) {
	s.dataMutex.Lock()
	defer s.dataMutex.Unlock()

	for k, v := range data {
		s.data[k] = v
	}

	log.Println("Put", len(data), "items in", s.name)
}

// Start starts asynchronous things if needed.
func (s *Storage) Start() {
	// The Context Listening Goroutine.
	go func() {
		<-s.ctx.Done()

		log.Println("In-memory storage", s.name, "done")

		s.doneChan <- struct{}{}
	}()
}

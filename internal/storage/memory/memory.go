package memory

import (
	"context"
	"log"
	"sync"
)

// Storage is an in-memory storage.
type Storage struct {
	ctx      context.Context
	doneChan chan struct{}
	name     string

	data      map[string]string
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
func (s *Storage) Get(key string) string {
	s.dataMutex.RLock()
	defer s.dataMutex.RUnlock()

	data, found := s.data[key]
	if !found {
		return "Not found"
	}

	return data
}

// GetDoneChan returns a channel which should be used to block execution
// until storage's routines are completed.
func (s *Storage) GetDoneChan() chan struct{} {
	return s.doneChan
}

// Initializes internal things.
func (s *Storage) initialize() {
	s.data = make(map[string]string)
}

// Put puts passed data into storage.
func (s *Storage) Put(data map[string]string) {
	s.dataMutex.Lock()
	defer s.dataMutex.Unlock()

	for k, v := range data {
		s.data[k] = v
	}
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

package datastore

import (
	"context"
	"log"
	"sync"

	"go.dev.pztrn.name/metricator/internal/configuration"
)

// DataStore is a data storage structure. It keeps all gathered metrics and gives
//  them away on request.
type DataStore struct {
	config   *configuration.Config
	ctx      context.Context
	doneChan chan struct{}

	datas      map[string]*applicationStorage
	datasMutex sync.RWMutex
}

// NewDataStore creates new data storage.
func NewDataStore(ctx context.Context, cfg *configuration.Config) (*DataStore, chan struct{}) {
	ds := &DataStore{
		config:   cfg,
		ctx:      ctx,
		doneChan: make(chan struct{}),
	}
	ds.initialize()

	return ds, ds.doneChan
}

// Internal things initialization.
func (ds *DataStore) initialize() {
	ds.datas = make(map[string]*applicationStorage)

	// Create applications defined in configuration.

	go func() {
		<-ds.ctx.Done()
		log.Println("Data storage stopped")

		ds.doneChan <- struct{}{}
	}()
}

// Start starts data storage asynchronous things.
func (ds *DataStore) Start() {
	log.Println("Starting data storage...")

	ds.datasMutex.RLock()
	for _, storage := range ds.datas {
		storage.start()
	}
}

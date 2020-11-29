package application

import (
	"context"
	"log"
	"net/http"
	"sync"

	"go.dev.pztrn.name/metricator/internal/storage"
	"go.dev.pztrn.name/metricator/internal/storage/memory"
)

// Application is a thing that responsible for all application-related
// actions like data fetching, storing, etc. on higher level.
type Application struct {
	config   *Config
	ctx      context.Context
	doneChan chan struct{}
	name     string

	storage     storage.Metrics
	storageDone chan struct{}

	fetchIsRunning      bool
	fetchIsRunningMutex sync.RWMutex

	httpClient *http.Client
}

// NewApplication creates new application.
func NewApplication(ctx context.Context, name string, config *Config) *Application {
	a := &Application{
		config:   config,
		ctx:      ctx,
		doneChan: make(chan struct{}),
		name:     name,
	}
	a.initialize()

	return a
}

// GetDoneChan returns a channel which should be used to block execution until
// application's routines are completed.
func (a *Application) GetDoneChan() chan struct{} {
	return a.doneChan
}

// Initializes internal things like storage, HTTP client, etc.
func (a *Application) initialize() {
	a.storage, a.storageDone = memory.NewStorage(a.ctx, a.name+" storage")

	log.Printf("Application '%s' initialized with configuration: %+v\n", a.name, a.config)
}

// Start starts asynchronous things like data fetching, storage cleanup, etc.
func (a *Application) Start() {
	a.storage.Start()

	go a.startFetcher()

	// The Context Listening Goroutine.
	go func() {
		<-a.ctx.Done()
		// We should wait until storage routines are also stopped.
		<-a.storage.GetDoneChan()

		log.Println("Application", a.name, "stopped")

		a.doneChan <- struct{}{}
	}()
}

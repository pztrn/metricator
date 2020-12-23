package application

import (
	"context"
	"net/http"
	"sync"

	"go.dev.pztrn.name/metricator/internal/common"
	"go.dev.pztrn.name/metricator/internal/logger"
	"go.dev.pztrn.name/metricator/internal/storage"
	"go.dev.pztrn.name/metricator/internal/storage/memory"
)

// Application is a thing that responsible for all application-related
// actions like data fetching, storing, etc. on higher level.
type Application struct {
	ctx                 context.Context
	storage             storage.Metrics
	config              *Config
	doneChan            chan struct{}
	logger              *logger.Logger
	storageDone         chan struct{}
	httpClient          *http.Client
	name                string
	fetchIsRunningMutex sync.RWMutex
	fetchIsRunning      bool
}

// NewApplication creates new application.
func NewApplication(ctx context.Context, name string, config *Config, logger *logger.Logger) *Application {
	a := &Application{
		config:   config,
		ctx:      ctx,
		doneChan: make(chan struct{}),
		logger:   logger,
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

// GetHandler returns HTTP requests handling function.
func (a *Application) GetHandler() common.HTTPHandlerFunc {
	return a.respond
}

// Initializes internal things like storage, HTTP client, etc.
func (a *Application) initialize() {
	a.storage, a.storageDone = memory.NewStorage(a.ctx, a.name+" storage", a.logger)

	a.logger.Debugf("Application '%s' initialized with configuration: %+v\n", a.name, a.config)
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

		a.logger.Infoln("Application", a.name, "stopped")

		a.doneChan <- struct{}{}
	}()
}

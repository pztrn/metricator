package httpserver

import (
	"context"
	"net"
	"net/http"
	"strings"
	"time"

	"go.dev.pztrn.name/metricator/internal/common"
	"go.dev.pztrn.name/metricator/internal/configuration"
	"go.dev.pztrn.name/metricator/internal/logger"
)

// HTTPServer is a controlling structure for HTTP server.
type HTTPServer struct {
	config   *configuration.Config
	ctx      context.Context
	doneChan chan struct{}
	logger   *logger.Logger
	handler  *handler
	server   *http.Server
}

// NewHTTPServer creates HTTP server and executes preliminary initialization
// (HTTP server structure initialized but it doesn't start).
func NewHTTPServer(ctx context.Context, cfg *configuration.Config, logger *logger.Logger) (*HTTPServer, chan struct{}) {
	// nolint:exhaustivestruct
	httpServer := &HTTPServer{
		config:   cfg,
		ctx:      ctx,
		doneChan: make(chan struct{}),
		logger:   logger,
	}
	httpServer.initialize()

	return httpServer, httpServer.doneChan
}

// Returns request's context based on main context of application.
// Basically it returns main context and does nothing more.
func (h *HTTPServer) getRequestContext(_ net.Listener) context.Context {
	return h.ctx
}

// Initializes handler and HTTP server structure.
func (h *HTTPServer) initialize() {
	h.logger.Debugln("Initializing HTTP server...")

	h.handler = &handler{
		handlers: make(map[string]common.HTTPHandlerFunc),
	}
	// We do not need to specify all possible parameters for HTTP server, so:
	// nolint:exhaustivestruct
	h.server = &http.Server{
		// ToDo: make it all configurable.
		Addr:           ":34421",
		BaseContext:    h.getRequestContext,
		Handler:        h.handler,
		ReadTimeout:    time.Second * 10,
		WriteTimeout:   time.Second * 10,
		MaxHeaderBytes: 1 << 20,
	}
}

// RegisterHandlerForApplication registers HTTP handler for application.
func (h *HTTPServer) RegisterHandlerForApplication(name string, handler common.HTTPHandlerFunc) {
	h.logger.Debugln("Registering handler for application", name)
	h.handler.register(name, handler)
}

// Start starts HTTP server in another goroutine and one more goroutine which
// is listening to main context's Cancel() call to stop HTTP server.
func (h *HTTPServer) Start() {
	go func() {
		err := h.server.ListenAndServe()
		if err != nil {
			if !strings.Contains(err.Error(), "Server closed") {
				h.logger.Infoln("HTTP server failed to listen:", err.Error())
			}
		}
	}()

	go func() {
		<-h.ctx.Done()
		h.logger.Infoln("Shutting down HTTP server")

		err := h.server.Shutdown(h.ctx)
		if err != nil && !strings.Contains(err.Error(), "context canceled") {
			h.logger.Infoln("Failed to stop HTTP server:", err.Error())
		}

		h.logger.Infoln("HTTP server stopped")

		h.doneChan <- struct{}{}
	}()
}

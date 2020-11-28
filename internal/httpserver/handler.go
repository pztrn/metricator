package httpserver

import (
	"net/http"

	"go.dev.pztrn.name/metricator/internal/common"
)

// HTTP requests handler.
type handler struct {
	handler common.HTTPHandlerFunc
}

// Registers request's handler.
func (h *handler) register(hndl common.HTTPHandlerFunc) {
	h.handler = hndl
}

// ServeHTTP handles every HTTP request.
func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {}

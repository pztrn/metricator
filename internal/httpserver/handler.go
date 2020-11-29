package httpserver

import (
	"context"
	"net/http"
	"strings"

	"go.dev.pztrn.name/metricator/internal/common"
)

// HTTP requests handler.
type handler struct {
	handlers map[string]common.HTTPHandlerFunc
}

// Registers request's handler.
func (h *handler) register(name string, hndl common.HTTPHandlerFunc) {
	h.handlers[name] = hndl
}

// ServeHTTP handles every HTTP request.
func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.URL.Path, "/api") {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("400 Bad Request - invalid path"))

		return
	}

	// Request validation.
	pathSplitted := strings.Split(r.URL.Path, "/")
	if len(pathSplitted) < 3 {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("400 Bad Request - invalid path"))
	}

	handler, found := h.handlers[pathSplitted[2]]
	if !found {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("400 Bad Request - invalid application name"))

		return
	}

	requestContext := r.Context()

	// Compose metric name.
	metricName := strings.Join(pathSplitted[3:], "/")
	ctx := context.WithValue(requestContext, common.ContextKeyMetric, metricName)

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(handler(ctx)))
}

package httpserver

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"go.dev.pztrn.name/metricator/internal/common"
	"go.dev.pztrn.name/metricator/internal/models"
)

// nolint:gochecknoglobals
var (
	errInvalidAPIVersion  = errors.New("invalid API version")
	errInvalidApplication = errors.New("invalid application")
	errInvalidPath        = errors.New("invalid path")
	errNoAppsRegistered   = errors.New("no applications registered")
	errNoData             = errors.New("no data")

	supportedAPIVersions = []int{1}
)

// HTTP requests handler.
type handler struct {
	handlers map[string]common.HTTPHandlerFunc
}

// Gets applications list from handlers map.
func (h *handler) getAppsList() ([]byte, error) {
	apps := make([]string, 0, len(h.handlers))

	for appName := range h.handlers {
		apps = append(apps, appName)
	}

	appsList, err := json.Marshal(apps)
	if err != nil {
		// ToDo: log error
		return nil, errNoAppsRegistered
	}

	return appsList, nil
}

// Gets request information from URL. Returns a structure with filled request
// info and error if it occurs.
func (h *handler) getRequestInfo(req *http.Request) (*models.RequestInfo, error) {
	// Request isn't for API or isn't versioned.
	if !strings.HasPrefix(req.URL.Path, "/api/v") {
		return nil, errInvalidPath
	}

	// Note: first element will always be empty!
	pathSplitted := strings.Split(req.URL.Path, "/")

	// Request is for API but not enough items in URL was passed.
	if len(pathSplitted) < 4 {
		return nil, errInvalidPath
	}

	var (
		appName     string
		metricName  string
		requestType string
	)

	// Parse API version.
	apiVersionRaw := strings.TrimLeft(pathSplitted[2], "v")

	apiVersion, err := strconv.Atoi(apiVersionRaw)
	if err != nil {
		// ToDo: log error
		return nil, errInvalidAPIVersion
	}

	// Check used API version.
	var supportedAPIVersionUsed bool

	for _, version := range supportedAPIVersions {
		if apiVersion == version {
			supportedAPIVersionUsed = true

			break
		}
	}

	if !supportedAPIVersionUsed {
		return nil, errInvalidAPIVersion
	}

	// Get request type and key.
	requestType = pathSplitted[3]

	if len(pathSplitted) >= 5 {
		appName = pathSplitted[4]
	}

	if len(pathSplitted) >= 6 {
		metricName = strings.Join(pathSplitted[5:], "/")
	}

	reqInfo := &models.RequestInfo{
		APIVersion:  apiVersion,
		Application: appName,
		Metric:      metricName,
		RequestType: requestType,
	}

	return reqInfo, nil
}

// Registers request's handler.
func (h *handler) register(appName string, hndl common.HTTPHandlerFunc) {
	h.handlers[appName] = hndl
}

// ServeHTTP handles every HTTP request.
func (h *handler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	startTime := time.Now()

	defer func() {
		requestDuration := time.Since(startTime)

		log.Printf("[HTTP Request] from %s to %s, duration %.4fs\n", req.RemoteAddr, req.URL.Path, requestDuration.Seconds())
	}()

	// Validate request and extract needed info.
	rInfo, err := h.getRequestInfo(req)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		_, _ = writer.Write([]byte("400 bad request - " + err.Error()))

		return
	}

	// Process request type. Here we process only known requests types,
	// all other requests will produce HTTP 400 error.
	switch rInfo.RequestType {
	// ToDo: move to constants.
	case "apps_list":
		appsList, err := h.getAppsList()
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			_, _ = writer.Write([]byte("400 bad request - " + err.Error()))

			return
		}

		writer.WriteHeader(http.StatusOK)
		_, _ = writer.Write(appsList)

		return
	case "info":
		infoData := struct {
			Branch     string
			Build      string
			CommitHash string
			Version    string
		}{
			Branch:     common.Branch,
			Build:      common.Build,
			CommitHash: common.CommitHash,
			Version:    common.Version,
		}

		// nolint:errchkjson
		infoBytes, _ := json.Marshal(infoData)

		writer.WriteHeader(http.StatusOK)
		_, _ = writer.Write(infoBytes)

		return
	case "metrics":
		handler, found := h.handlers[rInfo.Application]
		if !found {
			writer.WriteHeader(http.StatusBadRequest)
			_, _ = writer.Write([]byte("400 bad request - " + errInvalidApplication.Error()))

			return
		}

		// Get data from handler.
		data := handler(rInfo)
		if data == "" {
			writer.WriteHeader(http.StatusBadRequest)
			_, _ = writer.Write([]byte("400 bad request - " + errNoData.Error()))
		}

		writer.WriteHeader(http.StatusOK)
		_, _ = writer.Write([]byte(data))

		return
	}

	writer.WriteHeader(http.StatusBadRequest)
	_, _ = writer.Write([]byte("400 bad request - " + errInvalidPath.Error()))
}

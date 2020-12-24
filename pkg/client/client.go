package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"go.dev.pztrn.name/metricator/internal/logger"
	"go.dev.pztrn.name/metricator/pkg/schema"
)

// Client is a Metricator client that is ready to be used in other applications
// or libraries.
type Client struct {
	config *Config
	logger *logger.Logger

	httpClient *http.Client
}

// NewClient creates new Metricator client.
func NewClient(config *Config, logger *logger.Logger) *Client {
	c := &Client{
		config: config,
		logger: logger,
	}
	c.initialize()

	return c
}

// Executes request and parses it's contents.
func (c *Client) executeAndParse(req *http.Request, dest interface{}) error {
	c.logger.Debugf("Executing HTTP request to %s%s", c.config.Host, req.URL.RequestURI())

	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*time.Duration(c.config.Timeout))
	defer cancelFunc()

	req = req.WithContext(ctx)

	response, err := c.httpClient.Do(req)
	if err != nil {
		c.logger.Infoln("Failed to execute request to Metricator:", err.Error())

		return fmt.Errorf("metricator client: %w", err)
	}

	defer response.Body.Close()

	respData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		c.logger.Infoln("Failed to read response body:", err.Error())

		return fmt.Errorf("metricator client: %w", err)
	}

	err = json.Unmarshal(respData, dest)
	if err != nil {
		c.logger.Infoln("Failed to parse response:", err.Error())

		return fmt.Errorf("metricator client: %w", err)
	}

	return nil
}

// GetAppsList returns a slice with applications that was registered at Metricator.
func (c *Client) GetAppsList() schema.AppsList {
	address := fmt.Sprintf("%s/api/v1/apps_list", c.config.Host)

	// Request's context sets in c.executeAndParse, so:
	// nolint:noctx
	req, err := http.NewRequest("GET", address, nil)
	if err != nil {
		c.logger.Infoln("Failed to create HTTP request:", err.Error())

		return nil
	}

	appsList := make(schema.AppsList, 0)

	err = c.executeAndParse(req, &appsList)
	if err != nil {
		return nil
	}

	return appsList
}

// GetMetric returns value for metric.
func (c *Client) GetMetric(appName, metricName string) interface{} {
	address := fmt.Sprintf("%s/api/v1/metrics/%s/%s", c.config.Host, appName, metricName)

	// Request's context sets in c.executeAndParse, so:
	// nolint:noctx
	req, err := http.NewRequest("GET", address, nil)
	if err != nil {
		c.logger.Infoln("Failed to create HTTP request:", err.Error())

		return ""
	}

	var data interface{}

	err = c.executeAndParse(req, &data)
	if err != nil {
		return ""
	}

	return data
}

// GetMetricsList returns a slice with metrics names for passed application.
func (c *Client) GetMetricsList(appName string) schema.Metrics {
	address := fmt.Sprintf("%s/api/v1/metrics/%s", c.config.Host, appName)

	// Request's context sets in c.executeAndParse, so:
	// nolint:noctx
	req, err := http.NewRequest("GET", address, nil)
	if err != nil {
		c.logger.Infoln("Failed to create HTTP request:", err.Error())

		return nil
	}

	data := make(schema.Metrics, 0)

	err = c.executeAndParse(req, &data)
	if err != nil {
		return nil
	}

	return data
}

// Initializes internal states and storages.
func (c *Client) initialize() {
	// We do not need to set everything for client actually, so:
	// nolint:exhaustivestruct
	c.httpClient = &http.Client{
		Timeout: time.Second * time.Duration(c.config.Timeout),
	}
}

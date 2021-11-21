package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"go.dev.pztrn.name/metricator/internal/common"
	"go.dev.pztrn.name/metricator/internal/configuration"
	"go.dev.pztrn.name/metricator/internal/logger"
	"go.dev.pztrn.name/metricator/pkg/client"
	"go.dev.pztrn.name/metricator/pkg/schema"
)

// nolint:gochecknoglobals
var (
	application       = flag.String("application", "", "Application to query.")
	appsList          = flag.Bool("apps-list", false, "Show application's list registered at Metricator.")
	metricatorHost    = flag.String("metricator-host", "", "IP address or domain on which Metricator is available")
	metricatorTimeout = flag.Int("metricator-timeout", 5, "Timeout for requests sent to Metricator.")
	metricsList       = flag.Bool("metrics-list", false, "Show metrics list. Requires 'application' parameter.")
	metric            = flag.String("metric", "", "Metric data to retrieve. Requires 'application' parameter.")
	output            = flag.String("output", "json", "Output format. Can be 'json' or 'plain-by-line'.")
)

// This function uses fmt.Println to print lines without timestamps to make it easy
// to parse output, so:
// nolint:forbidigo
func main() {
	config := configuration.NewConfig()

	// Parse configuration.
	flag.Parse()

	err := config.Parse()
	if err != nil {
		log.Fatalln("Failed to parse configuration:", err.Error())
	}

	logger := logger.NewLogger(config.Logger)

	logger.Debugf("Starting Metricator client, version %s from branch %s (build #%s, commit hash %s)\n",
		common.Version,
		common.Branch,
		common.Build,
		common.CommitHash,
	)

	// Check configuration.
	// We cannot work at all if host isn't defined.
	if *metricatorHost == "" {
		logger.Infoln("Host isn't defined.")

		flag.PrintDefaults()
		os.Exit(1)
	}

	// If nothing is requested - show error message.
	if !*appsList && !*metricsList && *metric == "" {
		logger.Infoln("No action specified.")

		flag.PrintDefaults()
		os.Exit(1)
	}

	// When asking to metrics list we need application to be defined.
	if *metricsList && *application == "" {
		logger.Infoln("Getting metrics list requires 'application' parameter to be filled.")

		flag.PrintDefaults()
		os.Exit(1)
	}

	// When asking for specific metric we need application to be defined.
	if *metric != "" && *application == "" {
		logger.Infoln("Getting metric data requires 'application' parameter to be filled.")

		flag.PrintDefaults()
		os.Exit(1)
	}

	clientConfig := &client.Config{
		Host:    *metricatorHost,
		Timeout: *metricatorTimeout,
	}

	clnt := client.NewClient(clientConfig, logger)

	var data interface{}

	switch {
	case *appsList:
		data = clnt.GetAppsList()
	case *metricsList:
		data = clnt.GetMetricsList(*application)
	case *metric != "":
		data = clnt.GetMetric(*application, *metric)
	}

	switch *output {
	case "json":
		dataAsBytes, err := json.Marshal(data)
		if err != nil {
			logger.Infoln("Failed to marshal data from Metricator:", err.Error())

			os.Exit(2)
		}

		fmt.Println(string(dataAsBytes))
	case "plain-by-lines":
		// For plain mode if we request metric - we should just print it and exit.
		if *metric != "" {
			fmt.Println(data)
			os.Exit(0)
		}

		dataToPrint := []string{}

		switch {
		case *appsList:
			appsListData, ok := data.(schema.AppsList)
			if !ok {
				logger.Infoln("Failed to cast parsed data into schema.AppsList!")

				os.Exit(3)
			}

			for _, app := range appsListData {
				dataToPrint = append(dataToPrint, app)
			}
		case *metric != "":
			metricData, ok := data.(string)
			if !ok {
				logger.Infoln("Failed to cast parsed data into string!")

				os.Exit(3)
			}

			dataToPrint = append(dataToPrint, metricData)
		case *metricsList:
			metricsData, ok := data.(schema.Metrics)
			if !ok {
				logger.Infoln("Failed to cast parsed data into schema.Metrics!")

				os.Exit(3)
			}

			for _, metric := range metricsData {
				dataToPrint = append(dataToPrint, metric.Name)
			}
		}

		for _, line := range dataToPrint {
			fmt.Println(line)
		}
	}
}

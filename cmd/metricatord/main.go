package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"go.dev.pztrn.name/metricator/internal/application"
	"go.dev.pztrn.name/metricator/internal/common"
	"go.dev.pztrn.name/metricator/internal/configuration"
	"go.dev.pztrn.name/metricator/internal/httpserver"
	"go.dev.pztrn.name/metricator/internal/logger"
)

func main() {
	log.Printf("Starting Metricator daemon, version %s from branch %s (build #%s, commit hash %s)\n",
		common.Version,
		common.Branch,
		common.Build,
		common.CommitHash,
	)

	mainCtx, cancelFunc := context.WithCancel(context.Background())
	config := configuration.NewConfig()

	// Parse configuration.
	flag.Parse()

	if err := config.Parse(); err != nil {
		log.Fatalln("Failed to parse configuration:", err.Error())
	}

	log.Printf("Configuration parsed: %+v\n", config)

	// Create applications.
	apps := make([]*application.Application, 0, len(config.Applications))

	logger := logger.NewLogger(config.Logger)
	httpSrv, httpStopped := httpserver.NewHTTPServer(mainCtx, config, logger)

	for appName, appConfig := range config.Applications {
		app := application.NewApplication(mainCtx, appName, appConfig, logger)
		app.Start()

		httpSrv.RegisterHandlerForApplication(appName, app.GetHandler())

		apps = append(apps, app)
	}

	httpSrv.Start()

	log.Println("Metricator is started and ready to serve requests")

	// CTRL+C handler.
	signalHandler := make(chan os.Signal, 1)
	shutdownDone := make(chan bool, 1)

	signal.Notify(signalHandler, os.Interrupt, syscall.SIGTERM)

	go func(apps []*application.Application) {
		<-signalHandler
		cancelFunc()

		for _, app := range apps {
			<-app.GetDoneChan()
		}

		<-httpStopped

		shutdownDone <- true
	}(apps)

	<-shutdownDone
	log.Println("Metricator stopped")
	os.Exit(0)
}

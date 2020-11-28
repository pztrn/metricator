package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"go.dev.pztrn.name/metricator/internal/common"
	"go.dev.pztrn.name/metricator/internal/configuration"
	"go.dev.pztrn.name/metricator/internal/datastore"
	"go.dev.pztrn.name/metricator/internal/httpserver"
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

	httpSrv, httpStopped := httpserver.NewHTTPServer(mainCtx, config)
	dataStore, dataStoreStopped := datastore.NewDataStore(mainCtx, config)

	flag.Parse()
	err := config.Parse()
	if err != nil {
		log.Fatalln("Failed to parse configuration:", err.Error())
	}
	log.Printf("Configuration parsed: %+v\n", config)

	dataStore.Start()
	httpSrv.Start()

	log.Println("Metricator is started and ready to serve requests")

	// CTRL+C handler.
	signalHandler := make(chan os.Signal, 1)
	shutdownDone := make(chan bool, 1)

	signal.Notify(signalHandler, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-signalHandler
		cancelFunc()

		<-dataStoreStopped
		<-httpStopped

		shutdownDone <- true
	}()

	<-shutdownDone
	log.Println("Metricator stopped")
	os.Exit(0)
}

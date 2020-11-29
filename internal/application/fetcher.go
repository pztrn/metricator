package application

import (
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// Fetches data from remote endpoint, parses it and updates in storage.
func (a *Application) fetch() {
	// Do not do anything if fetching is running.
	// ToDo: maybe another approach?
	a.fetchIsRunningMutex.RLock()
	isFetching := a.fetchIsRunning
	a.fetchIsRunningMutex.RUnlock()

	if isFetching {
		return
	}

	a.fetchIsRunningMutex.Lock()
	a.fetchIsRunning = true
	a.fetchIsRunningMutex.Unlock()

	log.Println("Fetching data for", a.name)

	req, err := http.NewRequestWithContext(a.ctx, "GET", a.config.Endpoint, nil)
	if err != nil {
		log.Println("Failed to create request for", a.name, "metrics:", err.Error())

		return
	}

	for header, value := range a.config.Headers {
		req.Header.Add(header, value)
	}

	resp, err := a.httpClient.Do(req)
	if err != nil {
		log.Println("Failed to execute request for", a.name, "metrics:", err.Error())

		return
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Failed to read response body for", a.name, "metrics:", err.Error())

		return
	}

	data := a.parse(string(body))

	a.storage.Put(data)

	a.fetchIsRunningMutex.Lock()
	a.fetchIsRunning = false
	a.fetchIsRunningMutex.Unlock()
}

// Configures and starts Prometheus data fetching goroutine.
func (a *Application) startFetcher() {
	fetchTicker := time.NewTicker(a.config.TimeBetweenRequests)

	// nolint:exhaustivestruct
	a.httpClient = &http.Client{
		Timeout: time.Second * 5,
	}

	defer log.Println("Fetcher for", a.name, "completed")

	// First fetch should be executed ASAP.
	a.fetch()

	for {
		select {
		case <-a.ctx.Done():
			return
		case <-fetchTicker.C:
			a.fetch()
		}
	}
}

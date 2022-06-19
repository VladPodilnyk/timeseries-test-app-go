package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/VladPodilnyk/timeseries-test-app-go/internal/api"
	"github.com/VladPodilnyk/timeseries-test-app-go/internal/config"
	"github.com/VladPodilnyk/timeseries-test-app-go/internal/repo"
)

func main() {

	// init later :))))
	var repo repo.TimeSeries
	var config config.LimitsConfig

	// Initialize a new logger which writes messages to the standard out stream,
	// prefixed with the current date and time.
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	// Initialize an app
	app := api.New(repo, config)

	srv := &http.Server{
		Addr:         "localhost",
		Handler:      app.Routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	logger.Printf("starting %s server on %s", "kitty", srv.Addr)
	err := srv.ListenAndServe()
	logger.Fatal(err)
}

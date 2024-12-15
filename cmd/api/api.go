package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/edwrdc/digitally/docs" // docs is generated by Swag CLI, you have to import it.
	"github.com/edwrdc/digitally/internal/store"
)

type application struct {
	config config
	store  *store.Storage
}

type config struct {
	addr   string
	env    string
	db     dbConfig
	apiURL string
}

type dbConfig struct {
	dsn          string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  time.Duration
}

func (app *application) run() error {
	// Docs
	docs.SwaggerInfo.Version = version
	docs.SwaggerInfo.Host = app.config.apiURL
	docs.SwaggerInfo.BasePath = "/v1"

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", app.config.addr),
		Handler:      app.routes(),
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  10 * time.Second,
		IdleTimeout:  time.Minute,
	}

	return srv.ListenAndServe()
}

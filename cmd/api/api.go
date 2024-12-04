package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/edwrdc/digitally/internal/store"
)

type application struct {
	config config
	store  *store.Storage
}

type config struct {
	addr string
	env  string
	db   dbConfig
}

type dbConfig struct {
	dsn          string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  time.Duration
}

func (app *application) run() error {

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", app.config.addr),
		Handler:      app.routes(),
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  10 * time.Second,
		IdleTimeout:  time.Minute,
	}

	return srv.ListenAndServe()
}

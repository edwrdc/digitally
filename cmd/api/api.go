package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/edwrdc/digitally/docs" // docs is generated by Swag CLI, you have to import it.
	"github.com/edwrdc/digitally/internal/auth"
	"github.com/edwrdc/digitally/internal/mailer"
	"github.com/edwrdc/digitally/internal/store"
	"go.uber.org/zap"
)

type application struct {
	config        config
	store         *store.Storage
	logger        *zap.SugaredLogger
	mailer        mailer.Client
	authenticator auth.Authenticator
}

type config struct {
	addr        string
	env         string
	db          dbConfig
	apiURL      string
	frontendURL string
	mail        mailConfig
	auth        authConfig
}

type dbConfig struct {
	dsn          string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  time.Duration
}

type mailConfig struct {
	// sendGrid  sendGridConfig
	mailtrap  mailtrapConfig
	exp       time.Duration
	fromEmail string
}

type authConfig struct {
	basic basicAuthConfig
	token tokenAuthConfig
}

type basicAuthConfig struct {
	user string
	pass string
}

type tokenAuthConfig struct {
	secret string
	expiry time.Duration
	iss    string
}

// type sendGridConfig struct {
// 	apiKey string
// }

type mailtrapConfig struct {
	apiKey  string
	inboxID string
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

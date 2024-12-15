package main

import (
	"time"

	"github.com/edwrdc/digitally/internal/db"
	"github.com/edwrdc/digitally/internal/env"
	"github.com/edwrdc/digitally/internal/store"
	"go.uber.org/zap"
)

const version = "0.0.1"

//	@title			Digitally API
//	@description	API for Digitally, a platform for buying and selling digital products.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@BasePath					/v1
//	@securityDefinitions.apikey	ApiKeyAuth
//	@in							header
//	@name						Authorization
//	@description				JWT authorization header

func main() {

	cfg := config{
		addr:   env.Get("API_PORT", ":8080"),
		env:    env.Get("API_ENV", "development"),
		apiURL: env.Get("EXTERNAL_API_URL", "localhost:6969"),
		db: dbConfig{
			dsn:          env.Get("DB_DSN", "postgres://admin:adminpassword@localhost:5432/digitally?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 25),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 25),
			maxIdleTime:  time.Duration(env.GetInt("DB_MAX_IDLE_TIME", 15)) * time.Minute,
		},
		mail: mailConfig{
			exp: time.Duration(env.GetInt("MAIL_EXPIRY", 3)) * time.Hour,
		},
	}

	// Logger
	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()

	// Database
	db, err := db.New(cfg.db.dsn, cfg.db.maxOpenConns, cfg.db.maxIdleConns, cfg.db.maxIdleTime)
	if err != nil {
		logger.Panic(err)
	}
	defer db.Close()

	logger.Info("Established connection pool to database")

	store := store.New(db)

	app := &application{
		config: cfg,
		store:  store,
		logger: logger,
	}

	app.logger.Infow("Server Started", "env", app.config.env, "addr", app.config.addr)

	logger.Fatal(app.run())

}

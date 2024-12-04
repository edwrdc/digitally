package main

import (
	"log"
	"time"

	"github.com/edwrdc/digitally/internal/db"
	"github.com/edwrdc/digitally/internal/env"
	"github.com/edwrdc/digitally/internal/store"
)

func main() {

	cfg := config{
		addr: env.Get("API_PORT", ":8080"),
		env:  env.Get("API_ENV", "development"),
		db: dbConfig{
			dsn:          env.Get("DB_DSN", "postgres://admin:adminpassword@localhost:5432/digitally?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 25),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 25),
			maxIdleTime:  time.Duration(env.GetInt("DB_MAX_IDLE_TIME", 15)) * time.Minute,
		},
	}

	db, err := db.New(cfg.db.dsn, cfg.db.maxOpenConns, cfg.db.maxIdleConns, cfg.db.maxIdleTime)
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()

	log.Println("Established connection pool to database")

	store := store.New(db)

	app := &application{
		config: cfg,
		store:  store,
	}

	log.Printf("Starting %s server on %s", cfg.env, cfg.addr)

	log.Fatal(app.run())

}

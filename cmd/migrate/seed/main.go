package main

import (
	"log"
	"time"

	"github.com/edwrdc/digitally/internal/db"
	"github.com/edwrdc/digitally/internal/env"
	"github.com/edwrdc/digitally/internal/store"
)

func main() {

	addr := env.Get("DB_DSN", "postgres://admin:adminpassword@localhost:5432/digitally?sslmode=disable")
	conn, err := db.New(addr, 10, 10, 10*time.Second)
	if err != nil {
		log.Fatalf("error connecting to database: %v", err)
	}

	defer conn.Close()
	store := store.New(conn)

	if err := db.Seed(store, conn); err != nil {
		log.Fatalf("error seeding database: %v", err)
	}
}

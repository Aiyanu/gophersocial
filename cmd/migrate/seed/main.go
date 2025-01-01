package main

import (
	"log"

	"github.com/gophersocial/internal/db"
	"github.com/gophersocial/internal/env"
	"github.com/gophersocial/internal/store"
)

func main() {
	addr := env.GetString("DB_ADDR", "postgres://aiyanu:incorrect@localhost/social?sslmode=disable")
	conn, err := db.New(addr, 3, 3, "15m")
	if err != nil {
		log.Fatal(err)
	}
	store := store.NewStorage(conn)

	db.Seed(store)
}

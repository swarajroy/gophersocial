package main

import (
	"log"
	"time"

	"github.com/swarajroy/gophersocial/internal/db"
	"github.com/swarajroy/gophersocial/internal/store"
)

const (
	DB_ADDR = "postgres://admin:adminpassword@localhost/social?sslmode=disable"
)

func main() {
	maxIdleTime := 15 * time.Minute
	conn, err := db.New(DB_ADDR, 3, 3, maxIdleTime)
	if err != nil {
		log.Fatalf("error occurred while establishing connection to the DB!")
	}
	defer conn.Close()

	store := store.NewStorage(conn)
	db.Seed(store, conn)
}

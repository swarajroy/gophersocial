package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/swarajroy/gophersocial/internal/env"
)

func main() {
	cfg := config{
		addr: env.GetString("ADDR", ":8080"),
	}
	a := &application{
		config: cfg,
	}
	mux := a.mount()

	log.Fatal(a.run(mux))
}

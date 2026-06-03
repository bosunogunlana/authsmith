package main

import (
	"log"
	"net/http"

	"github.com/bosunogunlana/authsmith/internal/config"
	"github.com/bosunogunlana/authsmith/internal/httpserver"
)

func main() {
	config := config.Load()
	log.Printf("config Loaded %+v", config)

	mux := httpserver.ServerMux()

	if err := http.ListenAndServe(":"+config.Port, mux); err != nil {
		log.Fatalf("server failed to start: %v", err)
	}
}

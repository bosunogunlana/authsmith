package main

import (
	"log"
	"net/http"

	"github.com/bosunogunlana/authsmith/internal/config"
	"github.com/bosunogunlana/authsmith/internal/database"
	"github.com/bosunogunlana/authsmith/internal/httpserver"
)

func main() {
	config := config.Load()
	log.Printf("config loaded")

	_, err := database.Connect(config.DatabaseDSN)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}

	mux := httpserver.ServerMux()

	if err := http.ListenAndServe(":"+config.Port, mux); err != nil {
		log.Fatalf("server failed to start: %v", err)
	}
}

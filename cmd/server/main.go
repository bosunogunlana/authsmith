package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/bosunogunlana/authsmith/internal/config"
	"github.com/bosunogunlana/authsmith/internal/database"
	"github.com/bosunogunlana/authsmith/internal/handlers"
	"github.com/bosunogunlana/authsmith/internal/httpserver"
)

func main() {
	cfg := config.Load()
	log.Printf("config loaded")

	db, err := database.Connect(cfg.DatabaseDSN)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	if err := database.RunMigration(db, "migrations"); err != nil {
		log.Fatalf("db migration failed: %v", err)
	}

	tmpls, err := template.ParseGlob("templates/*.html")
	if err != nil {
		log.Fatalf("failed to load templates: %v", err)
	}

	h := &handlers.Handlers{
		DB:        db,
		Templates: tmpls,
	}

	mux := httpserver.ServerMux(h)

	log.Printf("server starting on :%s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, mux); err != nil {
		log.Fatalf("server failed to start: %v", err)
	}
}

package main

import (
	"log"

	"github.com/bosunogunlana/authsmith/internal/auth"
	"github.com/bosunogunlana/authsmith/internal/config"
	"github.com/bosunogunlana/authsmith/internal/database"
)

func main() {
	cfg := config.Load()

	db, err := database.Connect(cfg.DatabaseDSN)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	passwordDigest, err := auth.HashPassword("password123")
	if err != nil {
		log.Fatalf("failed to hash password: %v", err)
	}

	user, err := database.CreateUser(db, "test@example.com", passwordDigest, "Test User")
	if err != nil {
		log.Fatalf("failed to create user: %v", err)
	}

	log.Printf("created user: id=%s email=%s", user.ID, user.Email)
}

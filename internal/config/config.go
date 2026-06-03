package config

import (
	"log"
	"os"
)

type Config struct {
	Port        string
	DatabaseDSN string
}

func Load() Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	databaseDSN := os.Getenv("DATABASE_DSN")
	if databaseDSN == "" {
		log.Fatal("missing database dsn")
	}

	return Config{
		Port:        port,
		DatabaseDSN: databaseDSN,
	}
}

package httpserver

import (
	"net/http"

	"github.com/bosunogunlana/authsmith/internal/handlers"
)

func ServerMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", handlers.HealthHandler)

	return mux
}

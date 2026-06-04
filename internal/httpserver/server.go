package httpserver

import (
	"net/http"

	"github.com/bosunogunlana/authsmith/internal/handlers"
)

func ServerMux(h *handlers.Handlers) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", h.HealthHandler)
	mux.HandleFunc("GET /login", h.GetLoginHandler)
	mux.HandleFunc("POST /login", h.PostLoginHandler)
	mux.HandleFunc("POST /logout", h.LogoutHandler)
	mux.HandleFunc("GET /me", h.GetUserProfile)

	return mux
}

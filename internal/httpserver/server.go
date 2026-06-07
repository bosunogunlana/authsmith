package httpserver

import (
	"net/http"

	"github.com/bosunogunlana/authsmith/internal/handlers"
)

func ServerMux(h *handlers.Handlers) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", h.GetUserProfile)
	mux.HandleFunc("GET /health", h.HealthHandler)
	mux.HandleFunc("GET /login", h.GetLoginHandler)
	mux.HandleFunc("POST /login", h.PostLoginHandler)
	mux.HandleFunc("POST /logout", h.LogoutHandler)
	mux.HandleFunc("GET /me", h.GetUserProfile)

	mux.HandleFunc("GET /developer/clients/new", h.NewClientHandler)
	mux.HandleFunc("POST /developer/clients", h.CreateClientHandler)
	mux.HandleFunc("GET /developer/clients", h.ListClientsHandler)
	mux.HandleFunc("GET /developer/client/{client_id}", h.ShowClientHandler)

	mux.HandleFunc("GET /oauth/authorize", h.AuthorizeHandler)
	mux.HandleFunc("GET /oauth/consent", h.GetConsentHandler)
	mux.HandleFunc("POST /oauth/consent", h.PostConsentHandler)
	mux.HandleFunc("POST /oauth/token", h.PostTokenHandler)

	// APIs
	mux.HandleFunc("GET /api/user", h.GetAPIUserProfileHandler)

	return mux
}

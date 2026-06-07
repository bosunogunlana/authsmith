package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/bosunogunlana/authsmith/internal/database"
)

type APIUserResponse struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

func (h *Handlers) GetAPIUserProfileHandler(w http.ResponseWriter, r *http.Request) {
	accessToken, err := h.extractBearerToken(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if !hasScope(accessToken.Scopes, "profile:read") {
		http.Error(w, "insufficient_scope", http.StatusForbidden)
		return
	}

	user, err := database.GetUserByID(h.DB, accessToken.UserID)
	if err != nil {
		http.Error(w, "not_found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Pragma", "no-cache")
	w.WriteHeader(http.StatusOK)

	resp := APIUserResponse{
		ID:    user.ID,
		Email: user.Email,
		Name:  user.Name,
	}

	json.NewEncoder(w).Encode(resp)
}

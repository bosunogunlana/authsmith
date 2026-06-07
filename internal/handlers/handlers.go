package handlers

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"html/template"
	"net/http"
	"strings"

	"github.com/bosunogunlana/authsmith/internal/database"
	"github.com/bosunogunlana/authsmith/internal/models"
)

type Handlers struct {
	DB        *sql.DB
	Templates *template.Template
}

func (h *Handlers) isAuthenticated(r *http.Request) (models.Session, bool) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		return models.Session{}, false
	}

	hash := sha256.Sum256([]byte(cookie.Value))
	digest := hex.EncodeToString(hash[:])

	session, err := database.GetSessionByTokenDigest(h.DB, digest)
	if err != nil || session.ID == "" {
		return models.Session{}, false
	}

	return session, true
}

func (h *Handlers) requireAuth(w http.ResponseWriter, r *http.Request) (models.Session, bool) {
	session, ok := h.isAuthenticated(r)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
	return session, ok
}

func (h *Handlers) extractBearerToken(r *http.Request) (models.AccessToken, error) {
	authHeader := r.Header.Get("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return models.AccessToken{}, errors.New("missing bearer token")
	}

	rawToken := strings.TrimPrefix(authHeader, "Bearer ")
	hash := sha256.Sum256([]byte(rawToken))
	digest := hex.EncodeToString(hash[:])

	return database.GetAccessTokenByDigest(h.DB, digest)
}

func hasScope(tokenScopes, required string) bool {
	for _, s := range strings.Fields(tokenScopes) {
		if s == required {
			return true
		}
	}

	return false
}
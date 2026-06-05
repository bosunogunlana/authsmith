package handlers

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"html/template"
	"net/http"

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

package handlers

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"

	"github.com/bosunogunlana/authsmith/internal/database"
)

func (h *Handlers) GetUserProfile(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusUnauthorized)
		return
	}

	hash := sha256.Sum256([]byte(cookie.Value))
	digest := hex.EncodeToString(hash[:])

	session, err := database.GetSessionByTokenDigest(h.DB, digest)
	if err != nil || session.ID == "" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	user, err := database.GetUserByID(h.DB, session.UserID)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	h.Templates.ExecuteTemplate(w, "profile.html", user)
}
package handlers

import (
	"net/http"

	"github.com/bosunogunlana/authsmith/internal/database"
)

func (h *Handlers) GetUserProfile(w http.ResponseWriter, r *http.Request) {
	session, ok := h.requireAuth(w, r)
	if !ok {
		return
	}

	user, err := database.GetUserByID(h.DB, session.UserID)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	h.Templates.ExecuteTemplate(w, "profile.html", user)
}

package handlers

import (
	"net/http"
	"time"

	"github.com/bosunogunlana/authsmith/internal/auth"
	"github.com/bosunogunlana/authsmith/internal/database"
)

type loginPageData struct {
	Error string
}

func (h *Handlers) GetLoginHandler(w http.ResponseWriter, r *http.Request) {
	h.Templates.ExecuteTemplate(w, "login.html", loginPageData{})
}

func (h *Handlers) PostLoginHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")

	user, err := database.GetUserByEmail(h.DB, email)
	if err != nil {
		h.renderLoginError(w, "invalid credentials")
		return
	}

	if user.ID == "" {
		h.renderLoginError(w, "invalid credentials")
		return
	}

	if err := auth.CheckPassword(password, user.PasswordDigest); err != nil {
		h.renderLoginError(w, "invalid credentials")
		return
	}

	token, digest, err := auth.GenerateToken()
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	_, err = database.CreateSession(h.DB, user.ID, digest)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    token,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
	})

	http.Redirect(w, r, "/me", http.StatusSeeOther)
}

func (h *Handlers) renderLoginError(w http.ResponseWriter, message string) {
	h.Templates.ExecuteTemplate(w, "login.html", loginPageData{Error: message})
}

package handlers

import (
	"net/http"

	"github.com/bosunogunlana/authsmith/internal/database"
	"github.com/bosunogunlana/authsmith/internal/models"
	"github.com/bosunogunlana/authsmith/internal/oauth"
)

type clientNewPageData struct {
	Errors []string
}

func (h *Handlers) NewClientHandler(w http.ResponseWriter, r *http.Request) {
	h.Templates.ExecuteTemplate(w, "client_new.html", clientNewPageData{})
}

type clientDetailsPageData struct {
	Client    models.OAuthClient
	RawSecret string
}

func (h *Handlers) CreateClientHandler(w http.ResponseWriter, r *http.Request) {
	session, ok := h.requireAuth(w, r)
	if !ok {
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	name := r.FormValue("name")
	clientType := r.FormValue("client_type")
	redirectURI := r.FormValue("redirect_uris")
	allowedScopes := r.FormValue("allowed_scopes")

	errorsInfo := []string{}
	if name == "" {
		errorsInfo = append(errorsInfo, "Invalid Name")
	}
	if clientType != "confidential" && clientType != "public" {
		errorsInfo = append(errorsInfo, "Invalid Client Type - must be eith Confidential or Public")
	}
	if redirectURI == "" {
		errorsInfo = append(errorsInfo, "Invalid Redirect URL")
	}
	if allowedScopes == "" {
		allowedScopes = "all"
	}

	if len(errorsInfo) > 0 {
		h.renderCreateError(w, errorsInfo)
		return
	}

	clientID, err := oauth.GenerateClientID()
	if err != nil {
		h.renderCreateError(w, []string{"Failed to create Oauth Client"})
		return
	}

	var clientToken, clientDigest string
	if clientType == "confidential" {
		clientToken, clientDigest, err = oauth.GenerateClientSecret()
		if err != nil {
			h.renderCreateError(w, []string{"Failed to create Oauth Client"})
			return
		}
	}

	client, err := database.CreateOAuthClient(
		h.DB,
		session.UserID,
		name,
		clientID,
		clientDigest,
		clientType,
		redirectURI,
		allowedScopes,
	)
	if err != nil {
		h.renderCreateError(w, []string{"Failed to create Oauth Client"})
		return
	}

	h.Templates.ExecuteTemplate(w, "client_detail.html", clientDetailsPageData{Client: client, RawSecret: clientToken})
}

type clientsPageData struct {
	Clients []models.OAuthClient
}

func (h *Handlers) ListClientsHandler(w http.ResponseWriter, r *http.Request) {
	session, ok := h.requireAuth(w, r)
	if !ok {
		return
	}

	clients, err := database.GetOAuthClientsByOwnerID(h.DB, session.UserID)
	if err != nil {
		http.Redirect(w, r, "/me", http.StatusSeeOther)
		return
	}

	h.Templates.ExecuteTemplate(w, "clients.html", clientsPageData{Clients: clients})
}

func (h *Handlers) ShowClientHandler(w http.ResponseWriter, r *http.Request) {
	session, ok := h.requireAuth(w, r)
	if !ok {
		return
	}

	clientID := r.PathValue("client_id")
	client, err := database.GetOAuthClientByClientID(h.DB, clientID)
	if err != nil {
		http.Redirect(w, r, "/developer/clients", http.StatusSeeOther)
		return
	}

	if client.ID == "" {
		http.NotFound(w, r)
		return
	}

	if client.OwnerUserID != session.UserID {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	h.Templates.ExecuteTemplate(w, "client_detail.html", clientDetailsPageData{Client: client})
}

func (h *Handlers) renderCreateError(w http.ResponseWriter, messages []string) {
	h.Templates.ExecuteTemplate(w, "client_new.html", clientNewPageData{Errors: messages})
}

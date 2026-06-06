package handlers

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/bosunogunlana/authsmith/internal/database"
	"github.com/bosunogunlana/authsmith/internal/oauth"
)

type oauthErrorPageData struct {
	Message   string
	ErrorCode string
}

func (h *Handlers) AuthorizeHandler(w http.ResponseWriter, r *http.Request) {
	authReq, err := oauth.ParseAuthorizationRequest(h.DB, r)
	if err != nil {
		h.Templates.ExecuteTemplate(w, "oauth_error.html", oauthErrorPageData{
			Message:   "The authorization request is invalid",
			ErrorCode: err.Error(),
		})
		return
	}

	_, ok := h.isAuthenticated(r)
	if !ok {
		next := url.QueryEscape(r.URL.String())
		http.Redirect(w, r, "/login?next="+next, http.StatusSeeOther)
		return
	}

	params := url.Values{}
	params.Set("response_type", authReq.ResponseType)
	params.Set("client_id", authReq.ClientID)
	params.Set("redirect_uri", authReq.RedirectURI)
	params.Set("scope", authReq.Scope)
	params.Set("state", authReq.State)
	params.Set("code_challenge", authReq.CodeChallenge)
	params.Set("code_challenge_method", authReq.CodeChallengeMethod)

	http.Redirect(w, r, "/oauth/consent?"+params.Encode(), http.StatusSeeOther)
}

type consentPageData struct {
	ClientName string
	UserEmail  string
	Scopes     []string
	// authorization request fields (for hidden inputs)
	ResponseType        string
	ClientID            string
	RedirectURI         string
	Scope               string
	State               string
	CodeChallenge       string
	CodeChallengeMethod string
}

func (h *Handlers) GetConsentHandler(w http.ResponseWriter, r *http.Request) {
	session, ok := h.isAuthenticated(r)
	if !ok {
		next := url.QueryEscape(r.URL.String())
		http.Redirect(w, r, "/login?next="+next, http.StatusSeeOther)
		return
	}

	authReq, err := oauth.ParseAuthorizationRequest(h.DB, r)
	if err != nil {
		h.Templates.ExecuteTemplate(w, "oauth_error.html", oauthErrorPageData{
			Message:   "The authorization request is invalid",
			ErrorCode: err.Error(),
		})
		return
	}

	oauthClient, err := database.GetOAuthClientByClientID(h.DB, authReq.ClientID)
	if err != nil {
		h.Templates.ExecuteTemplate(w, "oauth_error.html", oauthErrorPageData{
			Message:   "The authorization request is invalid",
			ErrorCode: err.Error(),
		})
		return
	}

	if oauthClient.ID == "" {
		h.Templates.ExecuteTemplate(w, "oauth_error.html", oauthErrorPageData{
			Message: "The Oauth Client is Invalid",
		})
		return
	}

	user, err := database.GetUserByID(h.DB, session.UserID)
	if err != nil {
		h.Templates.ExecuteTemplate(w, "oauth_error.html", oauthErrorPageData{
			Message: "Invalid User",
		})
		return
	}

	h.Templates.ExecuteTemplate(w, "consent.html", consentPageData{
		ClientName:          oauthClient.Name,
		UserEmail:           user.Email,
		Scopes:              strings.Fields(authReq.Scope),
		Scope:               authReq.Scope,
		ResponseType:        authReq.ResponseType,
		ClientID:            authReq.ClientID,
		RedirectURI:         authReq.RedirectURI,
		State:               authReq.State,
		CodeChallenge:       authReq.CodeChallenge,
		CodeChallengeMethod: authReq.CodeChallengeMethod,
	})
}

func (h *Handlers) PostConsentHandler(w http.ResponseWriter, r *http.Request) {
	session, ok := h.isAuthenticated(r)
	if !ok {
		// TODO: Add redirect url with oauth params
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	clientID := r.FormValue("client_id")
	client, err := database.GetOAuthClientByClientID(h.DB, clientID)

	params := url.Values{}
	state := r.FormValue("state")
	params.Set("state", state)

	redirectURI := r.FormValue("redirect_uri")
	if err != nil || client.ID == "" {
		params.Set("error", "access_denied")
		http.Redirect(w, r, redirectURI+"?"+params.Encode(), http.StatusSeeOther)
		return
	}

	consent := r.FormValue("decision")
	if consent == "deny" || consent == "" {
		params.Set("error", "access_denied")
		http.Redirect(w, r, redirectURI+"?"+params.Encode(), http.StatusSeeOther)
		return
	}

	code, digest, err := oauth.GenerateAuthCode()
	if err != nil {
		http.Error(w, "internal_server_error", http.StatusInternalServerError)
		return
	}

	scope := r.FormValue("scope")
	codeChallenge := r.FormValue("code_challenge")
	codeChallengeMethod := r.FormValue("code_challenge_method")

	if scope == "" || codeChallenge == "" || codeChallengeMethod == "" {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	_, err = database.CreateAuthorizationCode(
		h.DB,
		session.UserID,
		clientID,
		redirectURI,
		scope,
		digest,
		codeChallenge,
		codeChallengeMethod,
	)

	if err != nil {
		http.Error(w, "internal_server_error", http.StatusInternalServerError)
		return
	}

	params.Set("code", code)
	http.Redirect(w, r, redirectURI+"?"+params.Encode(), http.StatusSeeOther)
}

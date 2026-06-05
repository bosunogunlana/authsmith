package handlers

import (
	"net/http"
	"net/url"

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
			Message: "The authorization request is invalid",
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

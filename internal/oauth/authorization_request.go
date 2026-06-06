package oauth

import (
	"database/sql"
	"net/http"
	"strings"

	"github.com/bosunogunlana/authsmith/internal/database"
)

type AuthorizationRequest struct {
	ResponseType        string
	ClientID            string
	RedirectURI         string
	Scope               string
	State               string
	CodeChallenge       string
	CodeChallengeMethod string
}

func ParseAuthorizationRequest(db *sql.DB, r *http.Request) (AuthorizationRequest, error) {
	query := r.URL.Query()
	responseType := query.Get("response_type")
	if responseType != "code" {
		return AuthorizationRequest{}, InvalidResponseCodeError
	}

	clientID := query.Get("client_id")
	if clientID == "" {
		return AuthorizationRequest{}, InvalidClientIDError
	}
	clientOauth, err := database.GetOAuthClientByClientID(db, clientID)
	if err != nil || clientOauth.ID == "" {
		return AuthorizationRequest{}, InvalidClientIDError
	}

	redirectURI := query.Get("redirect_uri")
	if clientOauth.RedirectURIs != redirectURI {
		return AuthorizationRequest{}, InvalidRedirectURIError
	}

	rawScope := query.Get("scope")
	requestedScope := strings.Fields(rawScope)
	clientScopes := strings.Fields(clientOauth.AllowedScopes)
	allowedSet := make(map[string]bool)
	for _, s := range clientScopes {
		allowedSet[s] = true
	}
	for _, s := range requestedScope {
		if !allowedSet[s] {
			return AuthorizationRequest{}, InvalidScopeError
		}
	}

	state := query.Get("state")
	if state == "" {
		return AuthorizationRequest{}, InvalidStateError
	}

	codeChallenge := query.Get("code_challenge")
	codeChallengeMethod := query.Get("code_challenge_method")
	if codeChallenge == "" || codeChallengeMethod != "S256" {
		return AuthorizationRequest{}, InvalidCodeChallenge
	}

	return AuthorizationRequest{
		ResponseType:        responseType,
		RedirectURI:         redirectURI,
		ClientID:            clientID,
		Scope:               rawScope,
		State:               state,
		CodeChallenge:       codeChallenge,
		CodeChallengeMethod: codeChallengeMethod,
	}, nil
}

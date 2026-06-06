package handlers

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/bosunogunlana/authsmith/internal/database"
	"github.com/bosunogunlana/authsmith/internal/oauth"
)

func (h *Handlers) PostTokenHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm(); if err != nil {
		tokenError(w, "invalid_request", http.StatusBadRequest)
		return
	}

	grantType := r.FormValue("grant_type")
	if grantType != "authorization_code" {
		tokenError(w, "unsupported_grant_type", http.StatusBadRequest)
		return
	}

	clientID := r.FormValue("client_id")
	client, err := database.GetOAuthClientByClientID(h.DB, clientID)
	if err != nil || client.ID == "" {
		tokenError(w, "invalid_client", http.StatusUnauthorized)
		return
	}

	if err := oauth.AuthenticateClient(r, client); err != nil {
		tokenError(w, "invalid_client", http.StatusUnauthorized)
		return
	}

	code := r.FormValue("code")
	hashedCode := sha256.Sum256([]byte(code))
	codeDigest := hex.EncodeToString(hashedCode[:])
	authCode, err := database.GetAuthorizationCodeByDigest(h.DB,codeDigest)
	if err != nil || authCode.ID == "" {
		tokenError(w, "invalid_grant", http.StatusBadRequest)
		return
	}

	redirectURI := r.FormValue("redirect_uri")
	if authCode.ClientID != clientID || authCode.RedirectURI != redirectURI {
		tokenError(w, "invalid_grant", http.StatusBadRequest)
		return
	}

	codeVerifier := r.FormValue("code_verifier")
	err = oauth.ValidatePKCE(codeVerifier, authCode.CodeChallenge)
	if err != nil {
		tokenError(w, "invalid_grant", http.StatusUnauthorized)
		return
	}

	// Mark code as used, generate access_token and refresh_token and return the tokens
	txn, err := h.DB.Begin()
	if err != nil {
		tokenError(w, "server_error", http.StatusInternalServerError)
    return
	}
	defer txn.Rollback()

	rowsAffected, err := database.MarkAuthorizationCodeUsed(txn, authCode.ID)
	if err != nil || rowsAffected == 0 {
		tokenError(w, "invalid_grant", http.StatusBadRequest)
		return
	}
	// Generate access token

	if err := txn.Commit(); err != nil {
		tokenError(w, "server_error", http.StatusInternalServerError)
    return
	}
}

func tokenError(w http.ResponseWriter, errCode string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	fmt.Fprintf(w, `{"error":%q}`, errCode)
}
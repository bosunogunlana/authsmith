package oauth

import (
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"net/http"

	"github.com/bosunogunlana/authsmith/internal/models"
)

type TokenRequest struct {
}

func AuthenticateClient(r *http.Request, client models.OAuthClient) error {
	if client.ClientType == "public" {
		return nil
	}

	clientID, clientSecret, ok := r.BasicAuth()
	if !ok {
		return MissingClientCredentialsError
	}

	if clientID != client.ClientID {
		return ClientIDMisMatchError
	}

	if !client.ClientSecretDigest.Valid {
		return NoClientSecretConfiguredError
	}

	secretHash := sha256.Sum256([]byte(clientSecret))
	incomingDigest := hex.EncodeToString(secretHash[:])

	if subtle.ConstantTimeCompare([]byte(incomingDigest), []byte(client.ClientSecretDigest.String)) == 0 {
		return InvalidClientError
	}

	return nil
}

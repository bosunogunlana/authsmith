package oauth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
)

func GenerateToken() (string, string, error) {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", "", err
	}

	token := base64.RawURLEncoding.EncodeToString(bytes)
	hash := sha256.Sum256([]byte(token))
	digest := hex.EncodeToString(hash[:])
	return token, digest, nil
}

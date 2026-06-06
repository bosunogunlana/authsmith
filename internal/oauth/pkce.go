package oauth

import (
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
)

func ValidatePKCE(verifier, challenge string) error {
	hash := sha256.Sum256([]byte(verifier))
	computed := base64.RawURLEncoding.EncodeToString(hash[:])
	if subtle.ConstantTimeCompare([]byte(challenge), []byte(computed)) != 1 {
		return InvalidPKCEError
	}

	return nil
} 
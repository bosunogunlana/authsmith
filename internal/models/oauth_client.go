package models

import (
	"database/sql"
	"time"
)

type OAuthClient struct {
	ID                 string
	OwnerUserID        string
	Name               string
	ClientID           string
	ClientSecretDigest sql.NullString
	ClientType         string
	RedirectURIs       string
	AllowedScopes      string
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

package models

import (
	"database/sql"
	"time"
)

type AuthorizationCode struct {
	ID                  string
	CodeDigest          string
	UserID              string
	ClientID            string
	RedirectURI         string
	Scopes              string
	CodeChallenge       string
	CodeChallengeMethod string
	ExpiresAt           time.Time
	UsedAt              sql.NullTime
	CreatedAt           time.Time
}

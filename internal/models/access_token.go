package models

import (
	"database/sql"
	"time"
)

type AccessToken struct {
	ID          string
	TokenDigest string
	UserID      string
	ClientID    string
	Scopes      string
	CreatedAt   time.Time
	ExpiresAt   time.Time
	RevokedAt   sql.NullTime
}

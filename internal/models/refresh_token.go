package models

import (
	"database/sql"
	"time"
)

type RefreshToken struct {
	ID          string
	TokenDigest string
	UserID      string
	ClientID    string
	Scopes      string
	CreatedAt   time.Time
	ExpiresAt   time.Time
	RevokedAt   sql.NullTime
	UsedAt      sql.NullTime
}

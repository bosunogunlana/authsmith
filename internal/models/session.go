package models

import "time"

type Session struct {
	ID string
	UserID string
	TokenDigest string
	CreatedAt time.Time
	ExpiresAt time.Time
}

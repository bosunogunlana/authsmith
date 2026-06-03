package models

import "time"

type User struct {
	ID             string
	Name           string
	Email          string
	PasswordDigest string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

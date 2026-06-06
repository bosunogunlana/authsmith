package database

import (
	"database/sql"
	"time"

	"github.com/bosunogunlana/authsmith/internal/models"
)

func CreateAccessToken(db *sql.Tx, tokenDigest, userID, clientID, scopes string) (models.AccessToken, error) {
	stmt := `
		INSERT INTO oauth_access_tokens (token_digest, user_id, client_id, scopes, expires_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, token_digest, user_id, client_id, scopes, expires_at, revoked_at, created_at
	`

	var acessToken models.AccessToken
	
	expiresAt := time.Now().Add(1 * time.Hour)
	row := db.QueryRow(stmt, tokenDigest, userID, clientID, scopes, expiresAt)
	err := row.Scan(
		&acessToken.ID,
		&acessToken.TokenDigest,
		&acessToken.UserID,
		&acessToken.ClientID,
		&acessToken.Scopes,
		&acessToken.ExpiresAt,
		&acessToken.RevokedAt,
		&acessToken.CreatedAt,
	)
	if err != nil {
		return models.AccessToken{}, err
	}

	return acessToken, nil
}

func GetAccessTokenBydigest(db *sql.DB, tokenDigest string) (models.AccessToken, error) {
	stmt := `
		SELECT id, token_digest, user_id, client_id, scopes, expires_at, revoked_at, created_at
		FROM oauth_access_tokens
		WHERE expires_at > $1 AND revoked_at is NULL
	`

	var acessToken models.AccessToken

	row := db.QueryRow(stmt, time.Now())
	err := row.Scan(
	  &acessToken.ID,
		&acessToken.TokenDigest,
		&acessToken.UserID,
		&acessToken.ClientID,
		&acessToken.Scopes,
		&acessToken.ExpiresAt,
		&acessToken.RevokedAt,
		&acessToken.CreatedAt,
	)
	if err != nil {
		return models.AccessToken{}, err
	}

	return acessToken, nil
}
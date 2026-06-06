package database

import (
	"database/sql"
	"time"

	"github.com/bosunogunlana/authsmith/internal/models"
)

func CreateRefreshToken(db *sql.Tx, tokenDigest, userID, clientID, scopes string) (models.RefreshToken, error) {
	stmt := `
		INSERT INTO oauth_refresh_tokens (token_digest, user_id, client_id, scopes, expires_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, token_digest, user_id, client_id, scopes, expires_at, used_at, revoked_at, created_at
	`

	var refreshToken models.RefreshToken

	expiresAt := time.Now().Add(30 * 24 * time.Hour)
	row := db.QueryRow(stmt, tokenDigest, userID, clientID, scopes, expiresAt)
	err := row.Scan(
		&refreshToken.ID,
		&refreshToken.TokenDigest,
		&refreshToken.UserID,
		&refreshToken.ClientID,
		&refreshToken.Scopes,
		&refreshToken.ExpiresAt,
		&refreshToken.UsedAt,
		&refreshToken.RevokedAt,
		&refreshToken.CreatedAt,
	)
	if err != nil {
		return models.RefreshToken{}, err
	}

	return refreshToken, nil
}

func GetRefreshTokenBydigest(db *sql.DB, tokenDigest string) (models.RefreshToken, error) {
	stmt := `
		SELECT id, token_digest, user_id, client_id, scopes, expires_at, used_at, revoked_at, created_at
		FROM oauth_refresh_tokens
		WHERE expires_at > $1 AND revoked_at IS NULL AND used_at IS NULL
	`

	var refreshToken models.RefreshToken

	row := db.QueryRow(stmt, time.Now())
	err := row.Scan(
	  &refreshToken.ID,
		&refreshToken.TokenDigest,
		&refreshToken.UserID,
		&refreshToken.ClientID,
		&refreshToken.Scopes,
		&refreshToken.ExpiresAt,
		&refreshToken.UsedAt,
		&refreshToken.RevokedAt,
		&refreshToken.CreatedAt,
	)
	if err != nil {
		return models.RefreshToken{}, err
	}

	return refreshToken, nil
}

func RevokeRefreshToken(db *sql.DB, id string) (bool, error) {
	stmt := `UPDATE oauth_refresh_tokens SET revoked_at = $2 where id = $1`

	result, err := db.Exec(stmt, id, time.Now())
	if err != nil {
		return false, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		return false, err
	}

	return true, nil
}
package database

import (
	"database/sql"
	"errors"
	"time"

	"github.com/bosunogunlana/authsmith/internal/models"
)

func CreateAuthorizationCode(db *sql.DB, userID, clientID, redirectURI, scopes, codeDigest, codeChallenge, codeChallengeMethod string) (models.AuthorizationCode, error) {
	stmt := `
		INSERT INTO oauth_authorization_codes (user_id, client_id, redirect_uri, scopes, code_digest, code_challenge, code_challenge_method, expires_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, user_id, client_id, redirect_uri, scopes, code_digest, code_challenge, code_challenge_method, expires_at, used_at, created_at
	`

	var authorizationCode models.AuthorizationCode

	expiresAt := time.Now().Add(10 * time.Minute)

	row := db.QueryRow(stmt, userID, clientID, redirectURI, scopes, codeDigest, codeChallenge, codeChallengeMethod, expiresAt)
	err := row.Scan(
		&authorizationCode.ID,
		&authorizationCode.UserID,
		&authorizationCode.ClientID,
		&authorizationCode.RedirectURI,
		&authorizationCode.Scopes,
		&authorizationCode.CodeDigest,
		&authorizationCode.CodeChallenge,
		&authorizationCode.CodeChallengeMethod,
		&authorizationCode.ExpiresAt,
		&authorizationCode.UsedAt,
		&authorizationCode.CreatedAt,
	)
	if err != nil {
		return models.AuthorizationCode{}, err
	}

	return authorizationCode, nil
}

func GetAuthorizationCodeByDigest(db *sql.DB, digest string) (models.AuthorizationCode, error) {
	stmt := `
		SELECT id, user_id, client_id, redirect_uri, scopes, code_digest, code_challenge, code_challenge_method, expires_at, used_at, created_at
		FROM oauth_authorization_codes
		WHERE code_digest = $1 AND used_at is null AND expires_at > $2
	`

	var authorizationCode models.AuthorizationCode

	row := db.QueryRow(stmt, digest, time.Now())

	err := row.Scan(
		&authorizationCode.ID,
		&authorizationCode.UserID,
		&authorizationCode.ClientID,
		&authorizationCode.RedirectURI,
		&authorizationCode.Scopes,
		&authorizationCode.CodeDigest,
		&authorizationCode.CodeChallenge,
		&authorizationCode.CodeChallengeMethod,
		&authorizationCode.ExpiresAt,
		&authorizationCode.UsedAt,
		&authorizationCode.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.AuthorizationCode{}, nil
		}

		return models.AuthorizationCode{}, err
	}

	return authorizationCode, nil
}

func MarkAuthorizationCodeUsed(db *sql.DB, id string) error {
	stmt := `UPDATE oauth_authorization_codes SET used_at = $2 WHERE id = $1`

	_, err := db.Exec(stmt, id, time.Now())
	return err
}

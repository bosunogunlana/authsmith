package database

import (
	"database/sql"
	"errors"
	"time"

	"github.com/bosunogunlana/authsmith/internal/models"
)

func CreateSession(db *sql.DB, userID, tokenDigest string) (models.Session, error) {
	stmt := `
		INSERT INTO sessions (user_id, token_digest, expires_at)
  	VALUES ($1, $2, $3)
    RETURNING id, user_id, token_digest, created_at, expires_at 
	`

	var session models.Session

	expresAt := time.Now().Add(24 * time.Hour)
	row := db.QueryRow(stmt, userID, tokenDigest, expresAt)
	err := row.Scan(
		&session.ID,
		&session.UserID,
		&session.TokenDigest,
		&session.CreatedAt,
		&session.ExpiresAt,
	)

	if err != nil {
		return models.Session{}, err
	}

	return session, nil
}

func DeleteSessionByTokenDigest(db *sql.DB, digest string) error {
	_, err := db.Exec(`DELETE FROM sessions WHERE token_digest = $1`, digest)
	return err
}

func GetSessionByTokenDigest(db *sql.DB, digest string) (models.Session, error) {
	stmt := `
		SELECT id, user_id, token_digest, created_at, expires_at
		FROM sessions
		WHERE token_digest = $1 AND expires_at > $2
	`

	var session models.Session

	row := db.QueryRow(stmt, digest, time.Now())

	err := row.Scan(
		&session.ID,
		&session.UserID,
		&session.TokenDigest,
		&session.CreatedAt,
		&session.ExpiresAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Session{}, nil
		}

		return models.Session{}, err
	}

	return session, nil
}
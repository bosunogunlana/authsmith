package database

import (
	"database/sql"
	"errors"

	"github.com/bosunogunlana/authsmith/internal/models"
)

func CreateUser(db *sql.DB, email, passwordDigest, name string) (models.User, error) {
	stmt := `
		INSERT INTO users (email, password_digest, name)
  	VALUES ($1, $2, $3)
    RETURNING id, email, password_digest, name, created_at, updated_at
	`

	var user models.User

	row := db.QueryRow(stmt, email, passwordDigest, name)
	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.PasswordDigest,
		&user.Name,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func GetUserByID(db *sql.DB, id string) (models.User, error) {
	stmt := `
		SELECT id, email, password_digest, name, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	var user models.User
	row := db.QueryRow(stmt, id)
	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.PasswordDigest,
		&user.Name,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, nil
		}

		return models.User{}, err
	}

	return user, nil
}

func GetUserByEmail(db *sql.DB, email string) (models.User, error) {
	stmt := `
		SELECT id, email, password_digest, name, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	var user models.User
	row := db.QueryRow(stmt, email)
	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.PasswordDigest,
		&user.Name,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, nil
		}

		return models.User{}, err
	}

	return user, nil
}
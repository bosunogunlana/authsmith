package database

import (
	"database/sql"
	"errors"

	"github.com/bosunogunlana/authsmith/internal/models"
)

func CreateOAuthClient(db *sql.DB, ownerUserID, name, clientID, clientSecretDigest, clientType, redirectURIs, allowedScopes string) (models.OAuthClient, error) {
	stmt := `
		INSERT INTO oauth_clients (owner_user_id, name, client_id, client_secret_digest, client_type, redirect_uris, allowed_scopes)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, owner_user_id, name, client_id, client_secret_digest, client_type, redirect_uris, allowed_scopes, created_at, updated_at
	`

	var client models.OAuthClient
	row := db.QueryRow(stmt, ownerUserID, name, clientID, clientSecretDigest, clientType, redirectURIs, allowedScopes)
	err := row.Scan(
		&client.ID,
		&client.OwnerUserID,
		&client.Name,
		&client.ClientID,
		&client.ClientSecretDigest,
		&client.ClientType,
		&client.RedirectURIs,
		&client.AllowedScopes,
		&client.CreatedAt,
		&client.UpdatedAt,
	)
	if err != nil {
		return models.OAuthClient{}, err
	}

	return client, nil
}

func GetOAuthClientsByOwnerID(db *sql.DB, ownerUserID string) ([]models.OAuthClient, error) {
	stmt := `
		SELECT id, owner_user_id, name, client_id, client_secret_digest, client_type, redirect_uris, allowed_scopes, created_at, updated_at
		FROM oauth_clients
		WHERE owner_user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := db.Query(stmt, ownerUserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var clients []models.OAuthClient
	for rows.Next() {
		var client models.OAuthClient
		if err := rows.Scan(
			&client.ID,
			&client.OwnerUserID,
			&client.Name,
			&client.ClientID,
			&client.ClientSecretDigest,
			&client.ClientType,
			&client.RedirectURIs,
			&client.AllowedScopes,
			&client.CreatedAt,
			&client.UpdatedAt,
		); err != nil {
			return nil, err
		}
		clients = append(clients, client)
	}

	return clients, rows.Err()
}

func GetOAuthClientByClientID(db *sql.DB, clientID string) (models.OAuthClient, error) {
	stmt := `
		SELECT id, owner_user_id, name, client_id, client_secret_digest, client_type, redirect_uris, allowed_scopes, created_at, updated_at
		FROM oauth_clients
		WHERE client_id = $1
	`

	var client models.OAuthClient
	row := db.QueryRow(stmt, clientID)
	err := row.Scan(
		&client.ID,
		&client.OwnerUserID,
		&client.Name,
		&client.ClientID,
		&client.ClientSecretDigest,
		&client.ClientType,
		&client.RedirectURIs,
		&client.AllowedScopes,
		&client.CreatedAt,
		&client.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.OAuthClient{}, nil
		}
		return models.OAuthClient{}, err
	}

	return client, nil
}

CREATE TABLE IF NOT EXISTS oauth_clients (
    id                   UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    owner_user_id        UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name                 TEXT NOT NULL,
    client_id            TEXT NOT NULL UNIQUE,
    client_secret_digest TEXT,
    client_type          TEXT NOT NULL CHECK (client_type IN ('confidential', 'public')),
    redirect_uris        TEXT NOT NULL,
    allowed_scopes       TEXT NOT NULL,
    created_at           TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at           TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

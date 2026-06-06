CREATE TABLE IF NOT EXISTS oauth_access_tokens (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    token_digest TEXT NOT NULL UNIQUE,
    user_id      UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    client_id    TEXT NOT NULL,
    scopes       TEXT NOT NULL,
    expires_at   TIMESTAMPTZ NOT NULL,
    revoked_at   TIMESTAMPTZ,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS oauth_refresh_tokens (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    token_digest TEXT NOT NULL UNIQUE,
    user_id      UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    client_id    TEXT NOT NULL,
    scopes       TEXT NOT NULL,
    expires_at   TIMESTAMPTZ NOT NULL,
    used_at      TIMESTAMPTZ,
    revoked_at   TIMESTAMPTZ,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

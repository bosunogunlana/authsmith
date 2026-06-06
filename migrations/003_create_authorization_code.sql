CREATE TABLE IF NOT EXISTS oauth_authorization_codes (
  id UUID                PRIMARY KEY DEFAULT gen_random_uuid(),
  code_digest            TEXT NOT NULL UNIQUE,
  user_id                UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  client_id              TEXT NOT NULL,
  redirect_uri           TEXT NOT NULL,
  scopes                 TEXT NOT NULL,
  code_challenge         TEXT NOT NULL,
  code_challenge_method  TEXT NOT NULL,
  expires_at             TIMESTAMPTZ NOT NULL,
  used_at                TIMESTAMPTZ,
  created_at             TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE TABLE auth(
  id BIGSERIAL PRIMARY KEY,

  email VARCHAR NOT NULL,
  password VARCHAR,

  email_verified_at TIMESTAMPTZ,
  email_changed_at TIMESTAMPTZ,
  password_changed_at TIMESTAMPTZ,

  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  deleted_at TIMESTAMPTZ
);

-- Enforce uniqueness of email for active (not soft-deleted) records
CREATE UNIQUE INDEX idx_auth_unique_email ON auth(email) WHERE deleted_at IS NULL;

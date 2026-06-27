CREATE SCHEMA authservice;

CREATE TABLE IF NOT EXISTS authservice.users (
    id UUID PRIMARY KEY,
    username VARCHAR(100) NOT NULL CHECK(char_length(username) BETWEEN 1 AND 50),
    phone_number VARCHAR(15) CHECK(
        phone_number ~ '^\+[0-9]+$'
        AND
        char_length(phone_number) BETWEEN 10 AND 15
    ),
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_users_username ON authservice.users(username);


CREATE TABLE IF NOT EXISTS authservice.refresh_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES authservice.users(id) ON DELETE CASCADE,
    token VARCHAR(255) UNIQUE NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    revoked BOOLEAN NOT NULL DEFAULT FALSE
);
-- Create indexes for fast lookups
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_token ON authservice.refresh_tokens(token);
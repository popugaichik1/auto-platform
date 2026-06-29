CREATE SCHEMA IF NOT EXISTS userservice;

CREATE TABLE IF NOT EXISTS userservice.users (
    id  UUID PRIMARY KEY,
    version INTEGER NOT NULL DEFAULT 1,
    username VARCHAR(125) NOT NULL,
    phone_number VARCHAR(15) CHECK(
        phone_number ~ '^\+[0-9]+$'
        AND
        char_length(phone_number) BETWEEN 10 AND 15
    ),
    bio VARCHAR(255),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_userservice_users_username ON userservice.users (username);
CREATE INDEX IF NOT EXISTS idx_userservice_users_phone_number ON userservice.users (phone_number);
CREATE INDEX IF NOT EXISTS idx_userservice_users_id ON userservice.users (id);

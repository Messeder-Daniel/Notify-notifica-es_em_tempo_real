CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS notifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title TEXT NOT NULL,
    message TEXT NOT NULL,
    is_read BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    read_at TIMESTAMPTZ
);

INSERT INTO users (name, email, password_hash)
VALUES
    ('Alice', 'alice@example.com', '$2a$10$inbLu0dHbOAO48GKLcGdEugGNfmGOwFQcapN7LMPwCIzDqJ7zTC5m'),
    ('Bob', 'bob@example.com', '$2a$10$inbLu0dHbOAO48GKLcGdEugGNfmGOwFQcapN7LMPwCIzDqJ7zTC5m'),
    ('Daniel', 'daniel@example.com', '$2a$10$inbLu0dHbOAO48GKLcGdEugGNfmGOwFQcapN7LMPwCIzDqJ7zTC5m')
ON CONFLICT (email)
DO UPDATE SET
    name = EXCLUDED.name,
    password_hash = EXCLUDED.password_hash;

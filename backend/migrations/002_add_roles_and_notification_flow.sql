CREATE EXTENSION IF NOT EXISTS "pgcrypto";

ALTER TABLE users
ADD COLUMN IF NOT EXISTS role TEXT NOT NULL DEFAULT 'user';

UPDATE users
SET role = 'user'
WHERE role IS NULL OR role NOT IN ('admin', 'user');

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM pg_constraint
        WHERE conname = 'users_role_check'
    ) THEN
        ALTER TABLE users
        ADD CONSTRAINT users_role_check CHECK (role IN ('admin', 'user'));
    END IF;
END $$;

INSERT INTO users (name, email, password_hash, role)
VALUES
    ('Daniel Messeder', 'messederdaniel@outlook.com', crypt('Teste@2026', gen_salt('bf')), 'admin'),
    ('Daniel Barreto', 'barretodaniel11971@hotmail.com', crypt('Teste@2026', gen_salt('bf')), 'user')
ON CONFLICT (email)
DO UPDATE SET
    name = EXCLUDED.name,
    role = EXCLUDED.role;

DO $$
BEGIN
    IF EXISTS (
        SELECT 1
        FROM information_schema.columns
        WHERE table_name = 'notifications'
          AND column_name = 'user_id'
    )
    AND NOT EXISTS (
        SELECT 1
        FROM information_schema.columns
        WHERE table_name = 'notifications'
          AND column_name = 'recipient_id'
    ) THEN
        ALTER TABLE notifications
        RENAME COLUMN user_id TO recipient_id;
    END IF;
END $$;

ALTER TABLE notifications
ADD COLUMN IF NOT EXISTS sender_id UUID;

ALTER TABLE notifications
ADD COLUMN IF NOT EXISTS parent_id UUID;

ALTER TABLE notifications
ADD COLUMN IF NOT EXISTS is_completed BOOLEAN NOT NULL DEFAULT FALSE;

ALTER TABLE notifications
ADD COLUMN IF NOT EXISTS completed_at TIMESTAMPTZ;

UPDATE notifications
SET sender_id = COALESCE(
    sender_id,
    (
        SELECT id
        FROM users
        WHERE email = 'messederdaniel@outlook.com'
        LIMIT 1
    ),
    recipient_id
)
WHERE sender_id IS NULL;

ALTER TABLE notifications
ALTER COLUMN sender_id SET NOT NULL;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM pg_constraint
        WHERE conname = 'notifications_sender_id_fkey'
    ) THEN
        ALTER TABLE notifications
        ADD CONSTRAINT notifications_sender_id_fkey
        FOREIGN KEY (sender_id) REFERENCES users(id) ON DELETE CASCADE;
    END IF;
END $$;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM pg_constraint
        WHERE conname = 'notifications_recipient_id_fkey'
    ) THEN
        ALTER TABLE notifications
        ADD CONSTRAINT notifications_recipient_id_fkey
        FOREIGN KEY (recipient_id) REFERENCES users(id) ON DELETE CASCADE;
    END IF;
END $$;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM pg_constraint
        WHERE conname = 'notifications_parent_id_fkey'
    ) THEN
        ALTER TABLE notifications
        ADD CONSTRAINT notifications_parent_id_fkey
        FOREIGN KEY (parent_id) REFERENCES notifications(id) ON DELETE CASCADE;
    END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_role ON users(role);
CREATE INDEX IF NOT EXISTS idx_notifications_sender_id ON notifications(sender_id);
CREATE INDEX IF NOT EXISTS idx_notifications_recipient_id ON notifications(recipient_id);
CREATE INDEX IF NOT EXISTS idx_notifications_parent_id ON notifications(parent_id);
CREATE INDEX IF NOT EXISTS idx_notifications_created_at ON notifications(created_at);

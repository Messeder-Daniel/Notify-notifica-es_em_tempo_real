ALTER TABLE users
ADD COLUMN IF NOT EXISTS role TEXT NOT NULL DEFAULT 'user';

UPDATE users
SET role = 'admin'
WHERE email = 'daniel@example.com';

UPDATE users
SET role = 'user'
WHERE email <> 'daniel@example.com';

DO $$
BEGIN
    IF EXISTS (
        SELECT 1
        FROM information_schema.columns
        WHERE table_name = 'notifications'
        AND column_name = 'user_id'
    ) AND NOT EXISTS (
        SELECT 1
        FROM information_schema.columns
        WHERE table_name = 'notifications'
        AND column_name = 'recipient_id'
    ) THEN
        ALTER TABLE notifications RENAME COLUMN user_id TO recipient_id;
    END IF;
END $$;

ALTER TABLE notifications
ADD COLUMN IF NOT EXISTS sender_id UUID REFERENCES users(id) ON DELETE SET NULL,
ADD COLUMN IF NOT EXISTS parent_id UUID REFERENCES notifications(id) ON DELETE CASCADE,
ADD COLUMN IF NOT EXISTS is_completed BOOLEAN NOT NULL DEFAULT FALSE,
ADD COLUMN IF NOT EXISTS completed_at TIMESTAMPTZ;

UPDATE notifications
SET sender_id = recipient_id
WHERE sender_id IS NULL;

ALTER TABLE notifications
ALTER COLUMN sender_id SET NOT NULL;

CREATE INDEX IF NOT EXISTS idx_notifications_recipient_id
ON notifications(recipient_id);

CREATE INDEX IF NOT EXISTS idx_notifications_sender_id
ON notifications(sender_id);

CREATE INDEX IF NOT EXISTS idx_notifications_parent_id
ON notifications(parent_id);

ALTER TABLE users
ADD COLUMN username_changed_at TIMESTAMP;

UPDATE users
SET username_changed_at = created_at
WHERE username_changed_at IS NULL;

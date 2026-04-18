DROP INDEX IF EXISTS idx_sessions_creator;

ALTER TABLE game_sessions
    DROP CONSTRAINT IF EXISTS game_sessions_creator_id_fkey,
    DROP COLUMN IF EXISTS creator_id;

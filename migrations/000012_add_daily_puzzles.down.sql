-- Drop indexes
DROP INDEX IF EXISTS idx_game_sessions_is_daily;
DROP INDEX IF EXISTS idx_daily_puzzles_date;

-- Drop user_daily_stats table
DROP TABLE IF EXISTS user_daily_stats;

-- Remove daily game fields from game_sessions
ALTER TABLE game_sessions
    DROP CONSTRAINT IF EXISTS fk_daily_puzzle,
    DROP COLUMN IF EXISTS daily_puzzle_id,
    DROP COLUMN IF EXISTS is_daily;

-- Drop daily_puzzles table
DROP TABLE IF EXISTS daily_puzzles;

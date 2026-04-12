DROP INDEX IF EXISTS idx_results_user_score;
DROP INDEX IF EXISTS idx_results_user_id;

ALTER TABLE game_results
    DROP CONSTRAINT IF EXISTS game_results_user_id_fkey,
    DROP COLUMN IF EXISTS user_id;

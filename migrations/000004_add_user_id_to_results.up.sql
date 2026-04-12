ALTER TABLE game_results
    ADD COLUMN user_id UUID NULL,
    ADD CONSTRAINT game_results_user_id_fkey
        FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL;

CREATE INDEX idx_results_user_id ON game_results(user_id);
CREATE INDEX idx_results_user_score ON game_results(user_id, score DESC);

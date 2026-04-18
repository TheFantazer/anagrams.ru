ALTER TABLE game_sessions
    ADD COLUMN creator_id UUID NULL,
    ADD CONSTRAINT game_sessions_creator_id_fkey
        FOREIGN KEY (creator_id) REFERENCES users(id) ON DELETE SET NULL;

CREATE INDEX idx_sessions_creator ON game_sessions(creator_id);

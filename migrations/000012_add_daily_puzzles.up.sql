CREATE TABLE IF NOT EXISTS daily_puzzles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    puzzle_date DATE NOT NULL UNIQUE,
    letters VARCHAR(20) NOT NULL,
    language VARCHAR(10) NOT NULL DEFAULT 'en',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

ALTER TABLE game_sessions
    ADD COLUMN is_daily BOOLEAN DEFAULT FALSE,
    ADD COLUMN daily_puzzle_id UUID,
    ADD CONSTRAINT fk_daily_puzzle FOREIGN KEY (daily_puzzle_id) REFERENCES daily_puzzles(id) ON DELETE SET NULL;

CREATE TABLE IF NOT EXISTS user_daily_stats (
    user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    current_streak INTEGER DEFAULT 0,
    longest_streak INTEGER DEFAULT 0,
    last_played_date DATE,
    total_daily_games INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_daily_puzzles_date ON daily_puzzles(puzzle_date);

CREATE INDEX idx_game_sessions_is_daily ON game_sessions(is_daily);

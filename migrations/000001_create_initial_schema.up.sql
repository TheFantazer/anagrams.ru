CREATE TABLE game_sessions
(
    id           UUID PRIMARY KEY     DEFAULT gen_random_uuid(),
    letters      VARCHAR(12) NOT NULL,
    language     VARCHAR(5)  NOT NULL DEFAULT 'ru',
    time_limit   SMALLINT    NOT NULL DEFAULT 60,
    letter_count SMALLINT    NOT NULL DEFAULT 7,
    valid_words  JSONB       NOT NULL DEFAULT '[]',
    max_score    INT         NOT NULL DEFAULT 0,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_sessions_created_at ON game_sessions(created_at);

CREATE TABLE game_results
(
    id                 UUID PRIMARY KEY     DEFAULT gen_random_uuid(),
    session_id         UUID        NOT NULL REFERENCES game_sessions (id) ON DELETE CASCADE,
    player_name        VARCHAR(64) NOT NULL DEFAULT 'Anonymous',
    player_fingerprint VARCHAR(64) NOT NULL,
    found_words        JSONB       NOT NULL DEFAULT '[]',
    word_count         SMALLINT    NOT NULL DEFAULT 0,
    score              INT                  DEFAULT 0,
    duration_ms        INT         NOT NULL,
    played_at          TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT uq_session_player UNIQUE (session_id, player_fingerprint)
);
CREATE INDEX idx_results_session_score ON game_results (session_id, score DESC);

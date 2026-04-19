CREATE TABLE session_invites
(
    id         UUID PRIMARY KEY     DEFAULT gen_random_uuid(),
    session_id UUID        NOT NULL REFERENCES game_sessions (id) ON DELETE CASCADE,
    user_id    UUID        NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT session_invite_unique UNIQUE (session_id, user_id)
);

CREATE INDEX idx_session_invites_session ON session_invites (session_id);
CREATE INDEX idx_session_invites_user ON session_invites (user_id);
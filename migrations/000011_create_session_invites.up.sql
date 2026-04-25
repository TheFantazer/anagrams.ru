DROP TABLE IF EXISTS session_invites;

CREATE TABLE session_invites (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    session_id UUID NOT NULL REFERENCES game_sessions(id) ON DELETE CASCADE,
    from_user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    to_user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE(session_id, to_user_id)
);

CREATE INDEX idx_session_invites_to_user ON session_invites(to_user_id);
CREATE INDEX idx_session_invites_session ON session_invites(session_id);

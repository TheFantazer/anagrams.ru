-- Таблица участников сессии (поддержка будущего 1vN режима)
CREATE TABLE IF NOT EXISTS session_participants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    session_id UUID NOT NULL REFERENCES game_sessions(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role VARCHAR(20) NOT NULL CHECK (role IN ('creator', 'opponent')),
    joined_at TIMESTAMP NOT NULL DEFAULT NOW(),
    started_at TIMESTAMP,
    UNIQUE(session_id, user_id)
);

CREATE INDEX idx_session_participants_session ON session_participants(session_id);
CREATE INDEX idx_session_participants_user ON session_participants(user_id);

-- Расширение game_sessions для поддержки разных режимов
ALTER TABLE game_sessions
ADD COLUMN IF NOT EXISTS max_opponents INT DEFAULT 1 CHECK (max_opponents >= 1),
ADD COLUMN IF NOT EXISTS invite_mode VARCHAR(20) DEFAULT 'link' CHECK (invite_mode IN ('link', 'friend'));

-- Миграция существующих данных: создать записи participant для всех создателей
INSERT INTO session_participants (session_id, user_id, role, joined_at, started_at)
SELECT id, creator_id, 'creator', created_at, created_at
FROM game_sessions
WHERE creator_id IS NOT NULL
ON CONFLICT (session_id, user_id) DO NOTHING;

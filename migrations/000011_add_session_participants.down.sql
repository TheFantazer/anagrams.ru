-- Откатить изменения в game_sessions
ALTER TABLE game_sessions
DROP COLUMN IF EXISTS invite_mode,
DROP COLUMN IF EXISTS max_opponents;

-- Удалить таблицу участников
DROP TABLE IF EXISTS session_participants;

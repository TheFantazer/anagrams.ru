CREATE TABLE IF NOT EXISTS users
(
    id             UUID PRIMARY KEY            DEFAULT gen_random_uuid(),
    username       VARCHAR(30) UNIQUE NOT NULL,
    email          VARCHAR(50) UNIQUE NULL,
    password       VARCHAR(255),
    oauth_provider VARCHAR(20),
    oauth_id       VARCHAR(100),
    created_at     TIMESTAMPTZ        NOT NULL DEFAULT now(),
    updated_at     TIMESTAMPTZ        NOT NULL DEFAULT now()
)
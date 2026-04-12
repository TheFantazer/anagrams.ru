ALTER TABLE users
    ADD COLUMN default_letter_count SMALLINT NOT NULL DEFAULT 7 CHECK (default_letter_count BETWEEN 6 AND 10),
    ADD COLUMN default_language VARCHAR(5) NOT NULL DEFAULT 'ru' CHECK (default_language IN ('ru', 'en')),
    ADD COLUMN default_time_limit SMALLINT NOT NULL DEFAULT 60 CHECK (default_time_limit > 0);

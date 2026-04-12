ALTER TABLE users
    DROP COLUMN IF EXISTS default_letter_count,
    DROP COLUMN IF EXISTS default_language,
    DROP COLUMN IF EXISTS default_time_limit;

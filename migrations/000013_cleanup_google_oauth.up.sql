-- Remove Google OAuth data from users table
UPDATE users
SET oauth_provider = NULL, oauth_id = NULL
WHERE oauth_provider = 'google';

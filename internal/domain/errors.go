package domain

import "errors"

// Session errors
var (
	ErrInvalidLetters      = errors.New("letters cannot be empty")
	ErrLetterCountMismatch = errors.New("letters count does not match letter_count")
	ErrUnsupportedLanguage = errors.New("unsupported language")
	ErrInvalidTimeLimit    = errors.New("time limit must be between 1 and 600 seconds")
	ErrNoValidWords        = errors.New("no valid words provided")
	ErrSessionExpired      = errors.New("session has expired")
)

// Result errors
var (
	ErrInvalidSessionID   = errors.New("invalid session ID")
	ErrMissingFingerprint = errors.New("player fingerprint is required")
	ErrInvalidDuration    = errors.New("duration must be positive")
	ErrInvalidWord        = errors.New("word is not valid for this session")
	ErrDuplicateResult    = errors.New("result already submitted for this session and player")
)

// Word errors
var (
	ErrWordTooShort = errors.New("word is too short")
	ErrWordTooLong  = errors.New("word is too long")
)

// User errors
var (
	ErrInvalidUsername    = errors.New("username must be 3-30 characters")
	ErrPasswordTooShort   = errors.New("password must be at least 6 characters")
	ErrUsernameTaken      = errors.New("username is already taken")
	ErrEmailTaken         = errors.New("email is already taken")
	ErrInvalidCredentials = errors.New("invalid username or password")
	ErrUserNotFound       = errors.New("user not found")
)

// Friend errors
var (
	ErrFriendRequestAlreadyExists = errors.New("friend request already exists")
	ErrFriendRequestNotPending    = errors.New("friend request is not pending")
	ErrAlreadyFriends             = errors.New("users are already friends")
	ErrNotFriends                 = errors.New("users are not friends")
)

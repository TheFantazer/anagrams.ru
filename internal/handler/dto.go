package handler

import (
	"time"

	"github.com/google/uuid"
)

type CreateSessionRequest struct {
	Language    string `json:"language"`
	LetterCount int    `json:"letter_count"`
	TimeLimit   int    `json:"time_limit"`
}

type SessionResponse struct {
	ID              uuid.UUID        `json:"id"`
	Letters         string           `json:"letters"`
	Language        string           `json:"language"`
	TimeLimit       int              `json:"time_limit"`
	LetterCount     int              `json:"letter_count"`
	MaxScore        int              `json:"max_score"`
	ValidWords      []string         `json:"valid_words"`
	CreatedAt       time.Time        `json:"created_at"`
	CreatorID       *uuid.UUID       `json:"creator_id,omitempty"`
	CreatorUsername *string          `json:"creator_username,omitempty"`
	Results         []ResultResponse `json:"results,omitempty"`
	Type            string           `json:"type,omitempty"` // "created" or "participated"
}

type PaginatedSessionsResponse struct {
	Sessions   []SessionResponse `json:"sessions"`
	Total      int               `json:"total"`
	Page       int               `json:"page"`
	PerPage    int               `json:"per_page"`
	TotalPages int               `json:"total_pages"`
}

type SubmitResultRequest struct {
	UserID      *string  `json:"user_id,omitempty"`
	PlayerName  string   `json:"player_name"`
	Fingerprint string   `json:"fingerprint"`
	FoundWords  []string `json:"found_words"`
	DurationMs  int      `json:"duration_ms"`
}

type ResultResponse struct {
	ID         uuid.UUID  `json:"id"`
	SessionID  uuid.UUID  `json:"session_id"`
	UserID     *uuid.UUID `json:"user_id,omitempty"`
	PlayerName string     `json:"player_name"`
	WordCount  int        `json:"word_count"`
	Score      int        `json:"score"`
	DurationMs int        `json:"duration_ms"`
	PlayedAt   time.Time  `json:"played_at"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserResponse struct {
	ID                 uuid.UUID `json:"id"`
	Username           string    `json:"username"`
	Email              *string   `json:"email,omitempty"`
	DefaultLetterCount int       `json:"default_letter_count"`
	DefaultLanguage    string    `json:"default_language"`
	DefaultTimeLimit   int       `json:"default_time_limit"`
	CreatedAt          time.Time `json:"created_at"`
}

type UpdateSettingsRequest struct {
	LetterCount int    `json:"letter_count"`
	Language    string `json:"language"`
	TimeLimit   int    `json:"time_limit"`
}

type UserStatsResponse struct {
	GamesPlayed  int     `json:"games_played"`
	BestScore    int     `json:"best_score"`
	LongestWord  string  `json:"longest_word"`
	TotalWords   int     `json:"total_words"`
	AverageScore float64 `json:"average_score"`
}

type SendFriendRequestRequest struct {
	ToUserID uuid.UUID `json:"to_user_id"`
}

type FriendRequestResponse struct {
	ID         uuid.UUID `json:"id"`
	FromUserID uuid.UUID `json:"from_user_id"`
	ToUserID   uuid.UUID `json:"to_user_id"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type DailyPuzzleResponse struct {
	ID         uuid.UUID `json:"id"`
	PuzzleDate time.Time `json:"puzzle_date"`
	Letters    string    `json:"letters"`
	Language   string    `json:"language"`
	CreatedAt  time.Time `json:"created_at"`
}

type DailySessionResponse struct {
	Session    SessionResponse     `json:"session"`
	HasPlayed  bool                `json:"has_played"`
	DailyStats *DailyStatsResponse `json:"daily_stats,omitempty"`
}

type DailyStatsResponse struct {
	UserID          uuid.UUID  `json:"user_id"`
	CurrentStreak   int        `json:"current_streak"`
	LongestStreak   int        `json:"longest_streak"`
	LastPlayedDate  *time.Time `json:"last_played_date,omitempty"`
	TotalDailyGames int        `json:"total_daily_games"`
}

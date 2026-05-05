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
	FoundWords []string   `json:"found_words,omitempty"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

type RegisterRequest struct {
	Username              string `json:"username"`
	Email                 string `json:"email"`
	Password              string `json:"password"`
	AcceptedPrivacyPolicy bool   `json:"accepted_privacy_policy"`
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

type UpdateUsernameRequest struct {
	NewUsername string `json:"new_username"`
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
	ID           uuid.UUID `json:"id"`
	FromUserID   uuid.UUID `json:"from_user_id"`
	FromUsername string    `json:"from_username"`
	ToUserID     uuid.UUID `json:"to_user_id"`
	ToUsername   string    `json:"to_username"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
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

type SessionInviteResponse struct {
	ID              uuid.UUID `json:"id"`
	SessionID       uuid.UUID `json:"session_id"`
	FromUserID      uuid.UUID `json:"from_user_id"`
	ToUserID        uuid.UUID `json:"to_user_id"`
	InvitedUsername string    `json:"invited_username"`
	Status          string    `json:"status"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// TelegramUpdate represents an incoming update from Telegram Bot API
type TelegramUpdate struct {
	UpdateID           int64               `json:"update_id"`
	Message            *TelegramMessage    `json:"message,omitempty"`
	InlineQuery        *InlineQuery        `json:"inline_query,omitempty"`
	ChosenInlineResult *ChosenInlineResult `json:"chosen_inline_result,omitempty"`
}

// TelegramMessage represents a message from Telegram
type TelegramMessage struct {
	MessageID int64         `json:"message_id"`
	From      *TelegramUser `json:"from,omitempty"`
	Chat      *TelegramChat `json:"chat,omitempty"`
	Text      string        `json:"text,omitempty"`
}

// TelegramUser represents a Telegram user
type TelegramUser struct {
	ID        int64  `json:"id"`
	IsBot     bool   `json:"is_bot"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name,omitempty"`
	Username  string `json:"username,omitempty"`
}

// TelegramChat represents a Telegram chat
type TelegramChat struct {
	ID   int64  `json:"id"`
	Type string `json:"type"`
}

// InlineQuery represents an inline query from Telegram
type InlineQuery struct {
	ID     string        `json:"id"`
	From   *TelegramUser `json:"from"`
	Query  string        `json:"query"`
	Offset string        `json:"offset"`
}

// ChosenInlineResult represents when user chooses an inline result
type ChosenInlineResult struct {
	ResultID string        `json:"result_id"`
	From     *TelegramUser `json:"from"`
	Query    string        `json:"query"`
}

// InlineQueryResultArticle represents an article result for inline query
type InlineQueryResultArticle struct {
	Type                string                   `json:"type"`
	ID                  string                   `json:"id"`
	Title               string                   `json:"title"`
	Description         string                   `json:"description,omitempty"`
	InputMessageContent *InputTextMessageContent `json:"input_message_content"`
	ReplyMarkup         *InlineKeyboardMarkup    `json:"reply_markup,omitempty"`
}

// InputTextMessageContent represents text message content
type InputTextMessageContent struct {
	MessageText string `json:"message_text"`
}

// InlineKeyboardMarkup represents an inline keyboard
type InlineKeyboardMarkup struct {
	InlineKeyboard [][]InlineKeyboardButton `json:"inline_keyboard"`
}

// InlineKeyboardButton represents a button in inline keyboard
type InlineKeyboardButton struct {
	Text   string      `json:"text"`
	WebApp *WebAppInfo `json:"web_app,omitempty"`
}

// WebAppInfo represents a Web App
type WebAppInfo struct {
	URL string `json:"url"`
}

// AnswerInlineQueryRequest represents request to answerInlineQuery
type AnswerInlineQueryRequest struct {
	InlineQueryID string        `json:"inline_query_id"`
	Results       []interface{} `json:"results"`
	CacheTime     int           `json:"cache_time,omitempty"`
}

package repository

import (
	"context"
	"time"

	"github.com/TheFantazer/anagrams.ru/internal/domain"
	"github.com/google/uuid"
)

type SessionWithResults struct {
	Session *domain.Session
	Results []*domain.Result
	Type    string // "created" or "participated"
}

type PaginatedSessions struct {
	Sessions   []*SessionWithResults
	Total      int
	Page       int
	PerPage    int
	TotalPages int
}

type SessionRepository interface {
	Create(ctx context.Context, session *domain.Session) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Session, error)
	GetByCreatorID(ctx context.Context, creatorID uuid.UUID, limit int) ([]*domain.Session, error)
	GetByParticipant(ctx context.Context, userID uuid.UUID, limit int) ([]*domain.Session, error)
	GetAllUserSessions(ctx context.Context, userID uuid.UUID, page int, perPage int) (*PaginatedSessions, error)
	Delete(ctx context.Context, id uuid.UUID) error
	DeleteExpired(ctx context.Context, before time.Time) (int64, error)
}

type ResultRepository interface {
	Create(ctx context.Context, result *domain.Result) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Result, error)
	GetBySessionID(ctx context.Context, sessionID uuid.UUID) ([]*domain.Result, error)
	GetTopBySessionID(ctx context.Context, sessionID uuid.UUID, limit int) ([]*domain.Result, error)
}

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
	GetByUsername(ctx context.Context, username string) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	GetByOAuthID(ctx context.Context, provider, oauthID string) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
	UpdateSettings(ctx context.Context, userID uuid.UUID, letterCount int, language string, timeLimit int) error
	LinkOAuth(ctx context.Context, userID uuid.UUID, provider, oauthID string) error
}

type FriendRepository interface {
	CreateRequest(ctx context.Context, fromUserId uuid.UUID, toUserId uuid.UUID) error
	GetPendingRequests(ctx context.Context, userId uuid.UUID) ([]*domain.FriendRequest, error)
	GetSendingRequests(ctx context.Context, userId uuid.UUID) ([]*domain.FriendRequest, error)
	AcceptRequest(ctx context.Context, requestId uuid.UUID) error
	RejectRequest(ctx context.Context, requestId uuid.UUID) error
	GetFriends(ctx context.Context, userId uuid.UUID) ([]*domain.User, error)
	RemoveFriend(ctx context.Context, userId uuid.UUID, friendId uuid.UUID) error
	SearchUsers(ctx context.Context, query string) ([]*domain.User, error)
	GetUserByID(ctx context.Context, userId uuid.UUID) (*domain.User, error)
	AreFriends(ctx context.Context, userId1 uuid.UUID, userId2 uuid.UUID) (bool, error)
}

type SessionInviteRepository interface {
	Create(ctx context.Context, invite *domain.SessionInvite) error
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.SessionInvite, error)
	GetBySessionID(ctx context.Context, sessionID uuid.UUID) ([]*domain.SessionInvite, error)
}

type SessionParticipantRepository interface {
	Create(ctx context.Context, participant *domain.SessionParticipant) error
	GetBySession(ctx context.Context, sessionID uuid.UUID) ([]*domain.SessionParticipant, error)
	GetBySessionAndUser(ctx context.Context, sessionID, userID uuid.UUID) (*domain.SessionParticipant, error)
	CountOpponents(ctx context.Context, sessionID uuid.UUID) (int, error)
	MarkAsStarted(ctx context.Context, sessionID, userID uuid.UUID) error
	HasStarted(ctx context.Context, sessionID, userID uuid.UUID) (bool, error)
}

type UserStats struct {
	GamesPlayed  int
	BestScore    int
	LongestWord  string
	TotalWords   int
	AverageScore float64
}

type LeaderboardEntry struct {
	Username   string
	BestScore  int
	TotalWords int
}

type StatsRepository interface {
	GetUserStats(ctx context.Context, userID uuid.UUID) (*UserStats, error)
	GetLeaderboard(ctx context.Context, period string, limit int) ([]*LeaderboardEntry, error)
}

type DailyPuzzleRepository interface {
	Create(ctx context.Context, puzzle *domain.DailyPuzzle) error
	GetByDate(ctx context.Context, date time.Time) (*domain.DailyPuzzle, error)
	GetToday(ctx context.Context) (*domain.DailyPuzzle, error)
}

type UserDailyStatsRepository interface {
	GetByUserID(ctx context.Context, userID uuid.UUID) (*domain.UserDailyStats, error)
	Upsert(ctx context.Context, stats *domain.UserDailyStats) error
	HasPlayedToday(ctx context.Context, userID uuid.UUID) (bool, error)
}

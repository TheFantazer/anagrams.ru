package repository

import (
	"context"
	"time"

	"github.com/TheFantazer/anagrams.ru/internal/domain"
	"github.com/google/uuid"
)

type SessionRepository interface {
	Create(ctx context.Context, session *domain.Session) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Session, error)
	GetByCreatorID(ctx context.Context, creatorID uuid.UUID, limit int) ([]*domain.Session, error)
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

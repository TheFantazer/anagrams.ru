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
	Update(ctx context.Context, user *domain.User) error
}

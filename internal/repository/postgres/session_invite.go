package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/TheFantazer/anagrams.ru/internal/domain"
	"github.com/TheFantazer/anagrams.ru/internal/repository"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type sessionInviteRepository struct {
	db *sqlx.DB
}

func NewSessionInviteRepository(db *sqlx.DB) repository.SessionInviteRepository {
	return &sessionInviteRepository{db: db}
}

type sessionInviteDB struct {
	ID         uuid.UUID `db:"id"`
	SessionID  uuid.UUID `db:"session_id"`
	FromUserID uuid.UUID `db:"from_user_id"`
	ToUserID   uuid.UUID `db:"to_user_id"`
	Status     string    `db:"status"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
}

func (s *sessionInviteDB) toDomain() *domain.SessionInvite {
	return &domain.SessionInvite{
		ID:         s.ID,
		SessionID:  s.SessionID,
		FromUserID: s.FromUserID,
		ToUserID:   s.ToUserID,
		Status:     s.Status,
		CreatedAt:  s.CreatedAt,
		UpdatedAt:  s.UpdatedAt,
	}
}

func fromDomainSessionInvite(s *domain.SessionInvite) *sessionInviteDB {
	return &sessionInviteDB{
		ID:         s.ID,
		SessionID:  s.SessionID,
		FromUserID: s.FromUserID,
		ToUserID:   s.ToUserID,
		Status:     s.Status,
		CreatedAt:  s.CreatedAt,
		UpdatedAt:  s.UpdatedAt,
	}
}

func (r *sessionInviteRepository) Create(ctx context.Context, invite *domain.SessionInvite) error {
	dbInvite := fromDomainSessionInvite(invite)
	query := `
		INSERT INTO session_invites (id, session_id, from_user_id, to_user_id, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := r.db.ExecContext(
		ctx,
		query,
		dbInvite.ID,
		dbInvite.SessionID,
		dbInvite.FromUserID,
		dbInvite.ToUserID,
		dbInvite.Status,
		dbInvite.CreatedAt,
		dbInvite.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to insert session invite: %w", err)
	}

	return nil
}

func (r *sessionInviteRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.SessionInvite, error) {
	query := `
		SELECT id, session_id, from_user_id, to_user_id, status, created_at, updated_at
		FROM session_invites
		WHERE to_user_id = $1
		ORDER BY created_at DESC
	`

	var dbInvites []sessionInviteDB
	err := r.db.SelectContext(ctx, &dbInvites, query, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []*domain.SessionInvite{}, nil
		}
		return nil, fmt.Errorf("failed to get session invites: %w", err)
	}

	invites := make([]*domain.SessionInvite, 0, len(dbInvites))
	for _, dbInvite := range dbInvites {
		invites = append(invites, dbInvite.toDomain())
	}

	return invites, nil
}

func (r *sessionInviteRepository) GetBySessionID(ctx context.Context, sessionID uuid.UUID) ([]*domain.SessionInvite, error) {
	query := `
		SELECT id, session_id, from_user_id, to_user_id, status, created_at, updated_at
		FROM session_invites
		WHERE session_id = $1
		ORDER BY created_at DESC
	`

	var dbInvites []sessionInviteDB
	err := r.db.SelectContext(ctx, &dbInvites, query, sessionID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []*domain.SessionInvite{}, nil
		}
		return nil, fmt.Errorf("failed to get session invites: %w", err)
	}

	invites := make([]*domain.SessionInvite, 0, len(dbInvites))
	for _, dbInvite := range dbInvites {
		invites = append(invites, dbInvite.toDomain())
	}

	return invites, nil
}

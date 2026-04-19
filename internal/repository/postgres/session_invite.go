package postgres

import (
	"context"
	"fmt"

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

func (r *sessionInviteRepository) CreateInvite(ctx context.Context, sessionID uuid.UUID, userID uuid.UUID) (*domain.SessionInvite, error) {
	invite := domain.NewSessionInvite(sessionID, userID)

	query := `
		INSERT INTO session_invites (id, session_id, user_id, created_at)
		VALUES ($1, $2, $3, $4)
	`

	_, err := r.db.ExecContext(ctx, query,
		invite.ID,
		invite.SessionID,
		invite.UserID,
		invite.CreatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create session invite: %w", err)
	}

	return invite, nil
}

func (r *sessionInviteRepository) GetInvites(ctx context.Context, userID uuid.UUID) ([]*domain.SessionInvite, error) {
	query := `
		SELECT id, session_id, user_id, created_at
		FROM session_invites
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	var invites []*domain.SessionInvite
	err := r.db.SelectContext(ctx, &invites, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get invites: %w", err)
	}

	return invites, nil
}

func (r *sessionInviteRepository) GetSessionInvites(ctx context.Context, sessionID uuid.UUID) ([]*domain.SessionInvite, error) {
	query := `
		SELECT id, session_id, user_id, created_at
		FROM session_invites
		WHERE session_id = $1
		ORDER BY created_at DESC
	`

	var invites []*domain.SessionInvite
	err := r.db.SelectContext(ctx, &invites, query, sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get session invites: %w", err)
	}

	return invites, nil
}

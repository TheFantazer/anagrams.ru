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

type postgresSessionParticipantRepo struct {
	db *sqlx.DB
}

func NewSessionParticipantRepository(db *sqlx.DB) repository.SessionParticipantRepository {
	return &postgresSessionParticipantRepo{db: db}
}

type sessionParticipantDB struct {
	ID        uuid.UUID  `db:"id"`
	SessionID uuid.UUID  `db:"session_id"`
	UserID    uuid.UUID  `db:"user_id"`
	Role      string     `db:"role"`
	JoinedAt  time.Time  `db:"joined_at"`
	StartedAt *time.Time `db:"started_at"`
}

func (p *sessionParticipantDB) toDomain() *domain.SessionParticipant {
	return &domain.SessionParticipant{
		ID:        p.ID,
		SessionID: p.SessionID,
		UserID:    p.UserID,
		Role:      p.Role,
		JoinedAt:  p.JoinedAt,
		StartedAt: p.StartedAt,
	}
}

func fromDomainParticipant(p *domain.SessionParticipant) *sessionParticipantDB {
	return &sessionParticipantDB{
		ID:        p.ID,
		SessionID: p.SessionID,
		UserID:    p.UserID,
		Role:      p.Role,
		JoinedAt:  p.JoinedAt,
		StartedAt: p.StartedAt,
	}
}

func (r *postgresSessionParticipantRepo) Create(ctx context.Context, participant *domain.SessionParticipant) error {
	dbParticipant := fromDomainParticipant(participant)
	query := `
		INSERT INTO session_participants (id, session_id, user_id, role, joined_at, started_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (session_id, user_id) DO NOTHING
	`
	_, err := r.db.ExecContext(
		ctx,
		query,
		dbParticipant.ID,
		dbParticipant.SessionID,
		dbParticipant.UserID,
		dbParticipant.Role,
		dbParticipant.JoinedAt,
		dbParticipant.StartedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to insert participant: %w", err)
	}
	return nil
}

func (r *postgresSessionParticipantRepo) GetBySession(ctx context.Context, sessionID uuid.UUID) ([]*domain.SessionParticipant, error) {
	query := `
		SELECT id, session_id, user_id, role, joined_at, started_at
		FROM session_participants
		WHERE session_id = $1
		ORDER BY joined_at ASC
	`
	var dbParticipants []sessionParticipantDB
	err := r.db.SelectContext(ctx, &dbParticipants, query, sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get participants by session: %w", err)
	}

	participants := make([]*domain.SessionParticipant, len(dbParticipants))
	for i, dbP := range dbParticipants {
		participants[i] = dbP.toDomain()
	}
	return participants, nil
}

func (r *postgresSessionParticipantRepo) GetBySessionAndUser(ctx context.Context, sessionID, userID uuid.UUID) (*domain.SessionParticipant, error) {
	query := `
		SELECT id, session_id, user_id, role, joined_at, started_at
		FROM session_participants
		WHERE session_id = $1 AND user_id = $2
	`
	var dbParticipant sessionParticipantDB
	err := r.db.GetContext(ctx, &dbParticipant, query, sessionID, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get participant: %w", err)
	}
	return dbParticipant.toDomain(), nil
}

func (r *postgresSessionParticipantRepo) CountOpponents(ctx context.Context, sessionID uuid.UUID) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM session_participants
		WHERE session_id = $1 AND role = 'opponent'
	`
	var count int
	err := r.db.GetContext(ctx, &count, query, sessionID)
	if err != nil {
		return 0, fmt.Errorf("failed to count opponents: %w", err)
	}
	return count, nil
}

func (r *postgresSessionParticipantRepo) MarkAsStarted(ctx context.Context, sessionID, userID uuid.UUID) error {
	query := `
		UPDATE session_participants
		SET started_at = NOW()
		WHERE session_id = $1 AND user_id = $2 AND started_at IS NULL
	`
	result, err := r.db.ExecContext(ctx, query, sessionID, userID)
	if err != nil {
		return fmt.Errorf("failed to mark as started: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return repository.ErrNotFound
	}
	return nil
}

func (r *postgresSessionParticipantRepo) HasStarted(ctx context.Context, sessionID, userID uuid.UUID) (bool, error) {
	query := `
		SELECT started_at IS NOT NULL
		FROM session_participants
		WHERE session_id = $1 AND user_id = $2
	`
	var hasStarted bool
	err := r.db.GetContext(ctx, &hasStarted, query, sessionID, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("failed to check if started: %w", err)
	}
	return hasStarted, nil
}

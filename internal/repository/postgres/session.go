package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/TheFantazer/anagrams.ru/internal/domain"
	"github.com/TheFantazer/anagrams.ru/internal/repository"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type postgresSessionRepo struct {
	db *sqlx.DB
}

func NewSessionRepository(db *sqlx.DB) repository.SessionRepository {
	return &postgresSessionRepo{db: db}
}

type sessionDB struct {
	ID          uuid.UUID       `db:"id"`
	Letters     string          `db:"letters"`
	Language    string          `db:"language"`
	TimeLimit   int             `db:"time_limit"`
	LetterCount int             `db:"letter_count"`
	ValidWords  json.RawMessage `db:"valid_words"`
	MaxScore    int             `db:"max_score"`
	CreatedAt   time.Time       `db:"created_at"`
}

func (s *sessionDB) toDomain() (*domain.Session, error) {
	var validWords []string
	if err := json.Unmarshal(s.ValidWords, &validWords); err != nil {
		return nil, fmt.Errorf("failed to unmarshal valid_words: %w", err)
	}

	return &domain.Session{
		ID:          s.ID,
		Letters:     s.Letters,
		Language:    s.Language,
		TimeLimit:   s.TimeLimit,
		LetterCount: s.LetterCount,
		ValidWords:  validWords,
		MaxScore:    s.MaxScore,
		CreatedAt:   s.CreatedAt,
	}, nil
}
func fromDomainSession(s *domain.Session) (*sessionDB, error) {
	validWordsJSON, err := json.Marshal(s.ValidWords)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal valid_words: %w", err)
	}
	return &sessionDB{
		ID:          s.ID,
		Letters:     s.Letters,
		Language:    s.Language,
		TimeLimit:   s.TimeLimit,
		LetterCount: s.LetterCount,
		ValidWords:  validWordsJSON,
		MaxScore:    s.MaxScore,
		CreatedAt:   s.CreatedAt,
	}, nil
}
func (r *postgresSessionRepo) Create(ctx context.Context, s *domain.Session) error {
	dbSession, err := fromDomainSession(s)
	if err != nil {
		return fmt.Errorf("failed to convert session: %w", err)
	}
	query := `
				INSERT INTO game_sessions (id, letters, language, time_limit, letter_count, valid_words, max_score, created_at)
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
`
	_, err = r.db.ExecContext(
		ctx,
		query,
		dbSession.ID,
		dbSession.Letters,
		dbSession.Language,
		dbSession.TimeLimit,
		dbSession.LetterCount,
		dbSession.ValidWords,
		dbSession.MaxScore,
		dbSession.CreatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to insert session: %w", err)
	}

	return nil
}

func (r *postgresSessionRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Session, error) {
	query := `
  		SELECT id, letters, language, time_limit, letter_count, valid_words, max_score, created_at
  		FROM game_sessions
  		WHERE id = $1
  	`

	var dbSession sessionDB
	err := r.db.GetContext(ctx, &dbSession, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	return dbSession.toDomain()
}

func (r *postgresSessionRepo) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM game_sessions WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
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

func (r *postgresSessionRepo) DeleteExpired(ctx context.Context, before time.Time) (int64, error) {
	query := `DELETE FROM game_sessions WHERE created_at < $1`
	result, err := r.db.ExecContext(ctx, query, before)
	if err != nil {
		return 0, fmt.Errorf("failed to delete expired session: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get rows affected: %w", err)
	}
	return rowsAffected, nil
}

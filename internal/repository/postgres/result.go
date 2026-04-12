package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/TheFantazer/anagrams.ru/internal/domain"
	"github.com/TheFantazer/anagrams.ru/internal/repository"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type postgresResultRepo struct {
	db *sqlx.DB
}

func NewResultRepository(db *sqlx.DB) repository.ResultRepository {
	return &postgresResultRepo{db: db}
}

type resultDB struct {
	ID                uuid.UUID       `db:"id"`
	SessionID         uuid.UUID       `db:"session_id"`
	UserID            *uuid.UUID      `db:"user_id"`
	PlayerName        string          `db:"player_name"`
	PlayerFingerprint string          `db:"player_fingerprint"`
	FoundWords        json.RawMessage `db:"found_words"`
	WordCount         int             `db:"word_count"`
	Score             int             `db:"score"`
	DurationMs        int             `db:"duration_ms"`
	PlayedAt          time.Time       `db:"played_at"`
}

func (r *resultDB) toDomain() (*domain.Result, error) {
	var foundWords []string
	if err := json.Unmarshal(r.FoundWords, &foundWords); err != nil {
		return nil, fmt.Errorf("failed to unmarshal found_words: %w", err)
	}

	return &domain.Result{
		ID:                r.ID,
		SessionID:         r.SessionID,
		UserID:            r.UserID,
		PlayerName:        r.PlayerName,
		PlayerFingerprint: r.PlayerFingerprint,
		FoundWords:        foundWords,
		WordCount:         r.WordCount,
		Score:             r.Score,
		DurationMs:        r.DurationMs,
		PlayedAt:          r.PlayedAt,
	}, nil
}

func fromDomainResult(r *domain.Result) (*resultDB, error) {
	foundWordsJSON, err := json.Marshal(r.FoundWords)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal found_words: %w", err)
	}

	return &resultDB{
		ID:                r.ID,
		SessionID:         r.SessionID,
		UserID:            r.UserID,
		PlayerName:        r.PlayerName,
		PlayerFingerprint: r.PlayerFingerprint,
		FoundWords:        foundWordsJSON,
		WordCount:         r.WordCount,
		Score:             r.Score,
		DurationMs:        r.DurationMs,
		PlayedAt:          r.PlayedAt,
	}, nil
}

func (r *postgresResultRepo) Create(ctx context.Context, result *domain.Result) error {
	dbResult, err := fromDomainResult(result)
	if err != nil {
		return fmt.Errorf("failed to convert result: %w", err)
	}

	query := `
  		INSERT INTO game_results (id, session_id, user_id, player_name, player_fingerprint, found_words, word_count, score, duration_ms,
  played_at)
  		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
  	`

	_, err = r.db.ExecContext(
		ctx,
		query,
		dbResult.ID,
		dbResult.SessionID,
		dbResult.UserID,
		dbResult.PlayerName,
		dbResult.PlayerFingerprint,
		dbResult.FoundWords,
		dbResult.WordCount,
		dbResult.Score,
		dbResult.DurationMs,
		dbResult.PlayedAt,
	)

	if err != nil {
		// Проверяем на нарушение UNIQUE constraint
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" { // unique_violation
				if strings.Contains(pqErr.Constraint, "uq_session_player") {
					return repository.ErrDuplicateResult
				}
			}
			// Проверяем на нарушение FK constraint
			if pqErr.Code == "23503" { // foreign_key_violation
				return repository.ErrForeignKeyViolation
			}
		}
		return fmt.Errorf("failed to insert result: %w", err)
	}

	return nil
}

func (r *postgresResultRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Result, error) {
	query := `
  		SELECT id, session_id, user_id, player_name, player_fingerprint, found_words, word_count, score, duration_ms, played_at
  		FROM game_results
  		WHERE id = $1
  	`

	var dbResult resultDB
	err := r.db.GetContext(ctx, &dbResult, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get result: %w", err)
	}

	return dbResult.toDomain()
}

func (r *postgresResultRepo) GetBySessionID(ctx context.Context, sessionID uuid.UUID) ([]*domain.Result, error) {
	query := `
  		SELECT id, session_id, user_id, player_name, player_fingerprint, found_words, word_count, score, duration_ms, played_at
  		FROM game_results
  		WHERE session_id = $1
  		ORDER BY played_at DESC
  	`

	var dbResults []resultDB
	err := r.db.SelectContext(ctx, &dbResults, query, sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get results by session: %w", err)
	}

	results := make([]*domain.Result, 0, len(dbResults))
	for _, dbResult := range dbResults {
		result, err := dbResult.toDomain()
		if err != nil {
			return nil, err
		}
		results = append(results, result)
	}

	return results, nil
}
func (r *postgresResultRepo) GetTopBySessionID(ctx context.Context, sessionID uuid.UUID, limit int) ([]*domain.Result, error) {
	query := `
  		SELECT id, session_id, user_id, player_name, player_fingerprint, found_words, word_count, score, duration_ms, played_at
  		FROM game_results
  		WHERE session_id = $1
  		ORDER BY score DESC, duration_ms ASC
  		LIMIT $2
  	`

	var dbResults []resultDB
	err := r.db.SelectContext(ctx, &dbResults, query, sessionID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get top results: %w", err)
	}

	results := make([]*domain.Result, 0, len(dbResults))
	for _, dbResult := range dbResults {
		result, err := dbResult.toDomain()
		if err != nil {
			return nil, err
		}
		results = append(results, result)
	}

	return results, nil
}

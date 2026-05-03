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
	ID            uuid.UUID       `db:"id"`
	Letters       string          `db:"letters"`
	Language      string          `db:"language"`
	TimeLimit     int             `db:"time_limit"`
	LetterCount   int             `db:"letter_count"`
	ValidWords    json.RawMessage `db:"valid_words"`
	MaxScore      int             `db:"max_score"`
	CreatorID     *uuid.UUID      `db:"creator_id"`
	CreatedAt     time.Time       `db:"created_at"`
	IsDaily       bool            `db:"is_daily"`
	DailyPuzzleID *uuid.UUID      `db:"daily_puzzle_id"`
}

func (s *sessionDB) toDomain() (*domain.Session, error) {
	var validWords []string
	if err := json.Unmarshal(s.ValidWords, &validWords); err != nil {
		return nil, fmt.Errorf("failed to unmarshal valid_words: %w", err)
	}

	return &domain.Session{
		ID:            s.ID,
		Letters:       s.Letters,
		Language:      s.Language,
		TimeLimit:     s.TimeLimit,
		LetterCount:   s.LetterCount,
		ValidWords:    validWords,
		MaxScore:      s.MaxScore,
		CreatorID:     s.CreatorID,
		CreatedAt:     s.CreatedAt,
		IsDaily:       s.IsDaily,
		DailyPuzzleID: s.DailyPuzzleID,
	}, nil
}
func fromDomainSession(s *domain.Session) (*sessionDB, error) {
	validWordsJSON, err := json.Marshal(s.ValidWords)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal valid_words: %w", err)
	}
	return &sessionDB{
		ID:            s.ID,
		Letters:       s.Letters,
		Language:      s.Language,
		TimeLimit:     s.TimeLimit,
		LetterCount:   s.LetterCount,
		ValidWords:    validWordsJSON,
		MaxScore:      s.MaxScore,
		CreatorID:     s.CreatorID,
		CreatedAt:     s.CreatedAt,
		IsDaily:       s.IsDaily,
		DailyPuzzleID: s.DailyPuzzleID,
	}, nil
}
func (r *postgresSessionRepo) Create(ctx context.Context, s *domain.Session) error {
	dbSession, err := fromDomainSession(s)
	if err != nil {
		return fmt.Errorf("failed to convert session: %w", err)
	}
	query := `
				INSERT INTO game_sessions (id, letters, language, time_limit, letter_count, valid_words, max_score, creator_id, created_at, is_daily, daily_puzzle_id)
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
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
		dbSession.CreatorID,
		dbSession.CreatedAt,
		dbSession.IsDaily,
		dbSession.DailyPuzzleID,
	)

	if err != nil {
		return fmt.Errorf("failed to insert session: %w", err)
	}

	return nil
}

func (r *postgresSessionRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Session, error) {
	query := `
  		SELECT id, letters, language, time_limit, letter_count, valid_words, max_score, creator_id, created_at
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

func (r *postgresSessionRepo) GetByCreatorID(ctx context.Context, creatorID uuid.UUID, limit int) ([]*domain.Session, error) {
	query := `
		SELECT id, letters, language, time_limit, letter_count, valid_words, max_score, creator_id, created_at
		FROM game_sessions
		WHERE creator_id = $1
		ORDER BY created_at DESC
		LIMIT $2
	`

	var dbSessions []sessionDB
	err := r.db.SelectContext(ctx, &dbSessions, query, creatorID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get sessions by creator: %w", err)
	}

	sessions := make([]*domain.Session, 0, len(dbSessions))
	for _, dbSession := range dbSessions {
		session, err := dbSession.toDomain()
		if err != nil {
			return nil, err
		}
		sessions = append(sessions, session)
	}

	return sessions, nil
}

func (r *postgresSessionRepo) GetByParticipant(ctx context.Context, userID uuid.UUID, limit int) ([]*domain.Session, error) {
	query := `
		SELECT DISTINCT gs.id, gs.letters, gs.language, gs.time_limit, gs.letter_count,
		       gs.valid_words, gs.max_score, gs.creator_id, gs.created_at
		FROM game_sessions gs
		INNER JOIN game_results gr ON gs.id = gr.session_id
		WHERE gr.user_id = $1
		ORDER BY gs.created_at DESC
		LIMIT $2
	`

	var dbSessions []sessionDB
	err := r.db.SelectContext(ctx, &dbSessions, query, userID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get sessions by participant: %w", err)
	}

	sessions := make([]*domain.Session, 0, len(dbSessions))
	for _, dbSession := range dbSessions {
		session, err := dbSession.toDomain()
		if err != nil {
			return nil, err
		}
		sessions = append(sessions, session)
	}

	return sessions, nil
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

func (r *postgresSessionRepo) GetAllUserSessions(ctx context.Context, userID uuid.UUID, page int, perPage int) (*repository.PaginatedSessions, error) {
	countQuery := `
		SELECT COUNT(DISTINCT gs.id)
		FROM game_sessions gs
		LEFT JOIN game_results gr ON gs.id = gr.session_id AND gr.user_id = $1
		LEFT JOIN session_invites si ON gs.id = si.session_id AND si.to_user_id = $1
		WHERE gs.creator_id = $1 OR gr.user_id = $1 OR si.to_user_id = $1
	`

	var total int
	err := r.db.GetContext(ctx, &total, countQuery, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to count sessions: %w", err)
	}

	offset := (page - 1) * perPage
	totalPages := (total + perPage - 1) / perPage

	query := `
		WITH user_sessions AS (
			SELECT DISTINCT gs.id, gs.created_at,
			       CASE
			           WHEN gs.creator_id = $1 THEN 'created'
			           WHEN EXISTS(SELECT 1 FROM game_results WHERE session_id = gs.id AND user_id = $1) THEN 'participated'
			           WHEN EXISTS(SELECT 1 FROM session_invites WHERE session_id = gs.id AND to_user_id = $1) THEN 'invited'
			           ELSE 'participated'
			       END as session_type
			FROM game_sessions gs
			LEFT JOIN game_results gr ON gs.id = gr.session_id AND gr.user_id = $1
			LEFT JOIN session_invites si ON gs.id = si.session_id AND si.to_user_id = $1
			WHERE gs.creator_id = $1 OR gr.user_id = $1 OR si.to_user_id = $1
			ORDER BY gs.created_at DESC
			LIMIT $2 OFFSET $3
		)
		SELECT
			gs.id, gs.letters, gs.language, gs.time_limit, gs.letter_count,
			gs.valid_words, gs.max_score, gs.creator_id, gs.created_at,
			us.session_type,
			gr.id as result_id, gr.session_id as result_session_id, gr.user_id as result_user_id,
			gr.player_name, gr.player_fingerprint, gr.found_words, gr.word_count, gr.score, gr.duration_ms, gr.played_at
		FROM user_sessions us
		JOIN game_sessions gs ON us.id = gs.id
		LEFT JOIN game_results gr ON gs.id = gr.session_id
		ORDER BY gs.created_at DESC, gr.score DESC
	`

	type sessionResultRow struct {
		sessionDB
		SessionType     string           `db:"session_type"`
		ResultID        *uuid.UUID       `db:"result_id"`
		ResultSessionID *uuid.UUID       `db:"result_session_id"`
		ResultUserID    *uuid.UUID       `db:"result_user_id"`
		PlayerName      *string          `db:"player_name"`
		Fingerprint     *string          `db:"player_fingerprint"`
		FoundWords      *json.RawMessage `db:"found_words"`
		WordCount       *int             `db:"word_count"`
		Score           *int             `db:"score"`
		DurationMs      *int             `db:"duration_ms"`
		PlayedAt        *time.Time       `db:"played_at"`
	}

	var rows []sessionResultRow
	err = r.db.SelectContext(ctx, &rows, query, userID, perPage, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get user sessions: %w", err)
	}

	sessionsMap := make(map[uuid.UUID]*repository.SessionWithResults)
	var orderedSessionIDs []uuid.UUID

	for _, row := range rows {
		if _, exists := sessionsMap[row.ID]; !exists {
			session, err := row.toDomain()
			if err != nil {
				return nil, err
			}

			sessionsMap[row.ID] = &repository.SessionWithResults{
				Session: session,
				Results: []*domain.Result{},
				Type:    row.SessionType,
			}
			orderedSessionIDs = append(orderedSessionIDs, row.ID)
		}

		if row.ResultID != nil {
			var foundWords []string
			if row.FoundWords != nil {
				if err := json.Unmarshal(*row.FoundWords, &foundWords); err != nil {
					return nil, fmt.Errorf("failed to unmarshal found_words: %w", err)
				}
			}

			result := &domain.Result{
				ID:                *row.ResultID,
				SessionID:         *row.ResultSessionID,
				UserID:            row.ResultUserID,
				PlayerName:        *row.PlayerName,
				PlayerFingerprint: *row.Fingerprint,
				FoundWords:        foundWords,
				WordCount:         *row.WordCount,
				Score:             *row.Score,
				DurationMs:        *row.DurationMs,
				PlayedAt:          *row.PlayedAt,
			}
			sessionsMap[row.ID].Results = append(sessionsMap[row.ID].Results, result)
		}
	}

	sessions := make([]*repository.SessionWithResults, 0, len(orderedSessionIDs))
	for _, id := range orderedSessionIDs {
		sessions = append(sessions, sessionsMap[id])
	}

	return &repository.PaginatedSessions{
		Sessions:   sessions,
		Total:      total,
		Page:       page,
		PerPage:    perPage,
		TotalPages: totalPages,
	}, nil
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

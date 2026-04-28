package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/TheFantazer/anagrams.ru/internal/domain"
	"github.com/TheFantazer/anagrams.ru/internal/repository"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type postgresDailyPuzzleRepo struct {
	db *sqlx.DB
}

func NewDailyPuzzleRepository(db *sqlx.DB) repository.DailyPuzzleRepository {
	return &postgresDailyPuzzleRepo{db: db}
}

type dailyPuzzleDB struct {
	ID         uuid.UUID `db:"id"`
	PuzzleDate time.Time `db:"puzzle_date"`
	Letters    string    `db:"letters"`
	Language   string    `db:"language"`
	CreatedAt  time.Time `db:"created_at"`
}

func (dp *dailyPuzzleDB) toDomain() *domain.DailyPuzzle {
	return &domain.DailyPuzzle{
		ID:         dp.ID,
		PuzzleDate: dp.PuzzleDate,
		Letters:    dp.Letters,
		Language:   dp.Language,
		CreatedAt:  dp.CreatedAt,
	}
}

func (r *postgresDailyPuzzleRepo) Create(ctx context.Context, puzzle *domain.DailyPuzzle) error {
	query := `
		INSERT INTO daily_puzzles (id, puzzle_date, letters, language, created_at)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (puzzle_date) DO NOTHING
	`

	_, err := r.db.ExecContext(ctx, query,
		puzzle.ID,
		puzzle.PuzzleDate,
		puzzle.Letters,
		puzzle.Language,
		puzzle.CreatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *postgresDailyPuzzleRepo) GetByDate(ctx context.Context, date time.Time) (*domain.DailyPuzzle, error) {
	// Normalize to UTC date
	normalizedDate := time.Date(
		date.Year(),
		date.Month(),
		date.Day(),
		0, 0, 0, 0,
		time.UTC,
	)

	query := `
		SELECT id, puzzle_date, letters, language, created_at
		FROM daily_puzzles
		WHERE puzzle_date = $1
	`

	var dbPuzzle dailyPuzzleDB
	err := r.db.GetContext(ctx, &dbPuzzle, query, normalizedDate)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}

	return dbPuzzle.toDomain(), nil
}

func (r *postgresDailyPuzzleRepo) GetToday(ctx context.Context) (*domain.DailyPuzzle, error) {
	today := time.Now().UTC()
	return r.GetByDate(ctx, today)
}

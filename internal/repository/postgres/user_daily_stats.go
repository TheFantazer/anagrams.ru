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

type postgresUserDailyStatsRepo struct {
	db *sqlx.DB
}

func NewUserDailyStatsRepository(db *sqlx.DB) repository.UserDailyStatsRepository {
	return &postgresUserDailyStatsRepo{db: db}
}

type userDailyStatsDB struct {
	UserID          uuid.UUID  `db:"user_id"`
	CurrentStreak   int        `db:"current_streak"`
	LongestStreak   int        `db:"longest_streak"`
	LastPlayedDate  *time.Time `db:"last_played_date"`
	TotalDailyGames int        `db:"total_daily_games"`
	CreatedAt       time.Time  `db:"created_at"`
	UpdatedAt       time.Time  `db:"updated_at"`
}

func (uds *userDailyStatsDB) toDomain() *domain.UserDailyStats {
	return &domain.UserDailyStats{
		UserID:          uds.UserID,
		CurrentStreak:   uds.CurrentStreak,
		LongestStreak:   uds.LongestStreak,
		LastPlayedDate:  uds.LastPlayedDate,
		TotalDailyGames: uds.TotalDailyGames,
		CreatedAt:       uds.CreatedAt,
		UpdatedAt:       uds.UpdatedAt,
	}
}

func (r *postgresUserDailyStatsRepo) GetByUserID(ctx context.Context, userID uuid.UUID) (*domain.UserDailyStats, error) {
	query := `
		SELECT user_id, current_streak, longest_streak, last_played_date,
		       total_daily_games, created_at, updated_at
		FROM user_daily_stats
		WHERE user_id = $1
	`

	var dbStats userDailyStatsDB
	err := r.db.GetContext(ctx, &dbStats, query, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Return empty stats for new users
			return domain.NewUserDailyStats(userID), nil
		}
		return nil, err
	}

	return dbStats.toDomain(), nil
}

func (r *postgresUserDailyStatsRepo) Upsert(ctx context.Context, stats *domain.UserDailyStats) error {
	query := `
		INSERT INTO user_daily_stats (
			user_id, current_streak, longest_streak, last_played_date,
			total_daily_games, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (user_id)
		DO UPDATE SET
			current_streak = EXCLUDED.current_streak,
			longest_streak = EXCLUDED.longest_streak,
			last_played_date = EXCLUDED.last_played_date,
			total_daily_games = EXCLUDED.total_daily_games,
			updated_at = EXCLUDED.updated_at
	`

	_, err := r.db.ExecContext(ctx, query,
		stats.UserID,
		stats.CurrentStreak,
		stats.LongestStreak,
		stats.LastPlayedDate,
		stats.TotalDailyGames,
		stats.CreatedAt,
		stats.UpdatedAt,
	)

	return err
}

func (r *postgresUserDailyStatsRepo) HasPlayedToday(ctx context.Context, userID uuid.UUID) (bool, error) {
	today := time.Now().UTC()
	todayDate := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, time.UTC)

	query := `
		SELECT EXISTS(
			SELECT 1
			FROM game_sessions gs
			INNER JOIN game_results gr ON gs.id = gr.session_id
			WHERE gs.is_daily = true
			  AND gr.user_id = $1
			  AND DATE(gr.played_at) = $2
		)
	`

	var exists bool
	err := r.db.GetContext(ctx, &exists, query, userID, todayDate)
	return exists, err
}

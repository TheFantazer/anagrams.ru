package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/TheFantazer/anagrams.ru/internal/repository"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type statsRepository struct {
	db *sqlx.DB
}

func NewStatsRepository(db *sqlx.DB) repository.StatsRepository {
	return &statsRepository{db: db}
}

func (r *statsRepository) GetUserStats(ctx context.Context, userID uuid.UUID) (*repository.UserStats, error) {
	query := `
		SELECT
			COUNT(*) as games_played,
			COALESCE(MAX(score), 0) as best_score,
			COALESCE(SUM(word_count), 0) as total_words,
			COALESCE(AVG(score), 0) as average_score
		FROM game_results
		WHERE user_id = $1
	`

	stats := &repository.UserStats{}
	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&stats.GamesPlayed,
		&stats.BestScore,
		&stats.TotalWords,
		&stats.AverageScore,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return &repository.UserStats{}, nil
		}
		return nil, fmt.Errorf("failed to get user stats: %w", err)
	}
	longestWordQuery := `
		SELECT found_words
		FROM game_results
		WHERE user_id = $1 AND found_words != '[]'::jsonb
		ORDER BY played_at DESC
		LIMIT 100
	`

	rows, err := r.db.QueryContext(ctx, longestWordQuery, userID)
	if err != nil {
		return stats, nil
	}
	defer rows.Close()

	longestWord := ""
	for rows.Next() {
		var foundWordsJSON []byte
		if err := rows.Scan(&foundWordsJSON); err != nil {
			continue
		}

		var words []string
		if err := json.Unmarshal(foundWordsJSON, &words); err != nil {
			continue
		}

		for _, word := range words {
			if len(word) > len(longestWord) {
				longestWord = word
			}
		}
	}

	stats.LongestWord = longestWord
	return stats, nil
}

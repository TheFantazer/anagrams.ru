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
	defer func() {
		if cerr := rows.Close(); cerr != nil {
			err = cerr
		}
	}()

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

func (r *statsRepository) GetLeaderboard(ctx context.Context, period string, limit int) ([]*repository.LeaderboardEntry, error) {
	var timeFilter string
	switch period {
	case "day":
		timeFilter = "AND gr.played_at >= NOW() - INTERVAL '1 day'"
	case "week":
		timeFilter = "AND gr.played_at >= NOW() - INTERVAL '7 days'"
	case "month":
		timeFilter = "AND gr.played_at >= NOW() - INTERVAL '30 days'"
	case "all":
		timeFilter = ""
	default:
		timeFilter = "AND gr.played_at >= NOW() - INTERVAL '7 days'"
	}

	query := fmt.Sprintf(`
		SELECT
			u.username,
			MAX(gr.score) as best_score,
			SUM(gr.word_count) as total_words
		FROM game_results gr
		INNER JOIN users u ON gr.user_id = u.id
		WHERE gr.user_id IS NOT NULL %s
		GROUP BY u.id, u.username
		ORDER BY best_score DESC
		LIMIT $1
	`, timeFilter)

	rows, err := r.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get leaderboard: %w", err)
	}
	defer func() {
		if cerr := rows.Close(); cerr != nil && err == nil {
			err = fmt.Errorf("failed to close rows: %w", cerr)
		}
	}()

	var entries []*repository.LeaderboardEntry
	for rows.Next() {
		entry := &repository.LeaderboardEntry{}
		if err := rows.Scan(&entry.Username, &entry.BestScore, &entry.TotalWords); err != nil {
			return nil, fmt.Errorf("failed to scan leaderboard entry: %w", err)
		}
		entries = append(entries, entry)
	}

	return entries, nil
}

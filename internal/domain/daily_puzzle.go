package domain

import (
	"time"

	"github.com/google/uuid"
)

type DailyPuzzle struct {
	ID         uuid.UUID `json:"id"`
	PuzzleDate time.Time `json:"puzzle_date"`
	Letters    string    `json:"letters"`
	Language   string    `json:"language"`
	CreatedAt  time.Time `json:"created_at"`
}

type UserDailyStats struct {
	UserID          uuid.UUID  `json:"user_id"`
	CurrentStreak   int        `json:"current_streak"`
	LongestStreak   int        `json:"longest_streak"`
	LastPlayedDate  *time.Time `json:"last_played_date,omitempty"`
	TotalDailyGames int        `json:"total_daily_games"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

// NewDailyPuzzle creates a new daily puzzle for the given date
func NewDailyPuzzle(puzzleDate time.Time, letters, language string) *DailyPuzzle {
	// Normalize to UTC date (00:00:00)
	normalizedDate := time.Date(
		puzzleDate.Year(),
		puzzleDate.Month(),
		puzzleDate.Day(),
		0, 0, 0, 0,
		time.UTC,
	)

	return &DailyPuzzle{
		ID:         uuid.New(),
		PuzzleDate: normalizedDate,
		Letters:    letters,
		Language:   language,
		CreatedAt:  time.Now().UTC(),
	}
}

// NewUserDailyStats creates initial stats for a user
func NewUserDailyStats(userID uuid.UUID) *UserDailyStats {
	now := time.Now().UTC()
	return &UserDailyStats{
		UserID:          userID,
		CurrentStreak:   0,
		LongestStreak:   0,
		LastPlayedDate:  nil,
		TotalDailyGames: 0,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
}

// UpdateStreak updates the user's streak after completing a daily puzzle
func (s *UserDailyStats) UpdateStreak(playedDate time.Time) {
	now := time.Date(
		playedDate.Year(),
		playedDate.Month(),
		playedDate.Day(),
		0, 0, 0, 0,
		time.UTC,
	)

	if s.LastPlayedDate == nil {
		// First time playing
		s.CurrentStreak = 1
		s.LongestStreak = 1
	} else {
		lastDate := time.Date(
			s.LastPlayedDate.Year(),
			s.LastPlayedDate.Month(),
			s.LastPlayedDate.Day(),
			0, 0, 0, 0,
			time.UTC,
		)

		daysDiff := int(now.Sub(lastDate).Hours() / 24)

		switch daysDiff {
		case 0:
			return
		case 1:
			s.CurrentStreak++
			if s.CurrentStreak > s.LongestStreak {
				s.LongestStreak = s.CurrentStreak
			}
		default:
			s.CurrentStreak = 1
		}
	}

	s.LastPlayedDate = &now
	s.TotalDailyGames++
	s.UpdatedAt = time.Now().UTC()
}

package service

import (
	"context"
	"fmt"
	"time"
	"unicode/utf8"

	"github.com/TheFantazer/anagrams.ru/internal/dictionary"
	"github.com/TheFantazer/anagrams.ru/internal/domain"
	"github.com/TheFantazer/anagrams.ru/internal/repository"
	"github.com/google/uuid"
)

type DailyPuzzleService interface {
	GetOrCreateTodaysPuzzle(ctx context.Context, language string) (*domain.DailyPuzzle, error)
	GetTodaysSession(ctx context.Context, language string) (*domain.Session, error)
	GetUserDailyStats(ctx context.Context, userID uuid.UUID) (*domain.UserDailyStats, error)
	HasPlayedToday(ctx context.Context, userID uuid.UUID) (bool, error)
	SubmitDailyResult(ctx context.Context, puzzleID, userID uuid.UUID, playerName, fingerprint string, words []string, durationMs int) (*domain.Result, error)
}

type dailyPuzzleService struct {
	puzzleRepo   repository.DailyPuzzleRepository
	statsRepo    repository.UserDailyStatsRepository
	sessionRepo  repository.SessionRepository
	resultRepo   repository.ResultRepository
	dictionaries map[string]*dictionary.Trie
	letterGen    *dictionary.LetterGenerator
}

func NewDailyPuzzleService(
	puzzleRepo repository.DailyPuzzleRepository,
	statsRepo repository.UserDailyStatsRepository,
	sessionRepo repository.SessionRepository,
	resultRepo repository.ResultRepository,
	dictionaries map[string]*dictionary.Trie,
	letterGen *dictionary.LetterGenerator,
) DailyPuzzleService {
	return &dailyPuzzleService{
		puzzleRepo:   puzzleRepo,
		statsRepo:    statsRepo,
		sessionRepo:  sessionRepo,
		resultRepo:   resultRepo,
		dictionaries: dictionaries,
		letterGen:    letterGen,
	}
}

// GetOrCreateTodaysPuzzle returns today's puzzle or creates it if it doesn't exist
func (s *dailyPuzzleService) GetOrCreateTodaysPuzzle(ctx context.Context, language string) (*domain.DailyPuzzle, error) {
	today := time.Now().UTC()

	puzzle, err := s.puzzleRepo.GetToday(ctx)
	if err == nil {
		return puzzle, nil
	}

	if err != repository.ErrNotFound {
		return nil, fmt.Errorf("failed to get today's puzzle: %w", err)
	}

	// Create new puzzle
	dict, ok := s.dictionaries[language]
	if !ok {
		return nil, domain.ErrUnsupportedLanguage
	}

	// Daily puzzles always use 7 letters
	const dailyLetterCount = 7
	letters := s.letterGen.GenerateFromDictionary(dict, dailyLetterCount)

	if letters == "" {
		letters = s.letterGen.GenerateBalancedLetters(language, dailyLetterCount)
	}

	validWords := dict.FindAllWords(letters)
	if len(validWords) == 0 {
		return s.GetOrCreateTodaysPuzzle(ctx, language)
	}

	puzzle = domain.NewDailyPuzzle(today, letters, language)

	if err := s.puzzleRepo.Create(ctx, puzzle); err != nil {
		return nil, fmt.Errorf("failed to create daily puzzle: %w", err)
	}

	return puzzle, nil
}

// GetTodaysSession creates a session from today's puzzle
func (s *dailyPuzzleService) GetTodaysSession(ctx context.Context, language string) (*domain.Session, error) {
	puzzle, err := s.GetOrCreateTodaysPuzzle(ctx, language)
	if err != nil {
		return nil, err
	}

	dict, ok := s.dictionaries[language]
	if !ok {
		return nil, domain.ErrUnsupportedLanguage
	}

	validWords := dict.FindAllWords(puzzle.Letters)
	if len(validWords) == 0 {
		return nil, fmt.Errorf("no valid words found for daily puzzle")
	}

	const dailyTimeLimit = 60

	session, err := domain.NewSession(
		puzzle.Letters,
		puzzle.Language,
		dailyTimeLimit,
		utf8.RuneCountInString(puzzle.Letters),
		validWords,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	session.IsDaily = true
	session.DailyPuzzleID = &puzzle.ID

	return session, nil
}

// GetUserDailyStats returns user's daily game statistics
func (s *dailyPuzzleService) GetUserDailyStats(ctx context.Context, userID uuid.UUID) (*domain.UserDailyStats, error) {
	stats, err := s.statsRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user daily stats: %w", err)
	}
	return stats, nil
}

// HasPlayedToday checks if user has played today's daily puzzle
func (s *dailyPuzzleService) HasPlayedToday(ctx context.Context, userID uuid.UUID) (bool, error) {
	hasPlayed, err := s.statsRepo.HasPlayedToday(ctx, userID)
	if err != nil {
		return false, fmt.Errorf("failed to check if user played today: %w", err)
	}
	return hasPlayed, nil
}

// SubmitDailyResult submits a result for a daily puzzle and updates user stats
func (s *dailyPuzzleService) SubmitDailyResult(ctx context.Context, puzzleID, userID uuid.UUID, playerName, fingerprint string, words []string, durationMs int) (*domain.Result, error) {
	puzzle, err := s.puzzleRepo.GetToday(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get today's puzzle: %w", err)
	}

	if puzzle.ID != puzzleID {
		return nil, fmt.Errorf("puzzle ID mismatch")
	}

	session, err := s.GetTodaysSession(ctx, puzzle.Language)
	if err != nil {
		return nil, err
	}

	existingSession, err := s.sessionRepo.GetByID(ctx, session.ID)
	if err == repository.ErrNotFound {
		if err := s.sessionRepo.Create(ctx, session); err != nil {
			return nil, fmt.Errorf("failed to create daily session: %w", err)
		}
	} else if err != nil {
		return nil, fmt.Errorf("failed to check session: %w", err)
	} else {
		session = existingSession
	}

	validatedWords := []string{}
	score := 0

	for _, word := range words {
		if session.IsValid(word) {
			validatedWords = append(validatedWords, word)
			score += domain.ScoreWord(word)
		}
	}

	result := &domain.Result{
		ID:                uuid.New(),
		SessionID:         session.ID,
		UserID:            &userID,
		PlayerName:        playerName,
		PlayerFingerprint: fingerprint,
		FoundWords:        validatedWords,
		WordCount:         len(validatedWords),
		Score:             score,
		DurationMs:        durationMs,
		PlayedAt:          time.Now().UTC(),
	}

	if err := s.resultRepo.Create(ctx, result); err != nil {
		return nil, fmt.Errorf("failed to save result: %w", err)
	}

	stats, err := s.statsRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user stats: %w", err)
	}

	stats.UpdateStreak(time.Now().UTC())

	if err := s.statsRepo.Upsert(ctx, stats); err != nil {
		return nil, fmt.Errorf("failed to update user stats: %w", err)
	}

	return result, nil
}

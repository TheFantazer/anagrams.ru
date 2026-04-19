package service

import (
	"context"
	"fmt"

	"github.com/TheFantazer/anagrams.ru/internal/dictionary"
	"github.com/TheFantazer/anagrams.ru/internal/domain"
	"github.com/TheFantazer/anagrams.ru/internal/repository"
	"github.com/google/uuid"
)

type GameService interface {
	CreateSession(ctx context.Context, language string, letterCount, timeLimit int, creatorID *uuid.UUID) (*domain.Session, error)
	GetSession(ctx context.Context, sessionID uuid.UUID) (*domain.Session, error)
	GetUserSessions(ctx context.Context, userID uuid.UUID, limit int) ([]*domain.Session, error)
	GetParticipatedSessions(ctx context.Context, userID uuid.UUID, limit int) ([]*domain.Session, error)
	SubmitResult(ctx context.Context, sessionID uuid.UUID, userID *uuid.UUID, playerName, fingerprint string, words []string, durationMs int) (*domain.Result, error)
	GetSessionResults(ctx context.Context, sessionID uuid.UUID, topN int) ([]*domain.Result, error)
}

type gameService struct {
	sessionRepo  repository.SessionRepository
	resultRepo   repository.ResultRepository
	dictionaries map[string]*dictionary.Trie
	letterGen    *dictionary.LetterGenerator
}

func NewGameService(
	sessionRepo repository.SessionRepository,
	resultRepo repository.ResultRepository,
	dictionaries map[string]*dictionary.Trie,
	letterGen *dictionary.LetterGenerator) GameService {
	return &gameService{
		sessionRepo:  sessionRepo,
		resultRepo:   resultRepo,
		dictionaries: dictionaries,
		letterGen:    letterGen,
	}
}

func (s *gameService) CreateSession(ctx context.Context, language string, letterCount, timeLimit int, creatorID *uuid.UUID) (*domain.Session, error) {
	dict, ok := s.dictionaries[language]
	if !ok {
		return nil, domain.ErrUnsupportedLanguage
	}

	letters := s.letterGen.GenerateFromDictionary(dict, letterCount)

	if letters == "" {
		letters = s.letterGen.GenerateBalancedLetters(language, letterCount)
	}

	validWords := dict.FindAllWords(letters)
	if len(validWords) == 0 {
		return s.CreateSession(ctx, language, letterCount, timeLimit, creatorID)
	}

	session, err := domain.NewSession(letters, language, timeLimit, letterCount, validWords)
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	session.CreatorID = creatorID

	if err := s.sessionRepo.Create(ctx, session); err != nil {
		return nil, fmt.Errorf("failed to save session: %w", err)
	}

	return session, nil
}

func (s *gameService) GetSession(ctx context.Context, sessionID uuid.UUID) (*domain.Session, error) {
	session, err := s.sessionRepo.GetByID(ctx, sessionID)
	if err != nil {
		return nil, err
	}
	return session, nil
}

func (s *gameService) GetUserSessions(ctx context.Context, userID uuid.UUID, limit int) ([]*domain.Session, error) {
	return s.sessionRepo.GetByCreatorID(ctx, userID, limit)
}

func (s *gameService) GetParticipatedSessions(ctx context.Context, userID uuid.UUID, limit int) ([]*domain.Session, error) {
	return s.sessionRepo.GetByParticipant(ctx, userID, limit)
}

func (s *gameService) SubmitResult(ctx context.Context, sessionID uuid.UUID, userID *uuid.UUID, playerName, fingerprint string, words []string, durationMs int) (*domain.Result, error) {
	session, err := s.sessionRepo.GetByID(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	if session.IsExpired() {
		return nil, domain.ErrSessionExpired
	}

	result, err := domain.NewResult(sessionID, userID, playerName, fingerprint, words, durationMs)
	if err != nil {
		return nil, fmt.Errorf("failed to create result: %w", err)
	}

	if err := result.ValidateWords(session); err != nil {
		return nil, err
	}

	if err := s.resultRepo.Create(ctx, result); err != nil {
		return nil, fmt.Errorf("failed to save result: %w", err)
	}

	return result, nil
}

func (s *gameService) GetSessionResults(ctx context.Context, sessionID uuid.UUID, topN int) ([]*domain.Result, error) {
	if topN > 0 {
		return s.resultRepo.GetTopBySessionID(ctx, sessionID, topN)
	}
	return s.resultRepo.GetBySessionID(ctx, sessionID)
}

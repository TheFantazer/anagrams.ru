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
	CreateSessionWithMode(ctx context.Context, language string, letterCount, timeLimit int, creatorID *uuid.UUID, inviteMode string, maxOpponents int) (*domain.Session, error)
	GetSession(ctx context.Context, sessionID uuid.UUID) (*domain.Session, error)
	GetUserSessions(ctx context.Context, userID uuid.UUID, limit int) ([]*domain.Session, error)
	GetParticipatedSessions(ctx context.Context, userID uuid.UUID, limit int) ([]*domain.Session, error)
	GetAllUserSessionsPaginated(ctx context.Context, userID uuid.UUID, page int, perPage int) (*repository.PaginatedSessions, error)
	SubmitResult(ctx context.Context, sessionID uuid.UUID, userID *uuid.UUID, playerName, fingerprint string, words []string, durationMs int) (*domain.Result, error)
	GetSessionResults(ctx context.Context, sessionID uuid.UUID, topN int) ([]*domain.Result, error)
	JoinSession(ctx context.Context, sessionID, userID uuid.UUID) error
	StartGame(ctx context.Context, sessionID, userID uuid.UUID) error
	CanJoinSession(ctx context.Context, sessionID, userID uuid.UUID) (bool, error)
}

type gameService struct {
	sessionRepo     repository.SessionRepository
	resultRepo      repository.ResultRepository
	participantRepo repository.SessionParticipantRepository
	dictionaries    map[string]*dictionary.Trie
	letterGen       *dictionary.LetterGenerator
}

func NewGameService(
	sessionRepo repository.SessionRepository,
	resultRepo repository.ResultRepository,
	participantRepo repository.SessionParticipantRepository,
	dictionaries map[string]*dictionary.Trie,
	letterGen *dictionary.LetterGenerator) GameService {
	return &gameService{
		sessionRepo:     sessionRepo,
		resultRepo:      resultRepo,
		participantRepo: participantRepo,
		dictionaries:    dictionaries,
		letterGen:       letterGen,
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

func (s *gameService) GetAllUserSessionsPaginated(ctx context.Context, userID uuid.UUID, page int, perPage int) (*repository.PaginatedSessions, error) {
	return s.sessionRepo.GetAllUserSessions(ctx, userID, page, perPage)
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

func (s *gameService) CreateSessionWithMode(ctx context.Context, language string, letterCount, timeLimit int, creatorID *uuid.UUID, inviteMode string, maxOpponents int) (*domain.Session, error) {
	dict, ok := s.dictionaries[language]
	if !ok {
		return nil, domain.ErrUnsupportedLanguage
	}

	letters := s.letterGen.GenerateFromDictionary(dict, letterCount)
	validWords := dict.FindAllWords(letters)

	session, err := domain.NewSession(letters, language, timeLimit, letterCount, validWords)
	if err != nil {
		return nil, err
	}

	session.CreatorID = creatorID
	session.InviteMode = inviteMode
	session.MaxOpponents = maxOpponents

	if err := s.sessionRepo.Create(ctx, session); err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	// Create participant record for creator
	if creatorID != nil {
		participant := domain.NewSessionParticipant(session.ID, *creatorID, "creator")
		participant.MarkAsStarted() // Creator auto-starts
		if err := s.participantRepo.Create(ctx, participant); err != nil {
			return nil, fmt.Errorf("failed to create creator participant: %w", err)
		}
	}

	return session, nil
}

func (s *gameService) CanJoinSession(ctx context.Context, sessionID, userID uuid.UUID) (bool, error) {
	session, err := s.sessionRepo.GetByID(ctx, sessionID)
	if err != nil {
		return false, err
	}

	// Check if session is expired
	if session.IsExpired() {
		return false, domain.ErrSessionExpired
	}

	// Check if user is already a participant
	participant, err := s.participantRepo.GetBySessionAndUser(ctx, sessionID, userID)
	if err != nil && err != repository.ErrNotFound {
		return false, err
	}
	if participant != nil {
		return true, nil // Already joined
	}

	// For "link" mode, check if max opponents reached
	if session.InviteMode == "link" {
		opponentCount, err := s.participantRepo.CountOpponents(ctx, sessionID)
		if err != nil {
			return false, err
		}
		if opponentCount >= session.MaxOpponents {
			return false, fmt.Errorf("session is full: max opponents reached")
		}
	}

	return true, nil
}

func (s *gameService) JoinSession(ctx context.Context, sessionID, userID uuid.UUID) error {
	canJoin, err := s.CanJoinSession(ctx, sessionID, userID)
	if err != nil {
		return err
	}
	if !canJoin {
		return fmt.Errorf("cannot join session")
	}

	// Check if already joined
	existing, err := s.participantRepo.GetBySessionAndUser(ctx, sessionID, userID)
	if err != nil && err != repository.ErrNotFound {
		return err
	}
	if existing != nil {
		return nil // Already joined, no error
	}

	// Create participant record
	participant := domain.NewSessionParticipant(sessionID, userID, "opponent")
	return s.participantRepo.Create(ctx, participant)
}

func (s *gameService) StartGame(ctx context.Context, sessionID, userID uuid.UUID) error {
	// Ensure user is a participant
	participant, err := s.participantRepo.GetBySessionAndUser(ctx, sessionID, userID)
	if err != nil {
		if err == repository.ErrNotFound {
			// Not a participant yet, join first
			if err := s.JoinSession(ctx, sessionID, userID); err != nil {
				return err
			}
		} else {
			return err
		}
	}

	// Mark as started
	if participant != nil && participant.HasStarted() {
		return nil // Already started, no error
	}

	return s.participantRepo.MarkAsStarted(ctx, sessionID, userID)
}

package mocks

import (
	"context"
	"sync"

	"github.com/TheFantazer/anagrams.ru/internal/domain"
	"github.com/TheFantazer/anagrams.ru/internal/repository"
	"github.com/google/uuid"
)

type MockSessionParticipantRepository struct {
	mu           sync.RWMutex
	participants map[string]*domain.SessionParticipant // key: "sessionID-userID"

	CreateFunc              func(ctx context.Context, participant *domain.SessionParticipant) error
	GetBySessionFunc        func(ctx context.Context, sessionID uuid.UUID) ([]*domain.SessionParticipant, error)
	GetBySessionAndUserFunc func(ctx context.Context, sessionID, userID uuid.UUID) (*domain.SessionParticipant, error)
	CountOpponentsFunc      func(ctx context.Context, sessionID uuid.UUID) (int, error)
	MarkAsStartedFunc       func(ctx context.Context, sessionID, userID uuid.UUID) error
	HasStartedFunc          func(ctx context.Context, sessionID, userID uuid.UUID) (bool, error)
}

func NewMockSessionParticipantRepository() *MockSessionParticipantRepository {
	return &MockSessionParticipantRepository{
		participants: make(map[string]*domain.SessionParticipant),
	}
}

func (m *MockSessionParticipantRepository) Create(ctx context.Context, participant *domain.SessionParticipant) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, participant)
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	key := participant.SessionID.String() + "-" + participant.UserID.String()
	m.participants[key] = participant
	return nil
}

func (m *MockSessionParticipantRepository) GetBySession(ctx context.Context, sessionID uuid.UUID) ([]*domain.SessionParticipant, error) {
	if m.GetBySessionFunc != nil {
		return m.GetBySessionFunc(ctx, sessionID)
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	var result []*domain.SessionParticipant
	for _, p := range m.participants {
		if p.SessionID == sessionID {
			result = append(result, p)
		}
	}
	return result, nil
}

func (m *MockSessionParticipantRepository) GetBySessionAndUser(ctx context.Context, sessionID, userID uuid.UUID) (*domain.SessionParticipant, error) {
	if m.GetBySessionAndUserFunc != nil {
		return m.GetBySessionAndUserFunc(ctx, sessionID, userID)
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	key := sessionID.String() + "-" + userID.String()
	participant, ok := m.participants[key]
	if !ok {
		return nil, repository.ErrNotFound
	}
	return participant, nil
}

func (m *MockSessionParticipantRepository) CountOpponents(ctx context.Context, sessionID uuid.UUID) (int, error) {
	if m.CountOpponentsFunc != nil {
		return m.CountOpponentsFunc(ctx, sessionID)
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	count := 0
	for _, p := range m.participants {
		if p.SessionID == sessionID {
			count++
		}
	}
	return count, nil
}

func (m *MockSessionParticipantRepository) MarkAsStarted(ctx context.Context, sessionID, userID uuid.UUID) error {
	if m.MarkAsStartedFunc != nil {
		return m.MarkAsStartedFunc(ctx, sessionID, userID)
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	key := sessionID.String() + "-" + userID.String()
	participant, ok := m.participants[key]
	if !ok {
		return repository.ErrNotFound
	}
	participant.MarkAsStarted()
	return nil
}

func (m *MockSessionParticipantRepository) HasStarted(ctx context.Context, sessionID, userID uuid.UUID) (bool, error) {
	if m.HasStartedFunc != nil {
		return m.HasStartedFunc(ctx, sessionID, userID)
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	key := sessionID.String() + "-" + userID.String()
	participant, ok := m.participants[key]
	if !ok {
		return false, repository.ErrNotFound
	}
	return participant.HasStarted(), nil
}

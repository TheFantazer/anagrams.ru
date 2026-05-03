package mocks

import (
	"context"
	"sync"
	"time"

	"github.com/TheFantazer/anagrams.ru/internal/domain"
	"github.com/TheFantazer/anagrams.ru/internal/repository"
	"github.com/google/uuid"
)

type MockSessionRepository struct {
	mu       sync.RWMutex
	sessions map[uuid.UUID]*domain.Session

	CreateFunc             func(ctx context.Context, session *domain.Session) error
	GetByIDFunc            func(ctx context.Context, id uuid.UUID) (*domain.Session, error)
	GetByCreatorIDFunc     func(ctx context.Context, creatorID uuid.UUID, limit int) ([]*domain.Session, error)
	GetByParticipantFunc   func(ctx context.Context, userID uuid.UUID, limit int) ([]*domain.Session, error)
	GetAllUserSessionsFunc func(ctx context.Context, userID uuid.UUID, page int, perPage int) (*repository.PaginatedSessions, error)
	DeleteFunc             func(ctx context.Context, id uuid.UUID) error
	DeleteExpiredFunc      func(ctx context.Context, before time.Time) (int64, error)
}

func NewMockSessionRepository() *MockSessionRepository {
	return &MockSessionRepository{
		sessions: make(map[uuid.UUID]*domain.Session),
	}
}

func (m *MockSessionRepository) Create(ctx context.Context, session *domain.Session) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, session)
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	m.sessions[session.ID] = session
	return nil
}

func (m *MockSessionRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Session, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(ctx, id)
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	session, ok := m.sessions[id]
	if !ok {
		return nil, repository.ErrNotFound
	}
	return session, nil
}

func (m *MockSessionRepository) GetByCreatorID(ctx context.Context, creatorID uuid.UUID, limit int) ([]*domain.Session, error) {
	if m.GetByCreatorIDFunc != nil {
		return m.GetByCreatorIDFunc(ctx, creatorID, limit)
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	// Simple mock implementation - returns empty result
	return []*domain.Session{}, nil
}

func (m *MockSessionRepository) GetByParticipant(ctx context.Context, userID uuid.UUID, limit int) ([]*domain.Session, error) {
	if m.GetByParticipantFunc != nil {
		return m.GetByParticipantFunc(ctx, userID, limit)
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	// Simple mock implementation - returns empty result
	return []*domain.Session{}, nil
}

func (m *MockSessionRepository) GetAllUserSessions(ctx context.Context, userID uuid.UUID, page int, perPage int) (*repository.PaginatedSessions, error) {
	if m.GetAllUserSessionsFunc != nil {
		return m.GetAllUserSessionsFunc(ctx, userID, page, perPage)
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	// Simple mock implementation - returns empty paginated result
	return &repository.PaginatedSessions{
		Sessions:   []*repository.SessionWithResults{},
		Total:      0,
		Page:       page,
		PerPage:    perPage,
		TotalPages: 0,
	}, nil
}

func (m *MockSessionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id)
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.sessions[id]; !ok {
		return repository.ErrNotFound
	}
	delete(m.sessions, id)
	return nil
}

func (m *MockSessionRepository) DeleteExpired(ctx context.Context, before time.Time) (int64, error) {
	if m.DeleteExpiredFunc != nil {
		return m.DeleteExpiredFunc(ctx, before)
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	var count int64
	for id, session := range m.sessions {
		if session.CreatedAt.Before(before) {
			delete(m.sessions, id)
			count++
		}
	}
	return count, nil
}

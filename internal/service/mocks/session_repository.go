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

	CreateFunc        func(ctx context.Context, session *domain.Session) error
	GetByIDFunc       func(ctx context.Context, id uuid.UUID) (*domain.Session, error)
	DeleteFunc        func(ctx context.Context, id uuid.UUID) error
	DeleteExpiredFunc func(ctx context.Context, before time.Time) (int64, error)
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

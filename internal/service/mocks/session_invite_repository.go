package mocks

import (
	"context"
	"sync"

	"github.com/TheFantazer/anagrams.ru/internal/domain"
	"github.com/google/uuid"
)

type MockSessionInviteRepository struct {
	mu      sync.RWMutex
	invites map[uuid.UUID]*domain.SessionInvite

	CreateFunc         func(ctx context.Context, invite *domain.SessionInvite) error
	GetByUserIDFunc    func(ctx context.Context, userID uuid.UUID) ([]*domain.SessionInvite, error)
	GetBySessionIDFunc func(ctx context.Context, sessionID uuid.UUID) ([]*domain.SessionInvite, error)
}

func NewMockSessionInviteRepository() *MockSessionInviteRepository {
	return &MockSessionInviteRepository{
		invites: make(map[uuid.UUID]*domain.SessionInvite),
	}
}

func (m *MockSessionInviteRepository) Create(ctx context.Context, invite *domain.SessionInvite) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, invite)
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	m.invites[invite.ID] = invite
	return nil
}

func (m *MockSessionInviteRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.SessionInvite, error) {
	if m.GetByUserIDFunc != nil {
		return m.GetByUserIDFunc(ctx, userID)
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	var result []*domain.SessionInvite
	for _, invite := range m.invites {
		if invite.ToUserID == userID {
			result = append(result, invite)
		}
	}
	return result, nil
}

func (m *MockSessionInviteRepository) GetBySessionID(ctx context.Context, sessionID uuid.UUID) ([]*domain.SessionInvite, error) {
	if m.GetBySessionIDFunc != nil {
		return m.GetBySessionIDFunc(ctx, sessionID)
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	var result []*domain.SessionInvite
	for _, invite := range m.invites {
		if invite.SessionID == sessionID {
			result = append(result, invite)
		}
	}
	return result, nil
}

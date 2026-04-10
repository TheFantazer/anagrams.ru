package mocks

import (
	"context"
	"sort"
	"sync"

	"github.com/TheFantazer/anagrams.ru/internal/domain"
	"github.com/TheFantazer/anagrams.ru/internal/repository"
	"github.com/google/uuid"
)

type MockResultRepository struct {
	mu      sync.RWMutex
	results map[uuid.UUID]*domain.Result

	fingerprints map[string]bool

	// Для контроля поведения в тестах
	CreateFunc            func(ctx context.Context, result *domain.Result) error
	GetByIDFunc           func(ctx context.Context, id uuid.UUID) (*domain.Result, error)
	GetBySessionIDFunc    func(ctx context.Context, sessionID uuid.UUID) ([]*domain.Result, error)
	GetTopBySessionIDFunc func(ctx context.Context, sessionID uuid.UUID, limit int) ([]*domain.Result, error)
}

func NewMockResultRepository() *MockResultRepository {
	return &MockResultRepository{
		results:      make(map[uuid.UUID]*domain.Result),
		fingerprints: make(map[string]bool),
	}
}

func (m *MockResultRepository) Create(ctx context.Context, result *domain.Result) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, result)
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	key := result.SessionID.String() + ":" + result.PlayerFingerprint
	if m.fingerprints[key] {
		return repository.ErrDuplicateResult
	}

	m.results[result.ID] = result
	m.fingerprints[key] = true
	return nil
}

func (m *MockResultRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Result, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(ctx, id)
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	result, ok := m.results[id]
	if !ok {
		return nil, repository.ErrNotFound
	}
	return result, nil
}

func (m *MockResultRepository) GetBySessionID(ctx context.Context, sessionID uuid.UUID) ([]*domain.Result, error) {
	if m.GetBySessionIDFunc != nil {
		return m.GetBySessionIDFunc(ctx, sessionID)
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	var results []*domain.Result
	for _, result := range m.results {
		if result.SessionID == sessionID {
			results = append(results, result)
		}
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].PlayedAt.After(results[j].PlayedAt)
	})

	return results, nil
}

func (m *MockResultRepository) GetTopBySessionID(ctx context.Context, sessionID uuid.UUID, limit int) ([]*domain.Result,
	error) {
	if m.GetTopBySessionIDFunc != nil {
		return m.GetTopBySessionIDFunc(ctx, sessionID, limit)
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	var results []*domain.Result
	for _, result := range m.results {
		if result.SessionID == sessionID {
			results = append(results, result)
		}
	}

	sort.Slice(results, func(i, j int) bool {
		if results[i].Score != results[j].Score {
			return results[i].Score > results[j].Score
		}
		return results[i].DurationMs < results[j].DurationMs
	})

	if len(results) > limit {
		results = results[:limit]
	}

	return results, nil
}

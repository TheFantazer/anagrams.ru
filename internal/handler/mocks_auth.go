package handler

import (
	"context"

	"github.com/TheFantazer/anagrams.ru/internal/domain"
	"github.com/TheFantazer/anagrams.ru/internal/repository"
	"github.com/TheFantazer/anagrams.ru/internal/service"
	"github.com/google/uuid"
)

type MockAuthService struct{}

func newMockAuthService() service.AuthService {
	return &MockAuthService{}
}

func (m *MockAuthService) Register(ctx context.Context, username, email, password string) (*domain.User, error) {
	return nil, nil
}

func (m *MockAuthService) Login(ctx context.Context, username, password string) (*domain.User, error) {
	return nil, nil
}

func (m *MockAuthService) GetUserByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	return nil, nil
}

func (m *MockAuthService) UpdateSettings(ctx context.Context, userID uuid.UUID, letterCount int, language string, timeLimit int) error {
	return nil
}

func (m *MockAuthService) UpdateUsername(ctx context.Context, userID uuid.UUID, newUsername string) error {
	return nil
}

func (m *MockAuthService) GetUserStats(ctx context.Context, userID uuid.UUID) (*repository.UserStats, error) {
	return nil, nil
}

func (m *MockAuthService) GetLeaderboard(ctx context.Context, period string, limit int) ([]*repository.LeaderboardEntry, error) {
	return nil, nil
}

func (m *MockAuthService) LoginOrRegisterWithOAuth(ctx context.Context, provider, oauthID, email, username string) (*domain.User, error) {
	return nil, nil
}

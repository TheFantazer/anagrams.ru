package service

import (
	"context"
	"fmt"
	"time"

	"github.com/TheFantazer/anagrams.ru/internal/domain"
	"github.com/TheFantazer/anagrams.ru/internal/repository"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(ctx context.Context, username, email, password string) (*domain.User, error)
	Login(ctx context.Context, username, password string) (*domain.User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
	UpdateSettings(ctx context.Context, userID uuid.UUID, letterCount int, language string, timeLimit int) error
	UpdateUsername(ctx context.Context, userID uuid.UUID, newUsername string) error
	GetUserStats(ctx context.Context, userID uuid.UUID) (*repository.UserStats, error)
	GetLeaderboard(ctx context.Context, period string, limit int) ([]*repository.LeaderboardEntry, error)
	LoginOrRegisterWithOAuth(ctx context.Context, provider, oauthID, email, username string) (*domain.User, error)
}

type authService struct {
	userRepo  repository.UserRepository
	statsRepo repository.StatsRepository
}

func NewAuthService(userRepo repository.UserRepository, statsRepo repository.StatsRepository) AuthService {
	return &authService{
		userRepo:  userRepo,
		statsRepo: statsRepo,
	}
}

func (s *authService) Register(ctx context.Context, username, email, password string) (*domain.User, error) {
	user, err := domain.NewUser(username, email, password)
	if err != nil {
		return nil, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}
	user.Password = string(hashedPassword)

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	user.Password = ""
	return user, nil
}

func (s *authService) Login(ctx context.Context, username, password string) (*domain.User, error) {
	user, err := s.userRepo.GetByUsername(ctx, username)
	if err != nil {
		if err == repository.ErrNotFound {
			user, err = s.userRepo.GetByEmail(ctx, username)
			if err != nil {
				if err == repository.ErrNotFound {
					return nil, domain.ErrInvalidCredentials
				}
				return nil, fmt.Errorf("failed to get user: %w", err)
			}
		} else {
			return nil, fmt.Errorf("failed to get user: %w", err)
		}
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, domain.ErrInvalidCredentials
	}

	user.Password = ""
	return user, nil
}

func (s *authService) GetUserByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	user.Password = ""
	return user, nil
}

func (s *authService) UpdateSettings(ctx context.Context, userID uuid.UUID, letterCount int, language string, timeLimit int) error {
	return s.userRepo.UpdateSettings(ctx, userID, letterCount, language, timeLimit)
}

func (s *authService) UpdateUsername(ctx context.Context, userID uuid.UUID, newUsername string) error {
	// Validate username
	if len(newUsername) < 3 || len(newUsername) > 30 {
		return domain.ErrInvalidUsername
	}

	// Get current user
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	// Check if username is the same
	if user.Username == newUsername {
		return fmt.Errorf("new username is the same as current username")
	}

	// Check cooldown period (2 weeks = 14 days)
	if user.UsernameChangedAt != nil {
		timeSinceLastChange := time.Since(*user.UsernameChangedAt)
		cooldownPeriod := 14 * 24 * time.Hour
		if timeSinceLastChange < cooldownPeriod {
			return domain.ErrUsernameChangeCooldown
		}
	}

	// Update username
	return s.userRepo.UpdateUsername(ctx, userID, newUsername)
}

func (s *authService) GetUserStats(ctx context.Context, userID uuid.UUID) (*repository.UserStats, error) {
	return s.statsRepo.GetUserStats(ctx, userID)
}

func (s *authService) GetLeaderboard(ctx context.Context, period string, limit int) ([]*repository.LeaderboardEntry, error) {
	return s.statsRepo.GetLeaderboard(ctx, period, limit)
}

func (s *authService) LoginOrRegisterWithOAuth(ctx context.Context, provider string, oauthID string, email string, username string) (*domain.User, error) {
	user, err := s.userRepo.GetByOAuthID(ctx, provider, oauthID)
	if err == nil {
		return user, nil
	}

	if email != "" {
		user, err = s.userRepo.GetByEmail(ctx, email)
		if err == nil {
			err = s.userRepo.LinkOAuth(ctx, user.ID, provider, oauthID)
			if err != nil {
				return nil, fmt.Errorf("failed to link oauth: %w", err)
			}
			return user, nil
		}
	}

	if username == "" {
		username = "user_" + oauthID
	}

	emailPtr := (*string)(nil)
	if email != "" {
		emailPtr = &email
	}

	newUser, err := domain.NewUser(username, "", "")
	if err != nil {
		return nil, fmt.Errorf("failed to create base user: %w", err)
	}

	newUser.Email = emailPtr
	newUser.OAuthProvider = &provider
	newUser.OAuthID = &oauthID

	err = s.userRepo.Create(ctx, newUser)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return newUser, nil
}

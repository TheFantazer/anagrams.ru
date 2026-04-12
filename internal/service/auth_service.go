package service

import (
	"context"
	"fmt"

	"github.com/TheFantazer/anagrams.ru/internal/domain"
	"github.com/TheFantazer/anagrams.ru/internal/repository"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(ctx context.Context, username, email, password string) (*domain.User, error)
	Login(ctx context.Context, username, password string) (*domain.User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
}

type authService struct {
	userRepo repository.UserRepository
}

func NewAuthService(userRepo repository.UserRepository) AuthService {
	return &authService{
		userRepo: userRepo,
	}
}

func (s *authService) Register(ctx context.Context, username, email, password string) (*domain.User, error) {
	// Создаём пользователя (валидация в domain)
	user, err := domain.NewUser(username, email, password)
	if err != nil {
		return nil, err
	}

	// Хешируем пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}
	user.Password = string(hashedPassword)

	// Сохраняем в БД
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	// Очищаем пароль перед возвратом
	user.Password = ""
	return user, nil
}

func (s *authService) Login(ctx context.Context, username, password string) (*domain.User, error) {
	// Получаем пользователя по username
	user, err := s.userRepo.GetByUsername(ctx, username)
	if err != nil {
		if err == repository.ErrNotFound {
			return nil, domain.ErrInvalidCredentials
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Проверяем пароль
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, domain.ErrInvalidCredentials
	}

	// Очищаем пароль перед возвратом
	user.Password = ""
	return user, nil
}

func (s *authService) GetUserByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Очищаем пароль
	user.Password = ""
	return user, nil
}

package service

import (
	"context"
	"fmt"

	"github.com/TheFantazer/anagrams.ru/internal/domain"
	"github.com/TheFantazer/anagrams.ru/internal/repository"
	"github.com/google/uuid"
)

type FriendService interface {
	SendFriendRequest(ctx context.Context, fromUserID, toUserID uuid.UUID) error
	AcceptFriendRequest(ctx context.Context, userID, requestID uuid.UUID) error
	RejectFriendRequest(ctx context.Context, userID, requestID uuid.UUID) error
	GetPendingRequests(ctx context.Context, userID uuid.UUID) ([]*domain.FriendRequest, error)
	GetSentRequests(ctx context.Context, userID uuid.UUID) ([]*domain.FriendRequest, error)
	GetFriends(ctx context.Context, userID uuid.UUID) ([]*domain.User, error)
	RemoveFriend(ctx context.Context, userID, friendID uuid.UUID) error
	SearchUsers(ctx context.Context, query string) ([]*domain.User, error)
	GetUserByID(ctx context.Context, userID uuid.UUID) (*domain.User, error)
	AreFriends(ctx context.Context, userID1, userID2 uuid.UUID) (bool, error)
}

type friendService struct {
	friendRepo repository.FriendRepository
	userRepo   repository.UserRepository
}

func (s *friendService) GetUserByID(ctx context.Context, userID uuid.UUID) (*domain.User, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	user.Password = ""

	return user, nil
}

func NewFriendService(friendRepo repository.FriendRepository, userRepo repository.UserRepository) FriendService {
	return &friendService{
		friendRepo: friendRepo,
		userRepo:   userRepo,
	}
}

func (s *friendService) SendFriendRequest(ctx context.Context, fromUserID, toUserID uuid.UUID) error {
	// Проверяем что оба пользователя существуют
	_, err := s.userRepo.GetByID(ctx, fromUserID)
	if err != nil {
		if err == repository.ErrNotFound {
			return domain.ErrUserNotFound
		}
		return fmt.Errorf("failed to get from user: %w", err)
	}

	_, err = s.userRepo.GetByID(ctx, toUserID)
	if err != nil {
		if err == repository.ErrNotFound {
			return domain.ErrUserNotFound
		}
		return fmt.Errorf("failed to get to user: %w", err)
	}

	// Проверяем что пользователи еще не друзья
	areFriends, err := s.friendRepo.AreFriends(ctx, fromUserID, toUserID)
	if err != nil {
		return fmt.Errorf("failed to check friendship: %w", err)
	}

	if areFriends {
		return domain.ErrAlreadyFriends
	}

	// Создаем запрос
	return s.friendRepo.CreateRequest(ctx, fromUserID, toUserID)
}

func (s *friendService) AcceptFriendRequest(ctx context.Context, userID, requestID uuid.UUID) error {
	// Получаем все pending requests для проверки прав
	requests, err := s.friendRepo.GetPendingRequests(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get pending requests: %w", err)
	}

	// Проверяем что запрос адресован этому пользователю
	var found bool
	for _, req := range requests {
		if req.ID == requestID {
			found = true
			break
		}
	}

	if !found {
		return repository.ErrNotFound
	}

	return s.friendRepo.AcceptRequest(ctx, requestID)
}

func (s *friendService) RejectFriendRequest(ctx context.Context, userID, requestID uuid.UUID) error {
	// Получаем все pending requests для проверки прав
	requests, err := s.friendRepo.GetPendingRequests(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get pending requests: %w", err)
	}

	// Проверяем что запрос адресован этому пользователю
	var found bool
	for _, req := range requests {
		if req.ID == requestID {
			found = true
			break
		}
	}

	if !found {
		return repository.ErrNotFound
	}

	return s.friendRepo.RejectRequest(ctx, requestID)
}

func (s *friendService) GetPendingRequests(ctx context.Context, userID uuid.UUID) ([]*domain.FriendRequest, error) {
	return s.friendRepo.GetPendingRequests(ctx, userID)
}

func (s *friendService) GetSentRequests(ctx context.Context, userID uuid.UUID) ([]*domain.FriendRequest, error) {
	return s.friendRepo.GetSendingRequests(ctx, userID)
}

func (s *friendService) GetFriends(ctx context.Context, userID uuid.UUID) ([]*domain.User, error) {
	users, err := s.friendRepo.GetFriends(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Очищаем пароли
	for _, user := range users {
		user.Password = ""
	}

	return users, nil
}

func (s *friendService) RemoveFriend(ctx context.Context, userID, friendID uuid.UUID) error {
	// Проверяем что они друзья
	areFriends, err := s.friendRepo.AreFriends(ctx, userID, friendID)
	if err != nil {
		return fmt.Errorf("failed to check friendship: %w", err)
	}

	if !areFriends {
		return domain.ErrNotFriends
	}

	return s.friendRepo.RemoveFriend(ctx, userID, friendID)
}

func (s *friendService) SearchUsers(ctx context.Context, query string) ([]*domain.User, error) {
	users, err := s.friendRepo.SearchUsers(ctx, query)
	if err != nil {
		return nil, err
	}

	// Очищаем пароли
	for _, user := range users {
		user.Password = ""
	}

	return users, nil
}

func (s *friendService) AreFriends(ctx context.Context, userID1, userID2 uuid.UUID) (bool, error) {
	return s.friendRepo.AreFriends(ctx, userID1, userID2)
}

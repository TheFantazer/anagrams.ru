package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/TheFantazer/anagrams.ru/internal/domain"
	"github.com/TheFantazer/anagrams.ru/internal/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"
)

type friendRepository struct {
	db *sqlx.DB
}

func (r *friendRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	query := `
		SELECT id, username, email, password, oauth_provider, oauth_id,
		       default_letter_count, default_language, default_time_limit,
		       created_at, updated_at
		FROM users
		WHERE id = $1
		LIMIT 1
	`

	var user domain.User
	err := r.db.GetContext(ctx, &user, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}

	return &user, nil
}

func NewFriendRepository(db *sqlx.DB) repository.FriendRepository {
	return &friendRepository{db: db}
}

func (r *friendRepository) CreateRequest(ctx context.Context, fromUserID uuid.UUID, toUserID uuid.UUID) error {
	friendRequest, err := domain.NewFriendRequest(fromUserID, toUserID)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO friend_requests (id, from_user_id, to_user_id, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (from_user_id, to_user_id)
		DO UPDATE SET
			status = 'pending',
			updated_at = NOW()
	`

	_, err = r.db.ExecContext(ctx, query,
		friendRequest.ID,
		friendRequest.FromUserID,
		friendRequest.ToUserID,
		friendRequest.Status,
		friendRequest.CreatedAt,
		friendRequest.UpdatedAt,
	)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23503" { // Foreign key violation
				return domain.ErrUserNotFound
			}
		}
		return fmt.Errorf("failed to create friend request: %w", err)
	}

	return nil
}

func (r *friendRepository) GetPendingRequests(ctx context.Context, userID uuid.UUID) ([]*domain.FriendRequest, error) {
	query := `
		SELECT id, from_user_id, to_user_id, status, created_at, updated_at
		FROM friend_requests
		WHERE to_user_id = $1 AND status = $2
		ORDER BY created_at DESC
	`

	requests := make([]*domain.FriendRequest, 0)
	err := r.db.SelectContext(ctx, &requests, query, userID, domain.FriendRequestStatusPending)
	if err != nil {
		return nil, fmt.Errorf("failed to get pending requests: %w", err)
	}

	return requests, nil
}

func (r *friendRepository) GetSendingRequests(ctx context.Context, userID uuid.UUID) ([]*domain.FriendRequest, error) {
	query := `
		SELECT id, from_user_id, to_user_id, status, created_at, updated_at
		FROM friend_requests
		WHERE from_user_id = $1
		ORDER BY created_at DESC
	`

	var requests []*domain.FriendRequest
	err := r.db.SelectContext(ctx, &requests, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get sending requests: %w", err)
	}

	return requests, nil
}

func (r *friendRepository) AcceptRequest(ctx context.Context, requestID uuid.UUID) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	// 1. Получаем request
	query := `
		SELECT id, from_user_id, to_user_id, status
		FROM friend_requests
		WHERE id = $1
	`

	var req domain.FriendRequest
	err = tx.GetContext(ctx, &req, query, requestID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return repository.ErrNotFound
		}
		return fmt.Errorf("failed to get friend request: %w", err)
	}

	// 2. Проверяем статус
	if !req.IsStatusPending() {
		return domain.ErrFriendRequestNotPending
	}

	// 3. Создаём friendships (bidirectional)
	friendship1, err := domain.NewFriendship(req.ToUserID, req.FromUserID)
	if err != nil {
		return err
	}

	friendship2, err := domain.NewFriendship(req.FromUserID, req.ToUserID)
	if err != nil {
		return err
	}

	insertQuery := `
		INSERT INTO friendships (id, user_id, friend_id, created_at)
		VALUES ($1, $2, $3, $4)
	`

	_, err = tx.ExecContext(ctx, insertQuery,
		friendship1.ID,
		friendship1.UserID,
		friendship1.FriendID,
		friendship1.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create friendship 1: %w", err)
	}

	_, err = tx.ExecContext(ctx, insertQuery,
		friendship2.ID,
		friendship2.UserID,
		friendship2.FriendID,
		friendship2.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create friendship 2: %w", err)
	}

	deleteQuery := `
		DELETE FROM friend_requests
		WHERE id = $1
	`

	_, err = tx.ExecContext(ctx, deleteQuery, requestID)
	if err != nil {
		return fmt.Errorf("failed to delete friend request: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *friendRepository) RejectRequest(ctx context.Context, requestID uuid.UUID) error {
	// Получаем friend request
	query := `
		SELECT id, from_user_id, to_user_id, status, created_at, updated_at
		FROM friend_requests
		WHERE id = $1
	`

	var req domain.FriendRequest
	err := r.db.GetContext(ctx, &req, query, requestID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return repository.ErrNotFound
		}
		return fmt.Errorf("failed to get friend request: %w", err)
	}

	// Проверяем что статус pending
	if !req.IsStatusPending() {
		return domain.ErrFriendRequestNotPending
	}

	// Обновляем статус на rejected
	updateQuery := `
		UPDATE friend_requests
		SET status = $1, updated_at = $2
		WHERE id = $3
	`

	result, err := r.db.ExecContext(ctx, updateQuery, domain.FriendRequestStatusRejected, time.Now().UTC(), requestID)
	if err != nil {
		return fmt.Errorf("failed to update friend request status: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return repository.ErrNotFound
	}

	return nil
}

func (r *friendRepository) GetFriends(ctx context.Context, userID uuid.UUID) ([]*domain.User, error) {
	query := `
		SELECT u.id, u.username, u.email, u.password, u.oauth_provider, u.oauth_id,
		       u.default_letter_count, u.default_language, u.default_time_limit,
		       u.created_at, u.updated_at
		FROM users u
		INNER JOIN friendships f ON u.id = f.friend_id
		WHERE f.user_id = $1
		ORDER BY u.username
	`

	var users []*domain.User
	err := r.db.SelectContext(ctx, &users, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get friends: %w", err)
	}

	return users, nil
}

func (r *friendRepository) RemoveFriend(ctx context.Context, userID uuid.UUID, friendID uuid.UUID) error {
	// Удаляем обе дружеские связи (bidirectional)
	query := `
		DELETE FROM friendships
		WHERE (user_id = $1 AND friend_id = $2) OR (user_id = $2 AND friend_id = $1)
	`

	result, err := r.db.ExecContext(ctx, query, userID, friendID)
	if err != nil {
		return fmt.Errorf("failed to remove friend: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return repository.ErrNotFound
	}

	return nil
}

func (r *friendRepository) SearchUsers(ctx context.Context, query string) ([]*domain.User, error) {
	searchQuery := `
		SELECT id, username, email, password, oauth_provider, oauth_id,
		       default_letter_count, default_language, default_time_limit,
		       created_at, updated_at
		FROM users
		WHERE username ILIKE $1 OR email ILIKE $1
		ORDER BY username
		LIMIT 50
	`

	// Добавляем % для ILIKE поиска
	searchPattern := "%" + query + "%"

	var users []*domain.User
	err := r.db.SelectContext(ctx, &users, searchQuery, searchPattern)
	if err != nil {
		return nil, fmt.Errorf("failed to search users: %w", err)
	}

	return users, nil
}

func (r *friendRepository) AreFriends(ctx context.Context, userID1 uuid.UUID, userID2 uuid.UUID) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1
			FROM friendships
			WHERE user_id = $1 AND friend_id = $2
		)
	`

	var exists bool
	err := r.db.GetContext(ctx, &exists, query, userID1, userID2)
	if err != nil {
		return false, fmt.Errorf("failed to check friendship: %w", err)
	}

	return exists, nil
}

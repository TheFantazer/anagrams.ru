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
)

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) repository.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *domain.User) error {
	query := `
		INSERT INTO users (id, username, email, password, oauth_provider, oauth_id,
		                   default_letter_count, default_language, default_time_limit,
		                   created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`

	_, err := r.db.ExecContext(ctx, query,
		user.ID,
		user.Username,
		user.Email,
		user.Password,
		user.OAuthProvider,
		user.OAuthID,
		user.DefaultLetterCount,
		user.DefaultLanguage,
		user.DefaultTimeLimit,
		user.CreatedAt,
		user.UpdatedAt,
	)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				if pgErr.ConstraintName == "users_username_key" {
					return domain.ErrUsernameTaken
				}
				if pgErr.ConstraintName == "users_email_key" {
					return domain.ErrEmailTaken
				}
			}
		}
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	query := `
		SELECT id, username, email, password, oauth_provider, oauth_id,
		       default_letter_count, default_language, default_time_limit,
		       username_changed_at, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	user := &domain.User{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.OAuthProvider,
		&user.OAuthID,
		&user.DefaultLetterCount,
		&user.DefaultLanguage,
		&user.DefaultTimeLimit,
		&user.UsernameChangedAt,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}

	return user, nil
}

func (r *userRepository) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	query := `
		SELECT id, username, email, password, oauth_provider, oauth_id,
		       default_letter_count, default_language, default_time_limit,
		       username_changed_at, created_at, updated_at
		FROM users
		WHERE username = $1
	`

	user := &domain.User{}
	err := r.db.QueryRowContext(ctx, query, username).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.OAuthProvider,
		&user.OAuthID,
		&user.DefaultLetterCount,
		&user.DefaultLanguage,
		&user.DefaultTimeLimit,
		&user.UsernameChangedAt,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get user by username: %w", err)
	}

	return user, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `
		SELECT id, username, email, password, oauth_provider, oauth_id,
		       default_letter_count, default_language, default_time_limit,
		       username_changed_at, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	user := &domain.User{}
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.OAuthProvider,
		&user.OAuthID,
		&user.DefaultLetterCount,
		&user.DefaultLanguage,
		&user.DefaultTimeLimit,
		&user.UsernameChangedAt,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return user, nil
}

func (r *userRepository) GetByOAuthID(ctx context.Context, provider, oauthID string) (*domain.User, error) {
	query := `
		SELECT id, username, email, password, oauth_provider, oauth_id,
		       default_letter_count, default_language, default_time_limit,
		       username_changed_at, created_at, updated_at
		FROM users
		WHERE oauth_provider = $1 AND oauth_id = $2
	`

	user := &domain.User{}
	err := r.db.QueryRowContext(ctx, query, provider, oauthID).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.OAuthProvider,
		&user.OAuthID,
		&user.DefaultLetterCount,
		&user.DefaultLanguage,
		&user.DefaultTimeLimit,
		&user.UsernameChangedAt,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get user by oauth id: %w", err)
	}

	return user, nil
}

func (r *userRepository) LinkOAuth(ctx context.Context, userID uuid.UUID, provider, oauthID string) error {
	query := `
		UPDATE users
		SET oauth_provider = $1, oauth_id = $2, updated_at = $3
		WHERE id = $4
	`

	result, err := r.db.ExecContext(ctx, query, provider, oauthID, time.Now().UTC(), userID)
	if err != nil {
		return fmt.Errorf("failed to link oauth: %w", err)
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

func (r *userRepository) Update(ctx context.Context, user *domain.User) error {
	query := `
		UPDATE users
		SET username = $1, email = $2, password = $3, oauth_provider = $4, oauth_id = $5,
		    default_letter_count = $6, default_language = $7, default_time_limit = $8,
		    updated_at = $9
		WHERE id = $10
	`

	result, err := r.db.ExecContext(ctx, query,
		user.Username,
		user.Email,
		user.Password,
		user.OAuthProvider,
		user.OAuthID,
		user.DefaultLetterCount,
		user.DefaultLanguage,
		user.DefaultTimeLimit,
		user.UpdatedAt,
		user.ID,
	)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				if pgErr.ConstraintName == "users_username_key" {
					return domain.ErrUsernameTaken
				}
				if pgErr.ConstraintName == "users_email_key" {
					return domain.ErrEmailTaken
				}
			}
		}
		return fmt.Errorf("failed to update user: %w", err)
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

func (r *userRepository) UpdateSettings(ctx context.Context, userID uuid.UUID, letterCount int, language string, timeLimit int) error {
	query := `
		UPDATE users
		SET default_letter_count = $1, default_language = $2, default_time_limit = $3, updated_at = $4
		WHERE id = $5
	`

	result, err := r.db.ExecContext(ctx, query, letterCount, language, timeLimit, time.Now().UTC(), userID)
	if err != nil {
		return fmt.Errorf("failed to update settings: %w", err)
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

func (r *userRepository) UpdateUsername(ctx context.Context, userID uuid.UUID, username string) error {
	query := `
		UPDATE users
		SET username = $1, username_changed_at = $2, updated_at = $3
		WHERE id = $4
	`

	now := time.Now().UTC()
	result, err := r.db.ExecContext(ctx, query, username, now, now, userID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" && pgErr.ConstraintName == "users_username_key" {
				return domain.ErrUsernameTaken
			}
		}
		return fmt.Errorf("failed to update username: %w", err)
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

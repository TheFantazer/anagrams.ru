package postgres

import (
	"context"
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"github.com/TheFantazer/anagrams.ru/internal/domain"
	"github.com/TheFantazer/anagrams.ru/internal/repository"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

// setupTestDB создаёт testcontainer с PostgreSQL и применяет миграции
func setupTestDB(t *testing.T) *sqlx.DB {
	ctx := context.Background()

	// Создаём PostgreSQL container
	postgresContainer, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("testuser"),
		postgres.WithPassword("testpass"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(60*time.Second)),
	)
	require.NoError(t, err, "Failed to start postgres container")

	// Cleanup при завершении теста
	t.Cleanup(func() {
		if err := testcontainers.TerminateContainer(postgresContainer); err != nil {
			t.Logf("Failed to terminate container: %v", err)
		}
	})

	// Получаем connection string
	connStr, err := postgresContainer.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	// Применяем миграции
	migrationsPath := filepath.Join("..", "..", "..", "migrations")
	m, err := migrate.New(
		fmt.Sprintf("file://%s", migrationsPath),
		connStr,
	)
	require.NoError(t, err, "Failed to create migrate instance")

	err = m.Up()
	require.NoError(t, err, "Failed to run migrations")

	// Открываем соединение с БД
	db, err := sqlx.Connect("postgres", connStr)
	require.NoError(t, err, "Failed to connect to database")

	t.Cleanup(func() {
		db.Close()
	})

	return db
}

// createTestSession создаёт тестовую сессию с дефолтными значениями
func createTestSession() *domain.Session {
	return &domain.Session{
		ID:          uuid.New(),
		Letters:     "абвгдеж",
		Language:    "ru",
		TimeLimit:   60,
		LetterCount: 7,
		ValidWords:  []string{"еда", "баг", "еж"},
		MaxScore:    300,
		CreatedAt:   time.Now().UTC().Truncate(time.Microsecond), // truncate для сравнения
	}
}

func TestSessionRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	repo := NewSessionRepository(db)
	ctx := context.Background()

	session := createTestSession()

	// Создаём сессию
	err := repo.Create(ctx, session)
	require.NoError(t, err)

	// Проверяем что сессия создана
	retrieved, err := repo.GetByID(ctx, session.ID)
	require.NoError(t, err)
	require.NotNil(t, retrieved)

	// Проверяем все поля
	assert.Equal(t, session.ID, retrieved.ID)
	assert.Equal(t, session.Letters, retrieved.Letters)
	assert.Equal(t, session.Language, retrieved.Language)
	assert.Equal(t, session.TimeLimit, retrieved.TimeLimit)
	assert.Equal(t, session.LetterCount, retrieved.LetterCount)
	assert.Equal(t, session.MaxScore, retrieved.MaxScore)
	assert.ElementsMatch(t, session.ValidWords, retrieved.ValidWords)
	assert.WithinDuration(t, session.CreatedAt, retrieved.CreatedAt, time.Second)
}

func TestSessionRepository_GetByID(t *testing.T) {
	db := setupTestDB(t)
	repo := NewSessionRepository(db)
	ctx := context.Background()

	t.Run("existing session", func(t *testing.T) {
		session := createTestSession()
		err := repo.Create(ctx, session)
		require.NoError(t, err)

		retrieved, err := repo.GetByID(ctx, session.ID)
		require.NoError(t, err)
		assert.Equal(t, session.ID, retrieved.ID)
	})

	t.Run("non-existing session", func(t *testing.T) {
		nonExistentID := uuid.New()
		retrieved, err := repo.GetByID(ctx, nonExistentID)
		assert.ErrorIs(t, err, repository.ErrNotFound)
		assert.Nil(t, retrieved)
	})
}

func TestSessionRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	repo := NewSessionRepository(db)
	ctx := context.Background()

	t.Run("delete existing session", func(t *testing.T) {
		session := createTestSession()
		err := repo.Create(ctx, session)
		require.NoError(t, err)

		// Удаляем
		err = repo.Delete(ctx, session.ID)
		require.NoError(t, err)

		// Проверяем что сессия удалена
		retrieved, err := repo.GetByID(ctx, session.ID)
		assert.ErrorIs(t, err, repository.ErrNotFound)
		assert.Nil(t, retrieved)
	})

	t.Run("delete non-existing session", func(t *testing.T) {
		nonExistentID := uuid.New()
		err := repo.Delete(ctx, nonExistentID)
		assert.ErrorIs(t, err, repository.ErrNotFound)
	})
}

func TestSessionRepository_DeleteExpired(t *testing.T) {
	db := setupTestDB(t)
	repo := NewSessionRepository(db)
	ctx := context.Background()

	// Создаём 3 сессии с разным временем создания
	now := time.Now().UTC()

	oldSession1 := createTestSession()
	oldSession1.CreatedAt = now.Add(-2 * time.Hour)
	err := repo.Create(ctx, oldSession1)
	require.NoError(t, err)

	oldSession2 := createTestSession()
	oldSession2.CreatedAt = now.Add(-1 * time.Hour)
	err = repo.Create(ctx, oldSession2)
	require.NoError(t, err)

	recentSession := createTestSession()
	recentSession.CreatedAt = now.Add(-10 * time.Minute)
	err = repo.Create(ctx, recentSession)
	require.NoError(t, err)

	// Удаляем сессии старше 30 минут
	cutoff := now.Add(-30 * time.Minute)
	deleted, err := repo.DeleteExpired(ctx, cutoff)
	require.NoError(t, err)
	assert.Equal(t, int64(2), deleted, "Should delete 2 old sessions")

	// Проверяем что старые сессии удалены
	_, err = repo.GetByID(ctx, oldSession1.ID)
	assert.ErrorIs(t, err, repository.ErrNotFound)

	_, err = repo.GetByID(ctx, oldSession2.ID)
	assert.ErrorIs(t, err, repository.ErrNotFound)

	// Проверяем что недавняя сессия осталась
	retrieved, err := repo.GetByID(ctx, recentSession.ID)
	require.NoError(t, err)
	assert.Equal(t, recentSession.ID, retrieved.ID)
}

func TestSessionRepository_ValidWordsJSON(t *testing.T) {
	db := setupTestDB(t)
	repo := NewSessionRepository(db)
	ctx := context.Background()

	session := createTestSession()
	session.ValidWords = []string{"абв", "где", "жзи", "клм"}

	err := repo.Create(ctx, session)
	require.NoError(t, err)

	retrieved, err := repo.GetByID(ctx, session.ID)
	require.NoError(t, err)

	// Проверяем что JSON правильно сериализовался и десериализовался
	assert.ElementsMatch(t, session.ValidWords, retrieved.ValidWords)
}

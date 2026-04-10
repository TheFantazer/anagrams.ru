package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/TheFantazer/anagrams.ru/internal/domain"
	"github.com/TheFantazer/anagrams.ru/internal/repository"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// createTestResult создаёт тестовый результат
func createTestResult(sessionID uuid.UUID) *domain.Result {
	return &domain.Result{
		ID:                uuid.New(),
		SessionID:         sessionID,
		PlayerName:        "TestPlayer",
		PlayerFingerprint: "fingerprint123",
		FoundWords:        []string{"еда", "баг"},
		WordCount:         2,
		Score:             200,
		DurationMs:        45000,
		PlayedAt:          time.Now().UTC().Truncate(time.Microsecond),
	}
}

func TestResultRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	sessionRepo := NewSessionRepository(db)
	resultRepo := NewResultRepository(db)
	ctx := context.Background()

	// Создаём сессию
	session := createTestSession()
	err := sessionRepo.Create(ctx, session)
	require.NoError(t, err)

	// Создаём результат
	result := createTestResult(session.ID)
	err = resultRepo.Create(ctx, result)
	require.NoError(t, err)

	// Проверяем что результат создан
	retrieved, err := resultRepo.GetByID(ctx, result.ID)
	require.NoError(t, err)
	require.NotNil(t, retrieved)

	// Проверяем все поля
	assert.Equal(t, result.ID, retrieved.ID)
	assert.Equal(t, result.SessionID, retrieved.SessionID)
	assert.Equal(t, result.PlayerName, retrieved.PlayerName)
	assert.Equal(t, result.PlayerFingerprint, retrieved.PlayerFingerprint)
	assert.ElementsMatch(t, result.FoundWords, retrieved.FoundWords)
	assert.Equal(t, result.WordCount, retrieved.WordCount)
	assert.Equal(t, result.Score, retrieved.Score)
	assert.Equal(t, result.DurationMs, retrieved.DurationMs)
	assert.WithinDuration(t, result.PlayedAt, retrieved.PlayedAt, time.Second)
}

func TestResultRepository_Create_Idempotency(t *testing.T) {
	db := setupTestDB(t)
	sessionRepo := NewSessionRepository(db)
	resultRepo := NewResultRepository(db)
	ctx := context.Background()

	// Создаём сессию
	session := createTestSession()
	err := sessionRepo.Create(ctx, session)
	require.NoError(t, err)

	// Создаём результат
	result := createTestResult(session.ID)
	err = resultRepo.Create(ctx, result)
	require.NoError(t, err)

	// Пытаемся создать дубликат (тот же session_id + player_fingerprint)
	duplicate := createTestResult(session.ID)
	duplicate.PlayerFingerprint = result.PlayerFingerprint // тот же fingerprint
	err = resultRepo.Create(ctx, duplicate)

	// Должна вернуться ошибка ErrDuplicateResult
	assert.ErrorIs(t, err, repository.ErrDuplicateResult)
}

func TestResultRepository_Create_ForeignKeyViolation(t *testing.T) {
	db := setupTestDB(t)
	resultRepo := NewResultRepository(db)
	ctx := context.Background()

	// Пытаемся создать результат для несуществующей сессии
	nonExistentSessionID := uuid.New()
	result := createTestResult(nonExistentSessionID)

	err := resultRepo.Create(ctx, result)

	// Должна вернуться ошибка ErrForeignKeyViolation
	assert.ErrorIs(t, err, repository.ErrForeignKeyViolation)
}

func TestResultRepository_GetByID(t *testing.T) {
	db := setupTestDB(t)
	sessionRepo := NewSessionRepository(db)
	resultRepo := NewResultRepository(db)
	ctx := context.Background()

	session := createTestSession()
	err := sessionRepo.Create(ctx, session)
	require.NoError(t, err)

	t.Run("existing result", func(t *testing.T) {
		result := createTestResult(session.ID)
		err := resultRepo.Create(ctx, result)
		require.NoError(t, err)

		retrieved, err := resultRepo.GetByID(ctx, result.ID)
		require.NoError(t, err)
		assert.Equal(t, result.ID, retrieved.ID)
	})

	t.Run("non-existing result", func(t *testing.T) {
		nonExistentID := uuid.New()
		retrieved, err := resultRepo.GetByID(ctx, nonExistentID)
		assert.ErrorIs(t, err, repository.ErrNotFound)
		assert.Nil(t, retrieved)
	})
}

func TestResultRepository_GetBySessionID(t *testing.T) {
	db := setupTestDB(t)
	sessionRepo := NewSessionRepository(db)
	resultRepo := NewResultRepository(db)
	ctx := context.Background()

	// Создаём сессию
	session := createTestSession()
	err := sessionRepo.Create(ctx, session)
	require.NoError(t, err)

	// Создаём несколько результатов
	result1 := createTestResult(session.ID)
	result1.PlayerFingerprint = "player1"
	result1.PlayedAt = time.Now().UTC().Add(-10 * time.Minute)
	err = resultRepo.Create(ctx, result1)
	require.NoError(t, err)

	result2 := createTestResult(session.ID)
	result2.PlayerFingerprint = "player2"
	result2.PlayedAt = time.Now().UTC().Add(-5 * time.Minute)
	err = resultRepo.Create(ctx, result2)
	require.NoError(t, err)

	result3 := createTestResult(session.ID)
	result3.PlayerFingerprint = "player3"
	result3.PlayedAt = time.Now().UTC()
	err = resultRepo.Create(ctx, result3)
	require.NoError(t, err)

	// Получаем все результаты для сессии
	results, err := resultRepo.GetBySessionID(ctx, session.ID)
	require.NoError(t, err)
	assert.Len(t, results, 3)

	// Проверяем сортировку по played_at DESC (последний первым)
	assert.Equal(t, result3.ID, results[0].ID)
	assert.Equal(t, result2.ID, results[1].ID)
	assert.Equal(t, result1.ID, results[2].ID)
}

func TestResultRepository_GetTopBySessionID(t *testing.T) {
	db := setupTestDB(t)
	sessionRepo := NewSessionRepository(db)
	resultRepo := NewResultRepository(db)
	ctx := context.Background()

	// Создаём сессию
	session := createTestSession()
	err := sessionRepo.Create(ctx, session)
	require.NoError(t, err)

	// Создаём результаты с разными score и duration
	result1 := createTestResult(session.ID)
	result1.PlayerFingerprint = "player1"
	result1.Score = 300
	result1.DurationMs = 50000
	err = resultRepo.Create(ctx, result1)
	require.NoError(t, err)

	result2 := createTestResult(session.ID)
	result2.PlayerFingerprint = "player2"
	result2.Score = 500
	result2.DurationMs = 40000
	err = resultRepo.Create(ctx, result2)
	require.NoError(t, err)

	result3 := createTestResult(session.ID)
	result3.PlayerFingerprint = "player3"
	result3.Score = 500
	result3.DurationMs = 30000 // Тот же score, но быстрее
	err = resultRepo.Create(ctx, result3)
	require.NoError(t, err)

	result4 := createTestResult(session.ID)
	result4.PlayerFingerprint = "player4"
	result4.Score = 100
	result4.DurationMs = 20000
	err = resultRepo.Create(ctx, result4)
	require.NoError(t, err)

	// Получаем топ-3
	top, err := resultRepo.GetTopBySessionID(ctx, session.ID, 3)
	require.NoError(t, err)
	assert.Len(t, top, 3)

	// Проверяем правильность сортировки:
	// 1. result3 (score=500, duration=30000) - лучший
	// 2. result2 (score=500, duration=40000) - тот же score, но медленнее
	// 3. result1 (score=300, duration=50000)
	assert.Equal(t, result3.ID, top[0].ID)
	assert.Equal(t, result2.ID, top[1].ID)
	assert.Equal(t, result1.ID, top[2].ID)
}

func TestResultRepository_CascadeDelete(t *testing.T) {
	db := setupTestDB(t)
	sessionRepo := NewSessionRepository(db)
	resultRepo := NewResultRepository(db)
	ctx := context.Background()

	// Создаём сессию
	session := createTestSession()
	err := sessionRepo.Create(ctx, session)
	require.NoError(t, err)

	// Создаём результат
	result := createTestResult(session.ID)
	err = resultRepo.Create(ctx, result)
	require.NoError(t, err)

	// Удаляем сессию
	err = sessionRepo.Delete(ctx, session.ID)
	require.NoError(t, err)

	// Проверяем что результат тоже удалился (CASCADE DELETE)
	retrieved, err := resultRepo.GetByID(ctx, result.ID)
	assert.ErrorIs(t, err, repository.ErrNotFound)
	assert.Nil(t, retrieved)
}

func TestResultRepository_FoundWordsJSON(t *testing.T) {
	db := setupTestDB(t)
	sessionRepo := NewSessionRepository(db)
	resultRepo := NewResultRepository(db)
	ctx := context.Background()

	session := createTestSession()
	err := sessionRepo.Create(ctx, session)
	require.NoError(t, err)

	result := createTestResult(session.ID)
	result.FoundWords = []string{"абв", "где", "жзи", "клм", "ноп"}

	err = resultRepo.Create(ctx, result)
	require.NoError(t, err)

	retrieved, err := resultRepo.GetByID(ctx, result.ID)
	require.NoError(t, err)

	// Проверяем что JSON правильно сериализовался и десериализовался
	assert.ElementsMatch(t, result.FoundWords, retrieved.FoundWords)
}

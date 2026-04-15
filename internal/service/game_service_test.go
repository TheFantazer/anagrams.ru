package service

import (
	"context"
	"testing"
	"time"
	"unicode/utf8"

	"github.com/TheFantazer/anagrams.ru/internal/dictionary"
	"github.com/TheFantazer/anagrams.ru/internal/domain"
	"github.com/TheFantazer/anagrams.ru/internal/repository"
	"github.com/TheFantazer/anagrams.ru/internal/service/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupTestDictionaries создает словари для тестов
func setupTestDictionaries() map[string]*dictionary.Trie {
	dictionaries := make(map[string]*dictionary.Trie)

	// Русский словарь
	ruTrie := dictionary.NewTrie()
	ruTrie.Insert("еда")
	ruTrie.Insert("еж")
	ruTrie.Insert("баг")
	ruTrie.Insert("вода")
	ruTrie.Insert("год")
	dictionaries["ru"] = ruTrie

	// Английский словарь
	enTrie := dictionary.NewTrie()
	enTrie.Insert("cat")
	enTrie.Insert("act")
	enTrie.Insert("bat")
	enTrie.Insert("tab")
	dictionaries["en"] = enTrie

	return dictionaries
}

// setupTestLetterGenerator создает генератор букв для тестов
func setupTestLetterGenerator() *dictionary.LetterGenerator {
	return dictionary.NewLetterGenerator()
}

func TestGameService_CreateSession(t *testing.T) {
	sessionRepo := mocks.NewMockSessionRepository()
	resultRepo := mocks.NewMockResultRepository()
	dictionaries := setupTestDictionaries()
	letterGen := setupTestLetterGenerator()

	service := NewGameService(sessionRepo, resultRepo, dictionaries, letterGen)
	ctx := context.Background()

	t.Run("success - creates session with valid words", func(t *testing.T) {
		session, err := service.CreateSession(ctx, "ru", 7, 60)
		require.NoError(t, err)
		require.NotNil(t, session)

		assert.NotEqual(t, uuid.Nil, session.ID)
		assert.Equal(t, "ru", session.Language)
		assert.Equal(t, 60, session.TimeLimit)
		assert.Equal(t, 7, session.LetterCount)
		assert.Equal(t, 7, utf8.RuneCountInString(session.Letters))
		assert.Greater(t, len(session.ValidWords), 0, "Should have at least one valid word")
		assert.Greater(t, session.MaxScore, 0, "Should have calculated max score")
	})

	t.Run("success - english dictionary", func(t *testing.T) {
		session, err := service.CreateSession(ctx, "en", 5, 120)
		require.NoError(t, err)
		require.NotNil(t, session)

		assert.Equal(t, "en", session.Language)
		assert.Equal(t, 120, session.TimeLimit)
		assert.Equal(t, 5, session.LetterCount)
	})

	t.Run("error - unsupported language", func(t *testing.T) {
		session, err := service.CreateSession(ctx, "fr", 7, 60)
		assert.ErrorIs(t, err, domain.ErrUnsupportedLanguage)
		assert.Nil(t, session)
	})

	t.Run("repository create error", func(t *testing.T) {
		repoWithError := mocks.NewMockSessionRepository()
		repoWithError.CreateFunc = func(ctx context.Context, session *domain.Session) error {
			return assert.AnError
		}

		serviceWithError := NewGameService(repoWithError, resultRepo, dictionaries, letterGen)
		session, err := serviceWithError.CreateSession(ctx, "ru", 7, 60)

		assert.Error(t, err)
		assert.Nil(t, session)
		assert.Contains(t, err.Error(), "failed to save session")
	})
}

func TestGameService_GetSession(t *testing.T) {
	sessionRepo := mocks.NewMockSessionRepository()
	resultRepo := mocks.NewMockResultRepository()
	dictionaries := setupTestDictionaries()
	letterGen := setupTestLetterGenerator()

	service := NewGameService(sessionRepo, resultRepo, dictionaries, letterGen)
	ctx := context.Background()

	t.Run("success - get existing session", func(t *testing.T) {
		// Создаем сессию
		created, err := service.CreateSession(ctx, "ru", 7, 60)
		require.NoError(t, err)

		// Получаем её
		retrieved, err := service.GetSession(ctx, created.ID)
		require.NoError(t, err)
		require.NotNil(t, retrieved)

		assert.Equal(t, created.ID, retrieved.ID)
		assert.Equal(t, created.Letters, retrieved.Letters)
		assert.Equal(t, created.Language, retrieved.Language)
	})

	t.Run("error - session not found", func(t *testing.T) {
		nonExistentID := uuid.New()
		session, err := service.GetSession(ctx, nonExistentID)

		assert.ErrorIs(t, err, repository.ErrNotFound)
		assert.Nil(t, session)
	})

	t.Run("repository error", func(t *testing.T) {
		repoWithError := mocks.NewMockSessionRepository()
		repoWithError.GetByIDFunc = func(ctx context.Context, id uuid.UUID) (*domain.Session, error) {
			return nil, assert.AnError
		}

		serviceWithError := NewGameService(repoWithError, resultRepo, dictionaries, letterGen)
		session, err := serviceWithError.GetSession(ctx, uuid.New())

		assert.Error(t, err)
		assert.Nil(t, session)
	})
}

func TestGameService_SubmitResult(t *testing.T) {
	sessionRepo := mocks.NewMockSessionRepository()
	resultRepo := mocks.NewMockResultRepository()
	dictionaries := setupTestDictionaries()
	letterGen := setupTestLetterGenerator()

	service := NewGameService(sessionRepo, resultRepo, dictionaries, letterGen)
	ctx := context.Background()

	t.Run("success - valid result", func(t *testing.T) {
		// Создаем сессию с известными буквами
		session := &domain.Session{
			ID:          uuid.New(),
			Letters:     "абвгдеж",
			Language:    "ru",
			TimeLimit:   60,
			LetterCount: 7,
			ValidWords:  []string{"еда", "баг", "еж"},
			MaxScore:    600,
			CreatedAt:   time.Now(),
		}
		err := sessionRepo.Create(ctx, session)
		require.NoError(t, err)

		// Отправляем результат
		result, err := service.SubmitResult(
			ctx,
			session.ID,
			nil, // userID для анонимного игрока
			"TestPlayer",
			"fingerprint123",
			[]string{"еда", "баг"},
			45000,
		)

		require.NoError(t, err)
		require.NotNil(t, result)

		assert.Equal(t, session.ID, result.SessionID)
		assert.Equal(t, "TestPlayer", result.PlayerName)
		assert.Equal(t, "fingerprint123", result.PlayerFingerprint)
		assert.Equal(t, 2, result.WordCount)
		assert.Equal(t, 200, result.Score) // еда (100) + баг (100)
		assert.Equal(t, 45000, result.DurationMs)
	})

	t.Run("error - session not found", func(t *testing.T) {
		nonExistentID := uuid.New()
		result, err := service.SubmitResult(
			ctx,
			nonExistentID,
			nil,
			"Player",
			"fp",
			[]string{"word"},
			1000,
		)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "failed to get session")
	})

	t.Run("error - session expired", func(t *testing.T) {
		// Создаем просроченную сессию (создана 2 часа назад, время игры 60 секунд)
		expiredSession := &domain.Session{
			ID:          uuid.New(),
			Letters:     "абвгдеж",
			Language:    "ru",
			TimeLimit:   60,
			LetterCount: 7,
			ValidWords:  []string{"еда"},
			MaxScore:    100,
			CreatedAt:   time.Now().Add(-2 * time.Hour),
		}
		err := sessionRepo.Create(ctx, expiredSession)
		require.NoError(t, err)

		result, err := service.SubmitResult(
			ctx,
			expiredSession.ID,
			nil,
			"Player",
			"fp",
			[]string{"еда"},
			1000,
		)

		assert.ErrorIs(t, err, domain.ErrSessionExpired)
		assert.Nil(t, result)
	})

	t.Run("error - invalid word (not in valid words)", func(t *testing.T) {
		session := &domain.Session{
			ID:          uuid.New(),
			Letters:     "абвгдеж",
			Language:    "ru",
			TimeLimit:   60,
			LetterCount: 7,
			ValidWords:  []string{"еда", "баг"},
			MaxScore:    200,
			CreatedAt:   time.Now(),
		}
		err := sessionRepo.Create(ctx, session)
		require.NoError(t, err)

		result, err := service.SubmitResult(
			ctx,
			session.ID,
			nil,
			"Player",
			"fp",
			[]string{"invalid"},
			1000,
		)

		assert.ErrorIs(t, err, domain.ErrInvalidWord)
		assert.Nil(t, result)
	})

	t.Run("error - duplicate result (same fingerprint)", func(t *testing.T) {
		session := &domain.Session{
			ID:          uuid.New(),
			Letters:     "абвгдеж",
			Language:    "ru",
			TimeLimit:   60,
			LetterCount: 7,
			ValidWords:  []string{"еда"},
			MaxScore:    100,
			CreatedAt:   time.Now(),
		}
		err := sessionRepo.Create(ctx, session)
		require.NoError(t, err)

		// Первая попытка - успешна
		_, err = service.SubmitResult(ctx, session.ID, nil, "Player1", "fp123", []string{"еда"}, 1000)
		require.NoError(t, err)

		// Вторая попытка с тем же fingerprint - ошибка
		result, err := service.SubmitResult(ctx, session.ID, nil, "Player1", "fp123", []string{"еда"}, 2000)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "failed to save result")
	})

	t.Run("error - empty fingerprint", func(t *testing.T) {
		session := &domain.Session{
			ID:          uuid.New(),
			Letters:     "абвгдеж",
			Language:    "ru",
			TimeLimit:   60,
			LetterCount: 7,
			ValidWords:  []string{"еда"},
			MaxScore:    100,
			CreatedAt:   time.Now(),
		}
		err := sessionRepo.Create(ctx, session)
		require.NoError(t, err)

		result, err := service.SubmitResult(ctx, session.ID, nil, "Player", "", []string{"еда"}, 1000)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "failed to create result")
	})
}

func TestGameService_GetSessionResults(t *testing.T) {
	sessionRepo := mocks.NewMockSessionRepository()
	resultRepo := mocks.NewMockResultRepository()
	dictionaries := setupTestDictionaries()
	letterGen := setupTestLetterGenerator()

	service := NewGameService(sessionRepo, resultRepo, dictionaries, letterGen)
	ctx := context.Background()

	t.Run("success - get all results", func(t *testing.T) {
		// Создаем сессию
		session := &domain.Session{
			ID:          uuid.New(),
			Letters:     "абвгдеж",
			Language:    "ru",
			TimeLimit:   60,
			LetterCount: 7,
			ValidWords:  []string{"еда", "баг", "еж"},
			MaxScore:    300,
			CreatedAt:   time.Now(),
		}
		err := sessionRepo.Create(ctx, session)
		require.NoError(t, err)

		// Добавляем 3 результата
		_, err = service.SubmitResult(ctx, session.ID, nil, "Player1", "fp1", []string{"еда"}, 30000)
		require.NoError(t, err)

		_, err = service.SubmitResult(ctx, session.ID, nil, "Player2", "fp2", []string{"еда", "баг"}, 40000)
		require.NoError(t, err)

		_, err = service.SubmitResult(ctx, session.ID, nil, "Player3", "fp3", []string{"еж"}, 20000)
		require.NoError(t, err)

		// Получаем все результаты (topN = 0)
		results, err := service.GetSessionResults(ctx, session.ID, 0)
		require.NoError(t, err)
		assert.Len(t, results, 3)

		// Проверяем сортировку по played_at DESC (последний первым)
		assert.Equal(t, "Player3", results[0].PlayerName)
		assert.Equal(t, "Player2", results[1].PlayerName)
		assert.Equal(t, "Player1", results[2].PlayerName)
	})

	t.Run("success - get top N results", func(t *testing.T) {
		// Создаем сессию
		session := &domain.Session{
			ID:          uuid.New(),
			Letters:     "абвгдеж",
			Language:    "ru",
			TimeLimit:   60,
			LetterCount: 7,
			ValidWords:  []string{"еда", "баг", "еж"},
			MaxScore:    300,
			CreatedAt:   time.Now(),
		}
		err := sessionRepo.Create(ctx, session)
		require.NoError(t, err)

		// Добавляем результаты с разными очками
		_, err = service.SubmitResult(ctx, session.ID, nil, "Player1", "fp1", []string{"еда"}, 30000)
		require.NoError(t, err)

		_, err = service.SubmitResult(ctx, session.ID, nil, "Player2", "fp2", []string{"еда", "баг"}, 40000)
		require.NoError(t, err)

		_, err = service.SubmitResult(ctx, session.ID, nil, "Player3", "fp3", []string{"еда", "баг", "еж"}, 20000)
		require.NoError(t, err)

		// Получаем топ-2
		results, err := service.GetSessionResults(ctx, session.ID, 2)
		require.NoError(t, err)
		assert.Len(t, results, 2)

		// Проверяем сортировку по score DESC, duration_ms ASC
		assert.Equal(t, "Player3", results[0].PlayerName) // 300 очков, 20000ms
		assert.Equal(t, "Player2", results[1].PlayerName) // 200 очков, 40000ms
	})

	t.Run("empty results - no submissions yet", func(t *testing.T) {
		// Создаем сессию без результатов
		session := &domain.Session{
			ID:          uuid.New(),
			Letters:     "абвгдеж",
			Language:    "ru",
			TimeLimit:   60,
			LetterCount: 7,
			ValidWords:  []string{"еда"},
			MaxScore:    100,
			CreatedAt:   time.Now(),
		}
		err := sessionRepo.Create(ctx, session)
		require.NoError(t, err)

		results, err := service.GetSessionResults(ctx, session.ID, 0)
		require.NoError(t, err)
		assert.Empty(t, results)
	})

	t.Run("repository error", func(t *testing.T) {
		repoWithError := mocks.NewMockResultRepository()
		repoWithError.GetBySessionIDFunc = func(ctx context.Context, sessionID uuid.UUID) ([]*domain.Result, error) {
			return nil, assert.AnError
		}

		serviceWithError := NewGameService(sessionRepo, repoWithError, dictionaries, letterGen)
		results, err := serviceWithError.GetSessionResults(ctx, uuid.New(), 0)

		assert.Error(t, err)
		assert.Nil(t, results)
	})
}

func TestGameService_CreateSession_WithShortWords(t *testing.T) {
	sessionRepo := mocks.NewMockSessionRepository()
	resultRepo := mocks.NewMockResultRepository()

	// Словарь с короткими словами - высокая вероятность найти что-то
	dictionaries := make(map[string]*dictionary.Trie)
	trie := dictionary.NewTrie()
	trie.Insert("а")
	trie.Insert("я")
	trie.Insert("и")
	trie.Insert("о")
	trie.Insert("у")
	trie.Insert("е")
	trie.Insert("ее")
	trie.Insert("еда")
	trie.Insert("он")
	trie.Insert("она")
	trie.Insert("оно")
	dictionaries["ru"] = trie

	letterGen := setupTestLetterGenerator()

	service := NewGameService(sessionRepo, resultRepo, dictionaries, letterGen)
	ctx := context.Background()

	// С таким богатым словарем должен найти валидные слова
	session, err := service.CreateSession(ctx, "ru", 7, 60)
	require.NoError(t, err)
	require.NotNil(t, session)

	// Проверяем что нашел хотя бы одно валидное слово
	assert.Greater(t, len(session.ValidWords), 0)
}

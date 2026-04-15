package handler

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/TheFantazer/anagrams.ru/internal/dictionary"
	"github.com/TheFantazer/anagrams.ru/internal/domain"
	"github.com/TheFantazer/anagrams.ru/internal/service"
	"github.com/TheFantazer/anagrams.ru/internal/service/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupTestService создает тестовый сервис с моками
func setupTestService() service.GameService {
	sessionRepo := mocks.NewMockSessionRepository()
	resultRepo := mocks.NewMockResultRepository()

	dictionaries := make(map[string]*dictionary.Trie)
	ruTrie := dictionary.NewTrie()
	ruTrie.Insert("еда")
	ruTrie.Insert("еж")
	ruTrie.Insert("баг")
	dictionaries["ru"] = ruTrie

	enTrie := dictionary.NewTrie()
	enTrie.Insert("cat")
	enTrie.Insert("act")
	enTrie.Insert("bat")
	dictionaries["en"] = enTrie

	letterGen := dictionary.NewLetterGenerator()

	return service.NewGameService(sessionRepo, resultRepo, dictionaries, letterGen)
}

func setupTestRouter() http.Handler {
	gameService := setupTestService()
	authService := newMockAuthService()
	logger := slog.New(slog.NewJSONHandler(io.Discard, nil))
	return NewRouter(gameService, authService, logger)
}

func TestCreateSession_Success(t *testing.T) {
	router := setupTestRouter()

	reqBody := CreateSessionRequest{
		Language:    "ru",
		LetterCount: 7,
		TimeLimit:   60,
	}

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/sessions", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)

	var resp SessionResponse
	err := json.NewDecoder(rec.Body).Decode(&resp)
	require.NoError(t, err)

	assert.NotEqual(t, uuid.Nil, resp.ID)
	assert.Equal(t, "ru", resp.Language)
	assert.Equal(t, 60, resp.TimeLimit)
	assert.Equal(t, 7, resp.LetterCount)
	assert.Greater(t, resp.MaxScore, 0)
}

func TestCreateSession_InvalidLanguage(t *testing.T) {
	router := setupTestRouter()

	reqBody := CreateSessionRequest{
		Language:    "fr",
		LetterCount: 7,
		TimeLimit:   60,
	}

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/sessions", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var errResp ErrorResponse
	err := json.NewDecoder(rec.Body).Decode(&errResp)
	require.NoError(t, err)
	assert.Equal(t, "invalid_language", errResp.Error)
}

func TestCreateSession_InvalidLetterCount(t *testing.T) {
	router := setupTestRouter()

	tests := []struct {
		name        string
		letterCount int
	}{
		{"too small", 4},
		{"too large", 11},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqBody := CreateSessionRequest{
				Language:    "ru",
				LetterCount: tt.letterCount,
				TimeLimit:   60,
			}

			body, _ := json.Marshal(reqBody)
			req := httptest.NewRequest(http.MethodPost, "/api/v1/sessions", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			assert.Equal(t, http.StatusBadRequest, rec.Code)

			var errResp ErrorResponse
			err := json.NewDecoder(rec.Body).Decode(&errResp)
			require.NoError(t, err)
			assert.Equal(t, "invalid_letter_count", errResp.Error)
		})
	}
}

func TestCreateSession_InvalidTimeLimit(t *testing.T) {
	router := setupTestRouter()

	tests := []struct {
		name      string
		timeLimit int
	}{
		{"too small", 29},
		{"too large", 301},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqBody := CreateSessionRequest{
				Language:    "ru",
				LetterCount: 7,
				TimeLimit:   tt.timeLimit,
			}

			body, _ := json.Marshal(reqBody)
			req := httptest.NewRequest(http.MethodPost, "/api/v1/sessions", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			assert.Equal(t, http.StatusBadRequest, rec.Code)

			var errResp ErrorResponse
			err := json.NewDecoder(rec.Body).Decode(&errResp)
			require.NoError(t, err)
			assert.Equal(t, "invalid_time_limit", errResp.Error)
		})
	}
}

func TestCreateSession_InvalidJSON(t *testing.T) {
	router := setupTestRouter()

	req := httptest.NewRequest(http.MethodPost, "/api/v1/sessions", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var errResp ErrorResponse
	json.NewDecoder(rec.Body).Decode(&errResp)
	assert.Equal(t, "invalid_request", errResp.Error)
}

func TestGetSession_Success(t *testing.T) {
	router := setupTestRouter()

	// Сначала создаем сессию
	createReq := CreateSessionRequest{
		Language:    "ru",
		LetterCount: 7,
		TimeLimit:   60,
	}
	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/sessions", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	var createResp SessionResponse
	json.NewDecoder(rec.Body).Decode(&createResp)

	// Теперь получаем сессию
	req = httptest.NewRequest(http.MethodGet, "/api/v1/sessions/"+createResp.ID.String(), nil)
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var getResp SessionResponse
	json.NewDecoder(rec.Body).Decode(&getResp)
	assert.Equal(t, createResp.ID, getResp.ID)
	assert.Equal(t, createResp.Letters, getResp.Letters)
}

func TestGetSession_NotFound(t *testing.T) {
	router := setupTestRouter()

	randomID := uuid.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/sessions/"+randomID.String(), nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)

	var errResp ErrorResponse
	json.NewDecoder(rec.Body).Decode(&errResp)
	assert.Equal(t, "not_found", errResp.Error)
}

func TestGetSession_InvalidUUID(t *testing.T) {
	router := setupTestRouter()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/sessions/invalid-uuid", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var errResp ErrorResponse
	json.NewDecoder(rec.Body).Decode(&errResp)
	assert.Equal(t, "invalid_uuid", errResp.Error)
}

func TestSubmitResult_Success(t *testing.T) {
	// Создаем service с известными словами
	sessionRepo := mocks.NewMockSessionRepository()
	resultRepo := mocks.NewMockResultRepository()

	dictionaries := make(map[string]*dictionary.Trie)
	ruTrie := dictionary.NewTrie()
	ruTrie.Insert("еда")
	ruTrie.Insert("баг")
	ruTrie.Insert("еж")
	dictionaries["ru"] = ruTrie

	letterGen := dictionary.NewLetterGenerator()
	gameService := service.NewGameService(sessionRepo, resultRepo, dictionaries, letterGen)
	authService := newMockAuthService()
	logger := slog.New(slog.NewJSONHandler(io.Discard, nil))
	router := NewRouter(gameService, authService, logger)

	// Создаем сессию с известными буквами напрямую в репозитории
	knownSession := &domain.Session{
		ID:          uuid.New(),
		Letters:     "абвгдеж",
		Language:    "ru",
		TimeLimit:   60,
		LetterCount: 7,
		ValidWords:  []string{"еда", "баг", "еж"},
		MaxScore:    300,
		CreatedAt:   time.Now(),
	}
	err := sessionRepo.Create(nil, knownSession)
	require.NoError(t, err)

	// Отправляем результат
	submitReq := SubmitResultRequest{
		PlayerName:  "TestPlayer",
		Fingerprint: "test-fp-123",
		FoundWords:  []string{"еда", "баг"},
		DurationMs:  45000,
	}

	body, _ := json.Marshal(submitReq)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/sessions/"+knownSession.ID.String()+"/results", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)

	var result ResultResponse
	json.NewDecoder(rec.Body).Decode(&result)
	assert.Equal(t, "TestPlayer", result.PlayerName)
	assert.Equal(t, 2, result.WordCount)
	assert.Equal(t, 200, result.Score)
}

func TestSubmitResult_InvalidWord(t *testing.T) {
	router := setupTestRouter()

	// Создаем сессию
	createReq := CreateSessionRequest{
		Language:    "ru",
		LetterCount: 7,
		TimeLimit:   60,
	}
	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/sessions", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	var session SessionResponse
	json.NewDecoder(rec.Body).Decode(&session)

	// Отправляем результат с невалидным словом
	submitReq := SubmitResultRequest{
		PlayerName:  "TestPlayer",
		Fingerprint: "test-fp-456",
		FoundWords:  []string{"invalid"},
		DurationMs:  45000,
	}

	body, _ = json.Marshal(submitReq)
	req = httptest.NewRequest(http.MethodPost, "/api/v1/sessions/"+session.ID.String()+"/results", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var errResp ErrorResponse
	json.NewDecoder(rec.Body).Decode(&errResp)
	assert.Equal(t, "invalid_word", errResp.Error)
}

func TestSubmitResult_DuplicateFingerprint(t *testing.T) {
	router := setupTestRouter()

	// Создаем сессию
	createReq := CreateSessionRequest{
		Language:    "ru",
		LetterCount: 7,
		TimeLimit:   60,
	}
	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/sessions", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	var session SessionResponse
	json.NewDecoder(rec.Body).Decode(&session)

	// Первая отправка
	submitReq := SubmitResultRequest{
		PlayerName:  "TestPlayer",
		Fingerprint: "duplicate-fp",
		FoundWords:  []string{"еда"},
		DurationMs:  45000,
	}

	body, _ = json.Marshal(submitReq)
	req = httptest.NewRequest(http.MethodPost, "/api/v1/sessions/"+session.ID.String()+"/results", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusCreated, rec.Code)

	// Вторая отправка с тем же fingerprint
	req = httptest.NewRequest(http.MethodPost, "/api/v1/sessions/"+session.ID.String()+"/results", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusConflict, rec.Code)

	var errResp ErrorResponse
	json.NewDecoder(rec.Body).Decode(&errResp)
	assert.Equal(t, "duplicate_result", errResp.Error)
}

func TestSubmitResult_SessionExpired(t *testing.T) {
	// Создаем service с mock репозиторием
	sessionRepo := mocks.NewMockSessionRepository()
	resultRepo := mocks.NewMockResultRepository()

	dictionaries := make(map[string]*dictionary.Trie)
	ruTrie := dictionary.NewTrie()
	ruTrie.Insert("еда")
	dictionaries["ru"] = ruTrie

	letterGen := dictionary.NewLetterGenerator()
	gameService := service.NewGameService(sessionRepo, resultRepo, dictionaries, letterGen)
	authService := newMockAuthService()
	logger := slog.New(slog.NewJSONHandler(io.Discard, nil))
	router := NewRouter(gameService, authService, logger)

	// Создаем просроченную сессию напрямую в репозитории
	expiredSession := &domain.Session{
		ID:          uuid.New(),
		Letters:     "абвгдеж",
		Language:    "ru",
		TimeLimit:   60,
		LetterCount: 7,
		ValidWords:  []string{"еда"},
		MaxScore:    100,
		CreatedAt:   time.Now().Add(-2 * time.Hour), // Сессия создана 2 часа назад
	}
	err := sessionRepo.Create(nil, expiredSession)
	require.NoError(t, err)

	// Пытаемся отправить результат
	submitReq := SubmitResultRequest{
		PlayerName:  "TestPlayer",
		Fingerprint: "test-fp",
		FoundWords:  []string{"еда"},
		DurationMs:  45000,
	}

	body, _ := json.Marshal(submitReq)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/sessions/"+expiredSession.ID.String()+"/results", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var errResp ErrorResponse
	json.NewDecoder(rec.Body).Decode(&errResp)
	assert.Equal(t, "session_expired", errResp.Error)
}

func TestGetSessionResults_Success(t *testing.T) {
	router := setupTestRouter()

	// Создаем сессию
	createReq := CreateSessionRequest{
		Language:    "ru",
		LetterCount: 7,
		TimeLimit:   60,
	}
	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/sessions", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	var session SessionResponse
	json.NewDecoder(rec.Body).Decode(&session)

	// Отправляем несколько результатов
	for i := 0; i < 3; i++ {
		submitReq := SubmitResultRequest{
			PlayerName:  "Player" + string(rune('1'+i)),
			Fingerprint: "fp-" + string(rune('1'+i)),
			FoundWords:  []string{"еда"},
			DurationMs:  45000,
		}

		body, _ := json.Marshal(submitReq)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/sessions/"+session.ID.String()+"/results", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)
	}

	// Получаем все результаты
	req = httptest.NewRequest(http.MethodGet, "/api/v1/sessions/"+session.ID.String()+"/results", nil)
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var results []ResultResponse
	json.NewDecoder(rec.Body).Decode(&results)
	assert.Len(t, results, 3)
}

func TestGetSessionResults_WithTopFilter(t *testing.T) {
	router := setupTestRouter()

	// Создаем сессию
	createReq := CreateSessionRequest{
		Language:    "ru",
		LetterCount: 7,
		TimeLimit:   60,
	}
	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/sessions", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	var session SessionResponse
	json.NewDecoder(rec.Body).Decode(&session)

	// Отправляем 5 результатов
	for i := 0; i < 5; i++ {
		submitReq := SubmitResultRequest{
			PlayerName:  "Player" + string(rune('1'+i)),
			Fingerprint: "fp-" + string(rune('1'+i)),
			FoundWords:  []string{"еда"},
			DurationMs:  45000,
		}

		body, _ := json.Marshal(submitReq)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/sessions/"+session.ID.String()+"/results", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)
	}

	// Получаем топ-2
	req = httptest.NewRequest(http.MethodGet, "/api/v1/sessions/"+session.ID.String()+"/results?top=2", nil)
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var results []ResultResponse
	json.NewDecoder(rec.Body).Decode(&results)
	assert.Len(t, results, 2)
}

func TestGetSessionResults_EmptyList(t *testing.T) {
	router := setupTestRouter()

	// Создаем сессию без результатов
	createReq := CreateSessionRequest{
		Language:    "ru",
		LetterCount: 7,
		TimeLimit:   60,
	}
	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/sessions", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	var session SessionResponse
	json.NewDecoder(rec.Body).Decode(&session)

	// Получаем результаты
	req = httptest.NewRequest(http.MethodGet, "/api/v1/sessions/"+session.ID.String()+"/results", nil)
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var results []ResultResponse
	json.NewDecoder(rec.Body).Decode(&results)
	assert.Empty(t, results)
}

func TestCORSHeaders(t *testing.T) {
	router := setupTestRouter()

	req := httptest.NewRequest(http.MethodOptions, "/api/v1/sessions", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNoContent, rec.Code)
	assert.Equal(t, "*", rec.Header().Get("Access-Control-Allow-Origin"))
	assert.Contains(t, rec.Header().Get("Access-Control-Allow-Methods"), "POST")
	assert.Contains(t, rec.Header().Get("Access-Control-Allow-Headers"), "Content-Type")
}

func TestRequestIDHeader(t *testing.T) {
	router := setupTestRouter()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/sessions/"+uuid.New().String(), nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	requestID := rec.Header().Get("X-Request-ID")
	assert.NotEmpty(t, requestID)

	// Проверяем что это валидный UUID
	_, err := uuid.Parse(requestID)
	assert.NoError(t, err)
}

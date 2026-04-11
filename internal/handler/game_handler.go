package handler

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/TheFantazer/anagrams.ru/internal/domain"
	"github.com/TheFantazer/anagrams.ru/internal/repository"
	"github.com/TheFantazer/anagrams.ru/internal/service"
	"github.com/google/uuid"
)

type GameHandler struct {
	service service.GameService
	logger  *slog.Logger
}

func NewGameHandler(service service.GameService, logger *slog.Logger) *GameHandler {
	return &GameHandler{
		service: service,
		logger:  logger,
	}
}

func (h *GameHandler) CreateSession(w http.ResponseWriter, r *http.Request) {
	var req CreateSessionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid_request", "Invalid JSON format")
		return
	}
	if req.Language != "ru" && req.Language != "en" {
		respondError(w, http.StatusBadRequest, "invalid_language", "Language must be 'ru' or 'en'")
		return
	}
	if req.LetterCount < 5 || req.LetterCount > 10 {
		respondError(w, http.StatusBadRequest, "invalid_letter_count", "Letter count must be between 5 and 10")
		return
	}
	if req.TimeLimit < 30 || req.TimeLimit > 300 {
		respondError(w, http.StatusBadRequest, "invalid_time_limit", "Time limit must be between 30 and 300 seconds")
		return
	}

	session, err := h.service.CreateSession(r.Context(), req.Language, req.LetterCount, req.TimeLimit)
	if err != nil {
		status, errCode, message := mapDomainError(err)
		h.logger.Error("Failed to create session", slog.String("error", err.Error()))
		respondError(w, status, errCode, message)
		return
	}

	response := SessionResponse{
		ID:          session.ID,
		Letters:     session.Letters,
		Language:    session.Language,
		TimeLimit:   session.TimeLimit,
		LetterCount: session.LetterCount,
		MaxScore:    session.MaxScore,
		CreatedAt:   session.CreatedAt,
	}

	respondJSON(w, http.StatusCreated, response)
}

func (h *GameHandler) GetSession(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	if idStr == "" {
		respondError(w, http.StatusBadRequest, "missing_id", "Session ID is required")
		return
	}

	sessionID, err := uuid.Parse(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid_uuid", "Invalid session ID format")
		return
	}

	session, err := h.service.GetSession(r.Context(), sessionID)
	if err != nil {
		status, errCode, message := mapDomainError(err)
		respondError(w, status, errCode, message)
		return
	}

	response := SessionResponse{
		ID:          session.ID,
		Letters:     session.Letters,
		Language:    session.Language,
		TimeLimit:   session.TimeLimit,
		LetterCount: session.LetterCount,
		MaxScore:    session.MaxScore,
		CreatedAt:   session.CreatedAt,
	}

	respondJSON(w, http.StatusOK, response)
}

func (h *GameHandler) SubmitResult(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	if idStr == "" {
		respondError(w, http.StatusBadRequest, "missing_id", "Session ID is required")
		return
	}

	sessionID, err := uuid.Parse(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid_uuid", "Invalid session ID format")
		return
	}

	var req SubmitResultRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid_request", "Invalid JSON format")
		return
	}

	result, err := h.service.SubmitResult(
		r.Context(),
		sessionID,
		req.PlayerName,
		req.Fingerprint,
		req.FoundWords,
		req.DurationMs,
	)

	if err != nil {
		status, errCode, message := mapDomainError(err)
		h.logger.Error("Failed to submit result",
			slog.String("error", err.Error()),
			slog.String("session_id", sessionID.String()),
		)
		respondError(w, status, errCode, message)
		return
	}

	response := ResultResponse{
		ID:         result.ID,
		SessionID:  result.SessionID,
		PlayerName: result.PlayerName,
		WordCount:  result.WordCount,
		Score:      result.Score,
		DurationMs: result.DurationMs,
		PlayedAt:   result.PlayedAt,
	}

	respondJSON(w, http.StatusCreated, response)
}

func (h *GameHandler) GetSessionResults(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	if idStr == "" {
		respondError(w, http.StatusBadRequest, "missing_id", "Session ID is required")
		return
	}

	sessionID, err := uuid.Parse(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid_uuid", "Invalid session ID format")
		return
	}

	topN := 0
	if topStr := r.URL.Query().Get("top"); topStr != "" {
		topN, err = strconv.Atoi(topStr)
		if err != nil || topN < 0 {
			respondError(w, http.StatusBadRequest, "invalid_top", "Parameter 'top' must be a non-negative integer")
			return
		}
	}

	results, err := h.service.GetSessionResults(r.Context(), sessionID, topN)
	if err != nil {
		status, errCode, message := mapDomainError(err)
		respondError(w, status, errCode, message)
		return
	}

	response := make([]ResultResponse, len(results))
	for i, r := range results {
		response[i] = ResultResponse{
			ID:         r.ID,
			SessionID:  r.SessionID,
			PlayerName: r.PlayerName,
			WordCount:  r.WordCount,
			Score:      r.Score,
			DurationMs: r.DurationMs,
			PlayedAt:   r.PlayedAt,
		}
	}

	respondJSON(w, http.StatusOK, response)
}

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		err := json.NewEncoder(w).Encode(data)
		if err != nil {
			status, errCode, message := mapDomainError(err)
			respondError(w, status, errCode, message)
			return
		}
	}
}

func respondError(w http.ResponseWriter, status int, err string, message string) {
	respondJSON(w, status, ErrorResponse{
		Error:   err,
		Message: message,
	})
}

func mapDomainError(err error) (int, string, string) {
	switch {
	case errors.Is(err, domain.ErrSessionExpired):
		return http.StatusBadRequest, "session_expired", "Session has expired"
	case errors.Is(err, domain.ErrInvalidWord):
		return http.StatusBadRequest, "invalid_word", "One or more words are invalid"
	case errors.Is(err, domain.ErrUnsupportedLanguage):
		return http.StatusBadRequest, "unsupported_language", "Language not supported"
	case errors.Is(err, repository.ErrNotFound):
		return http.StatusNotFound, "not_found", "Resource not found"
	case errors.Is(err, repository.ErrDuplicateResult):
		return http.StatusConflict, "duplicate_result", "Result already submitted for this session"
	default:
		return http.StatusInternalServerError, "internal_error", "An internal error occurred"
	}
}

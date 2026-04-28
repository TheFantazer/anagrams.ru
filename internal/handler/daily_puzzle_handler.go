package handler

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/TheFantazer/anagrams.ru/internal/service"
	"github.com/google/uuid"
)

type DailyPuzzleHandler struct {
	dailyService service.DailyPuzzleService
	gameService  service.GameService
	logger       *slog.Logger
}

func NewDailyPuzzleHandler(dailyService service.DailyPuzzleService, gameService service.GameService, logger *slog.Logger) *DailyPuzzleHandler {
	return &DailyPuzzleHandler{
		dailyService: dailyService,
		gameService:  gameService,
		logger:       logger,
	}
}

// GetTodaysPuzzle returns today's daily puzzle session
func (h *DailyPuzzleHandler) GetTodaysPuzzle(w http.ResponseWriter, r *http.Request) {
	language := r.URL.Query().Get("language")
	if language == "" {
		language = "en"
	}

	if language != "en" && language != "ru" {
		respondError(w, http.StatusBadRequest, "invalid_language", "Language must be 'en' or 'ru'")
		return
	}

	session, err := h.dailyService.GetTodaysSession(r.Context(), language)
	if err != nil {
		status, errCode, message := mapDomainError(err)
		h.logger.Error("Failed to get today's puzzle",
			slog.String("error", err.Error()),
			slog.String("type", fmt.Sprintf("%T", err)),
		)
		respondError(w, status, errCode, message)
		return
	}

	response := DailySessionResponse{
		Session: SessionResponse{
			ID:          session.ID,
			Letters:     session.Letters,
			Language:    session.Language,
			TimeLimit:   session.TimeLimit,
			LetterCount: session.LetterCount,
			MaxScore:    session.MaxScore,
			ValidWords:  session.ValidWords,
			CreatedAt:   session.CreatedAt,
		},
		HasPlayed: false,
	}

	if userIDStr := r.URL.Query().Get("user_id"); userIDStr != "" {
		userID, err := uuid.Parse(userIDStr)
		if err == nil {
			hasPlayed, _ := h.dailyService.HasPlayedToday(r.Context(), userID)
			response.HasPlayed = hasPlayed

			stats, err := h.dailyService.GetUserDailyStats(r.Context(), userID)
			if err == nil {
				response.DailyStats = &DailyStatsResponse{
					UserID:          stats.UserID,
					CurrentStreak:   stats.CurrentStreak,
					LongestStreak:   stats.LongestStreak,
					LastPlayedDate:  stats.LastPlayedDate,
					TotalDailyGames: stats.TotalDailyGames,
				}
			}
		}
	}

	respondJSON(w, http.StatusOK, response)
}

// GetDailyStats returns user's daily game statistics
func (h *DailyPuzzleHandler) GetDailyStats(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("user_id")
	if userIDStr == "" {
		respondError(w, http.StatusBadRequest, "missing_user_id", "User ID is required")
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid_user_id", "Invalid user ID format")
		return
	}

	stats, err := h.dailyService.GetUserDailyStats(r.Context(), userID)
	if err != nil {
		status, errCode, message := mapDomainError(err)
		h.logger.Error("Failed to get daily stats", slog.String("error", err.Error()))
		respondError(w, status, errCode, message)
		return
	}

	response := DailyStatsResponse{
		UserID:          stats.UserID,
		CurrentStreak:   stats.CurrentStreak,
		LongestStreak:   stats.LongestStreak,
		LastPlayedDate:  stats.LastPlayedDate,
		TotalDailyGames: stats.TotalDailyGames,
	}

	respondJSON(w, http.StatusOK, response)
}

// SubmitDailyResult submits a result for today's daily puzzle
func (h *DailyPuzzleHandler) SubmitDailyResult(w http.ResponseWriter, r *http.Request) {
	var req SubmitResultRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid_request", "Invalid JSON format")
		return
	}

	if req.PlayerName == "" {
		respondError(w, http.StatusBadRequest, "missing_player_name", "Player name is required")
		return
	}

	if req.Fingerprint == "" {
		respondError(w, http.StatusBadRequest, "missing_fingerprint", "Fingerprint is required")
		return
	}

	if req.UserID == nil {
		respondError(w, http.StatusUnauthorized, "unauthorized", "User ID required for daily puzzle")
		return
	}

	userID, err := uuid.Parse(*req.UserID)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid_user_id", "Invalid user ID format")
		return
	}

	hasPlayed, err := h.dailyService.HasPlayedToday(r.Context(), userID)
	if err != nil {
		h.logger.Error("Failed to check if user played today", slog.String("error", err.Error()))
		respondError(w, http.StatusInternalServerError, "internal_error", "Failed to check play status")
		return
	}

	if hasPlayed {
		respondError(w, http.StatusConflict, "already_played", "You have already played today's puzzle")
		return
	}

	// Get today's puzzle
	puzzle, err := h.dailyService.GetOrCreateTodaysPuzzle(r.Context(), "en") // TODO: get language from request
	if err != nil {
		status, errCode, message := mapDomainError(err)
		h.logger.Error("Failed to get today's puzzle", slog.String("error", err.Error()))
		respondError(w, status, errCode, message)
		return
	}

	result, err := h.dailyService.SubmitDailyResult(
		r.Context(),
		puzzle.ID,
		userID,
		req.PlayerName,
		req.Fingerprint,
		req.FoundWords,
		req.DurationMs,
	)
	if err != nil {
		status, errCode, message := mapDomainError(err)
		h.logger.Error("Failed to submit daily result", slog.String("error", err.Error()))
		respondError(w, status, errCode, message)
		return
	}

	response := ResultResponse{
		ID:         result.ID,
		SessionID:  result.SessionID,
		UserID:     result.UserID,
		PlayerName: result.PlayerName,
		WordCount:  result.WordCount,
		Score:      result.Score,
		DurationMs: result.DurationMs,
		PlayedAt:   result.PlayedAt,
	}

	respondJSON(w, http.StatusCreated, response)
}

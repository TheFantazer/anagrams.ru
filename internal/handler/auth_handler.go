package handler

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/TheFantazer/anagrams.ru/internal/domain"
	"github.com/TheFantazer/anagrams.ru/internal/repository"
	"github.com/TheFantazer/anagrams.ru/internal/service"
	"github.com/google/uuid"
)

type AuthHandler struct {
	authService service.AuthService
	logger      *slog.Logger
}

func NewAuthHandler(authService service.AuthService, logger *slog.Logger) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		logger:      logger,
	}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid_request", "Invalid JSON format")
		return
	}

	if req.Username == "" || req.Password == "" {
		respondError(w, http.StatusBadRequest, "missing_fields", "Username and password are required")
		return
	}

	if !req.AcceptedPrivacyPolicy {
		respondError(w, http.StatusBadRequest, "privacy_policy_required", "You must accept the privacy policy")
		return
	}

	user, err := h.authService.Register(r.Context(), req.Username, req.Email, req.Password)
	if err != nil {
		status, errCode, message := mapAuthError(err)
		h.logger.Error("Failed to register user", slog.String("error", err.Error()))
		respondError(w, status, errCode, message)
		return
	}

	response := UserResponse{
		ID:                 user.ID,
		Username:           user.Username,
		Email:              user.Email,
		DefaultLetterCount: user.DefaultLetterCount,
		DefaultLanguage:    user.DefaultLanguage,
		DefaultTimeLimit:   user.DefaultTimeLimit,
		CreatedAt:          user.CreatedAt,
	}

	respondJSON(w, http.StatusCreated, response)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid_request", "Invalid JSON format")
		return
	}

	if req.Username == "" || req.Password == "" {
		respondError(w, http.StatusBadRequest, "missing_fields", "Username and password are required")
		return
	}

	user, err := h.authService.Login(r.Context(), req.Username, req.Password)
	if err != nil {
		status, errCode, message := mapAuthError(err)
		h.logger.Error("Failed to login user", slog.String("error", err.Error()))
		respondError(w, status, errCode, message)
		return
	}

	response := UserResponse{
		ID:                 user.ID,
		Username:           user.Username,
		Email:              user.Email,
		DefaultLetterCount: user.DefaultLetterCount,
		DefaultLanguage:    user.DefaultLanguage,
		DefaultTimeLimit:   user.DefaultTimeLimit,
		CreatedAt:          user.CreatedAt,
	}

	respondJSON(w, http.StatusOK, response)
}

func (h *AuthHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("user_id")
	if userIDStr == "" {
		respondError(w, http.StatusUnauthorized, "unauthorized", "User ID is required")
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid_user_id", "Invalid user ID format")
		return
	}

	user, err := h.authService.GetUserByID(r.Context(), userID)
	if err != nil {
		status, errCode, message := mapAuthError(err)
		respondError(w, status, errCode, message)
		return
	}

	response := UserResponse{
		ID:                 user.ID,
		Username:           user.Username,
		Email:              user.Email,
		DefaultLetterCount: user.DefaultLetterCount,
		DefaultLanguage:    user.DefaultLanguage,
		DefaultTimeLimit:   user.DefaultTimeLimit,
		CreatedAt:          user.CreatedAt,
	}

	respondJSON(w, http.StatusOK, response)
}

func (h *AuthHandler) UpdateSettings(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("user_id")
	if userIDStr == "" {
		respondError(w, http.StatusUnauthorized, "unauthorized", "User ID is required")
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid_user_id", "Invalid user ID format")
		return
	}

	var req UpdateSettingsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid_request", "Invalid JSON format")
		return
	}

	if req.LetterCount < 6 || req.LetterCount > 10 {
		respondError(w, http.StatusBadRequest, "invalid_letter_count", "Letter count must be between 6 and 10")
		return
	}

	if req.Language != "ru" && req.Language != "en" {
		respondError(w, http.StatusBadRequest, "invalid_language", "Language must be 'ru' or 'en'")
		return
	}

	if req.TimeLimit <= 0 {
		respondError(w, http.StatusBadRequest, "invalid_time_limit", "Time limit must be positive")
		return
	}

	err = h.authService.UpdateSettings(r.Context(), userID, req.LetterCount, req.Language, req.TimeLimit)
	if err != nil {
		h.logger.Error("Failed to update settings", slog.String("error", err.Error()))
		respondError(w, http.StatusInternalServerError, "update_failed", "Failed to update settings")
		return
	}

	user, err := h.authService.GetUserByID(r.Context(), userID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "internal_error", "Failed to fetch updated user")
		return
	}

	response := UserResponse{
		ID:                 user.ID,
		Username:           user.Username,
		Email:              user.Email,
		DefaultLetterCount: user.DefaultLetterCount,
		DefaultLanguage:    user.DefaultLanguage,
		DefaultTimeLimit:   user.DefaultTimeLimit,
		CreatedAt:          user.CreatedAt,
	}

	respondJSON(w, http.StatusOK, response)
}

func (h *AuthHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("user_id")
	if userIDStr == "" {
		respondError(w, http.StatusUnauthorized, "unauthorized", "User ID is required")
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid_user_id", "Invalid user ID format")
		return
	}

	stats, err := h.authService.GetUserStats(r.Context(), userID)
	if err != nil {
		h.logger.Error("Failed to get user stats", slog.String("error", err.Error()))
		respondError(w, http.StatusInternalServerError, "stats_error", "Failed to get statistics")
		return
	}

	response := UserStatsResponse{
		GamesPlayed:  stats.GamesPlayed,
		BestScore:    stats.BestScore,
		LongestWord:  stats.LongestWord,
		TotalWords:   stats.TotalWords,
		AverageScore: stats.AverageScore,
	}

	respondJSON(w, http.StatusOK, response)
}

func (h *AuthHandler) GetLeaderboard(w http.ResponseWriter, r *http.Request) {
	period := r.URL.Query().Get("period")
	if period == "" {
		period = "week"
	}

	validPeriods := map[string]bool{"day": true, "week": true, "month": true, "all": true}
	if !validPeriods[period] {
		respondError(w, http.StatusBadRequest, "invalid_period", "Period must be day, week, month, or all")
		return
	}

	leaderboard, err := h.authService.GetLeaderboard(r.Context(), period, 20)
	if err != nil {
		h.logger.Error("Failed to get leaderboard", slog.String("error", err.Error()))
		respondError(w, http.StatusInternalServerError, "leaderboard_error", "Failed to get leaderboard")
		return
	}

	type LeaderboardResponse struct {
		Name  string `json:"name"`
		Score int    `json:"score"`
		Words int    `json:"words"`
	}

	response := make([]LeaderboardResponse, 0)
	for _, entry := range leaderboard {
		response = append(response, LeaderboardResponse{
			Name:  entry.Username,
			Score: entry.BestScore,
			Words: entry.TotalWords,
		})
	}

	respondJSON(w, http.StatusOK, response)
}

func mapAuthError(err error) (int, string, string) {
	switch {
	case errors.Is(err, domain.ErrInvalidUsername):
		return http.StatusBadRequest, "invalid_username", "Username must be 3-30 characters"
	case errors.Is(err, domain.ErrPasswordTooShort):
		return http.StatusBadRequest, "password_too_short", "Password must be at least 6 characters"
	case errors.Is(err, domain.ErrUsernameTaken):
		return http.StatusConflict, "username_taken", "Username is already taken"
	case errors.Is(err, domain.ErrEmailTaken):
		return http.StatusConflict, "email_taken", "Email is already taken"
	case errors.Is(err, domain.ErrInvalidCredentials):
		return http.StatusUnauthorized, "invalid_credentials", "Invalid username or password"
	case errors.Is(err, domain.ErrPrivacyPolicyNotAccepted):
		return http.StatusBadRequest, "privacy_policy_required", "You must accept the privacy policy"
	case errors.Is(err, repository.ErrNotFound):
		return http.StatusNotFound, "not_found", "User not found"
	default:
		return http.StatusInternalServerError, "internal_error", "An internal error occurred"
	}
}

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

	user, err := h.authService.Register(r.Context(), req.Username, req.Email, req.Password)
	if err != nil {
		status, errCode, message := mapAuthError(err)
		h.logger.Error("Failed to register user", slog.String("error", err.Error()))
		respondError(w, status, errCode, message)
		return
	}

	response := UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
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
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}

	respondJSON(w, http.StatusOK, response)
}

func (h *AuthHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	// Предполагается, что ID пользователя будет приходить из middleware после авторизации
	// Пока используем query parameter для демонстрации
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
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
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
	case errors.Is(err, repository.ErrNotFound):
		return http.StatusNotFound, "not_found", "User not found"
	default:
		return http.StatusInternalServerError, "internal_error", "An internal error occurred"
	}
}

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

type FriendHandler struct {
	friendService service.FriendService
	logger        *slog.Logger
}

func NewFriendHandler(friendService service.FriendService, logger *slog.Logger) *FriendHandler {
	return &FriendHandler{
		friendService: friendService,
		logger:        logger,
	}
}

// SendFriendRequest - POST /api/v1/friends/requests
func (h *FriendHandler) SendFriendRequest(w http.ResponseWriter, r *http.Request) {
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

	var req SendFriendRequestRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid_request", "Invalid JSON format")
		return
	}

	if err := h.friendService.SendFriendRequest(r.Context(), userID, req.ToUserID); err != nil {
		status, errCode, message := mapFriendError(err)
		h.logger.Error("Failed to send friend request", slog.String("error", err.Error()))
		respondError(w, status, errCode, message)
		return
	}

	respondJSON(w, http.StatusCreated, map[string]string{"message": "Friend request sent"})
}

// GetPendingRequests - GET /api/v1/friends/requests/pending
func (h *FriendHandler) GetPendingRequests(w http.ResponseWriter, r *http.Request) {
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

	requests, err := h.friendService.GetPendingRequests(r.Context(), userID)
	if err != nil {
		h.logger.Error("Failed to get pending requests", slog.String("error", err.Error()))
		respondError(w, http.StatusInternalServerError, "internal_error", "Failed to get pending requests")
		return
	}

	response := make([]FriendRequestResponse, 0)
	for _, req := range requests {
		response = append(response, FriendRequestResponse{
			ID:         req.ID,
			FromUserID: req.FromUserID,
			ToUserID:   req.ToUserID,
			Status:     req.Status,
			CreatedAt:  req.CreatedAt,
			UpdatedAt:  req.UpdatedAt,
		})
	}

	respondJSON(w, http.StatusOK, response)
}

// GetSentRequests - GET /api/v1/friends/requests/sent
func (h *FriendHandler) GetSentRequests(w http.ResponseWriter, r *http.Request) {
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

	requests, err := h.friendService.GetSentRequests(r.Context(), userID)
	if err != nil {
		h.logger.Error("Failed to get sent requests", slog.String("error", err.Error()))
		respondError(w, http.StatusInternalServerError, "internal_error", "Failed to get sent requests")
		return
	}

	response := make([]FriendRequestResponse, 0)
	for _, req := range requests {
		response = append(response, FriendRequestResponse{
			ID:         req.ID,
			FromUserID: req.FromUserID,
			ToUserID:   req.ToUserID,
			Status:     req.Status,
			CreatedAt:  req.CreatedAt,
			UpdatedAt:  req.UpdatedAt,
		})
	}

	respondJSON(w, http.StatusOK, response)
}

// AcceptFriendRequest - POST /api/v1/friends/requests/{id}/accept
func (h *FriendHandler) AcceptFriendRequest(w http.ResponseWriter, r *http.Request) {
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

	requestIDStr := r.PathValue("id")
	requestID, err := uuid.Parse(requestIDStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid_request_id", "Invalid request ID format")
		return
	}

	if err := h.friendService.AcceptFriendRequest(r.Context(), userID, requestID); err != nil {
		status, errCode, message := mapFriendError(err)
		h.logger.Error("Failed to accept friend request", slog.String("error", err.Error()))
		respondError(w, status, errCode, message)
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "Friend request accepted"})
}

// RejectFriendRequest - POST /api/v1/friends/requests/{id}/reject
func (h *FriendHandler) RejectFriendRequest(w http.ResponseWriter, r *http.Request) {
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

	requestIDStr := r.PathValue("id")
	requestID, err := uuid.Parse(requestIDStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid_request_id", "Invalid request ID format")
		return
	}

	if err := h.friendService.RejectFriendRequest(r.Context(), userID, requestID); err != nil {
		status, errCode, message := mapFriendError(err)
		h.logger.Error("Failed to reject friend request", slog.String("error", err.Error()))
		respondError(w, status, errCode, message)
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "Friend request rejected"})
}

// GetFriends - GET /api/v1/friends
func (h *FriendHandler) GetFriends(w http.ResponseWriter, r *http.Request) {
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

	friends, err := h.friendService.GetFriends(r.Context(), userID)
	if err != nil {
		h.logger.Error("Failed to get friends", slog.String("error", err.Error()))
		respondError(w, http.StatusInternalServerError, "internal_error", "Failed to get friends")
		return
	}

	response := make([]UserResponse, 0)
	for _, user := range friends {
		response = append(response, UserResponse{
			ID:                 user.ID,
			Username:           user.Username,
			Email:              user.Email,
			DefaultLetterCount: user.DefaultLetterCount,
			DefaultLanguage:    user.DefaultLanguage,
			DefaultTimeLimit:   user.DefaultTimeLimit,
			CreatedAt:          user.CreatedAt,
		})
	}

	respondJSON(w, http.StatusOK, response)
}

// RemoveFriend - DELETE /api/v1/friends/{id}
func (h *FriendHandler) RemoveFriend(w http.ResponseWriter, r *http.Request) {
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

	friendIDStr := r.PathValue("id")
	friendID, err := uuid.Parse(friendIDStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid_friend_id", "Invalid friend ID format")
		return
	}

	if err := h.friendService.RemoveFriend(r.Context(), userID, friendID); err != nil {
		status, errCode, message := mapFriendError(err)
		h.logger.Error("Failed to remove friend", slog.String("error", err.Error()))
		respondError(w, status, errCode, message)
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "Friend removed"})
}

// SearchUsers - GET /api/v1/users/search
func (h *FriendHandler) SearchUsers(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		respondError(w, http.StatusBadRequest, "missing_query", "Search query is required")
		return
	}

	users, err := h.friendService.SearchUsers(r.Context(), query)
	if err != nil {
		h.logger.Error("Failed to search users", slog.String("error", err.Error()))
		respondError(w, http.StatusInternalServerError, "internal_error", "Failed to search users")
		return
	}

	response := make([]UserResponse, 0)
	for _, user := range users {
		response = append(response, UserResponse{
			ID:                 user.ID,
			Username:           user.Username,
			Email:              user.Email,
			DefaultLetterCount: user.DefaultLetterCount,
			DefaultLanguage:    user.DefaultLanguage,
			DefaultTimeLimit:   user.DefaultTimeLimit,
			CreatedAt:          user.CreatedAt,
		})
	}

	respondJSON(w, http.StatusOK, response)
}

// GetUserByID - GET /api/v1/users/{id}
func (h *FriendHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid_user_id", "Invalid user ID format")
		return
	}

	user, err := h.friendService.GetUserByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			respondError(w, http.StatusNotFound, "not_found", "User not found")
			return
		}

		h.logger.Error("Failed to get user by id", slog.String("error", err.Error()))
		respondError(w, http.StatusInternalServerError, "internal_error", "Failed to get user")
		return
	}

	respondJSON(w, http.StatusOK, UserResponse{
		ID:                 user.ID,
		Username:           user.Username,
		Email:              user.Email,
		DefaultLetterCount: user.DefaultLetterCount,
		DefaultLanguage:    user.DefaultLanguage,
		DefaultTimeLimit:   user.DefaultTimeLimit,
		CreatedAt:          user.CreatedAt,
	})
}

func mapFriendError(err error) (int, string, string) {
	switch {
	case errors.Is(err, domain.ErrUserNotFound):
		return http.StatusNotFound, "user_not_found", "User not found"
	case errors.Is(err, domain.ErrFriendRequestAlreadyExists):
		return http.StatusConflict, "request_already_exists", "Friend request already exists"
	case errors.Is(err, domain.ErrAlreadyFriends):
		return http.StatusConflict, "already_friends", "Users are already friends"
	case errors.Is(err, domain.ErrNotFriends):
		return http.StatusBadRequest, "not_friends", "Users are not friends"
	case errors.Is(err, domain.ErrFriendRequestNotPending):
		return http.StatusBadRequest, "request_not_pending", "Friend request is not pending"
	case errors.Is(err, domain.ErrSelfFriendRequest):
		return http.StatusBadRequest, "self_request", "Cannot send friend request to yourself"
	case errors.Is(err, repository.ErrNotFound):
		return http.StatusNotFound, "not_found", "Resource not found"
	default:
		return http.StatusInternalServerError, "internal_error", "An internal error occurred"
	}
}

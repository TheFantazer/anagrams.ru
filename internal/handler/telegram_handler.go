package handler

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log/slog"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/TheFantazer/anagrams.ru/internal/service"
)

type TelegramHandler struct {
	authService service.AuthService
	jwtService  service.JWTService
	botToken    string
	frontendURL string
	logger      *slog.Logger
}

func NewTelegramHandler(
	authService service.AuthService,
	jwtService service.JWTService,
	botToken string,
	frontendURL string,
	logger *slog.Logger,
) *TelegramHandler {
	return &TelegramHandler{
		authService: authService,
		jwtService:  jwtService,
		botToken:    botToken,
		frontendURL: frontendURL,
		logger:      logger,
	}
}

// TelegramCallback handles Telegram Login Widget callback
func (h *TelegramHandler) TelegramCallback(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	telegramID := query.Get("id")
	firstName := query.Get("first_name")
	lastName := query.Get("last_name")
	username := query.Get("username")
	// photoURL := query.Get("photo_url")
	authDate := query.Get("auth_date")
	hash := query.Get("hash")

	if telegramID == "" || hash == "" {
		h.logger.Error("Missing required Telegram auth parameters")
		http.Redirect(w, r, h.frontendURL+"/auth?error=invalid_telegram_data", http.StatusTemporaryRedirect)
		return
	}

	// Verify Telegram auth data authenticity
	if !h.verifyTelegramAuth(query, hash) {
		h.logger.Error("Telegram auth data verification failed")
		http.Redirect(w, r, h.frontendURL+"/auth?error=telegram_auth_failed", http.StatusTemporaryRedirect)
		return
	}

	// Check auth_date is not too old (within 1 hour)
	authTimestamp, err := strconv.ParseInt(authDate, 10, 64)
	if err != nil {
		h.logger.Error("Invalid auth_date", slog.String("error", err.Error()))
		http.Redirect(w, r, h.frontendURL+"/auth?error=invalid_auth_date", http.StatusTemporaryRedirect)
		return
	}

	if time.Now().Unix()-authTimestamp > 3600 {
		h.logger.Error("Telegram auth data is too old")
		http.Redirect(w, r, h.frontendURL+"/auth?error=auth_expired", http.StatusTemporaryRedirect)
		return
	}

	// display name
	displayName := firstName
	if username != "" {
		displayName = username
	} else if lastName != "" {
		displayName = firstName + " " + lastName
	}

	// Login or register user via Telegram
	user, err := h.authService.LoginOrRegisterWithOAuth(
		r.Context(),
		"telegram",
		telegramID,
		"", // Telegram doesn't provide email
		displayName,
	)

	if err != nil {
		h.logger.Error("Failed to login/register via Telegram", slog.String("error", err.Error()))
		http.Redirect(w, r, h.frontendURL+"/auth?error=registration_failed", http.StatusTemporaryRedirect)
		return
	}

	accessToken, refreshToken, err := h.jwtService.GenerateTokenPair(
		user.ID,
		user.Username,
	)
	if err != nil {
		h.logger.Error("Failed to generate tokens", slog.String("error", err.Error()))
		http.Redirect(w, r, h.frontendURL+"/auth?error=token_generation_failed", http.StatusTemporaryRedirect)
		return
	}

	redirectURL := fmt.Sprintf(
		"%s/auth/callback?access_token=%s&refresh_token=%s",
		h.frontendURL,
		accessToken,
		refreshToken,
	)

	http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
}

// verifyTelegramAuth verifies the authenticity of Telegram Login Widget data
func (h *TelegramHandler) verifyTelegramAuth(query map[string][]string, receivedHash string) bool {
	dataCheckString := make([]string, 0)
	for key, values := range query {
		if key == "hash" {
			continue
		}
		if len(values) > 0 {
			dataCheckString = append(dataCheckString, fmt.Sprintf("%s=%s", key, values[0]))
		}
	}

	sort.Strings(dataCheckString)

	dataCheckStr := strings.Join(dataCheckString, "\n")

	secretKeyHash := sha256.Sum256([]byte(h.botToken))

	mac := hmac.New(sha256.New, secretKeyHash[:])
	mac.Write([]byte(dataCheckStr))
	expectedHash := hex.EncodeToString(mac.Sum(nil))

	// Compare hashes
	return hmac.Equal([]byte(expectedHash), []byte(receivedHash))
}

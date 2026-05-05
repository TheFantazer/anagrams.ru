package handler

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
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
func (h *TelegramHandler) verifyTelegramAuth(query url.Values, receivedHash string) bool {
	fields := []string{
		"id",
		"first_name",
		"last_name",
		"username",
		"photo_url",
		"auth_date",
	}

	var dataCheckString []string

	for _, key := range fields {
		value := query.Get(key)
		if value != "" {
			dataCheckString = append(dataCheckString, key+"="+value)
		}
	}

	sort.Strings(dataCheckString)

	dataCheckStr := strings.Join(dataCheckString, "\n")

	secretKey := sha256.Sum256([]byte(h.botToken))

	mac := hmac.New(sha256.New, secretKey[:])
	mac.Write([]byte(dataCheckStr))
	expectedHash := hex.EncodeToString(mac.Sum(nil))

	return hmac.Equal([]byte(expectedHash), []byte(receivedHash))
}

// BotWebhook handles incoming updates from Telegram Bot API
func (h *TelegramHandler) BotWebhook(w http.ResponseWriter, r *http.Request) {
	var update TelegramUpdate

	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		h.logger.Error("Invalid webhook payload", slog.String("error", err.Error()))
		respondError(w, http.StatusBadRequest, "invalid_payload", "Invalid webhook payload")
		return
	}

	h.logger.Info("Telegram webhook received",
		slog.Int64("update_id", update.UpdateID),
		slog.Bool("has_message", update.Message != nil),
		slog.Bool("has_inline_query", update.InlineQuery != nil),
	)

	if update.InlineQuery != nil {
		go h.handleInlineQuery(update.InlineQuery)
	}

	if update.Message != nil {
		go h.handleMessage(update.Message)
	}

	respondJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

// handleInlineQuery handles inline query from Telegram
func (h *TelegramHandler) handleInlineQuery(query *InlineQuery) {
	h.logger.Info("Handling inline query",
		slog.String("query_id", query.ID),
		slog.String("query_text", query.Query),
		slog.String("from_username", query.From.Username),
	)

	webAppURL := "https://anagrams.ru/tg"

	results := []interface{}{
		InlineQueryResultArticle{
			Type:        "article",
			ID:          "start_game",
			Title:       "🎮 Начать игру",
			Description: "Быстрый старт игры в анаграммы",
			InputMessageContent: &InputTextMessageContent{
				MessageText: "Играю в анаграммы! 🎮",
			},
			ReplyMarkup: &InlineKeyboardMarkup{
				InlineKeyboard: [][]InlineKeyboardButton{
					{
						{
							Text: "🎮 Играть",
							WebApp: &WebAppInfo{
								URL: webAppURL,
							},
						},
					},
				},
			},
		},
		InlineQueryResultArticle{
			Type:        "article",
			ID:          "open_app",
			Title:       "📱 Открыть приложение",
			Description: "Полная версия Mini App",
			InputMessageContent: &InputTextMessageContent{
				MessageText: "Открываю приложение анаграмм 📱",
			},
			ReplyMarkup: &InlineKeyboardMarkup{
				InlineKeyboard: [][]InlineKeyboardButton{
					{
						{
							Text: "📱 Открыть",
							WebApp: &WebAppInfo{
								URL: webAppURL,
							},
						},
					},
				},
			},
		},
	}

	err := h.answerInlineQuery(query.ID, results)
	if err != nil {
		h.logger.Error("Failed to answer inline query",
			slog.String("query_id", query.ID),
			slog.String("error", err.Error()),
		)
	}
}

// handleMessage handles incoming messages from Telegram
func (h *TelegramHandler) handleMessage(message *TelegramMessage) {
	h.logger.Info("Handling message",
		slog.Int64("message_id", message.MessageID),
		slog.String("text", message.Text),
	)

	// TODO: Implement message handling logic
}

// answerInlineQuery sends answer to inline query via Telegram Bot API
func (h *TelegramHandler) answerInlineQuery(queryID string, results []interface{}) error {
	payload := AnswerInlineQueryRequest{
		InlineQueryID: queryID,
		Results:       results,
		CacheTime:     300, // Cache for 5 minutes
	}

	return h.sendTelegramRequest("answerInlineQuery", payload)
}

// sendTelegramRequest sends a request to Telegram Bot API
func (h *TelegramHandler) sendTelegramRequest(method string, payload interface{}) error {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/%s", h.botToken, method)

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	resp, err := http.Post(url, "application/json", strings.NewReader(string(jsonData)))
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			h.logger.Error("failed to close response body",
				slog.Any("error", err),
			)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("telegram API returned non-200 status: %d", resp.StatusCode)
	}

	return nil
}

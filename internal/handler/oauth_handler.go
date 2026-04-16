package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/TheFantazer/anagrams.ru/internal/service"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type OAuthHandler struct {
	authService     service.AuthService
	jwtService      service.JWTService
	oauthConfig     *oauth2.Config
	frontendBaseURL string
	logger          *slog.Logger
}

func NewOAuthHandler(authService service.AuthService, jwtService service.JWTService, clientID, clientSecret, redirectURL, frontendBaseURL string, logger *slog.Logger) *OAuthHandler {
	return &OAuthHandler{
		authService: authService,
		jwtService:  jwtService,
		oauthConfig: &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  redirectURL,
			Scopes:       []string{"email", "profile"},
			Endpoint:     google.Endpoint,
		},
		frontendBaseURL: frontendBaseURL,
		logger:          logger,
	}
}

func (h *OAuthHandler) GoogleLogin(w http.ResponseWriter, r *http.Request) {
	url := h.oauthConfig.AuthCodeURL("state", oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (h *OAuthHandler) GoogleCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		respondError(w, http.StatusBadRequest, "invalid_request", "Code not found")
		return
	}

	token, err := h.oauthConfig.Exchange(context.Background(), code)
	if err != nil {
		h.logger.Error("failed to exchange token", slog.String("error", err.Error()))
		respondError(w, http.StatusInternalServerError, "oauth_error", "Failed to exchange token")
		return
	}

	client := h.oauthConfig.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		h.logger.Error("failed to get user info", slog.String("error", err.Error()))
		respondError(w, http.StatusInternalServerError, "oauth_error", "Failed to get user info")
		return
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			h.logger.Error("failed to close response body", slog.String("error", cerr.Error()))
		}
	}()

	var userInfo struct {
		ID    string `json:"id"`
		Email string `json:"email"`
		Name  string `json:"name"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		h.logger.Error("failed to decode user info", slog.String("error", err.Error()))
		respondError(w, http.StatusInternalServerError, "decode_error", "Failed to decode user info")
		return
	}

	user, err := h.authService.LoginOrRegisterWithOAuth(
		r.Context(),
		"google",
		userInfo.ID,
		userInfo.Email,
		userInfo.Name,
	)
	if err != nil {
		h.logger.Error("failed to login/register with oauth", slog.String("error", err.Error()))
		respondError(w, http.StatusInternalServerError, "auth_error", "Failed to authenticate")
		return
	}

	accessToken, refreshToken, err := h.jwtService.GenerateTokenPair(user.ID, user.Username)
	if err != nil {
		h.logger.Error("failed to generate tokens", slog.String("error", err.Error()))
		respondError(w, http.StatusInternalServerError, "token_error", "Failed to generate tokens")
		return
	}

	redirectURL := fmt.Sprintf("%s/auth/callback?access_token=%s&refresh_token=%s", h.frontendBaseURL, accessToken, refreshToken)
	http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
}

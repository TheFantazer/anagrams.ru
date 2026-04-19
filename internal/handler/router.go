package handler

import (
	"log/slog"
	"net/http"

	"github.com/TheFantazer/anagrams.ru/internal/config"
	"github.com/TheFantazer/anagrams.ru/internal/service"
)

func NewRouter(gameService service.GameService, authService service.AuthService, jwtService service.JWTService, cfg *config.Config, logger *slog.Logger) http.Handler {
	mux := http.NewServeMux()

	gameHandler := NewGameHandler(gameService, logger)
	authHandler := NewAuthHandler(authService, logger)
	oauthHandler := NewOAuthHandler(
		authService,
		jwtService,
		cfg.GoogleOAuth.ClientID,
		cfg.GoogleOAuth.ClientSecret,
		cfg.GoogleOAuth.RedirectURI,
		cfg.App.FrontendURL,
		logger,
	)

	// Game endpoints
	mux.HandleFunc("POST /api/v1/sessions", gameHandler.CreateSession)
	mux.HandleFunc("GET /api/v1/sessions/my", gameHandler.GetUserSessions)
	mux.HandleFunc("GET /api/v1/sessions/participated", gameHandler.GetParticipatedSessions)
	mux.HandleFunc("GET /api/v1/sessions/{id}", gameHandler.GetSession)
	mux.HandleFunc("POST /api/v1/sessions/{id}/results", gameHandler.SubmitResult)
	mux.HandleFunc("GET /api/v1/sessions/{id}/results", gameHandler.GetSessionResults)

	// Auth endpoints
	mux.HandleFunc("POST /api/v1/auth/register", authHandler.Register)
	mux.HandleFunc("POST /api/v1/auth/login", authHandler.Login)
	mux.HandleFunc("GET /api/v1/auth/me", authHandler.GetMe)
	mux.HandleFunc("PUT /api/v1/auth/settings", authHandler.UpdateSettings)
	mux.HandleFunc("GET /api/v1/auth/stats", authHandler.GetStats)
	mux.HandleFunc("GET /api/v1/leaderboard", authHandler.GetLeaderboard)

	// OAuth endpoints
	mux.HandleFunc("GET /api/v1/auth/google", oauthHandler.GoogleLogin)
	mux.HandleFunc("GET /api/v1/auth/google/callback", oauthHandler.GoogleCallback)

	handler := RecoveryMiddleware(logger)(
		RequestIDMiddleware(
			LoggingMiddleware(logger)(
				CORSMiddleware(mux),
			),
		),
	)

	return handler
}

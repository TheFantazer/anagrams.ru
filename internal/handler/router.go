package handler

import (
	"log/slog"
	"net/http"

	"github.com/TheFantazer/anagrams.ru/internal/config"
	"github.com/TheFantazer/anagrams.ru/internal/service"
)

func NewRouter(gameService service.GameService, authService service.AuthService, jwtService service.JWTService, friendService service.FriendService, cfg *config.Config, logger *slog.Logger) http.Handler {
	mux := http.NewServeMux()

	gameHandler := NewGameHandler(gameService, authService, logger)
	authHandler := NewAuthHandler(authService, logger)
	friendHandler := NewFriendHandler(friendService, logger)
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

	// Friend endpoints
	mux.HandleFunc("POST /api/v1/friends/requests", friendHandler.SendFriendRequest)
	mux.HandleFunc("GET /api/v1/friends/requests/pending", friendHandler.GetPendingRequests)
	mux.HandleFunc("GET /api/v1/friends/requests/sent", friendHandler.GetSentRequests)
	mux.HandleFunc("POST /api/v1/friends/requests/{id}/accept", friendHandler.AcceptFriendRequest)
	mux.HandleFunc("POST /api/v1/friends/requests/{id}/reject", friendHandler.RejectFriendRequest)
	mux.HandleFunc("GET /api/v1/friends", friendHandler.GetFriends)
	mux.HandleFunc("DELETE /api/v1/friends/{id}", friendHandler.RemoveFriend)

	// User search endpoint
	mux.HandleFunc("GET /api/v1/users/search", friendHandler.SearchUsers)

	handler := RecoveryMiddleware(logger)(
		RequestIDMiddleware(
			LoggingMiddleware(logger)(
				CORSMiddleware(mux),
			),
		),
	)

	return handler
}

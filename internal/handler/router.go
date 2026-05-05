package handler

import (
	"log/slog"
	"net/http"

	"github.com/TheFantazer/anagrams.ru/internal/config"
	"github.com/TheFantazer/anagrams.ru/internal/repository"
	"github.com/TheFantazer/anagrams.ru/internal/service"
)

func NewRouter(gameService service.GameService, authService service.AuthService, jwtService service.JWTService, friendService service.FriendService, dailyPuzzleService service.DailyPuzzleService, sessionInviteRepo repository.SessionInviteRepository, participantRepo repository.SessionParticipantRepository, cfg *config.Config, logger *slog.Logger) http.Handler {
	mux := http.NewServeMux()

	gameHandler := NewGameHandler(gameService, authService, sessionInviteRepo, participantRepo, logger)
	authHandler := NewAuthHandler(authService, logger)
	friendHandler := NewFriendHandler(friendService, logger)
	dailyPuzzleHandler := NewDailyPuzzleHandler(dailyPuzzleService, gameService, logger)
	telegramHandler := NewTelegramHandler(
		authService,
		jwtService,
		cfg.Telegram.BotToken,
		cfg.App.FrontendURL,
		logger,
	)

	// Game endpoints
	mux.HandleFunc("POST /api/v1/sessions", gameHandler.CreateSession)
	mux.HandleFunc("GET /api/v1/sessions/my", gameHandler.GetUserSessions)
	mux.HandleFunc("GET /api/v1/sessions/participated", gameHandler.GetParticipatedSessions)
	mux.HandleFunc("GET /api/v1/sessions/all", gameHandler.GetAllUserSessionsPaginated)
	mux.HandleFunc("GET /api/v1/sessions/{id}", gameHandler.GetSession)
	mux.HandleFunc("POST /api/v1/sessions/{id}/start", gameHandler.StartSession)
	mux.HandleFunc("POST /api/v1/sessions/{id}/results", gameHandler.SubmitResult)
	mux.HandleFunc("GET /api/v1/sessions/{id}/results", gameHandler.GetSessionResults)
	mux.HandleFunc("POST /api/v1/sessions/{id}/invites", gameHandler.CreateSessionInvite)
	mux.HandleFunc("GET /api/v1/sessions/{id}/invites", gameHandler.GetSessionInvites)
	// Auth endpoints
	mux.HandleFunc("POST /api/v1/auth/register", authHandler.Register)
	mux.HandleFunc("POST /api/v1/auth/login", authHandler.Login)
	mux.HandleFunc("GET /api/v1/auth/me", authHandler.GetMe)
	mux.HandleFunc("PUT /api/v1/auth/settings", authHandler.UpdateSettings)
	mux.HandleFunc("PUT /api/v1/auth/username", authHandler.UpdateUsername)
	mux.HandleFunc("GET /api/v1/auth/stats", authHandler.GetStats)
	mux.HandleFunc("GET /api/v1/leaderboard", authHandler.GetLeaderboard)

	// Telegram endpoints
	mux.HandleFunc("GET /api/v1/auth/telegram/callback", telegramHandler.TelegramCallback)
	mux.HandleFunc("POST /api/v1/bot-webhook", telegramHandler.BotWebhook)

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
	mux.HandleFunc("GET /api/v1/users/username/{username}", friendHandler.GetUserByUsername)
	mux.HandleFunc("GET /api/v1/users/{id}", friendHandler.GetUserByID)

	// Friend suggestions endpoint
	mux.HandleFunc("GET /api/v1/friends/suggestions", friendHandler.GetSuggestedFriends)

	// Daily puzzle endpoints
	mux.HandleFunc("GET /api/v1/daily-puzzle/today", dailyPuzzleHandler.GetTodaysPuzzle)
	mux.HandleFunc("GET /api/v1/daily-puzzle/stats", dailyPuzzleHandler.GetDailyStats)
	mux.HandleFunc("POST /api/v1/daily-puzzle/submit", dailyPuzzleHandler.SubmitDailyResult)

	handler := RecoveryMiddleware(logger)(
		RequestIDMiddleware(
			LoggingMiddleware(logger)(
				CORSMiddleware(mux),
			),
		),
	)

	return handler
}

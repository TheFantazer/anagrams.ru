package handler

import (
	"log/slog"
	"net/http"

	"github.com/TheFantazer/anagrams.ru/internal/service"
)

func NewRouter(gameService service.GameService, authService service.AuthService, logger *slog.Logger) http.Handler {
	mux := http.NewServeMux()

	gameHandler := NewGameHandler(gameService, logger)
	authHandler := NewAuthHandler(authService, logger)

	// Game endpoints
	mux.HandleFunc("POST /api/v1/sessions", gameHandler.CreateSession)
	mux.HandleFunc("GET /api/v1/sessions/{id}", gameHandler.GetSession)
	mux.HandleFunc("POST /api/v1/sessions/{id}/results", gameHandler.SubmitResult)
	mux.HandleFunc("GET /api/v1/sessions/{id}/results", gameHandler.GetSessionResults)

	// Auth endpoints
	mux.HandleFunc("POST /api/v1/auth/register", authHandler.Register)
	mux.HandleFunc("POST /api/v1/auth/login", authHandler.Login)
	mux.HandleFunc("GET /api/v1/auth/me", authHandler.GetMe)

	handler := RecoveryMiddleware(logger)(
		RequestIDMiddleware(
			LoggingMiddleware(logger)(
				CORSMiddleware(mux),
			),
		),
	)

	return handler
}

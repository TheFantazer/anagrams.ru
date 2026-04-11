package handler

import (
	"log/slog"
	"net/http"

	"github.com/TheFantazer/anagrams.ru/internal/service"
)

func NewRouter(gameService service.GameService, logger *slog.Logger) http.Handler {
	mux := http.NewServeMux()

	gameHandler := NewGameHandler(gameService, logger)

	mux.HandleFunc("POST /api/v1/sessions", gameHandler.CreateSession)
	mux.HandleFunc("GET /api/v1/sessions/{id}", gameHandler.GetSession)
	mux.HandleFunc("POST /api/v1/sessions/{id}/results", gameHandler.SubmitResult)
	mux.HandleFunc("GET /api/v1/sessions/{id}/results", gameHandler.GetSessionResults)

	// Применяем middleware в правильном порядке (от внешнего к внутреннему):
	// 1. Recovery (перехватывает панники от всех остальных)
	// 2. RequestID (добавляет ID в контекст)
	// 3. Logging (логирует с request ID)
	// 4. CORS (обрабатывает CORS headers)
	// 5. mux (роутинг)
	handler := RecoveryMiddleware(logger)(
		RequestIDMiddleware(
			LoggingMiddleware(logger)(
				CORSMiddleware(mux),
			),
		),
	)

	return handler
}

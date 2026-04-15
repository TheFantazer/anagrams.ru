package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/TheFantazer/anagrams.ru/internal/config"
	"github.com/TheFantazer/anagrams.ru/internal/dictionary"
	"github.com/TheFantazer/anagrams.ru/internal/handler"
	"github.com/TheFantazer/anagrams.ru/internal/repository/postgres"
	"github.com/TheFantazer/anagrams.ru/internal/service"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load config: %v\n", err)
		os.Exit(1)
	}

	logger := setupLogger(cfg.App.LogLevel)
	slog.SetDefault(logger)

	slog.Info("starting application",
		slog.String("env", cfg.App.Env),
		slog.String("port", cfg.App.Port),
	)

	db, err := connectDB(cfg.Postgres)
	if err != nil {
		slog.Error("failed to connect to database", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer func() {
		if err := db.Close(); err != nil {
			slog.Error("failed to close database", slog.String("error", err.Error()))
		}
	}()
	slog.Info("connected to database")

	dictionaries, err := loadDictionaries(cfg.Dict)
	if err != nil {
		slog.Error("failed to load dictionaries", slog.String("error", err.Error()))
		os.Exit(1)
	}
	slog.Info("dictionaries loaded", slog.Int("count", len(dictionaries)))

	sessionRepo := postgres.NewSessionRepository(db)
	resultRepo := postgres.NewResultRepository(db)
	userRepo := postgres.NewUserRepository(db.DB)
	statsRepo := postgres.NewStatsRepository(db)

	letterGen := dictionary.NewLetterGenerator()
	gameService := service.NewGameService(sessionRepo, resultRepo, dictionaries, letterGen)
	authService := service.NewAuthService(userRepo, statsRepo)

	router := handler.NewRouter(gameService, authService, logger)

	server := &http.Server{
		Addr:         ":" + cfg.App.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		slog.Info("starting HTTP server", slog.String("addr", server.Addr))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server error", slog.String("error", err.Error()))
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("shutting down server")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("server shutdown error", slog.String("error", err.Error()))
	}

	slog.Info("server stopped gracefully")
}

func setupLogger(level string) *slog.Logger {
	var logLevel slog.Level
	switch level {
	case "debug":
		logLevel = slog.LevelDebug
	case "info":
		logLevel = slog.LevelInfo
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{
		Level: logLevel,
	}

	return slog.New(slog.NewJSONHandler(os.Stdout, opts))
}

func connectDB(cfg config.PostgresConfig) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("connect to database: %w", err)
	}

	db.SetMaxOpenConns(cfg.MaxConns)
	db.SetMaxIdleConns(cfg.MaxConns / 2)
	db.SetConnMaxLifetime(time.Hour)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ping database: %w", err)
	}

	return db, nil
}

func loadDictionaries(cfg config.DictConfig) (map[string]*dictionary.Trie, error) {
	dictMap := make(map[string]*dictionary.Trie)

	for _, lang := range cfg.Languages {
		filePath := filepath.Join(cfg.Path, lang+".txt")

		trie, err := dictionary.LoadFromFile(filePath)
		if err != nil {
			return nil, fmt.Errorf("load dictionary %s: %w", lang, err)
		}

		dictMap[lang] = trie
	}

	return dictMap, nil
}

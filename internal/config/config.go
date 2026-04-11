package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	App      AppConfig
	Postgres PostgresConfig
	Redis    RedisConfig
	Dict     DictConfig
	Game     GameConfig
}

type AppConfig struct {
	Port     string
	Env      string
	LogLevel string
}

type PostgresConfig struct {
	Host     string
	Port     string
	Database string
	User     string
	Password string
	SSLMode  string
	MaxConns int
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

type DictConfig struct {
	Path      string
	Languages []string
}

type GameConfig struct {
	DefaultTimeLimit   int
	DefaultLetterCount int
	MinWordLength      int
}

func Load() (*Config, error) {
	cfg := &Config{
		App: AppConfig{
			Port:     getEnv("APP_PORT", "8080"),
			Env:      getEnv("APP_ENV", "development"),
			LogLevel: getEnv("APP_LOG_LEVEL", "info"),
		},
		Postgres: PostgresConfig{
			Host:     getEnv("POSTGRES_HOST", "localhost"),
			Port:     getEnv("POSTGRES_PORT", "5432"),
			Database: getEnv("POSTGRES_DB", "anagram"),
			User:     getEnv("POSTGRES_USER", "anagram"),
			Password: getEnv("POSTGRES_PASSWORD", "secret"),
			SSLMode:  getEnv("POSTGRES_SSLMODE", "disable"),
			MaxConns: getEnvInt("POSTGRES_MAX_CONNS", 20),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", "secret"),
			DB:       getEnvInt("REDIS_DB", 0),
		},
		Dict: DictConfig{
			Path:      getEnv("DICT_PATH", "./dictionaries"),
			Languages: getEnvSlice("DICT_LANGUAGES", "ru,en"),
		},
		Game: GameConfig{
			DefaultTimeLimit:   getEnvInt("GAME_DEFAULT_TIME_LIMIT", 60),
			DefaultLetterCount: getEnvInt("GAME_DEFAULT_LETTER_COUNT", 7),
			MinWordLength:      getEnvInt("GAME_MIN_WORD_LENGTH", 3),
		},
	}

	return cfg, nil
}

// PostgresDSN возвращает строку подключения к PostgreSQL
func (c *PostgresConfig) DSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.Database, c.SSLMode)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvSlice(key, defaultValue string) []string {
	value := getEnv(key, defaultValue)
	parts := strings.Split(value, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		if trimmed := strings.TrimSpace(part); trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

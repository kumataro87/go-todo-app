package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Server ServerConfig
	Log    LogConfig
	DB     DBConfig
}

type ServerConfig struct {
	Port            string
	ShutdownTimeout time.Duration
}

type LogConfig struct {
	Level  string
	Format string
}

type DBConfig struct {
	DSN             string
	MaxOpenConns    int
	MaxIdleConnes   int
	ConnMaxLifetime time.Duration
}

func Load() (*Config, error) {
	dsn := os.Getenv("DATABASE_URL")
	port := getEnv("PORT", "8080")
	shutdownTimeout, err := parseDuration(getEnv("SHUTDOWN_TIMEOUT", "30s"))

	maxOpenConns, err := parseInt(getEnv("DB_MAX_OPEN_CONNS", "25"))
	if err != nil {
		return nil, fmt.Errorf("invalid DB_MAX_OPEN_CONNS: %w", err)
	}

	maxIdleConnes, err := parseInt(getEnv("DB_MAX_IDLE_CONNS", "5"))
	if err != nil {
		return nil, fmt.Errorf("invalid DB_MAX_IDLE_CONNS: %w", err)
	}

	connMaxLifetime, err := parseDuration(getEnv("DB_CONN_MAX_LIFETIME", "5m"))
	if err != nil {
		return nil, fmt.Errorf("invalid DB_CONN_MAX_LIFETIME: %w", err)
	}

	return &Config{
		Server: ServerConfig{
			Port:            port,
			ShutdownTimeout: shutdownTimeout,
		},
		Log: LogConfig{
			Level:  getEnv("LOG_LEVEL", "info"),
			Format: getEnv("LOG_FORMAT", "json"),
		},
		DB: DBConfig{
			DSN:             dsn,
			MaxOpenConns:    maxOpenConns,
			MaxIdleConnes:   maxIdleConnes,
			ConnMaxLifetime: connMaxLifetime,
		},
	}, nil
}

func getEnv(key string, defaultValue string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultValue
}

func parseDuration(s string) (time.Duration, error) {
	return time.ParseDuration(s)
}

func parseInt(s string) (int, error) {
	return strconv.Atoi(s)
}

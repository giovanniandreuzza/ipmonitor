// Package config handles application configuration.
package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

// Config holds all application configuration.
type Config struct {
	TelegramBotToken string
	TelegramChatID   string
	LogLevel         string
}

// Load loads configuration from environment variables.
func Load() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
		TelegramBotToken: os.Getenv("TELEGRAM_BOT_TOKEN"),
		TelegramChatID:   os.Getenv("TELEGRAM_CHAT_ID"),
		LogLevel:         defaultValue(os.Getenv("LOG_LEVEL"), "INFO"),
	}

	if cfg.TelegramBotToken == "" {
		return nil, errors.New("TELEGRAM_BOT_TOKEN is required")
	}

	if cfg.TelegramChatID == "" {
		return nil, errors.New("TELEGRAM_CHAT_ID is required")
	}

	return cfg, nil
}

func defaultValue(value, fallback string) string {
	if value == "" {
		return fallback
	}
	return value
}

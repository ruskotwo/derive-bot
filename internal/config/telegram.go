package config

import (
	"os"
	"strconv"
)

type TelegramBotConfig struct {
	Token   string
	Timeout int
}

func NewTelegramConfig() *TelegramBotConfig {
	config := TelegramBotConfig{
		Token:   os.Getenv("TELEGRAM_BOT_TOKEN"),
		Timeout: 30,
	}

	if timeout, err := strconv.Atoi(os.Getenv("TELEGRAM_BOT_TIMEOUT")); err == nil && timeout > 0 {
		config.Timeout = timeout
	}

	return &config
}

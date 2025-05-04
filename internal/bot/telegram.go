package bot

import (
	"fmt"
	"log/slog"
	"runtime/debug"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/ruskotwo/derive-bot/internal/config"
	"github.com/ruskotwo/derive-bot/internal/derive"
	"github.com/ruskotwo/derive-bot/internal/domain/user"
	"github.com/ruskotwo/derive-bot/internal/infrastructure/localization"
	"github.com/ruskotwo/derive-bot/internal/infrastructure/logger"
)

type TelegramBot struct {
	api            *tgbotapi.BotAPI
	config         *config.TelegramBotConfig
	derive         *derive.Derive
	localize       *localization.Localize
	logger         *slog.Logger
	userRepository *user.Repository
}

func NewTelegramBot(
	config *config.TelegramBotConfig,
	derive *derive.Derive,
	localize *localization.Localize,
	logger *slog.Logger,
	userRepository *user.Repository,
) *TelegramBot {
	return &TelegramBot{
		config:         config,
		derive:         derive,
		localize:       localize,
		logger:         logger,
		userRepository: userRepository,
	}
}

func (b *TelegramBot) Start() error {
	bot, err := tgbotapi.NewBotAPI(b.config.Token)
	if err != nil {
		return err
	}

	b.api = bot

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = b.config.Timeout

	updates := b.api.GetUpdatesChan(updateConfig)

	for update := range updates {
		go func() {
			defer func() {
				if err := recover(); err != nil {
					b.logger.Error(
						fmt.Sprintf("panic: %v\n%s", err, debug.Stack()),
						slog.Any(logger.TelegramUserId, update.SentFrom().ID),
					)
				}
			}()

			if err := b.handleUpdate(update); err != nil {
				b.logger.Error(
					fmt.Sprintf("failed to handle update: %v", err),
					slog.Any(logger.TelegramUserId, update.SentFrom().ID),
				)
			}
		}()
	}

	return nil
}

package bot

import (
	"log/slog"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/nicksnyder/go-i18n/v2/i18n"

	"github.com/ruskotwo/derive-bot/internal/infrastructure/logger"
)

func (b *TelegramBot) handleUpdate(update tgbotapi.Update) error {
	telegramUser := update.SentFrom()
	localizer := b.localize.GetLocalizer(telegramUser.LanguageCode)

	userModel, err := b.userRepository.GetOneOrCreateUserByTelegramId(int(telegramUser.ID))
	if err != nil {
		return err
	}

	if userModel.Lang != telegramUser.LanguageCode {
		userModel.Lang = telegramUser.LanguageCode
		if err := b.userRepository.Update(userModel); err != nil {
			return err
		}
	}

	dto := handleDTO{
		user:      userModel,
		update:    &update,
		localizer: localizer,
	}

	switch {
	case update.Message != nil:
		return b.handleMessage(dto)
	case update.CallbackQuery != nil:
		err = b.handleCallbackQuery(dto)
		if err != nil {
			b.sendCallback(&tgbotapi.CallbackConfig{
				CallbackQueryID: dto.update.CallbackQuery.ID,
				Text:            "Invalid error",
			})
		}
		return err
	default:
		return nil
	}
}

func (b *TelegramBot) handleMessage(dto handleDTO) error {
	b.logger.Info(
		"handleMessage",
		slog.Any(logger.TelegramUserId, dto.user.TelegramId),
		slog.Any(logger.Lang, dto.update.Message.From.LanguageCode),
		slog.Any(logger.MessageText, dto.update.Message.Text),
	)

	switch dto.update.Message.Text {
	case "/start":
		return b.sendStartMessage(dto.update.Message.Chat.ID, dto.localizer)
	case "/help":
		return b.sendAboutMessage(dto.update.Message.Chat.ID, dto.localizer, true)
	default:
		return nil
	}
}

func (b *TelegramBot) handleCallbackQuery(dto handleDTO) error {
	b.logger.Info(
		"handleCallbackQuery",
		slog.Any(logger.TelegramUserId, dto.user.TelegramId),
		slog.Any(logger.Lang, dto.update.CallbackQuery.From.LanguageCode),
		slog.Any(logger.CallbackQueryData, dto.update.CallbackQuery.Data),
	)

	data, err := NewCallbackDataFromJson(dto.update.CallbackQuery.Data)
	if err != nil {
		text, err2 := dto.localizer.LocalizeMessage(&i18n.Message{ID: "invalid_callback_data"})
		if err2 != nil {
			return err2
		}

		b.sendCallback(&tgbotapi.CallbackConfig{
			CallbackQueryID: dto.update.CallbackQuery.ID,
			Text:            text,
			ShowAlert:       true,
		})

		return err
	}

	switch data.Action {
	case LetsStartBtn:
		b.sendCallback(&tgbotapi.CallbackConfig{CallbackQueryID: dto.update.CallbackQuery.ID})
		return b.letsStart(dto.update.CallbackQuery.From.ID, dto.localizer, dto.user)
	case ShowAboutBtn:
		b.sendCallback(&tgbotapi.CallbackConfig{CallbackQueryID: dto.update.CallbackQuery.ID})
		return b.sendAboutMessage(dto.update.CallbackQuery.From.ID, dto.localizer, true)
	case CompleteJourneyBtn:
		ok, err := b.completeJourney(
			dto.update.CallbackQuery,
			dto.localizer,
			dto.user,
			data.JourneyId,
		)
		if err != nil {
			return err
		}
		if !ok {
			return nil
		}

		return b.letsStart(dto.update.CallbackQuery.From.ID, dto.localizer, dto.user)
	case FinishJourneyBtn:
		ok, err := b.FinishJourney(
			dto.update.CallbackQuery,
			dto.localizer,
			dto.user,
			data.JourneyId,
		)
		if err != nil {
			return err
		}
		if !ok {
			return nil
		}

		return b.sendFinishMessage(dto.update.CallbackQuery.From.ID, dto.localizer)
	default:
		return nil
	}
}

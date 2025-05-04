package bot

import (
	"errors"
	"log/slog"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/nicksnyder/go-i18n/v2/i18n"

	"github.com/ruskotwo/derive-bot/internal/derive"
	"github.com/ruskotwo/derive-bot/internal/domain/journey"
	"github.com/ruskotwo/derive-bot/internal/domain/quest"
	"github.com/ruskotwo/derive-bot/internal/domain/user"
	"github.com/ruskotwo/derive-bot/internal/infrastructure/logger"
)

func (b *TelegramBot) sendStartMessage(chatID int64, localizer *i18n.Localizer) error {
	startMsgText, err := localizer.LocalizeMessage(&i18n.Message{ID: "start_msg"})
	if err != nil {
		return err
	}

	letsStartBtnText, err := localizer.LocalizeMessage(&i18n.Message{ID: "lets_start_btn"})
	if err != nil {
		return err
	}

	showAboutBtnText, err := localizer.LocalizeMessage(&i18n.Message{ID: "show_about_btn"})
	if err != nil {
		return err
	}

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(letsStartBtnText, CallbackData{
				Action: LetsStartBtn,
			}.ToJson()),
			tgbotapi.NewInlineKeyboardButtonData(showAboutBtnText, CallbackData{
				Action: ShowAboutBtn,
			}.ToJson()),
		),
	)

	msg := tgbotapi.NewMessage(chatID, startMsgText)
	msg.ReplyMarkup = keyboard

	_, err = b.api.Send(msg)

	return err
}

func (b *TelegramBot) sendAboutMessage(chatID int64, localizer *i18n.Localizer, needPrefix bool) error {
	aboutMsgText, err := localizer.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID: "about_msg",
		},
		TemplateData: map[string]interface{}{
			"CompleteTillSeconds": derive.CompleteTillSeconds,
		},
	})
	if err != nil {
		return err
	}

	if needPrefix {
		aboutMsgPrefixText, err := localizer.LocalizeMessage(&i18n.Message{ID: "about_msg_prefix"})
		if err != nil {
			return err
		}

		aboutMsgText = aboutMsgPrefixText + "\n" + aboutMsgText
	}

	letsStartBtnText, err := localizer.LocalizeMessage(&i18n.Message{ID: "lets_start_btn"})
	if err != nil {
		return err
	}

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(letsStartBtnText, CallbackData{
				Action: LetsStartBtn,
			}.ToJson()),
		),
	)

	msg := tgbotapi.NewMessage(chatID, aboutMsgText)
	msg.ReplyMarkup = keyboard

	_, err = b.api.Send(msg)

	return err
}

func (b *TelegramBot) letsStart(chatID int64, localizer *i18n.Localizer, userModel *user.Model) error {
	journeyModel, questModel, err := b.derive.LetsDerive(userModel)
	if err != nil {
		return err
	}

	return b.sendQuestMessage(chatID, journeyModel, questModel, localizer)
}

func (b *TelegramBot) completeJourney(
	callbackQuery *tgbotapi.CallbackQuery,
	localizer *i18n.Localizer,
	userModel *user.Model,
	journeyId int,
) (bool, error) {
	err := b.derive.CompleteJourney(journeyId, userModel.Id)
	if err != nil {
		var msgId string

		switch true {
		case errors.Is(err, derive.NotFoundJourneyError):
			fallthrough
		case errors.Is(err, derive.JourneyAlreadyCompletedError):
			msgId = "invalid_callback_data"
		case errors.Is(err, derive.JourneyCannotBeCompletedYetError):
			msgId = "journey_cannot_be_completed_yet"
		default:
			return false, err
		}

		b.logger.Info(
			"cant complete journey: "+msgId,
			slog.Int(logger.TelegramUserId, userModel.TelegramId),
			slog.Int(logger.UserId, userModel.Id),
			slog.Int(logger.JourneyId, journeyId),
		)

		text, err := localizer.LocalizeMessage(&i18n.Message{ID: msgId})
		if err != nil {
			return false, err
		}

		b.sendCallback(&tgbotapi.CallbackConfig{
			CallbackQueryID: callbackQuery.ID,
			Text:            text,
			ShowAlert:       true,
			CacheTime:       -1,
		})

		return false, nil
	}

	b.sendCallback(&tgbotapi.CallbackConfig{
		CallbackQueryID: callbackQuery.ID,
	})

	return true, nil
}

func (b *TelegramBot) sendQuestMessage(
	chatID int64,
	journeyModel *journey.Model,
	questModel *quest.Model,
	localizer *i18n.Localizer,
) error {
	letsStartBtnText, err := localizer.LocalizeMessage(&i18n.Message{ID: "next_journey_btn"})
	if err != nil {
		return err
	}

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(letsStartBtnText, CallbackData{
				Action:    CompleteJourneyBtn,
				JourneyId: journeyModel.Id,
			}.ToJson()),
		),
	)

	msg := tgbotapi.NewMessage(chatID, questModel.Description)
	msg.ReplyMarkup = keyboard

	_, err = b.api.Send(msg)

	return err
}

func (b *TelegramBot) sendCallback(config *tgbotapi.CallbackConfig) {
	if config.CacheTime == 0 {
		config.CacheTime = 600
	}

	if config.CacheTime < 0 {
		config.CacheTime = 0
	}

	_, _ = b.api.Request(config)
}

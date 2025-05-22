package bot

import (
	"errors"
	"fmt"
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
	if needPrefix {
		err := b.sendPartAboutMessage(chatID, localizer, "about_msg_prefix", false)
		if err != nil {
			b.logger.Error(
				fmt.Sprintf("failed to handle update: %v", err),
				slog.Int(logger.TelegramUserId, int(chatID)),
			)
		}
	}

	messages := []string{
		"about_msg_preface",
		"about_msg_whats_its_works",
		"about_msg_rules",
		"about_msg_conclusion",
	}

	for _, msg := range messages {
		err := b.sendPartAboutMessage(chatID, localizer, msg, msg == "about_msg_conclusion")
		if err != nil {
			b.logger.Error(
				fmt.Sprintf("failed to handle update: %v", err),
				slog.Int(logger.TelegramUserId, int(chatID)),
			)
		}
	}

	return nil
}

func (b *TelegramBot) sendPartAboutMessage(chatID int64, localizer *i18n.Localizer, msgId string, isConclusion bool) error {
	text, err := localizer.LocalizeMessage(&i18n.Message{ID: msgId})
	if err != nil {
		return err
	}

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = tgbotapi.ModeHTML

	if isConclusion {
		letsStartBtnText, err := localizer.LocalizeMessage(&i18n.Message{ID: "lets_start_btn"})
		if err != nil {
			return err
		}

		msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(letsStartBtnText, CallbackData{
					Action: LetsStartBtn,
				}.ToJson()),
			),
		)
	}

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
		case errors.Is(err, derive.JourneyAlreadyCanselError):
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

func (b *TelegramBot) finishJourney(
	callbackQuery *tgbotapi.CallbackQuery,
	localizer *i18n.Localizer,
	userModel *user.Model,
	journeyId int,
) (bool, error) {
	err := b.derive.CanselJourney(journeyId, userModel.Id)
	if err != nil {
		var msgId string

		switch true {
		case errors.Is(err, derive.NotFoundJourneyError):
			fallthrough
		case errors.Is(err, derive.JourneyAlreadyCanselError):
			fallthrough
		case errors.Is(err, derive.JourneyAlreadyCompletedError):
			msgId = "invalid_callback_data"
		default:
			return false, err
		}

		b.logger.Info(
			"cant cansel journey: "+msgId,
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
	completeJourneyBtnText, err := localizer.LocalizeMessage(&i18n.Message{ID: "complete_journey_btn"})
	if err != nil {
		return err
	}

	finishJourneyBtnText, err := localizer.LocalizeMessage(&i18n.Message{ID: "finish_journey_btn"})
	if err != nil {
		return err
	}

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(completeJourneyBtnText, CallbackData{
				Action:    CompleteJourneyBtn,
				JourneyId: journeyModel.Id,
			}.ToJson()),
			tgbotapi.NewInlineKeyboardButtonData(finishJourneyBtnText, CallbackData{
				Action:    FinishJourneyBtn,
				JourneyId: journeyModel.Id,
			}.ToJson()),
		),
	)

	switch true {
	case questModel.File != nil:
		msg := tgbotapi.NewPhoto(chatID, tgbotapi.FileURL(*questModel.File))
		msg.Caption = questModel.Description
		msg.ReplyMarkup = keyboard
		_, err = b.api.Send(msg)
		if err == nil {
			return nil
		}

		b.logger.Error(
			fmt.Sprintf("failed to send photo: %v", err),
			slog.Int(logger.TelegramUserId, int(chatID)),
			slog.Int(logger.JourneyId, journeyModel.Id),
			slog.Int(logger.QuestId, questModel.Id),
		)

		fallthrough
	default:
		msg := tgbotapi.NewMessage(chatID, questModel.Description)
		msg.ReplyMarkup = keyboard
		_, err = b.api.Send(msg)
	}

	return err
}

func (b *TelegramBot) sendFinishMessage(
	chatID int64,
	localizer *i18n.Localizer,
) error {
	finishJourneyMsgText, err := localizer.LocalizeMessage(&i18n.Message{ID: "finish_journey_msg"})
	if err != nil {
		return err
	}

	newJourneyBtnText, err := localizer.LocalizeMessage(&i18n.Message{ID: "new_journey_btn"})
	if err != nil {
		return err
	}

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(newJourneyBtnText, CallbackData{
				Action: LetsStartBtn,
			}.ToJson()),
		),
	)

	msg := tgbotapi.NewMessage(chatID, finishJourneyMsgText)
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

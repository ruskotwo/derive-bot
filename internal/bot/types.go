package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	jsoniter "github.com/json-iterator/go"
	"github.com/nicksnyder/go-i18n/v2/i18n"

	"github.com/ruskotwo/derive-bot/internal/domain/user"
)

const (
	LetsStartBtn       = "lets_start"
	ShowAboutBtn       = "show_about"
	CompleteJourneyBtn = "complete_journey"
)

var json = jsoniter.ConfigFastest

type handleDTO struct {
	user      *user.Model
	update    *tgbotapi.Update
	localizer *i18n.Localizer
}

type CallbackData struct {
	Action    string `json:"action"`
	JourneyId int    `json:"journey_id"`
}

func NewCallbackDataFromJson(data string) (CallbackData, error) {
	var callbackData CallbackData
	if err := json.UnmarshalFromString(data, &callbackData); err != nil {
		return CallbackData{}, err
	}
	return callbackData, nil
}

func (d CallbackData) ToJson() string {
	marshaled, _ := json.MarshalToString(d)
	return marshaled
}

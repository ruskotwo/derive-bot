//go:build wireinject
// +build wireinject

package factory

import (
	"github.com/google/wire"

	"github.com/ruskotwo/derive-bot/internal/bot"
)

func InitTelegramBot() (*bot.TelegramBot, func(), error) {
	panic(wire.Build(botSet))
}

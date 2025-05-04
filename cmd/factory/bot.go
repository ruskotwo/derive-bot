package factory

import (
	"github.com/google/wire"

	"github.com/ruskotwo/derive-bot/internal/bot"
	"github.com/ruskotwo/derive-bot/internal/config"
	"github.com/ruskotwo/derive-bot/internal/derive"
)

var botSet = wire.NewSet(
	infrastructureSet,
	databaseSet,
	config.NewTelegramConfig,
	derive.NewDerive,
	bot.NewTelegramBot,
)

package factory

import (
	"github.com/google/wire"

	"github.com/ruskotwo/derive-bot/internal/infrastructure/localization"
	"github.com/ruskotwo/derive-bot/internal/infrastructure/logger"
)

var infrastructureSet = wire.NewSet(
	logger.NewLogger,
	localization.NewLocalize,
)

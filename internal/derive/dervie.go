package derive

import (
	"database/sql"
	"errors"
	"log/slog"
	"time"

	"github.com/ruskotwo/derive-bot/internal/domain/journey"
	"github.com/ruskotwo/derive-bot/internal/domain/quest"
	"github.com/ruskotwo/derive-bot/internal/domain/user"
	"github.com/ruskotwo/derive-bot/internal/infrastructure/logger"
)

const defaultLang = "ru"

const CompleteTillSeconds = 15
const CompleteTill = CompleteTillSeconds * time.Second

var NotFoundNextQuestError = errors.New("not found next quest")
var NotFoundJourneyError = errors.New("not found journey")
var JourneyAlreadyCompletedError = errors.New("journey already completed")
var JourneyCannotBeCompletedYetError = errors.New("journey cannot be completed yet")

type Derive struct {
	logger            *slog.Logger
	journeyRepository *journey.Repository
	questRepository   *quest.Repository
}

func NewDerive(
	logger *slog.Logger,
	journeyRepository *journey.Repository,
	questRepository *quest.Repository,
) *Derive {
	return &Derive{
		logger:            logger,
		journeyRepository: journeyRepository,
		questRepository:   questRepository,
	}
}

func (d Derive) LetsDerive(userModel *user.Model) (*journey.Model, *quest.Model, error) {
	d.logger.Info("LetsDerive", slog.Int(logger.TelegramUserId, userModel.TelegramId))

	lastJourney, err := d.journeyRepository.GetLastForUserId(userModel.Id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return d.StartDerive(userModel, quest.CategoryDirection)
		}
		return nil, nil, err
	}

	lastQuest, err := d.questRepository.GetOneById(lastJourney.QuestId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return d.StartDerive(userModel, quest.CategoryDirection)
		}
		return nil, nil, err
	}

	if lastJourney.Progress == journey.ProgressNone {
		d.logger.Info(
			"non progress, return old",
			slog.Int(logger.TelegramUserId, userModel.TelegramId),
			slog.Int(logger.JourneyId, lastJourney.Id),
			slog.Int(logger.QuestId, lastJourney.QuestId),
		)
		// Если прогресса нет - значит можем вернуть этот же квест
		return lastJourney, lastQuest, nil
	}

	return d.StartDerive(userModel, d.getNextCategory(lastQuest.CategoryId))
}

func (d Derive) CompleteJourney(journeyId int, userId int) error {
	d.logger.Info(
		"CompleteJourney",
		slog.Int(logger.UserId, userId),
		slog.Int(logger.JourneyId, journeyId),
	)

	journeyModel, err := d.journeyRepository.GetOneByIdAndUserId(journeyId, userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return NotFoundJourneyError
		}
		return err
	}

	if journeyModel.Progress == journey.ProgressCompleted {
		return JourneyAlreadyCompletedError
	}

	if journeyModel.CompleteTillAt.After(time.Now()) {
		return JourneyCannotBeCompletedYetError
	}

	journeyModel.Progress = journey.ProgressCompleted
	err = d.journeyRepository.Save(journeyModel)
	if err != nil {
		return err
	}

	d.logger.Info(
		"completed journey",
		slog.Int(logger.UserId, userId),
		slog.Int(logger.JourneyId, journeyId),
	)

	return nil
}

func (d Derive) StartDerive(userModel *user.Model, categoryId int) (*journey.Model, *quest.Model, error) {
	d.logger.Info(
		"StartDerive",
		slog.Int(logger.TelegramUserId, userModel.TelegramId),
	)

	newQuest, err := d.getRandomQuest(userModel.Lang, categoryId)
	if err != nil {
		return nil, nil, err
	}

	err = d.journeyRepository.Save(&journey.Model{
		UserId:         userModel.Id,
		QuestId:        newQuest.Id,
		CompleteTillAt: time.Now().Add(CompleteTill),
	})
	if err != nil {
		return nil, nil, err
	}

	currentJourney, err := d.journeyRepository.GetLastForUserId(userModel.Id)
	if err != nil {
		return nil, nil, err
	}

	if currentJourney.QuestId != newQuest.Id {
		return nil, nil, errors.New("quest id mismatch")
	}

	d.logger.Info(
		"started derive",
		slog.Int(logger.TelegramUserId, userModel.TelegramId),
		slog.Int(logger.JourneyId, currentJourney.Id),
		slog.Int(logger.QuestId, currentJourney.QuestId),
	)

	return currentJourney, newQuest, nil
}

func (d Derive) getRandomQuest(
	lang string,
	categoryId int,
) (*quest.Model, error) {
	d.logger.Debug("getRandomQuest", slog.Any("lang", lang), slog.Any("categoryId", categoryId))

	newQuest, err := d.questRepository.GetRandomByLangAndCategoryId(lang, categoryId)

	if err == nil {
		return newQuest, nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	if lang != defaultLang {
		return d.getRandomQuest(defaultLang, categoryId)
	}

	return nil, NotFoundNextQuestError
}

func (d Derive) getNextCategory(lastCategory int) int {
	switch lastCategory {
	case quest.CategoryDirection:
		return quest.CategoryAction
	case quest.CategoryAction:
		return quest.CategoryCreative
	case quest.CategoryCreative:
		return quest.CategoryDirection
	default:
		return quest.CategoryDirection
	}
}

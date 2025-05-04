package factory

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/wire"
	"github.com/jmoiron/sqlx"

	"github.com/ruskotwo/derive-bot/internal/config"
	"github.com/ruskotwo/derive-bot/internal/domain/journey"
	"github.com/ruskotwo/derive-bot/internal/domain/quest"
	"github.com/ruskotwo/derive-bot/internal/domain/user"
)

var databaseSet = wire.NewSet(
	config.NewMysqlDatabaseConfig,
	NewDB,
	user.NewRepository,
	quest.NewRepository,
	journey.NewRepository,
)

func NewDB(
	config *config.MysqlDatabaseConfig,
) (*sqlx.DB, func(), error) {
	db, err := sqlx.Connect("mysql", config.GetDsn())
	if err != nil {
		return nil, func() {}, err
	}

	db.SetMaxOpenConns(config.MaxOpenConnections)
	db.SetMaxIdleConns(config.MaxIdleConnections)
	db.SetConnMaxLifetime(config.ConnectionMaxLifeTime)
	db.SetConnMaxIdleTime(config.ConnectionMaxIdleTime)

	return db, func() { _ = db.Close() }, nil
}

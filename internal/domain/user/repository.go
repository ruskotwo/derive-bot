package user

import (
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) GetOneOrCreateUserByTelegramId(telegramId int) (*Model, error) {
	user, err := r.GetOneByTelegramId(telegramId)
	if err == nil {
		return user, nil
	}

	if errors.Is(err, sql.ErrNoRows) {
		err = r.CreateUser(&Model{
			TelegramId: telegramId,
		})
		if err != nil {
			return nil, err
		}
	}

	return r.GetOneByTelegramId(telegramId)
}

func (r *Repository) GetOneByTelegramId(telegramId int) (*Model, error) {
	user := Model{}

	err := r.db.Get(&user, "SELECT * FROM users WHERE telegram_id=?", telegramId)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *Repository) CreateUser(user *Model) error {
	_, err := r.db.NamedExec(
		"INSERT IGNORE INTO users (telegram_id) VALUES (:telegram_id)",
		user,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) Update(user *Model) error {
	_, err := r.db.NamedExec(
		"UPDATE users SET lang=:lang WHERE id=:id",
		user,
	)
	if err != nil {
		return err
	}

	return nil
}

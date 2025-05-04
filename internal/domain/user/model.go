package user

import "time"

type Model struct {
	Id         int       `json:"id" db:"id"`
	TelegramId int       `json:"telegram_id" db:"telegram_id"`
	Lang       string    `json:"lang" db:"lang"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

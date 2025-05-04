package journey

import "time"

type Model struct {
	Id             int       `json:"id" db:"id"`
	UserId         int       `json:"user_id" db:"user_id"`
	QuestId        int       `json:"quest_id" db:"quest_id"`
	Progress       int       `json:"progress" db:"progress"`
	CompleteTillAt time.Time `json:"complete_till_at" db:"complete_till_at"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}

package quest_review

import "time"

type Model struct {
	Id        int       `json:"id" db:"id"`
	QuestId   int       `json:"quest_id" db:"quest_id"`
	UserId    int       `json:"user_id" db:"user_id"`
	Rating    int       `json:"rating" db:"rating"`
	Comment   string    `json:"comment" db:"comment"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

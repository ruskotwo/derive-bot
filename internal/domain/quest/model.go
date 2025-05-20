package quest

import "time"

type Model struct {
	Id          int       `json:"id" db:"id"`
	Title       string    `json:"title" db:"title"`
	Description string    `json:"description" db:"description"`
	File        *string   `json:"file" db:"file"`
	Lang        string    `json:"lang" db:"lang"`
	CategoryId  int       `json:"category_id" db:"category_id"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

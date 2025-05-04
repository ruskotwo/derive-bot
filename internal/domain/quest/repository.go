package quest

import (
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

func (r *Repository) GetOneById(id int) (*Model, error) {
	quest := Model{}

	err := r.db.Get(&quest, "SELECT * FROM quests WHERE id=?", id)
	if err != nil {
		return nil, err
	}

	return &quest, nil
}

func (r *Repository) GetRandomByLangAndCategoryId(lang string, categoryId int) (*Model, error) {
	quest := Model{}

	err := r.db.Get(&quest, "SELECT * FROM quests WHERE lang=? AND category_id=? ORDER BY RAND() LIMIT 1", lang, categoryId)
	if err != nil {
		return nil, err
	}

	return &quest, nil
}

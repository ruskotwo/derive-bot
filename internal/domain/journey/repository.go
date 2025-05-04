package journey

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

func (r *Repository) GetLastForUserId(userId int) (*Model, error) {
	j := Model{}

	err := r.db.Get(&j, "SELECT * FROM journeys WHERE user_id=? ORDER BY id DESC LIMIT 1", userId)
	if err != nil {
		return nil, err
	}

	return &j, nil
}

func (r *Repository) GetOneByIdAndUserId(id, userId int) (*Model, error) {
	j := Model{}

	err := r.db.Get(&j, "SELECT * FROM journeys WHERE id=? AND user_id=?", id, userId)
	if err != nil {
		return nil, err
	}

	return &j, nil
}

func (r *Repository) Save(j *Model) error {
	_, err := r.db.NamedExec("INSERT INTO journeys (user_id, quest_id, progress, complete_till_at) VALUES (:user_id, :quest_id, :progress, :complete_till_at)", j)
	if err != nil {
		return err
	}

	return nil
}

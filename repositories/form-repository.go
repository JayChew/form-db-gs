package repositories

import (
	"fmt"

	"github.com/JayChew/form-db-gs.git/models"
	"github.com/jmoiron/sqlx"
)

type FormRepo struct {
	DB *sqlx.DB
}

func (c *FormRepo) Create(form models.FormModel) (int64, error) {
	result, err := c.DB.Exec(
		`INSERT INTO forms (name, email) VALUES (?, ?)`,
		form.Name,
		form.Email,
	)

	if err != nil {
		return 0, fmt.Errorf("failed to insert form: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert id: %w", err)
	}
	
	return id, nil
}
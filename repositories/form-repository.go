package repositories

import (
	"fmt"

	"github.com/JayChew/form-db-gs.git/models"
	"github.com/jmoiron/sqlx"
)

type FormRepo struct {
	DB *sqlx.DB
}

func (c *FormRepo) Create(form models.FormModel) bool {
	_, err := c.DB.Exec(
		`INSERT INTO forms (
			name,
			email
		) VALUES (?,?)`,
		form.Name,
		form.Email,
	)

	if err != nil {
		fmt.Println(err)
		return false
	}
	
	return true;
}
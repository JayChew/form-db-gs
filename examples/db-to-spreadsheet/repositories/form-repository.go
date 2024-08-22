package repositories

import (
	"fmt"

	"github.com/JayChew/form-db-gs.git/examples/db-to-spreadsheet/helpers"
	"github.com/JayChew/form-db-gs.git/examples/db-to-spreadsheet/models"
	"github.com/jmoiron/sqlx"
)

type FormRepo struct {
	DB *sqlx.DB
}

func (c *FormRepo) Create(form models.FormModel) (int64, error) {
	query, values, err := helpers.GenerateInsertIntoQuery(form, "forms")
	if err != nil {
		fmt.Println("Error generating query:", err)
	}

	result, err := c.DB.Exec(
		query,
		values...,
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

func (c *FormRepo) GetAll() ([]models.FormModel, error) {
	var forms []models.FormModel
	var query = `SELECT * FROM forms`

	rows, err := c.DB.Queryx(query)
	if err != nil {
		return forms, fmt.Errorf("failed to select forms: %w", err)
	}
	defer func() {
		closeErr := rows.Close()
		if closeErr != nil {
			err = fmt.Errorf("failed to close rows: %w", closeErr)
		}
	}()

	for rows.Next() {
		form := models.FormModel{}
		err := rows.StructScan(&form)
		if err != nil {
			return forms, fmt.Errorf("failed to append form: %w", err)
		}
		forms = append(forms, form)
	}

	fmt.Println(forms)

	return forms, nil
}

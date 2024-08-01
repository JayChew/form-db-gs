package services

import (
	"github.com/JayChew/form-db-gs.git/models"
)

type FormService struct {
	IForm models.FormInterface
}

func (c *FormService) Create(name string, email string) (int64, error) {
	form := models.FormModel{
		Name: name,
		Email: email,
	}

	return c.IForm.Create(form);
}
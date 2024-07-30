package services

import (
	"github.com/JayChew/form-db-gs.git/models"
)

type FormService struct {
	IForm models.FormInterface
}

func (c *FormService) Create(name string, email string) (bool, models.FormModel) {
	form := models.FormModel{
		Name: name,
		Email: email,
	}

	c.IForm.Create(form);
	return true, form
}
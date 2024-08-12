package models

type FormModel struct {
	ID            int    `db:"id" json:"-" col:"Id"`
	Name          string `db:"name" json:"name" col:"Name"`
	Email         string `db:"email" json:"email" col:"Email"`
	ContactNumber string `db:"contact_number" json:"contact_number" col:"Contact number"`
}

type FormInterface interface {
	Create(FormModel) (int64, error)
	GetAll() ([]FormModel, error)
}

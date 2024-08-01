package models

type FormModel struct {
	ID int `db:"id" json:"-"`
	Name string `db:"name" json:"name"`
	Email string `db:"email" json:"email"`
}

type FormInterface interface {
	Create(FormModel) (int64, error)
}
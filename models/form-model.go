package models

import "time"

type FormModel struct {
	ID              int       `db:"id" json:"id" col:"Id"`
	ShortCode       string    `db:"short_code" json:"shortCode" col:"Short Code"`
	FormName        string    `db:"form_name" json:"formName" col:"Form Name"`
	Description     string    `db:"description" json:"description" col:"Description"`
	Notes           string    `db:"notes" json:"notes" col:"Notes"`
	History         string    `db:"history" json:"history" col:"History"`
	Status          int8      `db:"status" json:"status" col:"Status"`
	Priority        int16     `db:"priority" json:"priority" col:"Priority"`
	SubmissionCount int64     `db:"submission_count" json:"submissionCount" col:"Submission Count"`
	FormID          int64     `db:"form_id" json:"formId" col:"Form Id"`
	Cost            float64   `db:"cost" json:"cost" col:"Cost"`
	Rating          float32   `db:"rating" json:"rating" col:"Rating"`
	AverageScore    float64   `db:"average_score" json:"averageScore" col:"Average Score"`
	CreatedAt       time.Time `db:"created_at" json:"createdAt" col:"CreatedAt" mysql_format:"DATE"`
	UpdatedAt       time.Time `db:"updated_at" json:"updatedAt" col:"UpdatedAt" mysql_format:"DATETIME"`
	LastSubmission  time.Time `db:"last_submission" json:"lastSubmission" col:"LastSubmission" mysql_format:"TIMESTAMP"`
	FormTime        string    `db:"form_time" json:"formTime" col:"FormTime" mysql_format:"TIME"`

	Category       string `db:"category" json:"category" col:"Category"`
	Tags           string `db:"tags" json:"tags" col:"Tags"`
	StrongURL      string `db:"strong_url" json:"strongUrl" col:"Strong Url"`
	ContactNumber  string `db:"contact_number" json:"contactNumber" col:"Contact Number"`
	Email          string `db:"email" json:"email" col:"Email"`
	NRIC           string `db:"nric" json:"nric" col:"NRIC"`
	CurrencyCode   string `db:"currency_code" json:"currencyCode" col:"Currency Code"`
	CurrencySymbol string `db:"currency_symbol" json:"currencySymbol" col:"Currency Symbol"`

	AddressLine   string  `db:"address_line" json:"addressLine" col:"Address Line"`
	City          string  `db:"city" json:"city" col:"City"`
	StateProvince string  `db:"state_province" json:"stateProvince" col:"State Province"`
	PostalCode    string  `db:"postal_code" json:"postalCode" col:"Postal Code"`
	Country       string  `db:"country" json:"country" col:"Country"`
	AddressType   string  `db:"address_type" json:"addressType" col:"Address Type"`
	Latitude      float64 `db:"latitude" json:"latitude" col:"Latitude"`
	Longitude     float64 `db:"longitude" json:"longitude" col:"Longitude"`
}

type FormInterface interface {
	// Create(FormModel) (int64, error)
	GetAll() ([]FormModel, error)
}

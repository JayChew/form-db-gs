package main

import (
	"log"
	"time"

	"github.com/JayChew/form-db-gs.git/examples/db-to-spreadsheet/core"
	"github.com/JayChew/form-db-gs.git/examples/db-to-spreadsheet/models"
	"github.com/JayChew/form-db-gs.git/examples/db-to-spreadsheet/repositories"
	"github.com/JayChew/form-db-gs.git/examples/db-to-spreadsheet/services"
	"github.com/joho/godotenv"
)

func main() {
	db := core.InitalizeDB(
		"127.0.0.1",
		"3307",
		"root",
		"root",
		"form_db_gs",
	)

	formRepo := &repositories.FormRepo{DB: db}
	formService := services.FormService{IForm: formRepo}

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env config file")
	}

	formReq := models.CreateFormRequest{
		ShortCode:       "SC001",
		FormName:        "Sample Form",
		Description:     "This is a sample form description.",
		Notes:           "Some notes about the form.",
		History:         "Form creation history.",
		Status:          1, // Example status
		Priority:        2, // Example priority
		SubmissionCount: 100,
		FormID:          12345,
		Cost:            29.99,
		Rating:          4.5,
		AverageScore:    4.2,
		CreatedAt:       time.Now().AddDate(0, 0, -10), // 10 days ago
		UpdatedAt:       time.Now(),
		LastSubmission:  time.Now(),
		FormTime:        "15:30:00", // Example time
		Category:        "Support",
		Tags:            "Urgent,Important",
		StrongURL:       "http://example.com",
		ContactNumber:   "+60129533753",
		Email:           "example@example.com",
		NRIC:            "S1234567A",
		CurrencyCode:    "USD",
		CurrencySymbol:  "$",
		AddressLine:     "123 Main St",
		City:            "Sample City",
		StateProvince:   "Sample State",
		PostalCode:      "12345",
		Country:         "Sample Country",
		AddressType:     "Home",
		Latitude:        37.7749,
		Longitude:       -122.4194,
	}
	_, err = formService.Create(formReq)
	if err != nil {
		log.Fatalf("%v", err)
	}

	formService.SyncToGoogleSpreadSheet()
}

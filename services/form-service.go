package services

import (
	"fmt"
	"log"

	"github.com/JayChew/form-db-gs.git/models"
	"github.com/JayChew/form-db-gs.git/utils"
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

func (c *FormService) GetAll() ([]models.FormModel, error) {
	return c.IForm.GetAll();
}

func (c *FormService) SyncToGoogleSpreadSheet() {
	forms, err := c.GetAll();
	if err != nil {
		log.Fatalf("Unable to get forms. %v", err)
	}
	
	clearSheet := true;
	srv := utils.GoogleSpreadSheetSRV();
	SpreadsheetId := "16Fm__SoBsDr9WQenGt206G8CSiFzt33Cgm1LHkmnCr4"
	response, err := utils.AppendValuesToGoogleSpreadSheet(srv, "Sheet1", SpreadsheetId, utils.GenerateRows(forms), clearSheet)
	if err != nil {
		log.Fatalf("Unable to append values to the sheet: %v", err)
	}

	fmt.Printf("Appended values to the range: %s\n", response)
}
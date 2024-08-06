package services

import (
	"fmt"
	"log"
	"reflect"

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

	var rows [][]interface{}
	
	if len(forms) > 0 {
		val := reflect.ValueOf(forms[0])
		typ := val.Type()
		headers := []string{}
		for i := 0; i < val.NumField(); i++ {
			headers = append(headers, typ.Field(i).Name)
		}
		headerRow := make([]interface{}, len(headers))
		for i, h := range headers {
			headerRow[i] = h
		}
		rows = append(rows, headerRow)

		for _, form := range forms {
			formVal := reflect.ValueOf(form)
			formRow := []interface{}{}
			for i := 0; i < formVal.NumField(); i++ {
				formRow = append(formRow, formVal.Field(i).Interface())
			}
			rows = append(rows, formRow)
		}
	}
	// rows = append(rows, []interface{}{"ID", "Name", "Email"})
	// for _, form := range forms {
	// 	rows = append(rows, []interface{}{form.ID, form.Name, form.Email})
	// }

	clearSheet := true;
	srv := utils.GoogleSpreadSheetSRV();
	SpreadsheetId := "16Fm__SoBsDr9WQenGt206G8CSiFzt33Cgm1LHkmnCr4"
	response, err := utils.AppendValuesToGoogleSpreadSheet(srv, "Sheet1", SpreadsheetId, rows, clearSheet)
	if err != nil {
		log.Fatalf("Unable to append values to the sheet: %v", err)
	}

	fmt.Printf("Appended values to the range: %s\n", response)
}
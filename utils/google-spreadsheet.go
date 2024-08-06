package utils

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"strings"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

func GoogleSpreadSheetSRV() *sheets.Service {
	ctx := context.Background()

	// get bytes from base64 encoded google service accounts key
	credBytes, err := base64.StdEncoding.DecodeString(os.Getenv("GOOGLE_SERVICE_ACCOUNT_KEY"))
	if err != nil {
		log.Fatal(err)
	}

	// authenticate and get configuration
	config, err := google.JWTConfigFromJSON(credBytes, "https://www.googleapis.com/auth/spreadsheets")
	if err != nil {
		log.Fatal(err)
	}

	// create client with config and context
	client := config.Client(ctx)

	// create new service using client
	srv, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatal(err)
	}

	return srv;
}

// retrieves the name of a sheet by its ID from a Google Spreadsheet.
func GetGoogleSpreadSheetNameById(srv *sheets.Service, sheetId int64, spreadsheetId string) (string, error) {
	// Validate inputs
	if srv == nil {
		return "", fmt.Errorf("sheets service is nil")
	}
	if sheetId <= 0 {
		return "", fmt.Errorf("invalid sheet ID: %d", sheetId)
	}
	if strings.TrimSpace(spreadsheetId) == "" {
		return "", fmt.Errorf("spreadsheet ID is empty")
	}

	// Retrieve the properties of all sheets in the spreadsheet
	response, err := srv.Spreadsheets.Get(spreadsheetId).Fields("sheets(properties(sheetId,title))").Do()
	if err != nil {
		return "", fmt.Errorf("unable to retrieve spreadsheet data: %v", err)
	}

	// Iterate through the sheets to find the one with the matching sheetId
	for _, sheet := range response.Sheets {
		prop := sheet.Properties
		if prop.SheetId == sheetId {
			log.Printf("Found sheet: %s (ID: %d)", prop.Title, prop.SheetId)
			return prop.Title, nil
		}
	}

	// Return an error if no matching sheet is found
	return "", fmt.Errorf("no sheet found with ID: %d", sheetId)
}

// retrieves the ID of a sheet by its name from a Google Spreadsheet.
func GetGoogleSpreadSheetIdByName(srv *sheets.Service, sheetName string, spreadsheetId string) (int64, error) {
	// Validate inputs
	if srv == nil {
		return -1, fmt.Errorf("sheets service is nil")
	}
	if strings.TrimSpace(sheetName) == "" {
		return -1, fmt.Errorf("sheet name is empty")
	}
	if strings.TrimSpace(spreadsheetId) == "" {
		return -1, fmt.Errorf("spreadsheet ID is empty")
	}

	// Retrieve the properties of all sheets in the spreadsheet
	response, err := srv.Spreadsheets.Get(spreadsheetId).Fields("sheets(properties(sheetId,title))").Do()
	if err != nil {
		return -1, fmt.Errorf("unable to retrieve spreadsheet data: %v", err)
	}

	// Iterate through the sheets to find the one with the matching sheetName
	for _, sheet := range response.Sheets {
		prop := sheet.Properties
		if prop.Title == sheetName {
			log.Printf("Found sheet: %s (ID: %d)", prop.Title, prop.SheetId)
			return prop.SheetId, nil
		}
	}

	// Return an error if no matching sheet is found
	return -1, fmt.Errorf("no sheet found with name: %s", sheetName)
}

func ClearGoogleSpreadSheet(srv *sheets.Service, sheetName string, spreadsheetId string) error {
	clearRange := sheetName + "!A:Z" // Adjust range as necessary
	_, err := srv.Spreadsheets.Values.Clear(spreadsheetId, clearRange, &sheets.ClearValuesRequest{}).Context(context.Background()).Do()
	if err != nil {
		return fmt.Errorf("unable to clear spreadsheet data: %v", err)
	}
	return nil
}

func AppendValuesToGoogleSpreadSheet(srv *sheets.Service, sheetName string, spreadsheetId string, rows [][]interface{}, clearSheet bool) (string, error) {
	_, err := GetGoogleSpreadSheetIdByName(srv, sheetName, spreadsheetId)
	if err != nil {
		return "", fmt.Errorf("unable to retrieve sheet id: %v", err)
	}

	// Clear existing data if clearSheet is true
	if clearSheet {
		err := ClearGoogleSpreadSheet(srv, sheetName, spreadsheetId)
		if err != nil {
			return "", fmt.Errorf("unable to clear existing data: %v", err)
		}
	}

	row := &sheets.ValueRange{
		Values: rows,
	}

	// Append the values to the sheet
	response, err := srv.Spreadsheets.Values.Append(spreadsheetId, sheetName, row).
		ValueInputOption("USER_ENTERED").
		InsertDataOption("INSERT_ROWS").
		Context(context.Background()).Do()
	if err != nil {
		return "", fmt.Errorf("unable to save spreadsheet data: %v", err)
	}

	return response.Updates.UpdatedRange, nil
}
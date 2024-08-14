package utils

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"

	"github.com/JayChew/form-db-gs.git/helpers"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

func GoogleSpreadSheetSRV() *sheets.Service {
	ctx := context.Background()

	// get bytes from base64 encoded google service accounts key
	credBytes, err := base64.StdEncoding.DecodeString(os.Getenv("GOOGLE_SERVICE_ACCOUNT_KEY"))
	if err != nil {
		log.Printf("get bytes from base64 encoded google service accounts key failed: %v", err)
	}

	// authenticate and get configuration
	config, err := google.JWTConfigFromJSON(credBytes, "https://www.googleapis.com/auth/spreadsheets")
	if err != nil {
		log.Printf("authenticate and get configuration failed: %v", err)
	}

	// create client with config and context
	client := config.Client(ctx)

	// create new service using client
	srv, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Printf("create new service using client failed: %v", err)
	}

	return srv
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
	// clearRange := sheetName + "!A:Z" // Adjust range as necessary
	// _, err := srv.Spreadsheets.Values.Clear(spreadsheetId, clearRange, &sheets.ClearValuesRequest{}).Context(context.Background()).Do()
	// if err != nil {
	// 	return fmt.Errorf("unable to clear spreadsheet data: %v", err)
	// }
	// return nil

	sheetId, err := GetGoogleSpreadSheetIdByName(srv, sheetName, spreadsheetId)
	if err != nil {
		return fmt.Errorf("unable to retrieve sheet id: %v", err)
	}

	// Create a request to clear the entire sheet
	requests := []*sheets.Request{
		{
			UpdateCells: &sheets.UpdateCellsRequest{
				Range: &sheets.GridRange{
					SheetId: sheetId,
				},
				Fields: "*",
			},
		},
	}

	// Create the batch update request
	batchUpdateRequest := &sheets.BatchUpdateSpreadsheetRequest{
		Requests: requests,
	}

	// Execute the request
	_, err = srv.Spreadsheets.BatchUpdate(spreadsheetId, batchUpdateRequest).Context(context.Background()).Do()
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

func GenerateRows(tableData interface{}) [][]interface{} {
	var rows [][]interface{}

	// Get the reflection value of tableData
	v := reflect.ValueOf(tableData)

	// Check if tableData is a slice and if it has elements
	if v.Kind() != reflect.Slice || v.Len() == 0 {
		return rows
	}

	// Get the type of the elements in the slice
	elemType := v.Index(0).Type()
	headers := []string{}
	for i := 0; i < elemType.NumField(); i++ {
		// Extract the `col` tag value for human-readable headers
		colTag := elemType.Field(i).Tag.Get("col")
		if colTag == "" {
			colTag = elemType.Field(i).Name // Fallback to the field name if `col` tag is not present
		}
		headers = append(headers, colTag)
	}
	headerRow := make([]interface{}, len(headers))
	for i, h := range headers {
		headerRow[i] = h
	}
	rows = append(rows, headerRow)

	// Populate rows with data from each element in the slice
	for i := 0; i < v.Len(); i++ {
		elem := v.Index(i)
		dataRow := make([]interface{}, elem.NumField())
		for j := 0; j < elem.NumField(); j++ {
			field := elem.Field(j)

			switch field.Kind() {
			case reflect.String:
				fieldString := field.String()

				fieldString = helpers.FormatDateTime(fieldString)

				// Escape special characters
				fieldString = escapeSpecialCharacters(fieldString)

				dataRow[j] = fieldString
			default:
				// Directly append the field value to the dataRow
				dataRow[j] = field.Interface()
			}
		}
		rows = append(rows, dataRow)
	}

	return rows
}

// Function to escape special characters
func escapeSpecialCharacters(input string) string {
	// Check if the string starts with '+'
	if len(input) >= 2 && input[:2] == "'+" {
		// If the string starts with '+', return it as-is
		return input
	}

	// Replace special characters with their escaped versions (this is a basic example)
	input = strings.ReplaceAll(input, "'", "''")    // Escape single quotes
	input = strings.ReplaceAll(input, "\"", "\\\"") // Escape double quotes
	// Add more escaping rules as needed

	return input
}

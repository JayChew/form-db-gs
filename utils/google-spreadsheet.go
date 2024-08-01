package utils

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

func GetGoogleSpreadSheetConfig() *oauth2.Config {

	b, err := os.ReadFile("form-db-gs-dce6cdecba59.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	config, err := google.ConfigFromJSON(b, sheets.SpreadsheetsReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	return config;
}

func SaveGoogleSpreadsheetToken(token *oauth2.Token) {
	path := "token.json";
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func ReplaceGoogleSpreadsheetToken(authCode string) {
	config := GetGoogleSpreadSheetConfig();

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}

	SaveGoogleSpreadsheetToken(tok)
	fmt.Println("Token successfully replaced and saved.")
}

func GetGoogleSpreadsheetTokenFromFile() (*oauth2.Token, error) {
	f, err := os.Open("token.json")
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

func GetGoogleSpreadsheetTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

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

// GetGoogleSpreadSheetNameById retrieves the name of a sheet by its ID from a Google Spreadsheet.
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

// GetGoogleSpreadSheetIdByName retrieves the ID of a sheet by its name from a Google Spreadsheet.
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

func CreateNewGoogleSpreadSheet(srv *sheets.Service, sheetName string) string {
	spreadsheet := &sheets.Spreadsheet{
    Properties: &sheets.SpreadsheetProperties{
      Title: sheetName,
    },
	}

	spreadsheet, err := srv.Spreadsheets.Create(spreadsheet).Do()
	if err != nil {
		log.Fatalf("Unable to create spreadsheet. %v", err)
	}
	
	fmt.Printf("Spreadsheet ID: %s\n", spreadsheet.SpreadsheetId)
	return spreadsheet.SpreadsheetId
}

func AppendValueToTheGoogleSpreadSheet(srv *sheets.Service, sheetName string, spreadsheetId string, id string, name string, email string) (string, error) {
	_, err := GetGoogleSpreadSheetIdByName(srv, sheetName, spreadsheetId)
	if err != nil {
		return "", fmt.Errorf("unable to retrieve sheet id: %v", err)
	}

	row := &sheets.ValueRange{
		Values: [][]interface{}{{id, name, email}},
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
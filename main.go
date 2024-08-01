package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/JayChew/form-db-gs.git/core"
	"github.com/JayChew/form-db-gs.git/repositories"
	"github.com/JayChew/form-db-gs.git/services"
	"github.com/JayChew/form-db-gs.git/utils"
	"github.com/go-chi/chi"
	"github.com/joho/godotenv"
	"github.com/justinas/nosurf"
)

func main() {
	db := core.InitalizeDB(
		"127.0.0.1",
		"3307",
		"root",
		"root",
		"form_db_gs",
	);

	formRepo := &repositories.FormRepo{DB: db}
	formService := services.FormService{IForm: formRepo}

	name := "Jay Chew";
	email := "jaychew.3753@gmail.com"
	lastInsertId, err := formService.Create(name, email)
	if err != nil {
		log.Fatalf("%v", err)
	}

	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env config file")
	}

	srv := utils.GoogleSpreadSheetSRV();
	SpreadsheetId := "16Fm__SoBsDr9WQenGt206G8CSiFzt33Cgm1LHkmnCr4"
	var rows [][]interface{}
	rows = append(rows, []interface{}{lastInsertId, name, email})
	response, err := utils.AppendValuesToGoogleSpreadSheet(srv, "Sheet1", SpreadsheetId, rows)
	if err != nil {
		log.Fatalf("Unable to append values to the sheet: %v", err)
	}

	fmt.Printf("Appended values to the range: %s\n", response)

	r := chi.NewRouter()
	
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {});

	http.ListenAndServe(":8085", nosurf.New(r));
}
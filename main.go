package main

import (
	"log"
	"net/http"

	"github.com/JayChew/form-db-gs.git/utils"
	"github.com/go-chi/chi"
	"github.com/joho/godotenv"
	"github.com/justinas/nosurf"
)

func main() {
	// db := core.InitalizeDB(
	// 	"127.0.0.1",
	// 	"3307",
	// 	"root",
	// 	"root",
	// 	"form_db_gs",
	// );

	// formRepo := &repositories.FormRepo{DB: db}
	// formService := services.FormService{IForm: formRepo}

	// name := "Jay Chew";
	// email := "jaychew.3753@gmail.com"
	// _, form := formService.Create(name, email)
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env config file")
	}
	
	utils.GoogleSpreadSheetSRV();
	
	r := chi.NewRouter()
	
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {});

	http.ListenAndServe(":8085", nosurf.New(r));
}
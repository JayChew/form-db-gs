package main

import (
	"log"
	"net/http"

	"github.com/JayChew/form-db-gs.git/core"
	"github.com/JayChew/form-db-gs.git/repositories"
	"github.com/JayChew/form-db-gs.git/services"
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

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env config file")
	}

	name := "Jay Chew";
	email := "jaychew.3753@gmail.com"
	contact_number := "+60129533753"
	_, err = formService.Create(name, email, contact_number)
	if err != nil {
		log.Fatalf("%v", err)
	}

	formService.SyncToGoogleSpreadSheet()

	r := chi.NewRouter()
	
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {});

	http.ListenAndServe(":8085", nosurf.New(r));
}
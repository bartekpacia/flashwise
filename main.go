package main

import (
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB

func main() {
	router := setUpRouter()
	db = setUpDB()
	http.ListenAndServe(":8080", router)
}

func setUpRouter() http.Handler {
	router := mux.NewRouter()

	router.HandleFunc("/register", CreateUser).Methods("POST")

	router.HandleFunc("/flashcards", GetFlashcards).Methods("GET")
	router.HandleFunc("/flashcards", CreateFlashcard).Methods("POST")
	router.HandleFunc("/flashcards/{id}", UpdateFlashcard).Methods("PATCH")
	router.HandleFunc("/flashcards", DeleteFlashcard).Methods("DELETE")

	return handlers.LoggingHandler(os.Stdout, router)
}

func setUpDB() *sqlx.DB {
	db, err := sqlx.Open("mysql", "root:@(localhost:3306)/flashwise?parseTime=true")
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	return db
}

package main

import (
	"context"
	"net/http"
	"os"
	"strings"

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

	router.HandleFunc("/flashcards", AuthHandler(GetFlashcards)).Methods("GET")
	router.HandleFunc("/flashcards", AuthHandler(CreateFlashcard)).Methods("POST")
	router.HandleFunc("/flashcards/{id}", AuthHandler(UpdateFlashcard)).Methods("PATCH")
	router.HandleFunc("/flashcards", AuthHandler(DeleteFlashcard)).Methods("DELETE")

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

type ContextKey string

const ContextUserKey ContextKey = "user_id"

func AuthHandler(next func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			http.Error(w, "No Authorization header provided", http.StatusUnauthorized)
			return
		}

		splitToken := strings.Split(token, "Token")
		if len(splitToken) != 2 {
			http.Error(w, "Bearer token not in proper format", http.StatusUnauthorized)
			return
		}

		token = strings.TrimSpace(splitToken[1])

		var user User
		err := db.Get(&user, "SELECT * FROM users WHERE token = ?", token)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), ContextUserKey, user.ID)
		next(w, r.WithContext(ctx))
	})
}

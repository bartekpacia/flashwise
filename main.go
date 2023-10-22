package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB

func main() {
	var err error

	router := setUpRouter()
	db, err = setUpDB()
	if err != nil {
		log.Fatalln("failed to set up database:", err)
	}

	http.ListenAndServe(":8080", router)
}

func setUpRouter() http.Handler {
	router := mux.NewRouter()

	router.HandleFunc("/register", CreateUser).Methods("POST")

	router.HandleFunc("/flashcards", AuthHandler(GetFlashcards)).Methods("GET")
	router.HandleFunc("/flashcards", AuthHandler(CreateFlashcard)).Methods("POST")
	router.HandleFunc("/flashcards/{id}", AuthHandler(UpdateFlashcard)).Methods("PATCH")
	router.HandleFunc("/flashcards", AuthHandler(DeleteFlashcard)).Methods("DELETE")

	router.HandleFunc("/sets", AuthHandler(GetFlashcardSet)).Methods("GET")
	router.HandleFunc("/sets", AuthHandler(CreateFlashcardSet)).Methods("POST")

	router.HandleFunc("/quiz/generate", AuthHandler(GenerateQuiz)).Methods("POST")
	router.HandleFunc("/quiz/check", AuthHandler(CheckQuiz)).Methods("PUT")

	return handlers.LoggingHandler(os.Stdout, router)
}

func setUpDB() (*sqlx.DB, error) {
	host, ok := os.LookupEnv("MYSQL_HOST")
	if !ok {
		return nil, errors.New("MYSQL_HOST env var not set")
	}

	user, ok := os.LookupEnv("MYSQL_USER")
	if !ok {
		return nil, errors.New("MYSQL_USER env var not set")
	}

	password, ok := os.LookupEnv("MYSQL_PASSWORD")
	if !ok {
		return nil, errors.New("MYSQL_PASSWORD env var not set")
	}

	dbName, ok := os.LookupEnv("MYSQL_DB")
	if !ok {
		return nil, errors.New("MYSQL_DB env var not set")
	}

	connString := fmt.Sprintf("%s:%s@(%s:3306)/%s?parseTime=true", user, password, host, dbName)

	var database *sqlx.DB
	var err error
	fails := 0
	maxFails := 10
	for {
		if fails >= maxFails {
			return nil, fmt.Errorf("failed to connect to database after %d fails", maxFails)
		}

		if fails > 0 {
			time.Sleep(1 * time.Second)
		}

		database, err = sqlx.Open("mysql", connString)
		if err != nil {
			log.Println("failed to connect to database:", err)
			fails++
			continue
		}

		err = database.Ping()
		if err != nil {
			log.Println("failed to ping database:", err)
			fails++
			continue
		}

		break
	}

	log.Println("successfully connected to database")

	return database, nil
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

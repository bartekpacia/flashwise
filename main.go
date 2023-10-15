package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

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

func GetFlashcards(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Queryx("SELECT * FROM flashcards")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error while executing query: %v\n", err)
		return
	}

	flashcards := make([]Flashcard, 0)
	for rows.Next() {
		var flashcard Flashcard
		err := rows.StructScan(&flashcard)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error while iterating over rows: %v\n", err)
			return
		}
		flashcards = append(flashcards, flashcard)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(flashcards)
}

func CreateFlashcard(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error while reading request body: %v\n", err)
		return
	}

	var f Flashcard
	err = json.Unmarshal(body, &f)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error while unmarshalling request body: %v\n", err)
		return
	}

	stmt := "INSERT INTO flashcards (front, back, author_id, set_id) VALUES (?, ?, ?, ?)"
	_, err = db.Exec(stmt, f.Front, f.Back, f.AuthorID, f.SetID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error while executing query: %v\n", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func UpdateFlashcard(w http.ResponseWriter, r *http.Request) {
	/*
		params := mux.Vars(r)
		id := params["id"]

		for index, flashcard := range flashcards {
			if flashcard.ID == id {
				var updatedFlashcard Flashcard
				_ = json.NewDecoder(r.Body).Decode(&updatedFlashcard)
				flashcards[index] = updatedFlashcard
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(updatedFlashcard)
				return
			}
		}

		http.NotFound(w, r)
	*/
}

func DeleteFlashcard(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("flashcard_id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Validation error: flashcard_id query parameter is required\n")
		return
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Validation error: flashcard_id query parameter must be an integer\n")
		return
	}

	result, err := db.Exec("DELETE FROM flashcards WHERE id = ?", idInt)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error while executing query: %v\n", err)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error while getting rows affected: %v\n", err)
		return
	}

	if rowsAffected == 0 {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Error: flashcard with id %d not found\n", idInt)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

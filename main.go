package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var flashcards []Flashcard

func main() {
	router := mux.NewRouter()

	// Dummy data for demonstration purposes
	flashcards = append(
		flashcards,
		Flashcard{
			ID:       "1",
			Question: "What is Go?",
			Answer:   "Go is a programming language.",
		},
	)

	// Route handling for /flashcards
	router.HandleFunc("/flashcards", GetFlashcards).Methods("GET")
	router.HandleFunc("/flashcards", CreateFlashcard).Methods("POST")
	router.HandleFunc("/flashcards/{id}", UpdateFlashcard).Methods("PATCH")
	router.HandleFunc("/flashcards/{id}", DeleteFlashcard).Methods("DELETE")

	// Start the server
	http.ListenAndServe(":8080", router)
	log.Println("server exited")
}

func GetFlashcards(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(flashcards)
}

func CreateFlashcard(w http.ResponseWriter, r *http.Request) {
	var flashcard Flashcard
	body, _ := io.ReadAll(r.Body)
	json.Unmarshal(body, &flashcard)
	flashcards = append(flashcards, flashcard)

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusCreated)
}

func UpdateFlashcard(w http.ResponseWriter, r *http.Request) {
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
}

func DeleteFlashcard(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	for index, flashcard := range flashcards {
		if flashcard.ID == id {
			flashcards = append(flashcards[:index], flashcards[index+1:]...)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}

	http.NotFound(w, r)
}

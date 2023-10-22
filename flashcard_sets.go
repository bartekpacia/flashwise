package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func GetFlashcardSet(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(ContextUserKey).(uint64)
	if !ok {
		http.Error(w, "user ID is not present in context", http.StatusInternalServerError)
		return
	}

	var sets []FlashcardSet
	err := db.Select(&sets, "SELECT * FROM flashcard_sets WHERE author_id = ?", userID)
	if err != nil {
		http.Error(w, fmt.Sprintln("failed to execute query:", err), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(sets)
	if err != nil {
		http.Error(w, fmt.Sprintln("failed to encode response", err), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
}

func CreateFlashcardSet(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(ContextUserKey).(uint64)
	if !ok {
		http.Error(w, "user ID is not present in context", http.StatusInternalServerError)
		return
	}

	var body CreateFlashcardSetRequest
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to decode request body: %v\n", err), http.StatusBadRequest)
		return
	}

	stmt := "INSERT INTO flashcard_sets (author_id, description) VALUES (?, ?)"
	_, err = db.Exec(stmt, userID, body.Description)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error while executing query: %v\n", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func CreateFlashcardSet(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(ContextUserKey).(uint64)
	if !ok {
		http.Error(w, "user ID is not present in context", http.StatusInternalServerError)
		return
	}

	var body CreateFlashcardSetRequest
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, "Error while decoding request body: %v\n", http.StatusBadRequest)
		return
	}

	// TODO: Check if set_id belongs to author_id

	stmt := "INSERT INTO flashcard_sets (author_id, description) VALUES (?, ?)"
	_, err = db.Exec(stmt, userID, body.Description)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error while executing query: %v\n", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

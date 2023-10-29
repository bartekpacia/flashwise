package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func GetFlashcardSet(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(ContextUserKey).(uint64)
	if !ok {
		http.Error(w, "user ID is not present in context", http.StatusInternalServerError)
		return
	}

	sets := make([]FlashcardSet, 0)
	err := db.Select(&sets, "SELECT * FROM flashcard_sets WHERE author_id = ?", userID)
	if err != nil {
		http.Error(w, fmt.Sprintln("failed to execute query:", err), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(sets)
	if err != nil {
		http.Error(w, fmt.Sprintln("failed to encode response", err), http.StatusInternalServerError)
		return
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
		http.Error(w, fmt.Sprintln("failed to decode request body:", err), http.StatusBadRequest)
		return
	}

	stmt := "INSERT INTO flashcard_sets (author_id, title, is_public, category_id) VALUES (?, ?, ?, ?)"
	result, err := db.Exec(stmt, userID, body.Title, body.Public, body.CategoryID)
	if err != nil {
		http.Error(w, fmt.Sprintln("failed to execute query:", err), http.StatusInternalServerError)
		return
	}

	id, _ := result.LastInsertId()
	err = json.NewEncoder(w).Encode(map[string]uint64{"id": uint64(id)})
	if err != nil {
		http.Error(w, fmt.Sprintln("failed to encode response", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

func DeleteFlashcardSet(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(ContextUserKey).(uint64)
	if !ok {
		http.Error(w, "user ID is not present in context", http.StatusInternalServerError)
		return
	}

	// TODO: delete all flashcards belonging to user

	id, err := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		http.Error(w, "Validation error: id route variable is missing or not uint64", http.StatusBadRequest)
		return
	}

	stmt := "DELETE FROM flashcard_sets WHERE id = ? AND author_id = ?"
	_, err = db.Exec(stmt, id, userID)
	if err != nil {
		http.Error(w, fmt.Sprintln("failed to execute query:", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

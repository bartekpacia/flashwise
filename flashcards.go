package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func GetFlashcards(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(ContextUserKey).(uint64)
	if !ok {
		http.Error(w, "user ID is not present in context", http.StatusInternalServerError)
		return
	}

	flashcardID := r.URL.Query().Get("flashcard_id")
	setID := r.URL.Query().Get("set_id")

	if flashcardID != "" && setID != "" {
		http.Error(w, "Only one of flashcard_id and set_id can be specified", http.StatusBadRequest)
		return
	} else if flashcardID != "" {
		// Return flashcard for the user with the specified ID
		flashcardIDInt, err := strconv.ParseUint(flashcardID, 10, 64)
		if err != nil {
			http.Error(w, "Validation error: flashcard_id query parameter must be an integer\n", http.StatusBadRequest)
			return
		}

		var flashcard Flashcard
		err = db.Get(&flashcard, "SELECT * FROM flashcards WHERE author_id = ? AND id = ?", userID, flashcardIDInt)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error while executing query: %v\n", err), http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(flashcard)
	} else if setID != "" {
		// Return all flashcards for the user with the specified set ID
		setIDInt, err := strconv.ParseUint(setID, 10, 64)
		if err != nil {
			http.Error(w, "Validation error: set_id query parameter must be an integer\n", http.StatusBadRequest)
			return
		}

		// Check if set exists
		var exists bool
		row := db.QueryRow("SELECT EXISTS(SELECT 1 FROM flashcard_sets WHERE id = ?)", setIDInt)
		err = row.Scan(&exists)
		if err != nil && err != sql.ErrNoRows {
			http.Error(w, fmt.Sprintf("Error while executing query: %v\n", err), http.StatusInternalServerError)
			return
		}

		if !exists {
			http.Error(w, fmt.Sprintf("Error: set with id %d not found\n", setIDInt), http.StatusNotFound)
			return
		}

		rows, err := db.Queryx("SELECT * FROM flashcards WHERE author_id = ? AND set_id = ?", userID, setIDInt)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error while executing query: %v\n", err), http.StatusInternalServerError)
		}

		flashcards := make([]Flashcard, 0)
		for rows.Next() {
			var flashcard Flashcard
			err := rows.StructScan(&flashcard)
			if err != nil {
				http.Error(w, fmt.Sprintf("Error while iterating over rows: %v\n", err), http.StatusInternalServerError)
				return
			}
			flashcards = append(flashcards, flashcard)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(flashcards)
	} else {
		// Return all flashcards for the user

		rows, err := db.Queryx("SELECT * FROM flashcards WHERE author_id = ?", userID)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error while executing query: %v\n", err), http.StatusInternalServerError)
			return
		}

		flashcards := make([]Flashcard, 0)
		for rows.Next() {
			var flashcard Flashcard
			err := rows.StructScan(&flashcard)
			if err != nil {
				http.Error(w, fmt.Sprintf("Error while iterating over rows: %v\n", err), http.StatusInternalServerError)
				return
			}
			flashcards = append(flashcards, flashcard)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(flashcards)
	}
}

func CreateFlashcard(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(ContextUserKey).(uint64)
	if !ok {
		http.Error(w, "user ID is not present in context", http.StatusInternalServerError)
		return
	}

	var body CreateFlashcardRequest
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, "Error while decoding request body: %v\n", http.StatusBadRequest)
		return
	}

	// TODO: Check if set_id belongs to author_id

	stmt := "INSERT INTO flashcards (front, back, author_id, set_id) VALUES (?, ?, ?, ?)"
	_, err = db.Exec(stmt, body.Front, body.Back, userID, body.SetID)
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
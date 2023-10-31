package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/bartekpacia/flashwise/internal/domain"
	"github.com/gorilla/mux"
)

type createFlashcardRequest struct {
	Front string `json:"front"`
	Back  string `json:"back"`
	SetID uint64 `json:"flashcard_set"`
}

func (a *api) getFlashcards(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(uint64)
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

func (a *api) createFlashcard(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(uint64)
	if !ok {
		http.Error(w, "user ID is not present in context", http.StatusInternalServerError)
		return
	}

	var body createFlashcardRequest
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, fmt.Sprintln("failed to decode request body:", err), http.StatusBadRequest)
		return
	}

	// var set FlashcardSet
	// err = db.Get(&set, "SELECT * FROM flashcard_sets WHERE id = ?", body.SetID)
	// if err != nil {
	// 	if err == sql.ErrNoRows {
	// 		http.Error(w, fmt.Sprintf("set with ID %d does not exist\n", body.SetID), http.StatusNotFound)
	// 		return
	// 	} else {
	// 		http.Error(w, fmt.Sprintln("failed to execute query:", err), http.StatusInternalServerError)
	// 		return
	// 	}
	// }

	// if set.AuthorID != userID {
	// 	http.Error(w, fmt.Sprintf("set with ID %d does not belong to user %d\n", body.SetID, userID), http.StatusNotFound)
	// 	return
	// }

	// stmt := `
	// 	INSERT INTO flashcards
	// 		(front, back, author_id, set_id)
	// 	VALUES
	// 		(?, ?, ?, ?)`
	// _, err = db.Exec(stmt, body.Front, body.Back, userID, body.SetID)
	// if err != nil {
	// 	http.Error(w, fmt.Sprintln("failed to execute query:", err), http.StatusInternalServerError)
	// 	return
	// }

	w.WriteHeader(http.StatusCreated)
}

func (a *api) updateFlashcard(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	id, err := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		http.Error(w, "Validation error: id route variable is missing or not uint64", http.StatusBadRequest)
		return
	}

	var body createFlashcardRequest
	err = json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, fmt.Sprintln("failed to decode request body:", err), http.StatusBadRequest)
		return
	}

	a.flashcardRepo.Update(ctx, id, body.Front, body.Back, body.SetID)

	w.WriteHeader(http.StatusOK)
}

func (a *api) deleteFlashcard(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	id, err := strconv.ParseUint(r.URL.Query().Get("flashcard_id"), 10, 64)
	if err != nil {
		http.Error(w, "Validation error: flashcard_id query parameter is missing or not uint64", http.StatusBadRequest)
		return
	}

	err = a.flashcardRepo.Delete(ctx, id)
	if err != nil {
		if err == domain.ErrNotFound {
			http.Error(w, fmt.Sprintf("error: flashcard with id %d not found\n", id), http.StatusNotFound)
			return
		} else {
			http.Error(w, fmt.Sprintf("error while executing query: %v\n", err), http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}
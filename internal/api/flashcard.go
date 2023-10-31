package api

import (
	"context"
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
	_, ok := r.Context().Value("user_id").(uint64)
	if !ok {
		http.Error(w, "user ID is not present in context", http.StatusInternalServerError)
		return
	}

	flashcardID := r.URL.Query().Get("flashcard_id")
	setID := r.URL.Query().Get("set_id")

	if flashcardID != "" && setID != "" {
		http.Error(w, "Only one of flashcard_id and set_id can be specified", http.StatusBadRequest)
		return
	}

	if flashcardID != "" {
		id, err := strconv.ParseUint(flashcardID, 10, 64)
		if err != nil {
			http.Error(w, "flashcard_id query parameter must be an integer\n", http.StatusBadRequest)
			return
		}

		flashcard, err := a.flashcardRepo.GetByID(r.Context(), id)
		if err != nil {
			if err == domain.ErrNotFound {
				w.WriteHeader(http.StatusNotFound)
				return
			}

			http.Error(w, fmt.Sprintln("failed to execute query:", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(flashcard)
	}

	if setID != "" {
		// Return all flashcards for the user with the specified set ID
		setIDInt, err := strconv.ParseUint(setID, 10, 64)
		if err != nil {
			http.Error(w, "set_id query parameter must be an integer\n", http.StatusBadRequest)
			return
		}

		flashcards, err := a.flashcardRepo.GetBySetID(r.Context(), setIDInt)
		if err != nil {
			http.Error(w, fmt.Sprintln("failed to execute query:", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(flashcards)
	} else {
		flashcards, err := a.flashcardRepo.GetAll(r.Context())
		if err != nil {
			http.Error(w, fmt.Sprintln("failed to execute query:", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(flashcards)
	}
}

func (a *api) createFlashcard(w http.ResponseWriter, r *http.Request) {
	_, ok := r.Context().Value("user_id").(uint64)
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

	id, err := a.flashcardRepo.Create(r.Context(), body.Front, body.Back, body.SetID)
	if err != nil {
		http.Error(w, fmt.Sprintln("failed to create flashcard:", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]uint64{"id": *id})
	w.WriteHeader(http.StatusCreated)
}

func (a *api) updateFlashcard(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	id, err := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		http.Error(w, "id route variable is missing or not uint64", http.StatusBadRequest)
		return
	}

	var body createFlashcardRequest
	err = json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, fmt.Sprintln("failed to decode request body:", err), http.StatusBadRequest)
		return
	}

	err = a.flashcardRepo.Update(ctx, id, body.Front, body.Back, body.SetID)
	if err != nil {
		http.Error(w, fmt.Sprintln("failed to update flashcard:", err), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
}

func (a *api) deleteFlashcard(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	id, err := strconv.ParseUint(r.URL.Query().Get("flashcard_id"), 10, 64)
	if err != nil {
		http.Error(w, "flashcard_id query parameter is missing or not uint64", http.StatusBadRequest)
		return
	}

	err = a.flashcardRepo.Delete(ctx, id)
	if err != nil {
		if err == domain.ErrNotFound {
			w.WriteHeader(http.StatusNotFound)
			return
		} else {
			http.Error(w, fmt.Sprintln("failed to execute query:", err), http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}

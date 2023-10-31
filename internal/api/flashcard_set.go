package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/bartekpacia/flashwise/internal/domain"
	"github.com/gorilla/mux"
)

type createFlashcardSetRequest struct {
	Title      string `json:"name"`
	Public     bool   `json:"is_public"`
	CategoryID uint64 `json:"category"`
}

func (a *api) getFlashcardSet(w http.ResponseWriter, r *http.Request) {
	// userID, ok := r.Context().Value("user_id").(uint64)
	// if !ok {
	// 	http.Error(w, "user ID is not present in context", http.StatusInternalServerError)
	// 	return
	// }

	// setID := r.URL.Query().Get("flashcard_set_id")
	// authorID := r.URL.Query().Get("author_id")

	// // a.flashcardSetRepo.

	// if err != nil {
	// 	http.Error(w, fmt.Sprintln("failed to execute query:", err), http.StatusInternalServerError)
	// 	return
	// }

	// err = json.NewEncoder(w).Encode(sets)
	// if err != nil {
	// 	http.Error(w, fmt.Sprintln("failed to encode response", err), http.StatusInternalServerError)
	// 	return
	// }

	// w.Header().Set("Content-Type", "application/json")
}

func (a *api) createFlashcardSet(w http.ResponseWriter, r *http.Request) {
	_, ok := r.Context().Value("user_id").(uint64)
	if !ok {
		http.Error(w, domain.ErrNoUserID.Error(), http.StatusInternalServerError)
		return
	}

	var body createFlashcardSetRequest
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, fmt.Sprintln("failed to decode request body:", err), http.StatusBadRequest)
		return
	}

	id, err := a.flashcardSetRepo.Create(r.Context(), body.Title, body.Public, body.CategoryID)
	if err != nil {
		http.Error(w, fmt.Sprintln("failed to execute query:", err), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(map[string]uint64{"id": *id})
	if err != nil {
		http.Error(w, fmt.Sprintln("failed to encode response", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

func (a *api) deleteFlashcardSet(w http.ResponseWriter, r *http.Request) {
	_, ok := r.Context().Value("user_id").(uint64)
	if !ok {
		http.Error(w, "user ID is not present in context", http.StatusInternalServerError)
		return
	}

	// TODO: delete all flashcards belonging to user?

	id, err := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		http.Error(w, "id route variable is missing or not uint64", http.StatusBadRequest)
		return
	}

	err = a.flashcardSetRepo.Delete(r.Context(), id)
	if err != nil {
		http.Error(w, fmt.Sprintln("failed to execute query:", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

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
	_, ok := r.Context().Value("user_id").(uint64)
	if !ok {
		http.Error(w, "user ID is not present in context", http.StatusInternalServerError)
		return
	}

	if len(r.URL.Query()) == 0 {
		flashcards, err := a.flashcardSetRepo.GetAll(r.Context(), false)
		if err != nil {
			http.Error(w, fmt.Sprintln("failed to execute query:", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(flashcards)
		return
	}

	if r.URL.Query().Has("include_private") {
		includePrivate := r.URL.Query().Get("include_private") == "true"

		flashcards, err := a.flashcardSetRepo.GetAll(r.Context(), includePrivate)
		if err != nil {
			http.Error(w, fmt.Sprintln("failed to execute query:", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(flashcards)
		return
	}

	if r.URL.Query().Has("flashcard_set_id") {
		setID, err := strconv.ParseUint(r.URL.Query().Get("flashcard_set_id"), 10, 64)
		if err != nil {
			http.Error(w, "flashcard_set_id query parameter must be an integer\n", http.StatusBadRequest)
			return
		}

		flashcardSet, err := a.flashcardSetRepo.GetByID(r.Context(), setID)
		if err != nil {
			if err == domain.ErrNotFound {
				w.WriteHeader(http.StatusNotFound)
				return
			}

			if err == domain.ErrNoAccess {
				w.WriteHeader(http.StatusForbidden)
			}

			http.Error(w, fmt.Sprintln("failed to execute query:", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(flashcardSet)
		return
	}

	if r.URL.Query().Has("category_id") {
		categoryID, err := strconv.ParseUint(r.URL.Query().Get("category_id"), 10, 64)
		if err != nil {
			http.Error(w, "category_id query parameter must be an integer\n", http.StatusBadRequest)
			return
		}

		flashcardSets, err := a.flashcardSetRepo.GetByCategory(r.Context(), categoryID)
		if err != nil {
			http.Error(w, fmt.Sprintln("failed to execute query:", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(flashcardSets)
		return
	}

	if r.URL.Query().Has("name") {
		name := r.URL.Query().Get("name")

		flashcardSets, err := a.flashcardSetRepo.GetByNameContains(r.Context(), name)
		if err != nil {
			http.Error(w, fmt.Sprintln("failed to execute query:", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(flashcardSets)
		return
	}

	http.Error(w, "invalid query parameters", http.StatusBadRequest)
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]uint64{"id": id})
}

func (a *api) updateFlashcardSet(w http.ResponseWriter, r *http.Request) {
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

	id, err := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		http.Error(w, "id route variable is missing or not uint64", http.StatusBadRequest)
		return
	}

	err = a.flashcardSetRepo.Update(r.Context(), id, body.Title, body.Public, body.CategoryID)
	if err != nil {
		http.Error(w, fmt.Sprintln("failed to execute query:", err), http.StatusInternalServerError)
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

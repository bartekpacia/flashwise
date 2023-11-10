package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/bartekpacia/flashwise/internal/domain"
)

type generateQuizRequest struct {
	SetID uint64 `json:"flashcard_set_id"`
}

func (a *api) generateQuiz(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(uint64)
	if !ok {
		http.Error(w, "user ID is not present in context", http.StatusInternalServerError)
		return
	}
	_ = userID

	var body generateQuizRequest
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, "failed to decode request body", http.StatusBadRequest)
		return
	}

	quiz, err := a.quizRepo.Generate(r.Context(), body.SetID)
	if err != nil {
		if err == domain.ErrNotFound {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		if err == domain.ErrNoAccess {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(quiz)
}

type checkQuizRequest struct {
	Answers map[string]string `json:"answers"`
}

func (a *api) checkQuiz(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(uint64)
	if !ok {
		http.Error(w, "user ID is not present in context", http.StatusInternalServerError)
		return
	}
	_ = userID

	var body checkQuizRequest
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, "failed to decode request body", http.StatusBadRequest)
		return
	}

	fmt.Println("answers", body.Answers)

	quiz, err := a.quizRepo.Check(r.Context(), body.Answers)
	if err != nil {
		if err == domain.ErrNotFound {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		if err == domain.ErrNoAccess {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(quiz)
}

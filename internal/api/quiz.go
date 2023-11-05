package api

import (
	"encoding/json"
	"net/http"
)

type generateQuizRequest struct {
	SetID uint64 `json:"flashcard_set"`
}

func (a *api) generateQuiz(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(uint64)
	if !ok {
		http.Error(w, "user ID is not present in context", http.StatusInternalServerError)
		return
	}

	var body generateQuizRequest
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, "failed to decode request body", http.StatusBadRequest)
		return
	}

	_ = userID

	// TODO: Check if flashcard set is public. If is private, check if it belongs to userID.

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

	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

func (a *api) checkQuiz(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(uint64)
	if !ok {
		http.Error(w, "user ID is not present in context", http.StatusInternalServerError)
		return
	}

	_ = userID

	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

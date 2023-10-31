package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type createUserRequest struct {
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password1 string `json:"password1"`
	Password2 string `json:"password2"`
}

func (a *api) createUser(w http.ResponseWriter, r *http.Request) {
	var body createUserRequest

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if body.Password1 != body.Password2 {
		http.Error(w, "Passwords do not match", http.StatusBadRequest)
		return
	}

	if len(body.Password1) < 8 {
		http.Error(w, fmt.Sprintf("Password must be at least 8 bytes long (got %d)", len(body.Password1)), http.StatusBadRequest)
		return
	}

	if len(body.Password1) > 72 {
		http.Error(w, "Password must be at most 72 bytes long", http.StatusBadRequest)
		return
	}

	user, err := a.userRepo.Register(r.Context(), body.Username, body.Email, body.Password1)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	tokenJSON := map[string]string{"token": user.Token}
	json.NewEncoder(w).Encode(tokenJSON)
}

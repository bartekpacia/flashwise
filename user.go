package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var body CreateUserRequest

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
		http.Error(w, "Password must be at least 8 bytes long", http.StatusBadRequest)
		return
	}

	if len(body.Password1) > 72 {
		http.Error(w, "Password must be at most 72 bytes long", http.StatusBadRequest)
		return
	}

	passwordHash, _ := hashPassword(body.Password1)
	token := generateToken()
	stmt := "INSERT INTO users (username, email, password_hash, token) VALUES (?, ?, ?, ?)"
	_, err = db.Exec(stmt, body.Username, body.Email, passwordHash, token)
	if err != nil {
		http.Error(w, fmt.Sprint("Failed to insert user in db:", err.Error()), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	tokenJSON := map[string]string{"token": token}
	json.NewEncoder(w).Encode(tokenJSON)
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func generateToken() string {
	randomBytes := make([]byte, 20)
	rand.Read(randomBytes)

	return hex.EncodeToString(randomBytes)
}

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func GetCategories(w http.ResponseWriter, r *http.Request) {
	categories := make([]Category, 0)
	err := db.Select(&categories, "SELECT * FROM categories")
	if err != nil {
		http.Error(w, fmt.Sprintln("failed to execute query:", err), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(categories)
	if err != nil {
		http.Error(w, fmt.Sprintln("failed to encode response", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
}

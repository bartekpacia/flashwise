package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (a *api) getCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := a.categoryRepo.GetAll(r.Context())
	if err != nil {
		http.Error(w, fmt.Sprintln("failed to execute query", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categories)
}

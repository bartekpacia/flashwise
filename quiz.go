package main

import "net/http"

func GenerateQuiz(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(ContextUserKey).(uint64)
	if !ok {
		http.Error(w, "user ID is not present in context", http.StatusInternalServerError)
		return
	}

	_ = userID

	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

func CheckQuiz(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(ContextUserKey).(uint64)
	if !ok {
		http.Error(w, "user ID is not present in context", http.StatusInternalServerError)
		return
	}

	_ = userID

	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

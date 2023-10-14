package main

import "time"

type Flashcard struct {
	ID           string       `json:"id"` // Primary key
	Front        string       `json:"front"`
	Back         string       `json:"back"`
	LastModified *time.Time   `json:"last_modified"`
	Set          FlashcardSet `json:"set"`    // Foreign key to FlashcardSet
	Author       User         `json:"author"` // Foreign key to User
}

type FlashcardSet struct {
	ID     string `json:"id"`     // Primary key
	Author User   `json:"author"` // Foreign key to User
}

type User struct {
	ID   string `json:"id"` // Primary key
	Name string `json:"name"`
}

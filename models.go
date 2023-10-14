package main

import "time"

type Flashcard struct {
	ID         uint64     `json:"id"` // Primary key
	Front      string     `json:"front"`
	Back       string     `json:"back"`
	CreatedAt  time.Time  `json:"created_at" db:"created_at"`
	ModifiedAt *time.Time `json:"modified_at" db:"modified_at"`
	SetID      uint64     `json:"set" db:"set_id"`       // Foreign key to FlashcardSet
	AuthorID   uint64     `json:"author" db:"author_id"` // Foreign key to User
}

type FlashcardSet struct {
	ID     string `json:"id"`     // Primary key
	Author User   `json:"author"` // Foreign key to User
}

type User struct {
	ID   string `json:"id"` // Primary key
	Name string `json:"name"`
}

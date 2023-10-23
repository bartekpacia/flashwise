package main

import "time"

type Flashcard struct {
	ID         uint64     `json:"id"` // Primary key
	Front      string     `json:"front"`
	Back       string     `json:"back"`
	CreatedAt  time.Time  `json:"created_at" db:"created_at"`
	ModifiedAt *time.Time `json:"modified_at" db:"modified_at"`
	SetID      uint64     `json:"flashcard_set" db:"set_id"` // Foreign key to FlashcardSet
	AuthorID   uint64     `json:"author_id" db:"author_id"`  // Foreign key to User
}

type FlashcardSet struct {
	ID          uint64     `json:"id" db:"id" ` // Primary key
	Status      string     `json:"stats" db:"status" `
	Description string     `json:"name" db:"description" `
	CreatedAt   time.Time  `json:"created_at" db:"created_at" `
	ModifiedAt  *time.Time `json:"modified_at" db:"modified_at" `
	AuthorID    uint64     `json:"author_id" db:"author_id" ` // Foreign key to User
}

type User struct {
	ID           uint64    `db:"id"` // Primary key
	Username     string    `db:"username"`
	Email        string    `db:"email"`
	CreatedAt    time.Time `db:"created_at"`
	Admin        bool      `db:"is_admin"`
	PasswordHash string    `db:"password_hash"`
	Token        string    `db:"token"`
}

type Quiz struct {
	ID             uint64 `db:"id"`            // Primary key
	FlashcardSetID uint64 `db:"flashcard_set"` // Foreign key to FlashcardSet
	AuthorID       uint64 `db:"author_id"`     // Foreign key to User
}

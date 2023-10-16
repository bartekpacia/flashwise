package main

import "time"

type CreateUserRequest struct {
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password1 string `json:"password1"`
	Password2 string `json:"password2"`
}

type CreateFlashcardRequest struct {
	Front string `json:"front"`
	Back  string `json:"back"`
	SetID uint64 `json:"set_id"`
}

type CreateFlashcardSetRequest struct {
	Status      bool   `json:"status"`
	Description string `json:"description"`
}

type Flashcard struct {
	ID         uint64     `json:"id"` // Primary key
	Front      string     `json:"front"`
	Back       string     `json:"back"`
	CreatedAt  time.Time  `json:"created_at" db:"created_at"`
	ModifiedAt *time.Time `json:"modified_at" db:"modified_at"`
	SetID      uint64     `json:"set_id" db:"set_id"`       // Foreign key to FlashcardSet
	AuthorID   uint64     `json:"author_id" db:"author_id"` // Foreign key to User
}

type FlashcardSet struct {
	ID          string     `json:"id"`     // Primary key
	Author      User       `json:"author"` // Foreign key to User
	Status      string     `json:"status"`
	Description string     `json:"description"`
	CreatedAt   time.Time  `json:"created_at"`
	ModifiedAt  *time.Time `json:"modified_at"`
	AuthorID    uint64     `json:"author_id"`
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

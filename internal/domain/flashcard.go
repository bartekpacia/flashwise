package domain

import (
	"context"
	"time"
)

type Flashcard struct {
	ID         uint64     `json:"id"` // Primary key
	Front      string     `json:"front"`
	Back       string     `json:"back"`
	CreatedAt  time.Time  `json:"created_at" db:"created_at"`
	ModifiedAt *time.Time `json:"modified_at" db:"modified_at"`
	SetID      uint64     `json:"flashcard_set" db:"set_id"` // Foreign key to FlashcardSet
	AuthorID   uint64     `json:"author_id" db:"author_id"`  // Foreign key to User
}

// All contexts must contain userID.

type FlashcardRepository interface {
	GetAll(ctx context.Context) ([]Flashcard, error)
	GetByID(ctx context.Context, id uint64) (*Flashcard, error)
	GetBySetID(ctx context.Context, setID uint64) ([]Flashcard, error)

	Create(ctx context.Context, front string, back string, setID uint64) (*uint64, error)
	Update(ctx context.Context, id uint64, front string, back string, setID uint64) error

	Delete(ctx context.Context, id uint64) error
	// DeleteBySetID(ctx context.Context, setID uint64) error
}

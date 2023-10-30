package domain

import (
	"context"
	"time"
)

type FlashcardSet struct {
	ID         uint64     `json:"id" db:"id"` // Primary key
	Public     string     `json:"is_public" db:"is_public"`
	Title      string     `json:"name" db:"title"`
	CreatedAt  time.Time  `json:"created_at" db:"created_at"`
	ModifiedAt *time.Time `json:"modified_at" db:"modified_at"`
	AuthorID   uint64     `json:"author_id" db:"author_id"`  // Foreign key to User
	CategoryID uint64     `json:"category" db:"category_id"` // Foreign key to Category
}

type FlashcardSetRepository interface {
	GetAll(ctx context.Context) ([]FlashcardSet, error)
	GetById(ctx context.Context) (*FlashcardSet, error)

	Create(ctx context.Context, title string, public bool, categoryID uint64) (*uint64, error)

	Delete(ctx context.Context, id uint64) error
}

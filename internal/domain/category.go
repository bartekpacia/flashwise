package domain

import "context"

type Category struct {
	ID    uint64 `json:"id" db:"id" ` // Primary key
	Title string `json:"name" db:"title" `
	Slug  string `json:"slug" db:"slug" `
}

type CategoryRepository interface {
	GetAll(ctx context.Context) ([]Category, error)
}

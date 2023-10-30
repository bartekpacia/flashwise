package database

import (
	"context"

	"github.com/bartekpacia/flashwise/internal/domain"
)

type categoryRepo struct {
	db Database
}

func NewCategoryRepository(db Database) domain.CategoryRepository {
	return &categoryRepo{db: db}
}

func (r *categoryRepo) GetAll(ctx context.Context) ([]domain.Category, error) {
	categories := make([]domain.Category, 0)
	err := r.db.Select(&categories, "SELECT * FROM categories")

	return categories, err
}

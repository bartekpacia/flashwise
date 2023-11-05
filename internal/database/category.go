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
	err := r.db.SelectContext(ctx, &categories, "SELECT * FROM categories")
	if err != nil {
		return nil, err
	}

	return categories, err
}

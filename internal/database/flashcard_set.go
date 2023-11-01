package database

import (
	"context"

	"github.com/bartekpacia/flashwise/internal/domain"
)

type flashcardSetRepo struct {
	db Database
}

func NewFlashcardSetRepository(db Database) domain.FlashcardSetRepository {
	return &flashcardSetRepo{db: db}
}

func (r *flashcardSetRepo) GetAll(ctx context.Context, includePrivate bool) ([]domain.FlashcardSet, error) {
	userID, ok := ctx.Value("user_id").(uint64)
	if !ok {
		return nil, domain.ErrNoUserID
	}

	flashcardSets := make([]domain.FlashcardSet, 0)
	stmt := "SELECT * FROM flashcard_sets WHERE author_id = ? AND (is_public = true OR is_public = ?)"
	err := r.db.SelectContext(ctx, &flashcardSets, stmt, userID, !includePrivate)
	if err != nil {
		return nil, err
	}

	return flashcardSets, nil
}

func (r *flashcardSetRepo) GetByID(ctx context.Context, id uint64) (*domain.FlashcardSet, error) {
	userID, ok := ctx.Value("user_id").(uint64)
	if !ok {
		return nil, domain.ErrNoUserID
	}

	var flashcardSet domain.FlashcardSet
	err := r.db.GetContext(ctx, &flashcardSet, "SELECT * FROM flashcard_sets WHERE id = ? AND author_id = ?", id, userID)
	if err != nil {
		return nil, err
	}

	return &flashcardSet, nil
}

func (r *flashcardSetRepo) GetByCategory(ctx context.Context, categoryID uint64) ([]domain.FlashcardSet, error) {
	flashcardSets := make([]domain.FlashcardSet, 0)
	stmt := "SELECT * FROM flashcard_sets WHERE is_public = true AND category_id = ?"
	err := r.db.SelectContext(ctx, &flashcardSets, stmt, categoryID)
	if err != nil {
		return nil, err
	}

	return flashcardSets, nil
}

func (r *flashcardSetRepo) GetByNameContains(ctx context.Context, name string) ([]domain.FlashcardSet, error) {
	flashcardSets := make([]domain.FlashcardSet, 0)
	stmt := "SELECT * FROM flashcard_sets WHERE is_public = true AND title LIKE ?"
	err := r.db.SelectContext(ctx, &flashcardSets, stmt, "%"+name+"%")
	if err != nil {
		return nil, err
	}

	return flashcardSets, nil
}

func (r *flashcardSetRepo) Create(ctx context.Context, title string, public bool, categoryID uint64) (uint64, error) {
	userID, ok := ctx.Value("user_id").(uint64)
	if !ok {
		return 0, domain.ErrNoUserID
	}

	stmt := "INSERT INTO flashcard_sets (author_id, title, is_public, category_id) VALUES (?, ?, ?, ?)"
	result, err := r.db.ExecContext(ctx, stmt, userID, title, public, categoryID)
	if err != nil {
		return 0, err
	}

	id, _ := result.LastInsertId()
	return uint64(id), nil
}

func (r *flashcardSetRepo) Update(ctx context.Context, id uint64, title string, public bool, categoryID uint64) error {
	userID, ok := ctx.Value("user_id").(uint64)
	if !ok {
		return domain.ErrNoUserID
	}

	stmt := "UPDATE flashcard_sets SET title = ?, is_public = ? , category_id = ? WHERE id = ? AND author_id = ?"
	_, err := r.db.ExecContext(ctx, stmt, title, public, categoryID, id, userID)
	if err != nil {
		return err
	}

	return nil
}

func (r *flashcardSetRepo) Delete(ctx context.Context, id uint64) error {
	userID, ok := ctx.Value("user_id").(uint64)
	if !ok {
		return domain.ErrNoUserID
	}

	_, err := r.db.ExecContext(ctx, "DELETE FROM flashcard_sets WHERE id = ? AND author_id = ?", id, userID)
	if err != nil {
		return err
	}

	return nil
}

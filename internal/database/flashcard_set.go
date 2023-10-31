package database

import (
	"context"

	"github.com/bartekpacia/flashwise/internal/domain"
)

type flashcardSetRepository struct {
	db Database
}

func NewFlashcardSetRepository(db Database) domain.FlashcardSetRepository {
	return &flashcardSetRepository{db: db}
}

func (r *flashcardSetRepository) GetAll(ctx context.Context) ([]domain.FlashcardSet, error) {
	userID, ok := ctx.Value("user_id").(uint64)
	if !ok {
		return nil, domain.ErrNoUserID
	}

	flashcardSets := make([]domain.FlashcardSet, 0)
	err := r.db.SelectContext(ctx, &flashcardSets, "SELECT * FROM flashcard_sets WHERE author_id = ?", userID)
	if err != nil {
		return nil, err
	}

	return flashcardSets, nil
}

func (r *flashcardSetRepository) GetByID(ctx context.Context, id uint64) (*domain.FlashcardSet, error) {
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

func (r *flashcardSetRepository) Create(ctx context.Context, title string, public bool, categoryID uint64) (*uint64, error) {
	userID, ok := ctx.Value("user_id").(uint64)
	if !ok {
		return nil, domain.ErrNoUserID
	}

	stmt := "INSERT INTO flashcard_sets (author_id, title, is_public, category_id) VALUES (?, ?, ?, ?)"
	result, err := r.db.ExecContext(ctx, stmt, userID, title, public, categoryID)
	if err != nil {
		return nil, err
	}

	id, _ := result.LastInsertId()
	idUint := uint64(id)
	return &idUint, nil
}

func (r *flashcardSetRepository) Delete(ctx context.Context, id uint64) error {
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

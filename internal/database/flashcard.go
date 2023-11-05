package database

import (
	"context"
	"fmt"

	"github.com/bartekpacia/flashwise/internal/domain"
)

type flashcardRepo struct {
	db Database
}

func NewFlashcardRepository(db Database) domain.FlashcardRepository {
	return &flashcardRepo{db: db}
}

func (r *flashcardRepo) GetAll(ctx context.Context) ([]domain.Flashcard, error) {
	userID, ok := ctx.Value("user_id").(uint64)
	if !ok {
		return nil, domain.ErrNoUserID
	}

	flashcards := make([]domain.Flashcard, 0)
	err := r.db.GetContext(ctx, &flashcards, "SELECT * FROM flashcards WHERE author_id = ?", userID)
	if err != nil {
		// TODO: handle no rows error
		return nil, fmt.Errorf("failed to get all flashcards: %v", err)
	}

	return flashcards, nil
}

func (r *flashcardRepo) GetByID(ctx context.Context, id uint64) (*domain.Flashcard, error) {
	userID, ok := ctx.Value("user_id").(uint64)
	if !ok {
		return nil, domain.ErrNoUserID
	}

	var flashcard domain.Flashcard
	err := r.db.GetContext(ctx, &flashcard, "SELECT * FROM flashcards WHERE author_id = ? AND id = ?", userID, id)
	if err != nil {
		// TODO: handle no rows error
		return nil, fmt.Errorf("failed to get flashcard by id: %v", err)
	}

	return &flashcard, nil
}

func (r *flashcardRepo) GetBySetID(ctx context.Context, setID uint64) ([]domain.Flashcard, error) {
	userID, ok := ctx.Value("user_id").(uint64)
	if !ok {
		return nil, domain.ErrNoUserID
	}

	// Check if set exists
	var exists bool
	row := r.db.QueryRowxContext(ctx, "SELECT EXISTS(SELECT 1 FROM flashcard_sets WHERE id = ?)", setID)
	err := row.Scan(&exists)
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, fmt.Errorf("set with id %d not found", setID)
	}

	flashcards := make([]domain.Flashcard, 0)
	err = r.db.SelectContext(ctx, &flashcards, "SELECT * FROM flashcards WHERE author_id = ? AND set_id = ?", userID, setID)
	if err != nil {
		// TODO: handle no rows error
		return nil, fmt.Errorf("failed to get flashcards by set id: %v", err)
	}

	return flashcards, nil
}

func (r *flashcardRepo) Create(ctx context.Context, front string, back string, setID uint64) (uint64, error) {
	userID, ok := ctx.Value("user_id").(uint64)
	if !ok {
		return 0, domain.ErrNoUserID
	}

	// Verify that the flashcard set belongs to the user with userID
	var set domain.FlashcardSet
	err := r.db.GetContext(ctx, &set, "SELECT * FROM flashcard_sets WHERE id = ?", setID)
	if err != nil {
		return 0, err
	}
	if set.AuthorID != userID {
		return 0, fmt.Errorf("set with ID %d does not belong to user %d", setID, userID)
	}

	stmt := `
		INSERT INTO flashcards
			(front, back, author_id, set_id)
		VALUES
			(?, ?, ?, ?)`
	result, err := r.db.ExecContext(ctx, stmt, front, back, userID, setID)
	if err != nil {
		return 0, err
	}

	id, _ := result.LastInsertId()
	return uint64(id), nil
}

func (r *flashcardRepo) Update(ctx context.Context, id uint64, front string, back string, setID uint64) error {
	userID, ok := ctx.Value("user_id").(uint64)
	if !ok {
		return domain.ErrNoUserID
	}

	stmt := `
		UPDATE flashcards f
		SET
			f.front = ?,
			f.back = ?,
			f.set_id = ?,
			f.modified_at = NOW()
		WHERE
			f.id = ? AND
			f.author_id = ? AND
			EXISTS (
				SELECT 1
				FROM flashcard_sets s
				WHERE s.id = f.set_id AND s.author_id = f.author_id
			)`

	result, err := r.db.ExecContext(ctx, stmt, front, back, setID, id, userID)
	if err != nil {
		return fmt.Errorf("failed to execute query: %v", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return domain.ErrNotFound
	}

	return nil
}

func (r *flashcardRepo) Delete(ctx context.Context, id uint64) error {
	userID, ok := ctx.Value("user_id").(uint64)
	if !ok {
		return domain.ErrNoUserID
	}

	result, err := r.db.ExecContext(ctx, "DELETE FROM flashcards WHERE id = ? AND author_id = ?", id, userID)
	if err != nil {
		return fmt.Errorf("failed to execute query: %v", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return domain.ErrNotFound
	}

	return nil
}

func (r *flashcardRepo) DeleteBySetID(ctx context.Context, setID uint64) error {
	userID, ok := ctx.Value("user_id").(uint64)
	if !ok {
		return domain.ErrNoUserID
	}

	result, err := r.db.ExecContext(ctx, "DELETE FROM flashcards WHERE set_id = ? AND author_id = ?", setID, userID)
	if err != nil {
		return fmt.Errorf("failed to execute query: %v", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return domain.ErrNotFound
	}

	return nil
}

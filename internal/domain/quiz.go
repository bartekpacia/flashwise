package domain

type Quiz struct {
	ID             uint64 `db:"id"`            // Primary key
	FlashcardSetID uint64 `db:"flashcard_set"` // Foreign key to FlashcardSet
	AuthorID       uint64 `db:"author_id"`     // Foreign key to User
}

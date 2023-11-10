package database

import (
	"context"
	"database/sql"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/bartekpacia/flashwise/internal/domain"
)

var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

type quizRepo struct {
	db Database
}

func NewQuizRepository(db Database) domain.QuizRepository {
	return &quizRepo{db: db}
}

func (r *quizRepo) Generate(ctx context.Context, id uint64) (*domain.Quiz, error) {
	userID, ok := ctx.Value("user_id").(uint64)
	if !ok {
		return nil, domain.ErrNoUserID
	}

	var flashcardSet domain.FlashcardSet
	err := r.db.GetContext(ctx, &flashcardSet, "SELECT * FROM flashcard_sets WHERE id = ?", id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrNotFound
		}

		return nil, err
	}

	if flashcardSet.AuthorID != userID && !flashcardSet.Public {
		return nil, domain.ErrNoAccess
	}

	flashcards := make([]domain.Flashcard, 0)
	err = r.db.SelectContext(ctx, &flashcards, "SELECT * FROM flashcards WHERE set_id = ?", id)
	if err != nil {
		return nil, err
	}

	rng.Shuffle(len(flashcards), func(i, j int) {
		flashcards[i], flashcards[j] = flashcards[j], flashcards[i]
	})

	quiz := domain.Quiz{
		ID:        uint64(rng.Intn(10000)),
		Questions: generateQuestions(flashcards),
	}

	return &quiz, nil
}

func generateQuestions(flashcards []domain.Flashcard) []domain.Question {
	// Step 1. Create Question with single good answer from Flashcards

	questions := make([]domain.Question, 0)

	for _, flashcard := range flashcards {
		answers := make([]domain.Answer, 0)

		// Add the correct answer
		answers = append(answers, domain.Answer{
			Letter: "A",
			Text:   flashcard.Back,
		})

		// Add 3 other incorrect answers
		answers = append(answers, genRandomAnswers(flashcards)...)

		question := domain.Question{
			ID:          uint64(rng.Intn(10000)),
			FlashcardID: flashcard.ID,
			Text:        flashcard.Front,
			Answers:     answers,
		}

		questions = append(questions, question)
	}

	return questions
}

func genRandomAnswers(flashcards []domain.Flashcard) []domain.Answer {
	backs := make([]string, 0)
	for i := 0; i < 3; i++ {
		flashcard := flashcards[i]
		backs = append(backs, flashcard.Back)
	}

	rng.Shuffle(len(backs), func(i, j int) {
		backs[i], backs[j] = backs[j], backs[i]
	})

	answers := make([]domain.Answer, 0)
	for i, back := range backs {
		answers = append(answers, domain.Answer{
			Letter: strconv.Itoa(66 + i),
			Text:   back,
		})
	}

	return answers
}

func (r *quizRepo) Check(ctx context.Context, answers map[string]string) (*domain.QuizResult, error) {
	// TODO: Use userID to check for ownership of flashcards
	_, ok := ctx.Value("user_id").(uint64)
	if !ok {
		return nil, domain.ErrNoUserID
	}

	quizResult := domain.QuizResult{
		Results: make(map[string]string),
	}

	for flashcardID, gotAnswer := range answers {
		var goodAnswer string
		err := r.db.GetContext(ctx, &goodAnswer, "SELECT back FROM flashcards WHERE id = ?", flashcardID)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, domain.ErrNotFound
			}

			return nil, err
		}

		fmt.Println("got good answer for flashcard", flashcardID, ":", goodAnswer)

		if gotAnswer != goodAnswer {
			quizResult.Results[flashcardID] = "Incorrect"
		} else {
			quizResult.Results[flashcardID] = "Correct"
			quizResult.FinalScore++
		}

	}

	return &quizResult, nil
}

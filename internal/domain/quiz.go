package domain

import "context"

type Quiz struct {
	Questions []Question `json:"question"`
}

type Question struct {
	FlashcardID uint64   `json:"flashcard_id"`
	Text        string   `json:"text"`
	Answers     []Answer `json:"answers"`
}

type Answer struct {
	Letter string `json:"letter"`
	Text   string `json:"text"`
}

type QuizResult struct {
	Results    map[string]string `json:"results"`
	FinalScore uint              `json:"final_score"`
}

type QuizRepository interface {
	Generate(ctx context.Context, id uint64) (*Quiz, error)
	Check(ctx context.Context, answers map[string]string) (*QuizResult, error)
}

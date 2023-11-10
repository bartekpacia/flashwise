package domain

import "context"

type Quiz struct {
	ID        uint64     `json:"quiz_id"`
	Questions []Question `json:"question"`
}

type Question struct {
	ID          uint64   `json:"id"`
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

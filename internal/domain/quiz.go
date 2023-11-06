package domain

import "context"

type Quiz struct {
	ID        uint64     `json:"quiz_id"`
	Questions []Question `json:"question"`
}

type Question struct {
	ID      uint64   `json:"id"`
	Text    string   `json:"text"`
	Answers []Answer `json:"answers"`
}

type Answer struct {
	Letter string `json:"letter"`
	Text   string `json:"text"`
}

type QuizRepository interface {
	Generate(ctx context.Context, id uint64) (*Quiz, error)
	// Check(ctx context.Context, id uint64) (*Flashcard, error)
}

package api

import (
	"log/slog"
	"net/http"

	"github.com/bartekpacia/flashwise/internal/domain"
)

type api struct {
	logger     *slog.Logger
	httpClient *http.Client

	userRepo         domain.UserRepository
	flashcardRepo    domain.FlashcardRepository
	flashcardSetRepo domain.FlashcardSetRepository
	categoryRepo     domain.CategoryRepository
}

func NewAPI(logger *slog.Logger, httpClient *http.Client) *api {
	return &api{
		logger:     logger,
		httpClient: httpClient,
	}
}

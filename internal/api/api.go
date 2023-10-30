package api

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/bartekpacia/flashwise/internal/database"
	"github.com/bartekpacia/flashwise/internal/domain"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

type api struct {
	logger     *slog.Logger
	httpClient *http.Client

	userRepo         domain.UserRepository
	flashcardRepo    domain.FlashcardRepository
	flashcardSetRepo domain.FlashcardSetRepository
	categoryRepo     domain.CategoryRepository
}

func NewAPI(logger *slog.Logger, db *sqlx.DB) *api {
	httpClient := &http.Client{}

	flashcardRepo := database.NewFlashcardRepository(db)
	categoryRepo := database.NewCategoryRepository(db)

	return &api{
		logger:     logger,
		httpClient: httpClient,

		flashcardRepo: flashcardRepo,
		categoryRepo:  categoryRepo,
	}
}

func (a *api) CreateServer(port int) *http.Server {
	return &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: a.routes(),
	}
}

func (a *api) routes() http.Handler {
	router := mux.NewRouter()

	router.HandleFunc("/api/register", CreateUser).Methods("POST")

	router.HandleFunc("/api/flashcards", AuthHandler(GetFlashcards)).Methods("GET")
	router.HandleFunc("/api/flashcards", AuthHandler(CreateFlashcard)).Methods("POST")
	router.HandleFunc("/api/flashcards/{id}", AuthHandler(a.updateFlashcard)).Methods("PATCH", "PUT")
	router.HandleFunc("/api/flashcards", AuthHandler(a.deleteFlashcard)).Methods("DELETE")

	router.HandleFunc("/api/sets", AuthHandler(GetFlashcardSet)).Methods("GET")
	router.HandleFunc("/api/sets", AuthHandler(CreateFlashcardSet)).Methods("POST")
	router.HandleFunc("/api/sets/{id}", AuthHandler(DeleteFlashcardSet)).Methods("DELETE")

	router.HandleFunc("/api/category", AuthHandler(GetCategories)).Methods("GET")

	router.HandleFunc("/api/quiz/generate", AuthHandler(GenerateQuiz)).Methods("POST")
	router.HandleFunc("/api/quiz/check", AuthHandler(CheckQuiz)).Methods("PUT")

	return TrailingSlashHandler(LogHandler(CORSHandler(router.ServeHTTP)))
}

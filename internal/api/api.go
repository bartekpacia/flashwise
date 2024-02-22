package api

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/bartekpacia/flashwise/internal/api/middleware"
	"github.com/bartekpacia/flashwise/internal/database"
	"github.com/bartekpacia/flashwise/internal/domain"
	"github.com/jmoiron/sqlx"
)

type api struct {
	logger     *slog.Logger
	httpClient *http.Client

	userRepo         domain.UserRepository
	flashcardRepo    domain.FlashcardRepository
	flashcardSetRepo domain.FlashcardSetRepository
	categoryRepo     domain.CategoryRepository
	quizRepo         domain.QuizRepository
}

func NewAPI(logger *slog.Logger, db *sqlx.DB) *api {
	httpClient := &http.Client{}

	middleware.DB = db // This is hacky. See #5

	userRepo := database.NewUserRepository(db)
	flashcardRepo := database.NewFlashcardRepository(db)
	flashcardSetRepo := database.NewFlashcardSetRepository(db)
	categoryRepo := database.NewCategoryRepository(db)
	quizRepo := database.NewQuizRepository(db)

	return &api{
		logger:     logger,
		httpClient: httpClient,

		userRepo:         userRepo,
		flashcardRepo:    flashcardRepo,
		flashcardSetRepo: flashcardSetRepo,
		categoryRepo:     categoryRepo,
		quizRepo:         quizRepo,
	}
}

func (a *api) CreateServer(port int) *http.Server {
	return &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: a.routes(),
	}
}

func (a *api) routes() http.Handler {
	router := http.NewServeMux()

	router.HandleFunc("POST /auth/login", a.login)
	router.HandleFunc("POST /api/register", a.createUser)

	router.HandleFunc("GET /api/flashcards", middleware.AuthHandler(a.getFlashcards))
	router.HandleFunc("POST /api/flashcards", middleware.AuthHandler(a.createFlashcard))
	router.HandleFunc("PATCH /api/flashcards/{id}", middleware.AuthHandler(a.updateFlashcard))
	router.HandleFunc("PUT /api/flashcards/{id}", middleware.AuthHandler(a.updateFlashcard))
	router.HandleFunc("DELETE /api/flashcards/{id}", middleware.AuthHandler(a.deleteFlashcard))

	router.HandleFunc("GET /api/sets", middleware.AuthHandler(a.getFlashcardSet))
	router.HandleFunc("POST /api/sets", middleware.AuthHandler(a.createFlashcardSet))
	router.HandleFunc("PATCH /api/sets/{id}", middleware.AuthHandler(a.updateFlashcardSet))
	router.HandleFunc("PUT /api/sets/{id}", middleware.AuthHandler(a.updateFlashcardSet))
	router.HandleFunc("DELETE /api/sets/{id}", middleware.AuthHandler(a.deleteFlashcardSet))

	router.HandleFunc("GET /api/category", middleware.AuthHandler(a.getCategories))

	router.HandleFunc("POST /api/quiz/generate", middleware.AuthHandler(a.generateQuiz))
	router.HandleFunc("PUT /api/quiz/check", middleware.AuthHandler(a.checkQuiz))

	return middleware.TrailingSlashHandler(middleware.LogHandler(middleware.CORSHandler(router)))
}

package api

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/bartekpacia/flashwise/internal/api/middleware"
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

	middleware.DB = db // This is hacky. See #5

	userRepo := database.NewUserRepository(db)
	flashcardRepo := database.NewFlashcardRepository(db)
	flashcardSetRepo := database.NewFlashcardSetRepository(db)
	categoryRepo := database.NewCategoryRepository(db)

	return &api{
		logger:     logger,
		httpClient: httpClient,

		userRepo:         userRepo,
		flashcardRepo:    flashcardRepo,
		flashcardSetRepo: flashcardSetRepo,
		categoryRepo:     categoryRepo,
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

	router.HandleFunc("/auth/login", a.login).Methods("POST")
	router.HandleFunc("/api/register", a.createUser).Methods("POST")

	router.HandleFunc("/api/flashcards", middleware.AuthHandler(a.getFlashcards)).Methods("GET")
	router.HandleFunc("/api/flashcards", middleware.AuthHandler(a.createFlashcard)).Methods("POST")
	router.HandleFunc("/api/flashcards/{id}", middleware.AuthHandler(a.updateFlashcard)).Methods("PATCH", "PUT")
	router.HandleFunc("/api/flashcards/{id}", middleware.AuthHandler(a.deleteFlashcard)).Methods("DELETE")

	router.HandleFunc("/api/sets", middleware.AuthHandler(a.getFlashcardSet)).Methods("GET")
	router.HandleFunc("/api/sets", middleware.AuthHandler(a.createFlashcardSet)).Methods("POST")
	router.HandleFunc("/api/sets/{id}", middleware.AuthHandler(a.updateFlashcardSet)).Methods("PATCH", "PUT")
	router.HandleFunc("/api/sets/{id}", middleware.AuthHandler(a.deleteFlashcardSet)).Methods("DELETE")

	router.HandleFunc("/api/category", middleware.AuthHandler(a.getCategories)).Methods("GET")

	router.HandleFunc("/api/quiz/generate", middleware.AuthHandler(a.generateQuiz)).Methods("POST")
	router.HandleFunc("/api/quiz/check", middleware.AuthHandler(a.checkQuiz)).Methods("PUT")

	return middleware.TrailingSlashHandler(middleware.LogHandler(middleware.CORSHandler(router)))
}

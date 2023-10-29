package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/felixge/httpsnoop"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/lmittmann/tint"
)

var db *sqlx.DB

func main() {
	var err error

	setUpLogger()
	router := setUpRouter()
	db, err = connectDB()
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
	}

	http.ListenAndServe(":8080", router)
}

func setUpLogger() {
	opts := &tint.Options{TimeFormat: time.TimeOnly}
	handler := tint.NewHandler(os.Stdout, opts)
	logger := slog.New(handler)
	slog.SetDefault(logger)
}

func setUpRouter() http.Handler {
	router := mux.NewRouter()

	router.HandleFunc("/api/register", CreateUser).Methods("POST")

	router.HandleFunc("/api/flashcards", AuthHandler(GetFlashcards)).Methods("GET")
	router.HandleFunc("/api/flashcards", AuthHandler(CreateFlashcard)).Methods("POST")
	router.HandleFunc("/api/flashcards/{id}", AuthHandler(UpdateFlashcard)).Methods("PATCH", "PUT")
	router.HandleFunc("/api/flashcards", AuthHandler(DeleteFlashcard)).Methods("DELETE")

	router.HandleFunc("/api/sets", AuthHandler(GetFlashcardSet)).Methods("GET")
	router.HandleFunc("/api/sets", AuthHandler(CreateFlashcardSet)).Methods("POST")

	router.HandleFunc("/api/category", AuthHandler(GetCategories)).Methods("GET")

	router.HandleFunc("/api/quiz/generate", AuthHandler(GenerateQuiz)).Methods("POST")
	router.HandleFunc("/api/quiz/check", AuthHandler(CheckQuiz)).Methods("PUT")

	return TrailingSlashHandler(LogHandler(CORSHandler(router.ServeHTTP)))
}

func connectDB() (*sqlx.DB, error) {
	host, ok := os.LookupEnv("MYSQL_HOST")
	if !ok {
		return nil, errors.New("MYSQL_HOST env var not set")
	}

	user, ok := os.LookupEnv("MYSQL_USER")
	if !ok {
		return nil, errors.New("MYSQL_USER env var not set")
	}

	password, ok := os.LookupEnv("MYSQL_PASSWORD")
	if !ok {
		return nil, errors.New("MYSQL_PASSWORD env var not set")
	}

	dbName, ok := os.LookupEnv("MYSQL_DB")
	if !ok {
		return nil, errors.New("MYSQL_DB env var not set")
	}

	connString := fmt.Sprintf("%s:%s@(%s:3306)/%s?parseTime=true", user, password, host, dbName)

	var database *sqlx.DB
	var err error
	fails := 0
	maxFails := 60
	for {
		if fails >= maxFails {
			return nil, fmt.Errorf("failed to connect to database after %d fails", maxFails)
		}

		if fails > 0 {
			time.Sleep(1 * time.Second)
		}

		database, err = sqlx.Open("mysql", connString)
		if err != nil {
			slog.Warn("failed to open connection to database", "error", err)
			fails++
			continue
		}

		err = database.Ping()
		if err != nil {
			slog.Warn("failed to ping database", "error", err)
			fails++
			continue
		}

		break
	}

	slog.Info("connected to database", "host", host, "user", user, "db_name", dbName)

	return database, nil
}

type ContextKey string

const ContextUserKey ContextKey = "user_id"

func AuthHandler(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			http.Error(w, "No Authorization header provided", http.StatusUnauthorized)
			return
		}

		splitToken := strings.Split(token, "Token")
		if len(splitToken) != 2 {
			http.Error(w, "Bearer token not in proper format", http.StatusUnauthorized)
			return
		}

		token = strings.TrimSpace(splitToken[1])

		var user User
		err := db.Get(&user, "SELECT * FROM users WHERE token = ?", token)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), ContextUserKey, user.ID)
		next(w, r.WithContext(ctx))
	})
}

var (
	allowedOrigins = []string{"http://localhost:3000"}
	allowedMethods = []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"}
)

func CORSHandler(next http.HandlerFunc) http.HandlerFunc {
	isPreflight := func(r *http.Request) bool {
		return r.Method == "OPTIONS" &&
			r.Header.Get("Origin") != "" &&
			r.Header.Get("Access-Control-Request-Method") != ""
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		if isPreflight(r) {
			method := r.Header.Get("Access-Control-Request-Method")
			if slices.Contains(allowedOrigins, origin) && slices.Contains(allowedMethods, method) {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Methods", strings.Join(allowedMethods, ", "))
				w.Header().Set("Access-Control-Allow-Headers", "*, Authorization")
				w.Header().Set("Access-Control-Max-Age", "86400")
				w.Header().Set("Vary", "Origin")
			}

			return
		}

		if slices.Contains(allowedOrigins, origin) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Vary", "Origin")
		}

		next(w, r)
	})
}

func LogHandler(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Info(r.Method, "from", r.RemoteAddr, "url", r.URL.Path)
		m := httpsnoop.CaptureMetrics(next, w, r)
		slog.Info(r.Method, "from", r.RemoteAddr, "url", r.URL.Path, "status_code", m.Code, "duration", m.Duration.Milliseconds())
	})
}

func TrailingSlashHandler(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = strings.TrimSuffix(r.URL.Path, "/")
		next(w, r)
	})
}

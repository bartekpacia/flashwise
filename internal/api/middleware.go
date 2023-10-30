package api

import (
	"context"
	"log/slog"
	"net/http"
	"slices"
	"strings"

	"github.com/felixge/httpsnoop"
)

// type ContextKey string
// const ContextUserKey ContextKey = "user_id"

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

		ctx := context.WithValue(r.Context(), "user_id", user.ID)
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

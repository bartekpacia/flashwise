package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/bartekpacia/flashwise/internal/api"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/lmittmann/tint"
)

func main() {
	ctx := context.Background()

	logger := setUpLogger()

	db, err := connectDB(logger)
	if err != nil {
		logger.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}

	defer db.Close()

	server := api.NewAPI(logger, db).CreateServer(8080)
	go func() {
		server.ListenAndServe()
	}()

	logger.Info("server started", "port", 8080)

	<-ctx.Done()

	err = server.Shutdown(ctx)
	if err != nil {
		logger.Error("failed to shutdown server", "error", err)
	}
}

func setUpLogger() *slog.Logger {
	opts := &tint.Options{TimeFormat: time.TimeOnly}
	handler := tint.NewHandler(os.Stdout, opts)
	return slog.New(handler)
}

func connectDB(logger *slog.Logger) (*sqlx.DB, error) {
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
			logger.Warn("failed to open connection to database", "error", err)
			fails++
			continue
		}

		err = database.Ping()
		if err != nil {
			logger.Warn("failed to ping database", "error", err)
			fails++
			continue
		}

		break
	}

	logger.Info("connected to database", "host", host, "user", user, "db_name", dbName)

	return database, nil
}

package database

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type Database interface {
	sqlx.QueryerContext
	sqlx.ExecerContext

	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
}

package domain

import "errors"

var (
	ErrNotFound = errors.New("requested item was not found")
	ErrConflict = errors.New("item already exists")
	ErrNoUserID = errors.New("no user id found in context")
)

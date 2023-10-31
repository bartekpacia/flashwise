package domain

import "errors"

var (
	ErrNoUserID = errors.New("no user id found in context")
	ErrNotFound = errors.New("requested item was not found")
	ErrNoAccess = errors.New("no permission to access the requested item")
	ErrConflict = errors.New("item already exists")
)

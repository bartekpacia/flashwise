package domain

import "time"

type User struct {
	ID           uint64    `db:"id"` // Primary key
	Username     string    `db:"username"`
	Email        string    `db:"email"`
	CreatedAt    time.Time `db:"created_at"`
	Admin        bool      `db:"is_admin"`
	PasswordHash string    `db:"password_hash"`
	Token        string    `db:"token"`
}

type UserRepository interface {
	CreateUser(user *User) error
}

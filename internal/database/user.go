package database

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"

	"github.com/bartekpacia/flashwise/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

type userRepo struct {
	db Database
}

func NewUserRepository(db Database) domain.UserRepository {
	return &userRepo{db: db}
}

func (r *userRepo) Login(ctx context.Context, email string, password string) (token *string, err error) {
	return nil, errors.New("not implemented")
}

func (r *userRepo) Register(ctx context.Context, username string, email string, password string) (*domain.User, error) {
	passwordHash, _ := hashPassword(password)
	token := generateToken()
	stmt := "INSERT INTO users (username, email, password_hash, token) VALUES (?, ?, ?, ?)"
	result, err := r.db.ExecContext(ctx, stmt, username, email, passwordHash, token)
	if err != nil {
		return nil, err
	}

	id, _ := result.LastInsertId()

	var user domain.User
	err = r.db.GetContext(ctx, &user, "SELECT * FROM users WHERE id = ?", id)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func generateToken() string {
	randomBytes := make([]byte, 20)
	rand.Read(randomBytes)

	return hex.EncodeToString(randomBytes)
}

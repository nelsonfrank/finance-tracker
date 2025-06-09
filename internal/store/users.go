package store

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/nelsonfrank/finance-tracker/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrDuplicateEmail    = errors.New("a user with that email already exists")
	ErrDuplicateUsername = errors.New("a user with that username already exists")
)


type Password struct {
	text *string
	hash []byte
}

func (p *Password) Set(text string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(text), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	p.text = &text
	p.hash = hash

	return nil
}

func (p *Password) Compare(text string) error {
	return bcrypt.CompareHashAndPassword(p.hash, []byte(text))
}

type UsersStorage struct {
	db *sqlx.DB
	repo *repository.Queries
}

func (s *UsersStorage) Create(ctx context.Context, db *sqlx.DB, user *repository.CreateUserParams) error {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()


	repo := repository.New(db.DB)

	if _, err := repo.CreateUser(ctx, *user) ; err != nil {
		// Handle duplicate constraint errors
		if strings.Contains(err.Error(), `duplicate key value violates unique constraint "users_email_key"`) {
			return ErrDuplicateEmail
		}
		if strings.Contains(err.Error(), `duplicate key value violates unique constraint "users_username_key"`) {
			return ErrDuplicateUsername
		}
		return err
	}

	return nil
}

type UserInvitation struct {
	ID     int64     `gorm:"primaryKey"`
	Token  string    `gorm:"column:token"`
	UserID int64     `gorm:"column:user_id"`
	Expiry time.Time `gorm:"column:expiry"`
}





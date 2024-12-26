package store

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID        int64     `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `gorm:"uniqueIndex;not null" json:"email"`
	Password  string    `gorm:"not null" json:"-"`
	CreatedAt time.Time `gorm:"type:timestamp with time zone;not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `gorm:"type:timestamp with time zone;not null;default:CURRENT_TIMESTAMP"`
}

type UsersStorage struct {
	db *gorm.DB
}

func (s *UsersStorage) Create(ctx context.Context, user *User) error {

	return nil
}

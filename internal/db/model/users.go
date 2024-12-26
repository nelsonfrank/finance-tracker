package model

import (
	"time"

	"gorm.io/gorm"
)

// User model for database
type User struct {
	gorm.Model
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `gorm:"uniqueIndex;not null" json:"email"`
	Password  string    `gorm:"not null" json:"-"`
	CreatedAt time.Time `gorm:"type:timestamp with time zone;not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `gorm:"type:timestamp with time zone;not null;default:CURRENT_TIMESTAMP"`
}

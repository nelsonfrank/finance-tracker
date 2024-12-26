package store

import (
	"context"

	"gorm.io/gorm"
)

type Storage struct {
	Posts interface {
		Create(context.Context, *Post) error
	}
	Users interface {
		Create(context.Context, *User) error
	}
}

func NewStorage(db *gorm.DB) Storage {
	return Storage{
		Posts: &PostsStorage{db},
		Users: &UsersStorage{db},
	}
}

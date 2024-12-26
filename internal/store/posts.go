package store

import (
	"context"

	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	ID        int64    `json:"id"`
	Content   string   `json:"content"`
	Title     string   `json:"title"`
	UserId    int64    `json:"user_id"`
	Tags      []string `json:"tags"`
	CreatedAt string   `json:"created_at"`
	UpdatedAt string   `json:"updated_at"`
}
type PostsStorage struct {
	db *gorm.DB
}

func (s *PostsStorage) Create(ctx context.Context, post *Post) error {

	return nil
}

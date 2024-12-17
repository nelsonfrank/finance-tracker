package store

import (
	"context"
	"database/sql"
)

type PostsStorage struct {
	db *sql.DB
}

func (s *PostsStorage) Create(ctx context.Context) error {
	return nil
}

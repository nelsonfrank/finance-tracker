package db

import (
	"context"
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

func New(ctx context.Context, addr string) (*sqlx.DB, error) {
	db, err := sql.Open("pgx", addr) // using "pgx" driver via stdlib
	if err != nil {
		return nil, err
	}

	sqlxDB := sqlx.NewDb(db, "pgx")

	// Optional: Test the connection
	if err := sqlxDB.PingContext(ctx); err != nil {
		return nil, err
	}

	return sqlxDB, nil
}

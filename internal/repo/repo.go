package repo

import (
	"context"
	"database/sql"

	_ "github.com/lib/pq"
)

type Repo struct {
	DB *sql.DB
}

type Repository interface {
	AddFavorite(ctx context.Context, sid string, emojiID string) error
	RemoveFavorite(ctx context.Context, sid string, emojiID string) error
	GetFavorites(ctx context.Context, sid string) ([]string, error)
}

func New(dsn string) (*Repo, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &Repo{DB: db}, nil
}

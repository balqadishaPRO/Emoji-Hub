package repo

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type Repo struct {
	DB *sql.DB
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

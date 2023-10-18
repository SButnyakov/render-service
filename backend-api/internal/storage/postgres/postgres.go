package postgres

import (
	"database/sql"
	"fmt"
)

const (
	UniqueViolationErrorCode = "23505"
)

type PGStorage struct {
	Db *sql.DB
}

func New(storagePath string) (*PGStorage, error) {
	const fn = "storage.postgres.New"

	db, err := sql.Open("postgres", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	return &PGStorage{Db: db}, nil
}

package postgres

import (
	"backend-api/internal/config"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

const (
	UniqueViolationErrorCode = "23505"
)

type PGStorage struct {
	Db *sql.DB
}

func New(cfg *config.Config) (*PGStorage, error) {
	const fn = "storage.postgres.New"

	url := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable",
		cfg.DB.User,
		cfg.DB.Password,
		cfg.DB.Host,
		cfg.DB.Port,
		cfg.DB.Name)

	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	return &PGStorage{Db: db}, nil
}

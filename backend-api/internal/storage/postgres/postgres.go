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

	// TODO: implement migration mechanism
	stmt, err := db.Prepare(`
		CREATE TABLE IF NOT EXISTS users(
		    id int GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
		    login text UNIQUE NOT NULL,
		    email text UNIQUE NOT NULL,
		    password text NOT NULL
		)
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	stmt, err = db.Prepare(`
		CREATE TABLE IF NOT EXISTS orders(
			id int GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
			fileName text NOT NULL,
			user_id int NOT NULL,
			FOREIGN KEY (user_id) REFERENCES users (id)
		)
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	return &PGStorage{Db: db}, nil
}

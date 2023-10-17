package repos

import (
	"backend-api/internal/storage"
	"backend-api/internal/storage/postgres"
	"database/sql"
	"errors"
	"fmt"
	"github.com/lib/pq"
	"time"
)

type UserRepository struct {
	pg *postgres.PGStorage
}

func NewUserRepository(pg *postgres.PGStorage) *UserRepository {
	return &UserRepository{pg: pg}
}

func (u *UserRepository) Create(user storage.User) error {
	const fn = "postgres.repos.users.Create"

	stmt, err := u.pg.Db.Prepare("INSERT INTO users(login, email, password, sub_expire_date) values($1, $2, $3, $4)")
	if err != nil {
		return fmt.Errorf("%s: prepare statement: %w", fn, err)
	}

	_, err = stmt.Exec(user.Login, user.Email, user.Password, time.Now().Add(-time.Hour*24))
	if err != nil {
		if postgresErr, ok := err.(*pq.Error); ok && postgresErr.Code == postgres.UniqueViolationErrorCode {
			return fmt.Errorf("%s: %w", fn, storage.ErrUserExists)
		}

		return fmt.Errorf("%s: execute statement: %w", fn, err)
	}

	return nil
}

func (u *UserRepository) User(uid int) (storage.User, error) {
	const fn = "postgres.repos.users.User"

	var resUser storage.User

	stmt, err := u.pg.Db.Prepare("SELECT id, login, email FROM users WHERE id=$1")
	if err != nil {
		return resUser, fmt.Errorf("%s: prepare statement: %w", fn, err)
	}

	err = stmt.QueryRow(uid).Scan(&resUser.Id, &resUser.Login, &resUser.Email, &resUser.Password)
	if errors.Is(err, sql.ErrNoRows) {
		return resUser, storage.ErrUserNotFound
	}
	if err != nil {
		return resUser, fmt.Errorf("%s: execute statement: %w", fn, err)
	}

	return resUser, nil
}

func (u *UserRepository) CheckCredentials(loginOrEmail, password string) (int64, error) {
	const fn = "postgres.repos.users.CheckCredentials"

	stmt, err := u.pg.Db.Prepare("SELECT id FROM users WHERE (login=$1 OR email=$1) AND password=$2")
	if err != nil {
		return 0, fmt.Errorf("%s: prepare statement: %w", fn, err)
	}

	var uid int64

	err = stmt.QueryRow(loginOrEmail, password).Scan(&uid)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, storage.ErrInvalidCredentials
	}
	if err != nil {
		return 0, fmt.Errorf("%s: execute statement: %w", fn, err)
	}

	return uid, nil
}

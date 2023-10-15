package storage

import "errors"

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrUserExists         = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type User struct {
	Id       int64
	Login    string
	Email    string
	Password string
}

type Order struct {
	Id       int64
	FileName string
	UserId   int64
}

package storage

import (
	"errors"
	"time"
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrUserExists         = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type User struct {
	Id            int64
	Login         string
	Email         string
	Password      string
	SubExpireDate time.Time
}

type Order struct {
	Id       int64
	FileName string
	UserId   int64
}

type Payment struct {
	Id       int64
	DateTime time.Time
	TypeId   string
	UserID   string
}

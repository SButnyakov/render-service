package storage

import (
	"errors"
	"time"
)

var (
	ErrFailedToConnect               = errors.New("failed to connect to storage")
	ErrUserNotFound                  = errors.New("user not found")
	ErrUserExists                    = errors.New("user already exists")
	ErrInvalidCredentials            = errors.New("invalid credentials")
	ErrSubscriptionNotFound          = errors.New("subscription not found")
	ErrSubscriptionTypeDoesNotExists = errors.New("subscription type does not exists")
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

type Payment struct {
	Id       int64
	DateTime time.Time
	TypeId   int64
	UserID   int64
}

type PaymentType struct {
	Id   int64
	Name string
}

type Subscription struct {
	Id         int64
	UserId     int64
	TypeId     int64
	ExpireDate time.Time
}

type SubscriptionType struct {
	Id   int64
	Name string
}

// TODO: subscribtion info
type RedisData struct {
	Format     string
	Resolution string
	SavePath   string
}

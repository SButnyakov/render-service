package storage

import (
	"errors"
	"time"
)

var (
	ErrFailedToConnect = errors.New("failed to connect to storage")

	ErrUserNotFound         = errors.New("user not found")
	ErrSubscriptionNotFound = errors.New("subscription not found")

	ErrUserExists         = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")

	ErrNoSubscriptionTypes = errors.New("subscription types not found")
	ErrNoPaymentTypes      = errors.New("payment types not found")
	ErrNoOrderStatuses     = errors.New("order statuses not found")
)

type User struct {
	Id       int64
	Login    string
	Email    string
	Password string
}

type Order struct {
	Id           int64
	FileName     string
	StoringName  string
	CreationDate time.Time
	UserId       int64
	StatusId     int64
	DownloadLink int64
}

type OrderStatus struct {
	Id   int64
	Name string
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

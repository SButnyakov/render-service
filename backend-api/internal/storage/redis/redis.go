package redis

import (
	"backend-api/internal/storage"
	"context"
	"github.com/redis/go-redis/v9"
)

var (
	Nil = redis.Nil
)

func New(address string) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{Addr: address})
	_, err := client.Ping(context.Background()).Result()

	if err != nil {
		return nil, storage.ErrFailedToConnect
	}

	return client, nil
}

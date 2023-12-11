package redis

import (
	"backend-api/internal/config"
	"backend-api/internal/storage"
	"context"
	"github.com/redis/go-redis/v9"
)

func New(cfg *config.Config) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Address,
		Password: cfg.Redis.Password,
		DB:       0,
	})
	_, err := client.Ping(context.Background()).Result()

	if err != nil {
		return nil, storage.ErrFailedToConnect
	}

	return client, nil
}

package myredis

import (
	storage_config "LostAndFound/internal/config/storage_config"
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

func NewRedis(cfg storage_config.RedisConfig) (*redis.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	if err := redisClient.Ping(ctx).Err(); err != nil {
		return nil, err
	}
	return redisClient, nil
}

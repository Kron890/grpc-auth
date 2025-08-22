package storage

import (
	"auth-grpc/internal/config"
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
)

type DataBaseRedis struct {
	DB   *redis.Client
	logs *logrus.Logger
}

func NewRedis(cfg *config.Config, logs *logrus.Logger) (*DataBaseRedis, error) {
	redisAddr := fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort)

	client := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
		PoolSize: cfg.RedisPoolSize,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis connection failed: %w", err)
	}

	return &DataBaseRedis{DB: client, logs: logs}, nil
}

func (r *DataBaseRedis) Close() error {
	r.logs.Info("stopping Redis server")
	return r.DB.Close()
}

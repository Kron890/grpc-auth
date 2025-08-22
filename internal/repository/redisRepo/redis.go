package redisRepo

import (
	"auth-grpc/internal/domain"
	"auth-grpc/internal/infrastructure/storage"
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type RepositoryRedis struct {
	DB *storage.DataBaseRedis
}

func New(db *storage.DataBaseRedis) *RepositoryRedis {
	return &RepositoryRedis{DB: db}
}

// ШАБЛОН
func (r *RepositoryRedis) SaveToken(ctx context.Context, token domain.Token) error {
	ttl := 24 * time.Hour
	return r.DB.DB.Set(ctx, token.UserID, token.RefreshToken, ttl).Err()

}

// ШАБЛОН
func (r *RepositoryRedis) GetToken(ctx context.Context) (string, error) {
	userID, err := r.DB.DB.Get(ctx, "key").Result()
	if err == redis.Nil {
		//TODO: Создать ошибку
		return "", fmt.Errorf("key not found")
	} else if err != nil {
		return "", err
	}
	return userID, nil
}

// ШАБЛОН
func (r *RepositoryRedis) DeleteToken(ctx context.Context, login string) error {
	return r.DB.DB.Del(ctx, login).Err()

}

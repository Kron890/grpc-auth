package redisRepo

import (
	"auth-grpc/internal/infrastructure/storage"
	"context"
	"fmt"
	"time"
)

type RepositoryRedis struct {
	DB *storage.DataBaseRedis
}

func New(db *storage.DataBaseRedis) *RepositoryRedis {
	return &RepositoryRedis{DB: db}
}

// ШАБЛОН
func (r *RepositoryRedis) SaveRefreshToken(ctx context.Context, userID int64, token string, ttl time.Duration) error {
	key := fmt.Sprintf("refresh:%d", userID)
	return r.DB.Set(ctx, key, token, ttl).Err()
}

// ШАБЛОН
func (r *RepositoryRedis) GetRefreshToken(ctx context.Context, userID int64) (string, error) {
	key := fmt.Sprintf("refresh:%d", userID)
	return r.DB.Get(ctx, key).Result()
}

// ШАБЛОН
func (r *RepositoryRedis) DeleteRefreshToken(ctx context.Context, userID int64) error {
	key := fmt.Sprintf("refresh:%d", userID)
	return r.DB.Del(ctx, key).Err()
}

// DeleteToken удаляет токен по логину (для совместимости с интерфейсом)
func (r *RepositoryRedis) DeleteToken(ctx context.Context, login string) error {
	// Этот метод не используется
	return nil
}

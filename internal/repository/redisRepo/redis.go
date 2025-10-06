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

// New создает репозиторий поверх Redis-клиента.
func New(db *storage.DataBaseRedis) *RepositoryRedis {
	return &RepositoryRedis{DB: db}
}

// refreshKey формирует ключ хранения refresh-токена пользователя.
func refreshKey(userID int64) string {
	return fmt.Sprintf("refresh:%d", userID)
}

func (r *RepositoryRedis) SaveRefreshToken(ctx context.Context, userID int64, token string, ttl time.Duration) error {
	key := refreshKey(userID)
	return r.DB.Set(ctx, key, token, ttl).Err()
}

// GetRefreshToken возвращает сохранённый refresh-токен для пользователя.
func (r *RepositoryRedis) GetRefreshToken(ctx context.Context, userID int64) (string, error) {
	key := refreshKey(userID)
	return r.DB.Get(ctx, key).Result()
}

// DeleteRefreshToken удаляет refresh-токен пользователя
func (r *RepositoryRedis) DeleteRefreshToken(ctx context.Context, userID int64) error {
	key := refreshKey(userID)
	return r.DB.Del(ctx, key).Err()
}

// DeleteToken удаляет токен по логину (для совместимости с интерфейсом).
func (r *RepositoryRedis) DeleteToken(ctx context.Context, login string) error {
	// Этот метод не используется
	return nil
}

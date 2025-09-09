package internal

import (
	"auth-grpc/internal/domain/filters"
	"context"
	"time"
)

// TODO: передавать структуру
type User interface {
	GetUser(ctx context.Context, login string) (filters.UserDB, error)
	Create(ctx context.Context, user filters.UserDB) (int64, error)
	CheckUser(ctx context.Context, login string) (bool, error)
}

type RepoRedis interface {
	DeleteToken(ctx context.Context, login string) error
	GetRefreshToken(ctx context.Context, userID int64) (string, error)
	SaveRefreshToken(ctx context.Context, userID int64, token string, ttl time.Duration) error
}

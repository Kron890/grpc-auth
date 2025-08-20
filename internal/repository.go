package internal

import (
	"auth-grpc/internal/domain/filters"
	"context"
)

// TODO: передавать структуру
type User interface {
	GetUser(ctx context.Context, login string) (filters.UserDB, error)
	Create(ctx context.Context, user filters.UserDB) (int64, error)
}
type AppProvider interface {
	App(ctx context.Context, appID int)
}

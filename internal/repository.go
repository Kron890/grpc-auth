package internal

import (
	"auth-grpc/internal/domain"
	"context"
)

// TODO: передавать структуру
type User interface {
	Creat(ctx context.Context, login string, passHash []byte) (domain.User, error)
	Authenticate(ctx context.Context, login string, passHash []byte) (domain.User, error)
}
type AppProvider interface {
	App(ctx context.Context, appID int)
}

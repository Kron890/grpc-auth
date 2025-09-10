package internal

import (
	"auth-grpc/internal/domain"
	"context"
)

type Auth interface {
	Register(ctx context.Context, login, password string) (int64, error)
	Login(ctx context.Context, login, password string) (domain.Token, error)
	Verify(string) error
	Refresh(ctx context.Context, refreshToken string) (domain.Token, error)
	GenerateTokens(id int64, login string) (domain.Token, error)
}

package internal

import "context"

type Auth interface {
	Register(ctx context.Context, login, password string) (int64, error)
	Login(ctx context.Context, login, password string) (string, error)
}

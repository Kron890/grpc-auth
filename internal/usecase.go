package internal

import "context"

type Auth interface {
	Register(ctx context.Context, login, password string) (string, error)
	Login(ctx context.Context, login, password string) error
}

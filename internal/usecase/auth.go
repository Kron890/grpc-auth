package usecase

import (
	"auth-grpc/internal/domain"
	"context"
	"time"

	"github.com/sirupsen/logrus"
)

type Auth struct {
	logs        *logrus.Logger
	usr         User
	appProvider AppProvider
}

// TODO: передавать структуру
type User interface {
	SaveUser(ctx context.Context, login string, passHash []byte) (domain.User, error)
	UserAuth(ctx context.Context, login string, passHash []byte) (domain.User, error)
}
type AppProvider interface {
	App(ctx context.Context, appID int)
}

func New(logs logrus.Logger, user User, appProvider AppProvider, tokenTTL time.Duration) *Auth {
	return &Auth{
		logs:        &logs,
		usr:         user,
		appProvider: appProvider}
}

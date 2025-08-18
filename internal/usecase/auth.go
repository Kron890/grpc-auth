package usecase

import (
	"auth-grpc/internal"
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

type Auth struct {
	logs        *logrus.Logger
	user        internal.User
	appProvider internal.AppProvider
	tokenTTL    time.Duration
}

// TODO:...
func New(logs *logrus.Logger, user internal.User, appProvider internal.AppProvider, tokenTTL time.Duration) *Auth {
	return &Auth{
		logs:        logs,
		user:        user,
		appProvider: appProvider}
}

// Register TODO
func (a *Auth) Register(ctx context.Context, login, password string) (string, error) {
	//TODO: хеш + соль

	user, err := a.user.Creat(ctx, login, []byte(password))
	if err != nil {
		a.logs.Error("user creation failed:", err)
		return "", err
	}

	a.logs.Info("user has been successfully registered")
	return fmt.Sprintf("%d", user.ID), nil
}

// Login TODO
func (a *Auth) Login(ctx context.Context, login, password string) error {
	//TODO: хеш + соль

	user, err := a.user.Authenticate(ctx, login, []byte(password))
	if err != nil {
		return err
	}
	fmt.Println("user", user) //TODO:...
	return nil
}

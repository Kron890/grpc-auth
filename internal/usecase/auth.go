package usecase

import (
	"auth-grpc/internal"
	"auth-grpc/internal/lib/jwt"
	"auth-grpc/internal/lib/jwt/dto"
	"auth-grpc/internal/repository"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

var (
	errInvailidCredentials = errors.New("invalid credentials")
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
		appProvider: appProvider,
		tokenTTL:    tokenTTL,
	}
}

// Register TODO
func (a *Auth) Register(ctx context.Context, login, password string) (string, error) {
	//TODO: хеш + соль
	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		a.logs.Error("failed to generate password hash: ", err)
		return "", err
	}
	user, err := a.user.Creat(ctx, login, passHash)
	if errors.Is(err, repository.ErrUserExists) {
		a.logs.Error(repository.ErrUserExists)
		return "", repository.ErrUserExists
	}

	if err != nil {
		a.logs.Error("user creation failed:", err)
		return "", err
	}

	a.logs.Info("user registered")

	return fmt.Sprintf("%d", user.ID), nil
}

// Login TODO: возращать токен
func (a *Auth) Login(ctx context.Context, login, password string) error {
	//TODO: хеш + соль
	user, err := a.user.GetUser(ctx, login)
	if errors.Is(err, repository.ErrUserNotFound) {
		a.logs.Warn("user not found", err)
		return err
	} else if err != nil {
		a.logs.Error(err) // TODO:...
		return err
	}

	if err = bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		a.logs.Error("Invalid password", err) // TODO:...
		return errInvailidCredentials
	}
	userClaims := dto.MapFromUser(user)

	tokenJWT, err := jwt.NewToken(userClaims, a.tokenTTL)
	if err != nil {
		a.logs.Error(err)
		return err
	}
	fmt.Println(tokenJWT) // TODO:...

	a.logs.Info("user logged in")
	return nil
}

package usecase

import (
	"auth-grpc/internal"
	"auth-grpc/internal/domain"
	"auth-grpc/internal/domain/filters"
	"auth-grpc/internal/lib/jwt"
	"auth-grpc/internal/lib/jwt/dto"
	"context"
	"errors"
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
func (a *Auth) Register(ctx context.Context, login, password string) (int64, error) {
	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		a.logs.Error("failed to generate password hash: ", err)
		return 0, err
	}

	filterUser := filters.UserDB{
		Login:    login,
		PassHash: passHash,
	}

	var user domain.User
	user.ID, err = a.user.Create(ctx, filterUser)
	if err != nil {
		a.logs.Error("user has not been created: ", err)
		return 0, err
	}

	if err != nil {
		a.logs.Error("user creation failed:", err)
		return 0, err
	}

	a.logs.Info("user registered")

	return user.ID, nil
}

// Login TODO: возращать токен
func (a *Auth) Login(ctx context.Context, login, password string) (string, error) {
	var userFilter filters.UserDB

	userFilter, err := a.user.GetUser(ctx, login)
	if err != nil {
		a.logs.Error("GetUcser", err) // TODO:...
		return "", err
	}

	user := domain.User{
		ID:       userFilter.ID,
		Login:    login,
		PassHash: userFilter.PassHash,
	}

	if err = bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		a.logs.Error("Invalid password", err)
		return "", errInvailidCredentials
	}
	userClaims := dto.MapFromUser(user)

	tokenJWT, err := jwt.NewToken(userClaims, a.tokenTTL)
	if err != nil {

		a.logs.Error("NewToken:", err)
		return "", err
	}

	a.logs.Info("user logged in")
	return tokenJWT, nil
}

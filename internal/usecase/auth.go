package usecase

import (
	"auth-grpc/internal"
	"auth-grpc/internal/domain"
	"auth-grpc/internal/domain/filters"
	"auth-grpc/internal/jwt"
	"auth-grpc/internal/jwt/dto"
	"auth-grpc/internal/repository/redisRepo"
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
	logs       *logrus.Logger
	user       internal.User
	redis      *redisRepo.RepositoryRedis
	JwtManager *jwt.Manager
}

// TODO:...
func New(logs *logrus.Logger, user internal.User, redis *redisRepo.RepositoryRedis, jwtManager *jwt.Manager) *Auth {
	return &Auth{
		logs:       logs,
		user:       user,
		redis:      redis,
		JwtManager: jwtManager,
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

	exists, err := a.user.CheckUser(ctx, login)
	if err != nil {
		a.logs.Error("CheckUser", err)
		return 0, errors.New("") //TODO: добавить ошибку
	}
	if exists {
		return 0, errors.New("user already exists")
	}

	user.ID, err = a.user.Create(ctx, filterUser)

	if err != nil {
		a.logs.Error("Create: ", err)
		return 0, errors.New("user has not been created")
	}

	a.logs.Info("user registered")

	return user.ID, nil
}

func (a *Auth) Login(ctx context.Context, login, password string) (domain.Token, error) {
	var userFilter filters.UserDB

	userFilter, err := a.user.GetUser(ctx, login)
	if err != nil {
		a.logs.Error("GetUser", err) // TODO:Добавить ошибки на то что нет логина и другая ошибка
		return domain.Token{}, err
	}

	user := domain.User{
		ID:       userFilter.ID,
		Login:    login,
		PassHash: userFilter.PassHash,
	}

	if err = bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		a.logs.Error("Invalid password", err)
		return domain.Token{}, errInvailidCredentials
	}

	//TODO: надо вынести в отдельный метод или пакет, который будет работать с токенами.
	claims := dto.MapFromUser(user.ID, user.Login)

	accessToken, err := a.JwtManager.NewAccess(claims)
	if err != nil {
		a.logs.Error("NewToken:", err)
		return domain.Token{}, err
	}

	refreshTTL := 7 * 24 * time.Hour
	refreshToken, err := a.JwtManager.NewRefresh(claims)
	if err != nil {
		return domain.Token{}, err
	}
	//
	if err := a.redis.SaveRefreshToken(ctx, user.ID, refreshToken, refreshTTL); err != nil {
		return domain.Token{}, err
	}

	tokenUser := domain.Token{Access: accessToken, Refresh: refreshToken}
	a.logs.Info("user logged in")
	return tokenUser, nil
}

func (a *Auth) Verify(string) error {
	return nil

	// Здесь принимаешь accessToken и:

	// валидируешь его через jwt.Parse.

	// если expired → ошибка.

	// если подпись невалидна → ошибка.

	// иначе всё ок.
}

// Refresh TODO:...
func (a *Auth) Refresh(ctx context.Context, refreshToken string) (domain.Token, error) {
	// валидируем refresh
	claims, err := a.JwtManager.ValidateRefresh(refreshToken)
	if err != nil {
		return domain.Token{}, err
	}

	userID := int64(claims["uid"].(float64))

	// сверяем с Redis
	stored, err := a.redis.GetRefreshToken(ctx, userID)
	if err != nil || stored != refreshToken {
		return domain.Token{}, errors.New("invalid refresh")
	}

	userClaims := dto.UserClaims{
		ID:    userID,
		Login: claims["login"].(string),
	}

	// новые токены
	accessToken, err := a.JwtManager.NewAccess(userClaims)
	if err != nil {
		return domain.Token{}, err
	}

	newRefreshToken, err := a.JwtManager.NewRefresh(userClaims)
	if err != nil {
		return domain.Token{}, err
	}

	// сохраняем новый refresh в Redis
	if err := a.redis.SaveRefreshToken(ctx, userID, newRefreshToken, a.JwtManager.RefreshTTL); err != nil {
		return domain.Token{}, err
	}

	return domain.Token{
		Access:  accessToken,
		Refresh: newRefreshToken,
	}, nil
}

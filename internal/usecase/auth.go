package usecase

import (
	"auth-grpc/internal"
	"auth-grpc/internal/domain"
	"auth-grpc/internal/domain/mappers"
	"auth-grpc/internal/jwt"
	"auth-grpc/internal/jwt/dto"
	"auth-grpc/internal/repository/postgres"
	"auth-grpc/internal/repository/redisRepo"
	"context"
	"errors"
	"time"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

var refreshTTL = 7 * 24 * time.Hour

type Auth struct {
	logs       *logrus.Logger
	user       internal.User
	redis      *redisRepo.RepositoryRedis
	JwtManager jwt.JWTManager
}

// TODO:...
func New(logs *logrus.Logger, user internal.User, redis *redisRepo.RepositoryRedis, jwtManager jwt.JWTManager) *Auth {
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
		return 0, internal.ErrInternal
	}

	filterUser := mappers.UserDB{
		Login:    login,
		PassHash: passHash,
	}

	exists, err := a.user.CheckUser(ctx, login)
	if err != nil {
		a.logs.Error("CheckUser", err)
		return 0, internal.ErrInternal
	}
	if exists {
		return 0, errors.New("user already exists")
	}

	userID, err := a.user.Create(ctx, filterUser)
	if err != nil {
		a.logs.Error("Create: ", err)
		return 0, internal.ErrInternal
	}

	a.logs.Info("user registered")

	return userID, nil
}

func (a *Auth) Login(ctx context.Context, login, password string) (domain.Token, error) {
	var userFilter mappers.UserDB

	userFilter, err := a.user.GetUser(ctx, login)
	if err != nil {
		if errors.Is(err, postgres.ErrUserNotFound) {
			return domain.Token{}, err
		}
		a.logs.Error("GetUser", err)
		return domain.Token{}, internal.ErrInternal
	}

	user := domain.User{
		ID:       userFilter.ID,
		Login:    login,
		PassHash: userFilter.PassHash,
	}

	if err = bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		a.logs.Error("Invalid password", err)
		return domain.Token{}, internal.ErrInvailidCredentials
	}

	tokenUser, err := a.GenerateTokens(user.ID, user.Login)
	if err != nil {
		a.logs.Error("GenerateTokens: ", err)
		return domain.Token{}, internal.ErrInternal
	}

	if err := a.redis.SaveRefreshToken(ctx, user.ID, tokenUser.Refresh, refreshTTL); err != nil {
		a.logs.Error("SaveRefreshToken: ", err)
		return domain.Token{}, internal.ErrInternal
	}

	a.logs.Info("user logged in")
	return tokenUser, nil
}

func (a *Auth) GenerateTokens(id int64, login string) (domain.Token, error) {
	claims := dto.MapFromUser(id, login)

	accessToken, err := a.JwtManager.NewAccess(claims)
	if err != nil {
		a.logs.Error("NewToken:", err)
		return domain.Token{}, internal.ErrInternal
	}

	refreshToken, err := a.JwtManager.NewRefresh(claims)
	if err != nil {
		a.logs.Error("NewRefresh: ", err)
		return domain.Token{}, internal.ErrInternal
	}
	return domain.Token{Access: accessToken, Refresh: refreshToken}, nil
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
		if errors.Is(err, jwt.ErrInvailidToken) {
			return domain.Token{}, err
		}
		a.logs.Error("ValidateRefresh: ", err)
		return domain.Token{}, internal.ErrInternal
	}

	userID := int64(claims["uid"].(float64))

	// сверяем с Redis
	stored, err := a.redis.GetRefreshToken(ctx, userID)
	if err != nil || stored != refreshToken {
		return domain.Token{}, errors.New("invalid refresh")
	}

	// новые токены
	tokens, err := a.CreateTokens(userID, claims["login"].(string))
	if err != nil {
		return domain.Token{}, err
	}
	// сохраняем новый refresh в Redis
	if err := a.redis.SaveRefreshToken(ctx, userID, tokens.Refresh, a.JwtManager.GetRefreshTTL()); err != nil {
		a.logs.Error("SaveRefreshToken: ", err)
		return domain.Token{}, internal.ErrInternal
	}

	return tokens, nil
}

func (a *Auth) CreateTokens(id int64, login string) (domain.Token, error) {
	userClaims := dto.UserClaims{
		ID:    id,
		Login: login,
	}

	accessToken, err := a.JwtManager.NewAccess(userClaims)
	if err != nil {
		a.logs.Error("NewAccess: ", err)
		return domain.Token{}, internal.ErrInternal
	}

	newRefreshToken, err := a.JwtManager.NewRefresh(userClaims)
	if err != nil {
		a.logs.Error("NewRefresh: ", err)
		return domain.Token{}, internal.ErrInternal
	}
	return domain.Token{
		UserID:  id,
		Access:  accessToken,
		Refresh: newRefreshToken,
	}, nil
}

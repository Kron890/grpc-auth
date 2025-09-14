package jwt

import (
	"auth-grpc/internal/jwt/dto"
	"time"

	"github.com/golang-jwt/jwt"
)

// JWTManager интерфейс для работы с JWT токенами
type JWTManager interface {
	NewAccess(user dto.UserClaims) (string, error)
	NewRefresh(user dto.UserClaims) (string, error)
	ValidateAccess(tokenStr string) (jwt.MapClaims, error)
	ValidateRefresh(tokenStr string) (jwt.MapClaims, error)
	GetRefreshTTL() time.Duration
}

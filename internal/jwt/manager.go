package jwt

import (
	"auth-grpc/internal/jwt/dto"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

var (
	ErrInvailidToken = errors.New("invalid token")
)

type Manager struct {
	secret     string
	accessTTL  time.Duration
	RefreshTTL time.Duration
}

func NewManager(accessTTL, refreshTTL time.Duration) (*Manager, error) {
	secret := os.Getenv("SECRET")
	if secret == "" {
		return nil, fmt.Errorf("SECRET not set") //TODO:...
	}

	return &Manager{
		secret:     secret,
		accessTTL:  accessTTL,
		RefreshTTL: refreshTTL,
	}, nil
}

// NewAccess генерирует access токен
func (m *Manager) NewAccess(user dto.UserClaims) (string, error) {
	claims := jwt.MapClaims{
		"uid":   user.ID,
		"login": user.Login,
		"exp":   time.Now().Add(m.accessTTL).Unix(),
		"typ":   "access",
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).
		SignedString([]byte(m.secret))
}

// NewRefresh генерирует refresh токен
func (m *Manager) NewRefresh(user dto.UserClaims) (string, error) {
	claims := jwt.MapClaims{
		"uid":   user.ID,
		"login": user.Login,
		"exp":   time.Now().Add(m.RefreshTTL).Unix(),
		"typ":   "refresh",
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).
		SignedString([]byte(m.secret))
}

// ValidateAccess проверяет access токен
func (m *Manager) ValidateAccess(tokenStr string) (jwt.MapClaims, error) {
	return m.parse(tokenStr, "access")
}

// ValidateRefresh проверяет refresh токен
func (m *Manager) ValidateRefresh(tokenStr string) (jwt.MapClaims, error) {
	return m.parse(tokenStr, "refresh")
}

func (m *Manager) parse(tokenStr, expectedType string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(m.secret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, ErrInvailidToken
	}

	if claims["typ"] != expectedType {
		return nil, fmt.Errorf("wrong token type: %v", claims["typ"])
	}

	return claims, nil
}

// GetRefreshTTL возвращает время жизни refresh токена
func (m *Manager) GetRefreshTTL() time.Duration {
	return m.RefreshTTL
}

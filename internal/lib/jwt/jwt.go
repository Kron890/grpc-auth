package jwt

import (
	"auth-grpc/internal/lib/jwt/dto"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

func NewToken(user dto.UserClaims, duration time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"login": user.Login,
		"uid":   user.ID,
		"exp":   time.Now().Add(duration).Unix(),
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims)

	secret := os.Getenv("SECRET")
	if secret == "" {
		return "", fmt.Errorf("secret not found")
	}
	signetToken, err := jwtToken.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return signetToken, nil
}

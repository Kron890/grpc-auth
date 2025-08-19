package jwt

import (
	"auth-grpc/internal/lib/jwt/dto"
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
		panic("secret nil") //TODO:...
	}

	signetToken, err := jwtToken.SignedString(secret)
	if err != nil {
		return "", err
	}
	return signetToken, nil
}

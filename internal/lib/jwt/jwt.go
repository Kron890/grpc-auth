package jwt

// import (
// 	"auth-grpc/internal/domain"
// 	"auth-grpc/internal/lib/jwt/dto"
// 	"fmt"
// 	"os"
// 	"time"

// 	"github.com/golang-jwt/jwt"
// 	"google.golang.org/grpc/codes"
// 	"google.golang.org/grpc/status"
// )

// // FIX: сделать структуру
// func NewAccess(user dto.UserClaims, duration time.Duration) (string, error) {
// 	claims := jwt.MapClaims{
// 		"login": user.Login,
// 		"uid":   user.ID,
// 		"exp":   time.Now().Add(duration).Unix(),
// 	}

// 	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims)

// 	secret := os.Getenv("SECRET")
// 	if secret == "" {
// 		return "", fmt.Errorf("secret not found") //TODO:...
// 	}
// 	signetToken, err := jwtToken.SignedString([]byte(secret))
// 	if err != nil {
// 		return "", err
// 	}
// 	return signetToken, nil
// }

// // FIX:
// func NewRefresh(user dto.UserClaims, duration time.Duration) (domain.Token, error) {
// 	return domain.Token{}, nil
// }

// func ParseToken(tokenString string) (bool, error) {
// 	// парсим токен
// 	secret := os.Getenv("SECRET")

// 	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
// 		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
// 			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
// 		}
// 		return secret, nil
// 	})
// 	if err != nil || !token.Valid {
// 		return false, status.Error(codes.Unauthenticated, "invalid token")
// 	}
// 	return true, nil
// }

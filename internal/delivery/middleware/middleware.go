package middleware

import (
	"auth-grpc/internal/jwt"
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// NewAuthInterceptor создает interceptor с внедренной зависимостью JWTManager
func NewAuthInterceptor(jwtManager jwt.JWTManager) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// Пропускаем аутентификацию для регистрации и логина
		if info.FullMethod == "/auth.AuthService/Register" || info.FullMethod == "/auth.AuthService/Login" {
			return handler(ctx, req)
		}

		// Получаем метаданные из контекста
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "metadata is not provided")
		}

		// Извлекаем токен из заголовков
		values := md.Get("authorization")
		if len(values) == 0 {
			return nil, status.Error(codes.Unauthenticated, "authorization token is not provided")
		}
		token := values[0]

		// Валидируем токен через JWTManager
		claims, err := jwtManager.ValidateAccess(token)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, "invalid token")
		}

		// Добавляем информацию о пользователе в контекст для дальнейшего использования
		type userIDKey struct{}
		type userLoginKey struct{}
		ctx = context.WithValue(ctx, userIDKey{}, claims["uid"])
		ctx = context.WithValue(ctx, userLoginKey{}, claims["login"])

		return handler(ctx, req)
	}
}

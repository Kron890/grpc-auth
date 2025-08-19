package dto

import (
	"auth-grpc/internal/domain"
)

type UserClaims struct {
	ID    int64
	Login string
}

func MapFromUser(user domain.User) UserClaims {
	return UserClaims{
		ID:    user.ID,
		Login: user.Login,
	}
}

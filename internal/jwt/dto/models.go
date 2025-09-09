package dto

type UserClaims struct {
	ID    int64
	Login string
}

func MapFromUser(id int64, login string) UserClaims {
	return UserClaims{
		ID:    id,
		Login: login,
	}
}

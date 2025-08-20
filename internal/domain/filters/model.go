package filters

import "auth-grpc/internal/domain"

type UserDB struct {
	ID       int64
	Login    string
	PassHash []byte
}

func UserToUserDB(u domain.User) UserDB {
	return UserDB{Login: u.Login, PassHash: u.PassHash}
}

package domain

type User struct {
	ID       int64
	Login    string
	passHash []byte
}

package domain

type Token struct {
	UserID  int64
	Refresh string
	Access  string
}

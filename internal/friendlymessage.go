package internal

import "errors"

var (
	ErrInvailidCredentials = errors.New("invalid credentials")
	ErrInternal            = errors.New("internal error")
)

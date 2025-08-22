package domain

import "time"

type RefreshToken struct {
	UserID    string        `json:"user_id"`
	CreatedAt int64         `json:"created_at"`
	ExpiresIn time.Duration `json:"expires_in"`
}

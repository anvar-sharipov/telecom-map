package domain

import (
	"time"
)

type RefreshToken struct {
	ID        int64
	UserID    int64
	TokenHash string
	UserAgent string
	IPAddress string
	ExpiresAt time.Time
	CreatedAt time.Time
}

package domain

import (
	"time"
)

type User struct {
	ID        int64
	Username  string
	Password  string
	FullName  string
	Phone     *string
	IsActive  bool
	CreatedAt time.Time
}

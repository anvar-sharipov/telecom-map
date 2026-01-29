package users

import "time"

type AdminGroupDTO struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type AdminUserDTO struct {
	ID        int64           `json:"id"`
	Username  string          `json:"username"`
	FullName  string          `json:"full_name"`
	IsActive  bool            `json:"is_active"`
	CreatedAt time.Time       `json:"created_at"`
	Groups    []AdminGroupDTO `json:"groups"`
}

type CreateUserRequest struct {
	Username string  `json:"username"`
	Password string  `json:"password"`
	FullName string  `json:"full_name"`
	GroupIDs []int64 `json:"group_ids"`
	IsActive bool    `json:"is_active"`
}

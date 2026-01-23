package admin

import (
	"time"

	"github.com/anvar-sharipov/telecom-map/internal/domain"
)

type GroupDTO struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	IsActive    bool   `json:"is_active"`
}

type UserDTO struct {
	ID        int64      `json:"id"`
	Username  string     `json:"username"`
	FullName  string     `json:"full_name"`
	IsActive  bool       `json:"is_active"`
	CreatedAt string     `json:"created_at"`
	Groups    []GroupDTO `json:"groups"`
}

func mapUserToDTO(u *domain.User) UserDTO {
	dto := UserDTO{
		ID:        u.ID,
		Username:  u.Username,
		FullName:  u.FullName,
		IsActive:  u.IsActive,
		CreatedAt: u.CreatedAt.Format(time.RFC3339),
		Groups:    []GroupDTO{},
	}

	for _, g := range u.Groups {
		dto.Groups = append(dto.Groups, GroupDTO{
			ID:          g.ID,
			Name:        g.Name,
			Description: g.Description,
			IsActive:    g.IsActive,
		})
	}

	return dto
}

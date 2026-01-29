package users

import "github.com/anvar-sharipov/telecom-map/internal/domain"

func ToAdminDTO(u *domain.User) *AdminUserDTO {
	groups := make([]AdminGroupDTO, 0, len(u.Groups))
	for _, g := range u.Groups {
		groups = append(groups, AdminGroupDTO{
			ID:   g.ID,
			Name: g.Name,
		})
	}

	return &AdminUserDTO{
		ID:        u.ID,
		Username:  u.Username,
		FullName:  u.FullName,
		IsActive:  u.IsActive,
		CreatedAt: u.CreatedAt,
		Groups:    groups,
	}
}

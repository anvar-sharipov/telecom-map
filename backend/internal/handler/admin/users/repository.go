package users

import (
	"context"

	"github.com/anvar-sharipov/telecom-map/internal/domain"
)

type UserRepository interface {
	ListWithGroups(ctx context.Context) ([]*domain.User, error)
	CreateWithGroups(
		ctx context.Context,
		user *domain.User,
		groupIDs []int64,
	) error
}

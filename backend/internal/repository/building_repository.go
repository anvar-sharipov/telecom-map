package repository

import (
	"context"

	"github.com/anvar-sharipov/telecom-map/internal/domain"
)

// BuildingRepository — контракт для работы со зданиями
type BuildingRepository interface {
	Create(ctx context.Context, building *domain.Building) error
}

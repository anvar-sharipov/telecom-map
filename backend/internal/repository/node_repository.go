package repository

import (
	"context"

	"github.com/anvar-sharipov/telecom-map/internal/domain"
	"github.com/google/uuid"
)

type NodeRepository interface {
	Create(ctx context.Context, node *domain.Node) error
	GetAll(ctx context.Context) ([]domain.Node, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Node, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

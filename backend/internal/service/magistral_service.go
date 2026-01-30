package service

import (
	"context"

	"github.com/anvar-sharipov/telecom-map/internal/domain"
	"github.com/anvar-sharipov/telecom-map/internal/repository"
)

type MagistralService struct {
	repo *repository.MagistralRepository
}

func NewMagistralService(repo *repository.MagistralRepository) *MagistralService {
	return &MagistralService{repo: repo}
}

func (s *MagistralService) Create(
	ctx context.Context,
	name, description string,
) (*domain.Magistral, error) {
	m := &domain.Magistral{
		Name:        name,
		Description: description,
	}

	if err := s.repo.Create(ctx, m); err != nil {
		return nil, err
	}

	return m, nil
}

func (s *MagistralService) List(
	ctx context.Context,
) ([]domain.Magistral, error) {
	return s.repo.List(ctx)
}

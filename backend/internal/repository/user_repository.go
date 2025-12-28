package repository

import (
	"github.com/anvar-sharipov/telecom-map/internal/domain"
)

type UserRepository interface {
	Create(user *domain.User) error
	GetByID(id int64) (*domain.User, error)
	GetByPhone(phone string) (*domain.User, error)
	List() ([]*domain.User, error)
}

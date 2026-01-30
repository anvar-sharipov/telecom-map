package repository

import (
	"context"

	"github.com/anvar-sharipov/telecom-map/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MagistralRepository struct {
	db *pgxpool.Pool
}

func NewMagistralRepository(db *pgxpool.Pool) *MagistralRepository {
	return &MagistralRepository{db: db}
}

func (r *MagistralRepository) Create(
	ctx context.Context,
	m *domain.Magistral,
) error {
	query := `
		INSERT INTO magistrals (name, description)
		VALUES ($1, $2)
		RETURNING id, created_at
	`

	return r.db.QueryRow(
		ctx,
		query,
		m.Name,
		m.Description,
	).Scan(&m.ID, &m.CreatedAt)
}

func (r *MagistralRepository) List(
	ctx context.Context,
) ([]domain.Magistral, error) {
	rows, err := r.db.Query(
		ctx,
		`SELECT id, name, description, created_at
		 FROM magistrals
		 ORDER BY created_at DESC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []domain.Magistral

	for rows.Next() {
		var m domain.Magistral
		if err := rows.Scan(
			&m.ID,
			&m.Name,
			&m.Description,
			&m.CreatedAt,
		); err != nil {
			return nil, err
		}
		result = append(result, m)
	}

	return result, nil
}

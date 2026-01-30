package postgres

import (
	"context"

	"github.com/anvar-sharipov/telecom-map/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type NodeRepository struct {
	pool *pgxpool.Pool
}

func NewNodeRepository(pool *pgxpool.Pool) *NodeRepository {
	return &NodeRepository{pool: pool}
}

func (r *NodeRepository) Create(ctx context.Context, node *domain.Node) error {
	query := `
		INSERT INTO nodes (name, type, lon, lat)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at
	`
	return r.pool.QueryRow(
		ctx,
		query,
		node.Name,
		node.Type,
		node.Lon,
		node.Lat,
	).Scan(&node.ID, &node.CreatedAt)
}

func (r *NodeRepository) GetAll(ctx context.Context) ([]domain.Node, error) {
	query := `
		SELECT id, name, type, lon, lat, created_at
		FROM nodes
		ORDER BY created_at DESC
	`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var nodes []domain.Node
	for rows.Next() {
		var n domain.Node
		err := rows.Scan(
			&n.ID,
			&n.Name,
			&n.Type,
			&n.Lon,
			&n.Lat,
			&n.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		nodes = append(nodes, n)
	}

	return nodes, nil
}

func (r *NodeRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Node, error) {
	query := `
		SELECT id, name, type, lon, lat, created_at
		FROM nodes
		WHERE id = $1
	`

	var n domain.Node
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&n.ID,
		&n.Name,
		&n.Type,
		&n.Lon,
		&n.Lat,
		&n.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &n, nil
}

func (r *NodeRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM nodes WHERE id = $1`, id)
	return err
}

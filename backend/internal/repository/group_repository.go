package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Group struct {
	ID          int64
	Name        string
	Description string
	IsActive    bool
	CreatedAt   time.Time
}

type GroupRepository struct {
	db *pgxpool.Pool
}

func NewGroupRepository(db *pgxpool.Pool) *GroupRepository {
	return &GroupRepository{db: db}
}

func (r *GroupRepository) Create(ctx context.Context, name, description string) (int64, error) {
	query := `INSERT INTO groups (name, description) VALUES ($1, $2) RETURNING id`
	var id int64

	err := r.db.QueryRow(ctx, query, name, description).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to create group: %w", err)
	}

	return id, nil
}

func (r *GroupRepository) GetById(ctx context.Context, id int64) (*Group, error) {
	var g Group
	query := `SELECT id, name, description, is_active, created_at FROM groups WHERE id = $1`

	err := r.db.QueryRow(ctx, query, id).Scan(
		&g.ID,
		&g.Name,
		&g.Description,
		&g.IsActive,
		&g.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get group by id: %w", err)
	}

	return &g, nil
}

func (r *GroupRepository) ListAll(ctx context.Context) ([]Group, error) {
	query := `SELECT id, name, description, is_active, created_at FROM groups`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list groups: %w", err)
	}

	defer rows.Close()

	var groups []Group
	for rows.Next() {
		var g Group
		if err := rows.Scan(
			&g.ID,
			&g.Name,
			&g.Description,
			&g.IsActive,
			&g.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan group: %w", err)
		}
		groups = append(groups, g)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return groups, nil
}

func (r *GroupRepository) ListAllActive(ctx context.Context) ([]Group, error) {
	query := `
		SELECT id, name, description, is_active, created_at 
		FROM groups
		WHERE is_active = true
	`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list groups: %w", err)
	}

	defer rows.Close()

	var groups []Group
	for rows.Next() {
		var g Group
		if err := rows.Scan(
			&g.ID,
			&g.Name,
			&g.Description,
			&g.IsActive,
			&g.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan group: %w", err)
		}
		groups = append(groups, g)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return groups, nil
}

func (r *GroupRepository) Disable(ctx context.Context, id int64) error {
	query := `UPDATE groups SET is_active = false WHERE id = $1`

	_, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to disable group: %w", err)
	}
	return nil
}

func (r *GroupRepository) Enable(ctx context.Context, id int64) error {
	query := `UPDATE groups SET is_active = true WHERE id = $1`

	_, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to enable group: %w", err)
	}
	return nil
}

func (r *GroupRepository) Update(
	ctx context.Context,
	id int64,
	name string,
	description string,
) error {
	query := `
		UPDATE groups
		SET name = $1,
			description = $2
		WHERE id = $3
	`
	cmd, err := r.db.Exec(ctx, query, name, description, id)
	if err != nil {
		return fmt.Errorf("failed to update group: %w", err)
	}

	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("group not found")
	}
	return nil
}

func (r *GroupRepository) ListByActive(ctx context.Context, active bool) ([]Group, error) {
	const query = `
		SELECT id, name, description, is_active, created_at
		FROM groups
		WHERE is_active = $1
		ORDER BY name
	`

	rows, err := r.db.Query(ctx, query, active)
	if err != nil {
		return nil, fmt.Errorf("failed to list groups by active: %w", err)
	}
	defer rows.Close()

	var groups []Group
	for rows.Next() {
		var g Group
		if err := rows.Scan(
			&g.ID,
			&g.Name,
			&g.Description,
			&g.IsActive,
			&g.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan group: %w", err)
		}
		groups = append(groups, g)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return groups, nil
}

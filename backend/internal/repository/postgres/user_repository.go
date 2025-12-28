package postgres

import (
	"context"
	"time"

	"github.com/anvar-sharipov/telecom-map/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) Create(user *domain.User) error {
	query := `
		INSERT INTO users (username, full_name, phone, password, is_active)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at
	`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := r.db.QueryRow(
		ctx,
		query,
		user.Username,
		user.FullName,
		user.Phone,
		user.Password,
		user.IsActive,
	).Scan(&user.ID, &user.CreatedAt)

	return err
}

func (r *UserRepository) GetByID(id int64) (*domain.User, error) {
	query := `
	SELECT id, full_name, phone, is_active, created_at
	From users
	WHERE id = $1
	`

	row := r.db.QueryRow(context.Background(), query, id)

	var user domain.User
	err := row.Scan(
		&user.ID,
		&user.FullName,
		&user.Phone,
		&user.IsActive,
		&user.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetByPhone(phone string) (*domain.User, error) {
	query := `
		SELECT id, full_name, phone, is_active, created_at
		FROM users
		WHERE phone = $1
	`

	row := r.db.QueryRow(context.Background(), query, phone)

	var user domain.User
	err := row.Scan(
		&user.ID,
		&user.FullName,
		&user.Phone,
		&user.IsActive,
		&user.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) List() ([]*domain.User, error) {
	query := `
	SELECT id, full_name, phone, is_active, created_at
	FROM users
	ORDER BY id
	`

	rows, err := r.db.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var users []*domain.User

	for rows.Next() {
		var user domain.User
		if err := rows.Scan(
			&user.ID,
			&user.FullName,
			&user.Phone,
			&user.IsActive,
			&user.CreatedAt,
		); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	return users, nil
}

package repository

import (
	"context"
	"time"

	"github.com/anvar-sharipov/telecom-map/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RefreshTokenRepository struct {
	db *pgxpool.Pool
}

func NewRefreshTokenRepository(db *pgxpool.Pool) *RefreshTokenRepository {
	return &RefreshTokenRepository{db: db}
}

func (r *RefreshTokenRepository) Create(ctx context.Context, userID int64, token string, expires time.Time) error {
	query := `
	INSERT INTO refresh_tokens (user_id, token, expires_at)
	VALUES ($1, $2, $3)
	`

	_, err := r.db.Exec(
		ctx,
		query,
		userID,
		token,
		expires,
	)
	return err
}

// 3️⃣ Добавим получение refresh token (важно!)
// Когда пользователь обновляет access token, мы должны найти refresh token:
func (r *RefreshTokenRepository) GetByToken(ctx context.Context, token string) (*domain.RefreshToken, error) {
	rt := &domain.RefreshToken{}
	query := `
		SELECT id, user_id, token, expires_at, created_at
		FROM refresh_tokens
		WHERE token = $1
	`

	err := r.db.QueryRow(ctx, query, token).Scan(&rt.ID, &rt.UserID, &rt.Token, &rt.ExpiresAt)
	if err != nil {
		return nil, err
	}

	return rt, nil
}

// 4️⃣ Удаление refresh token (logout / rotation)
// Очень важно для безопасности:
func (r *RefreshTokenRepository) Delete(ctx context.Context, token string) error {
	query := `
		DELETE FROM refresh_tokens 
		WHERE token = $1
	`
	_, err := r.db.Exec(ctx, query, token)
	return err
}

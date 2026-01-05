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

func (r *RefreshTokenRepository) Create(
	ctx context.Context,
	userID int64,
	tokenHash string,
	userAgent string,
	ipAddress string,
	expires time.Time,
) error {
	query := `
	INSERT INTO refresh_tokens (
		user_id, 
		token_hash, 
		user_agent,
		ip_address,
		expires_at
	)
	VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.db.Exec(
		ctx,
		query,
		userID,
		tokenHash,
		userAgent,
		ipAddress,
		expires,
	)
	return err
}

// 3️⃣ Добавим получение refresh token (важно!)
// Когда пользователь обновляет access token, мы должны найти refresh token:
func (r *RefreshTokenRepository) GetByTokenHash(
	ctx context.Context,
	tokenHash string,
) (*domain.RefreshToken, error) {
	rt := &domain.RefreshToken{}
	query := `
		SELECT 
		id, 
		user_id, 
		token_hash,
		user_agent,
		ip_address, 
		expires_at, 
		created_at
	FROM refresh_tokens
	WHERE token_hash = $1
	`

	err := r.db.QueryRow(ctx, query, tokenHash).Scan(
		&rt.ID,
		&rt.UserID,
		&rt.TokenHash,
		&rt.UserAgent,
		&rt.IPAddress,
		&rt.ExpiresAt,
		&rt.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return rt, nil
}

// 4️⃣ Удаление refresh token (logout / rotation)
// Очень важно для безопасности:
func (r *RefreshTokenRepository) Delete(ctx context.Context, tokenHash string) error {
	query := `
		DELETE FROM refresh_tokens 
		WHERE token_hash = $1
	`
	_, err := r.db.Exec(ctx, query, tokenHash)
	return err
}

// Один refresh token на устройство (РЕКОМЕНДУЮ)
// 👉 Новый логин НЕ удаляет другие,
// 👉 но перезаписывает токен для того же устройства
// Как определить устройство:
// user_agent
// ip_address (опционально)
func (r *RefreshTokenRepository) DeleteByUserAndAgent(
	ctx context.Context,
	userID int64,
	userAgent string,
) error {
	_, err := r.db.Exec(ctx, `
	DELETE FROM refresh_tokens
	WHERE user_id = $1 AND user_agent = $2
	`, userID, userAgent)

	return err
}

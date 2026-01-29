package service

import (
	"context"
	"log"
	"time"

	"github.com/anvar-sharipov/telecom-map/internal/repository"
	"github.com/anvar-sharipov/telecom-map/internal/utils"
)

type AuthService struct {
	RefreshTokenRepo *repository.RefreshTokenRepository
}

// SessionData содержит все данные для создания сессии
type SessionData struct {
	UserID    int64
	UserAgent string
	IP        string
}

// TokenPair содержит оба токена и срок действия
type TokenPair struct {
	AccessToken  string
	RefreshToken string
	ExpiresAt    time.Time
}

// CreateSession создаёт access + refresh токены, сохраняет в БД
func (s *AuthService) CreateSession(ctx context.Context, data SessionData) (*TokenPair, error) {
	// 1️⃣ Генерируем ACCESS TOKEN
	accessToken, err := utils.GenerateToken(data.UserID)
	if err != nil {
		return nil, err
	}
	// 2️⃣ Генерируем REFRESH TOKEN
	refreshToken, err := utils.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	// 3️⃣ Хешируем refresh token
	refreshTokenHash := utils.HashToken(refreshToken)

	// 4️⃣ Срок действия (30 секунд для теста, потом 7 дней)
	// expiresAt := time.Now().Add(10 * time.Second)
	expiresAt := time.Now().Add(7 * 24 * time.Hour) // продакшн

	// 5️⃣ Удаляем старый токен для этого устройства (если есть)
	err = s.RefreshTokenRepo.DeleteByUserAndAgent(ctx, data.UserID, data.UserAgent)
	if err != nil {
		log.Println("warning: failed to cleanup old refresh token:", err)
		// не прерываем выполнение
	}

	// 6️⃣ Сохраняем новый refresh token
	err = s.RefreshTokenRepo.Create(
		ctx,
		data.UserID,
		refreshTokenHash,
		data.UserAgent,
		data.IP,
		expiresAt,
	)
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
	}, nil
}

// RotateSession удаляет старый refresh token и создаёт новую пару токенов
func (s *AuthService) RotateSession(ctx context.Context, data SessionData) (*TokenPair, error) {
	// Точно такая же логика, как CreateSession
	// (rotation = delete old + create new)
	return s.CreateSession(ctx, data)
}

// InvalidateSession удаляет refresh token из БД
func (s *AuthService) InvalidateSession(ctx context.Context, userID int64, userAgent string) error {
	return s.RefreshTokenRepo.DeleteByUserAndAgent(ctx, userID, userAgent)
}

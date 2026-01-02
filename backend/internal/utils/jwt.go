package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(userID int64) (string, error) {
	secret := os.Getenv("JWT_SECRET") // секрет берём из .env
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims) // создаём JWT
	return token.SignedString([]byte(secret))                  // подписываем
}

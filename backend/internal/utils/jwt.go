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
		// "exp": time.Now().Add(10 * time.Second).Unix(), // 🔑 Access token — 10 секунд
		"iat": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims) // создаём JWT
	return token.SignedString([]byte(secret))                  // подписываем
}

func ParseToken(tokenString string) (int64, error) {
	secret := os.Getenv("JWT_SECRET")

	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(secret), nil
	})
	if err != nil || !token.Valid {
		return 0, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, jwt.ErrTokenInvalidClaims
	}

	// 🔐 exp
	exp, ok := claims["exp"].(float64)
	if !ok {
		return 0, jwt.ErrTokenExpired
	}
	if time.Now().Unix() > int64(exp) {
		return 0, jwt.ErrTokenExpired
	}

	// 👤 user_id
	userID, ok := claims["user_id"].(float64)
	if !ok {
		return 0, jwt.ErrTokenInvalidClaims
	}

	return int64(userID), nil
}

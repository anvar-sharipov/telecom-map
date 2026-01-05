package utils

import (
	"crypto/rand"
	"encoding/base64"
)

// GenerateRefreshToken генерирует криптостойкий refresh token
func GenerateRefreshToken() (string, error) {
	// 32 байта = 256 бит безопасности
	b := make([]byte, 32)

	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	// URL-safe строка (можно класть в cookie / header)
	return base64.RawStdEncoding.EncodeToString(b), nil
}

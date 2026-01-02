package utils

import (
	"crypto/rand"
	"encoding/base64"
)

// GenerateRefreshToken creates a secure random string
func GenerateRefreshToken() (string, error) {
	b := make([]byte, 64) // 512 бит
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(b), nil
}

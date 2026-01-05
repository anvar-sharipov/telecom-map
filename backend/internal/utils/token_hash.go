package utils

import (
	"crypto/sha256"
	"encoding/hex"
)

// HashToken хеширует refresh token перед сохранением в БД
func HashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

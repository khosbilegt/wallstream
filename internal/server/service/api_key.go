package service

import (
	"crypto/rand"
	"encoding/hex"
)

// GenerateAPIKey generates a secure random API key
func GenerateAPIKey() (string, error) {
	bytes := make([]byte, 32) // 256 bits
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

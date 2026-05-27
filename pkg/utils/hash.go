package utils

import (
	"crypto/sha256"
	"encoding/hex"
)

func HashToken(token string) string {
	hash := sha256.New()
	hash.Write([]byte(token))
	return hex.EncodeToString(hash.Sum(nil))
}
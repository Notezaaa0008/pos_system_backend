package utils

import (
	"crypto/rand"
	"encoding/hex"
)

func RandomToken() string{
	b := make([]byte, 32)
	rand.Read(b)
	token := hex.EncodeToString(b)

	return token
}
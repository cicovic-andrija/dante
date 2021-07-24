package util

import (
	"crypto/rand"
	"encoding/hex"
)

func RandHexString(n int) (string, error) {
	if n < 1 {
		return "", nil
	}

	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

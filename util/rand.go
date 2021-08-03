package util

import (
	"crypto/rand"
	"encoding/hex"
)

// RandHexString returns a string of a specified length,
// that consists of hexadecimal digits.
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

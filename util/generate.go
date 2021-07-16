package util

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"
)

func GenerateHexString(n int) (string, error) {
	if n < 1 {
		return "", nil
	}

	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func GenerateLogName(server string, rotation int, logType string) string {
	return fmt.Sprintf(
		"%s_%s_%03d_%s.log",
		server,
		time.Now().Format("D20060102_T150405"),
		rotation,
		logType,
	)
}

package util

import (
	"regexp"
)

var uuidv4Regex = regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")

// IsValidUUIDv4 returns a flag indicating whether a string
// is of a valid UUID v4 format.
func IsValidUUIDv4(str string) bool {
	return uuidv4Regex.MatchString(str)
}

package util

// SearchForString returns a flag indicating whether a string value
// is found in a string slice.
func SearchForString(str string, strings ...string) (found bool) {
	for _, s := range strings {
		if str == s {
			found = true
			return
		}
	}
	return
}

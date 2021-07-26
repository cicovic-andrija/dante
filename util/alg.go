package util

func SearchForString(str string, strings ...string) (found bool) {
	for _, s := range strings {
		if str == s {
			found = true
			return
		}
	}
	return
}

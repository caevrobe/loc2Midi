package helpers

// returns string stripped of newline characters
func StripNL(s string) string {
	last := len(s) - 1
	for x := last; x > -1; x-- {
		if s[x] == '\r' || s[x] == '\n' {
			last = x
		} else {
			break
		}
	}

	return s[:last]
}

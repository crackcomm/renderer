package helpers

// Contain - Returns true if slice contains string.
func Contain(list []string, str string) bool {
	for _, s := range list {
		if s == str {
			return true
		}
	}
	return false
}

// MergeUnique - Merges `source` string slice into `dest` and returns result.
// Inserts from `source` only when `dest` does not `Contain` given string.
func MergeUnique(dest, source []string) []string {
	for _, str := range source {
		if !Contain(dest, str) {
			dest = append(dest, str)
		}
	}
	return dest
}

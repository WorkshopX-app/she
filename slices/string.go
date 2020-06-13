package slices

//Contains reports whether str or modifier(str) is within slice.
func Contains(slice []string, str string, modifier func(s string) string) bool {
	for _, she := range slice {
		if she == str {
			return true
		}
		if modifier != nil && modifier(str) == str {
			return true
		}
	}
	return false
}

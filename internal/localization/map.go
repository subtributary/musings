package localization

// FindMatch finds the value with the best locale match in a map of locale to value.
func FindMatch[T any](haystack map[string]T, locale string) (T, bool) {
	return FindMatchCb(haystack, locale, func(T) bool {
		return true
	})
}

// FindMatchCb finds the value with the best locale match in a map of locale to value.
//
// The `foundCb` callback is called for every match, and it can return false to continue searching.
func FindMatchCb[T any](haystack map[string]T, locale string, foundCb func(T) bool) (value T, isFound bool) {
	// todo: get the priorities in a smarter way
	priority := []string{locale, ""}

	for _, p := range priority {
		if value, isFound = haystack[p]; isFound && foundCb(value) {
			break
		}
	}

	return
}

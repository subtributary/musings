package localization

import "strings"

// ExtractLocale parses the locale out of the first segment of a path.
// It returns the parsed locale and the remaining path after that.
// If the locale is invalid or missing, then it returns (UndLocale, path).
func ExtractLocale(locales []Locale, path string) (Locale, string) {
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	segments := strings.SplitN(path, "/", 3)

	// Search for the locale in the supported locales and return it if found.
	for _, locale := range locales {
		if strings.EqualFold(locale.Tag, segments[1]) {
			if len(segments) == 3 {
				return locale, "/" + segments[2]
			}
			return locale, "/"
		}
	}

	return UndLocale, path
}

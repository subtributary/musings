package localization

import (
	"strings"
)

// ExtractLocale parses the locale and trailing path out of a localized URL.
// If the locale is invalid or missing, then UndLocale is returned.
// The trailing path is always returned with a "/" prefix.
// The name argument must use "/" for a path separator.
func ExtractLocale(locales []Locale, name string) (Locale, string) {
	name, _ = strings.CutPrefix(name, "/")

	if len(locales) == 1 && locales[0].Tag == Und {
		return UndLocale, "/" + name
	}

	segments := strings.SplitN(name, "/", 2)

	// Search for the locale in the supported locales
	for _, locale := range locales {
		if strings.EqualFold(locale.Tag, segments[0]) {
			if len(segments) == 2 {
				return locale, "/" + segments[1]
			}
			return locale, "/"
		}
	}

	return UndLocale, "/" + name
}

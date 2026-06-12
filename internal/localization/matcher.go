package localization

import (
	"log"

	"golang.org/x/text/language"
)

// LocaleMatcher matches an Accept-Language value to a configured locale.
type LocaleMatcher struct {
	locales []Locale
	matcher language.Matcher
	tags    []language.Tag
}

func NewLocaleMatcher(locales []Locale) LocaleMatcher {
	tags := make([]language.Tag, len(locales))
	for i, locale := range locales {
		tag, err := language.Parse(locale.Tag)
		if err != nil {
			// This shouldn't happen. Tags are validated at config load.
			log.Fatalf("Error: unable to parse tag %q", locale.Tag)
		}
		tags[i] = tag
	}

	return LocaleMatcher{
		locales: locales,
		matcher: language.NewMatcher(tags),
		tags:    tags,
	}
}

// Choose chooses the best locale match for an accept-language header.
// If no best match is found, then the first configured locale is returned.
// If no locales are configured, then UndLocale is returned.
func (m *LocaleMatcher) Choose(acceptLanguage string) Locale {
	if len(m.locales) == 0 {
		return UndLocale
	}

	_, i := language.MatchStrings(m.matcher, acceptLanguage)
	return m.locales[i]
}

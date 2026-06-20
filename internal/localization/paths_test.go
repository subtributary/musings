package localization_test

import (
	"testing"

	"github.com/subtributary/musings/internal/localization"
)

func TestExtractLocale(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		locales      []localization.Locale
		urlPath      string
		wantTag      string
		wantTrailing string
	}{
		/* no locales */
		{
			name:         "no locales, empty route",
			urlPath:      "",
			wantTag:      "und",
			wantTrailing: "/",
		},
		{
			name:         "no locales, root path",
			urlPath:      "/",
			wantTag:      "und",
			wantTrailing: "/",
		},
		{
			name:         "no locales, path with slashes",
			urlPath:      "/test/",
			wantTag:      "und",
			wantTrailing: "/test/",
		},
		{
			name:         "no locales, path without leading slash",
			urlPath:      "test/",
			wantTag:      "und",
			wantTrailing: "/test/",
		},
		{
			name:         "no locales, path without trailing slash",
			urlPath:      "/test",
			wantTag:      "und",
			wantTrailing: "/test",
		},
		{
			name:         "no locales, path without slashes",
			urlPath:      "test",
			wantTag:      "und",
			wantTrailing: "/test",
		},

		/* one locale */
		{
			name:         "one locale, empty route",
			locales:      []localization.Locale{en},
			urlPath:      "",
			wantTag:      "und",
			wantTrailing: "/",
		},
		{
			name:         "one locale, root path",
			locales:      []localization.Locale{en},
			urlPath:      "/",
			wantTag:      "und",
			wantTrailing: "/",
		},
		{
			name:         "one locale, path with slashes",
			locales:      []localization.Locale{en},
			urlPath:      "/test/",
			wantTag:      "und",
			wantTrailing: "/test/",
		},
		{
			name:         "one locale, path without leading slash",
			locales:      []localization.Locale{en},
			urlPath:      "test/",
			wantTag:      "und",
			wantTrailing: "/test/",
		},
		{
			name:         "one locale, path without trailing slash",
			locales:      []localization.Locale{en},
			urlPath:      "/test",
			wantTag:      "und",
			wantTrailing: "/test",
		},
		{
			name:         "one locale, path without slashes",
			locales:      []localization.Locale{en},
			urlPath:      "test",
			wantTag:      "und",
			wantTrailing: "/test",
		},
		{
			name:         "one locale, localized root path",
			locales:      []localization.Locale{en},
			urlPath:      "/en/",
			wantTag:      "en",
			wantTrailing: "/",
		},
		{
			name:         "one locale, localized root without leading slash",
			locales:      []localization.Locale{en},
			urlPath:      "en/",
			wantTag:      "en",
			wantTrailing: "/",
		},
		{
			name:         "one locale, localized root without trailing slash",
			locales:      []localization.Locale{en},
			urlPath:      "/en",
			wantTag:      "en",
			wantTrailing: "/",
		},
		{
			name:         "one locale, localized root without slashes",
			locales:      []localization.Locale{en},
			urlPath:      "en",
			wantTag:      "en",
			wantTrailing: "/",
		},
		{
			name:         "one locale, localized path",
			locales:      []localization.Locale{en},
			urlPath:      "/en/test",
			wantTag:      "en",
			wantTrailing: "/test",
		},
		{
			name:         "one locale, nested path",
			locales:      []localization.Locale{en},
			urlPath:      "/en/sub/test",
			wantTag:      "en",
			wantTrailing: "/sub/test",
		},
		{
			// Not supported!
			name:         "one locale, Windows file path",
			locales:      []localization.Locale{en},
			urlPath:      `\en\test`,
			wantTag:      "und",
			wantTrailing: `/\en\test`,
		},

		/* two locales */
		{
			name:         "two locales, unlocalized path",
			locales:      []localization.Locale{ar, en},
			urlPath:      "/test",
			wantTag:      "und",
			wantTrailing: "/test",
		},
		{
			name:         "two locales, localized path",
			locales:      []localization.Locale{ar, en},
			urlPath:      "/en/test",
			wantTag:      "en",
			wantTrailing: "/test",
		},

		/* special */
		{
			name:         "invalid locale in path",
			locales:      []localization.Locale{en},
			urlPath:      "/xx/test",
			wantTag:      "und",
			wantTrailing: "/xx/test",
		},
		{
			name:         "und locale",
			locales:      []localization.Locale{und},
			urlPath:      "/und/test",
			wantTag:      "und",
			wantTrailing: "/und/test", // und does not change the path.
		},
		{
			name:         "locale with region",
			locales:      []localization.Locale{en, zhHans},
			urlPath:      "/zh-Hans/test",
			wantTag:      "zh-Hans",
			wantTrailing: "/test",
		},
		{
			name:         "locale with lowercase region",
			locales:      []localization.Locale{en, zhHans},
			urlPath:      "/zh-hans/test",
			wantTag:      "zh-Hans",
			wantTrailing: "/test",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			locale, trailing := localization.ExtractLocale(tt.locales, tt.urlPath)

			if tt.wantTag != locale.Tag {
				t.Errorf("ExtractLocale(%q): locale = %v, want %v",
					tt.urlPath, locale.Tag, tt.wantTag)
			}

			if tt.wantTrailing != trailing {
				t.Errorf("ExtractLocale(%q): trailing = %v, want %v",
					tt.urlPath, trailing, tt.wantTrailing)
			}
		})
	}
}

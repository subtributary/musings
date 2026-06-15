package localization_test

import (
	"testing"

	"github.com/subtributary/musings/internal/localization"
)

func TestExtractLocale(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		urlPath   string
		localeTag string
		trailing  string
	}{
		{
			name:      "empty path",
			urlPath:   "",
			localeTag: "und",
			trailing:  "/",
		},
		{
			name:      "root path",
			urlPath:   "/",
			localeTag: "und",
			trailing:  "/",
		},
		{
			name:      "localized root path",
			urlPath:   "/en/",
			localeTag: "en",
			trailing:  "/",
		},
		{
			name:      "localized root path without trailing slash",
			urlPath:   "/en",
			localeTag: "en",
			trailing:  "/",
		},
		{
			name:      "localized root path without leading slash",
			urlPath:   "en/",
			localeTag: "en",
			trailing:  "/",
		},
		{
			name:      "localized root path without slashes",
			urlPath:   "en",
			localeTag: "en",
			trailing:  "/",
		},
		{
			name:      "localized file path",
			urlPath:   "/en/index.html",
			localeTag: "en",
			trailing:  "/index.html",
		},
		{
			name:      "localized nested file path",
			urlPath:   "/en/sub/index.html",
			localeTag: "en",
			trailing:  "/sub/index.html",
		},
		{
			name:      "invalid locale",
			urlPath:   "/xx/index.html",
			localeTag: "und",
			trailing:  "/xx/index.html",
		},
		{
			name:      "locale with region",
			urlPath:   "/zh-Hans/index.html",
			localeTag: "zh-Hans",
			trailing:  "/index.html",
		},
		{
			name:      "locale with lowercase region",
			urlPath:   "/zh-hans/index.html",
			localeTag: "zh-Hans",
			trailing:  "/index.html",
		},
	}

	locales := []localization.Locale{en, zhHans}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			locale, trailing := localization.ExtractLocale(locales, tt.urlPath)

			if tt.localeTag != locale.Tag {
				t.Errorf("ExtractLocale(%q): locale = %v, want %v",
					tt.urlPath, locale.Tag, tt.localeTag)
			}

			if tt.trailing != trailing {
				t.Errorf("ExtractLocale(%q): trailing = %v, want %v",
					tt.urlPath, trailing, tt.trailing)
			}
		})
	}
}

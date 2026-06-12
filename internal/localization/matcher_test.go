package localization_test

import (
	"testing"

	"github.com/subtributary/musings/internal/localization"
)

func TestLocaleMatcher_Choose(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		locales []localization.Locale
		accept  string
		want    localization.Locale
	}{
		{
			name:    "no accept, no configured locale",
			locales: []localization.Locale{},
			want:    und,
		},
		{
			name:    "no accept, has configured locales",
			locales: []localization.Locale{ar, en},
			want:    ar,
		},
		{
			name:    "no match",
			locales: []localization.Locale{ar, en},
			accept:  "fr",
			want:    ar,
		},
		{
			name:    "tag with just language",
			locales: []localization.Locale{ar, en},
			accept:  "en",
			want:    en,
		},
		{
			name:    "tag with language and region",
			locales: []localization.Locale{ar, en, enGB},
			accept:  "en-GB",
			want:    enGB,
		},
		{
			name:    "tag for child of supported",
			locales: []localization.Locale{ar, en},
			accept:  "en-GB",
			want:    en,
		},
		{
			name:    "tag for parent of supported",
			locales: []localization.Locale{ar, enGB},
			accept:  "en",
			want:    enGB,
		},
		{
			name:    "locales specified with weights",
			locales: []localization.Locale{ar, en, enGB},
			accept:  "en;q=0.5, en-GB;q=0.9",
			want:    enGB,
		},
		{
			name:    "prefer zh-Hans over zh-Hant",
			locales: []localization.Locale{zhHant, zhHans},
			accept:  "zh",
			want:    zhHans,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			m := localization.NewLocaleMatcher(tt.locales)

			got := m.Choose(tt.accept)
			if got != tt.want {
				t.Errorf("Choose = %v, want %v", got, tt.want)
			}
		})
	}
}

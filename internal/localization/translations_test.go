package localization

import (
	"testing"

	"golang.org/x/text/language"
)

func TestLoadFor(t *testing.T) {
	t.Parallel()

	// Ensure the catalog is populated before any test uses LoadFor.
	// This must happen before the subtests run in parallel.
	InitTranslations()

	en := LoadFor(language.English)
	zhHans := LoadFor(language.SimplifiedChinese)
	zhHant := LoadFor(language.TraditionalChinese)

	tests := []struct {
		name string
		in   language.Tag
		want Strings
	}{
		{
			name: "en-GB falls back to en",
			in:   language.BritishEnglish,
			want: en,
		},
		{
			name: "zh falls back to zh-Hans",
			in:   language.Chinese,
			want: zhHans,
		},
		{
			name: "zh-Hant uses zh-Hant",
			in:   language.TraditionalChinese,
			want: zhHant,
		},
		{
			name: "und falls back to en",
			in:   language.Und,
			want: en,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := LoadFor(tt.in)
			if got != tt.want {
				t.Errorf("LoadFor(%q) = %+v, want %+v", tt.in, got, tt.want)
			}
		})
	}
}

package localization_test

import (
	"encoding/json"
	"testing"

	"github.com/subtributary/musings/internal/localization"
)

var (
	ar = localization.Locale{
		Tag:         "ar",
		NativeName:  "العربية",
		Direction:   "rtl",
		WritingMode: "horizontal-tb",
	}
	en = localization.Locale{
		Tag:         "en",
		NativeName:  "English",
		Direction:   "auto",
		WritingMode: "horizontal-tb",
	}
	enGB = localization.Locale{
		Tag:         "en-GB",
		NativeName:  "British English",
		Direction:   "auto",
		WritingMode: "horizontal-tb",
	}
	mnMong = localization.Locale{
		Tag:         "mn-Mong",
		NativeName:  "ᠮᠣᠩᠭᠣᠯ",
		Direction:   "ltr",
		WritingMode: "vertical-lr",
	}
	und    = localization.Locale{Tag: "und"}
	zhHans = localization.Locale{Tag: "zh-Hans"}
	zhHant = localization.Locale{Tag: "zh-Hant"}
)

func TestLocale_JSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		json    string
		want    localization.Locale
		wantErr bool
	}{
		{
			name:    "empty",
			wantErr: true,
		},
		{
			name:    "invalid locale",
			json:    `{"tag":"invalid"}`,
			wantErr: true,
		},
		{
			name: "undefined locale",
			json: `{"tag":"und","native_name":"und"}`,
			want: und,
		},
		{
			name: "tag only",
			json: `{"tag":"en"}`,
			want: en,
		},
		{
			name: "tag with region",
			json: `{"tag":"en-GB"}`,
			want: enGB,
		},
		{
			name: "tag region is lowercased",
			json: `{"tag":"en-gb"}`,
			want: enGB,
		},
		{
			name:    "tag only but native name not known",
			json:    `{"tag":"mn-Mong"}`,
			wantErr: true,
		},
		{
			name: "tag and direction only",
			json: `{"tag":"ar", "direction":"rtl"}`,
			want: ar,
		},
		{
			name:    "digit count not 10",
			json:    `{"tag:"en", digits: "abcdef"}`,
			wantErr: true,
		},
		{
			name: "locale fully defined",
			json: `{"tag":"mn-Mong",
				"date_format": "2006.01.02",
				"digits": "0123456789",
				"native_name":"ᠮᠣᠩᠭᠣᠯ",
				"direction":"ltr", 
				"writing_mode":"vertical-lr"}`,
			want: mnMong,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var locale localization.Locale
			err := json.Unmarshal([]byte(tt.json), &locale)

			if tt.wantErr && err == nil {
				t.Fatalf("want error, got none")
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("want no error, got: %v", err)
			}

			if locale.Tag != tt.want.Tag {
				t.Errorf("locale = %v, want %v", locale, tt.want)
			}
		})
	}
}

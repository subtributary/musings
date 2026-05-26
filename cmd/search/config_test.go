package main

import (
	"testing"

	"github.com/subtributary/musings/internal/config"
	"golang.org/x/text/language"
)

func TestLoadAppConfig(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		args    []string
		base    config.Global
		want    Config
		wantErr bool
	}{
		{
			name:    "empty",
			args:    []string{},
			base:    config.Global{Locales: []language.Tag{language.Und}},
			wantErr: true,
		},
		{
			name: "locale in global config",
			args: []string{"--locale", "en", "query"},
			base: config.Global{Locales: []language.Tag{language.English}},
			want: Config{Locale: language.English, Query: "query"},
		},
		{
			name:    "locale not in global config",
			args:    []string{"--locale", "en", "query"},
			base:    config.Global{Locales: []language.Tag{language.Und}},
			wantErr: true,
		},
		{
			name: "uppercase locale",
			args: []string{"--locale", "EN", "query"},
			base: config.Global{Locales: []language.Tag{language.English}},
			want: Config{Locale: language.English, Query: "query"},
		},
		{
			name: "query without locale",
			args: []string{"query"},
			base: config.Global{Locales: []language.Tag{language.English}},
			want: Config{Locale: language.English, Query: "query"},
		},
		{
			name:    "locale without query",
			args:    []string{"--locale", "en"},
			base:    config.Global{Locales: []language.Tag{language.English}},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			cfg := Config{Global: tt.base}
			err := cfg.loadAppConfig(tt.args)

			if tt.wantErr && err == nil {
				t.Fatalf("expected error, got none")
			} else if !tt.wantErr && err != nil {
				t.Fatalf("expected no error, got: %v", err)
			} else if tt.wantErr {
				return
			}

			if tt.want.Query != cfg.Query {
				t.Errorf("Query = %s, want %s", cfg.Query, tt.want.Query)
			}

			if tt.want.Locale != cfg.Locale {
				t.Errorf("Locale = %s, want %s", cfg.Locale, tt.want.Locale)
			}
		})
	}
}

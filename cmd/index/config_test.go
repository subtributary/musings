package main

import (
	"slices"
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
			name: "empty",
			args: []string{},
			base: config.Global{Locales: []language.Tag{language.Und}},
			want: Config{Locales: []language.Tag{language.Und}},
		},
		{
			name: "locale in global config",
			args: []string{"--locale", "en"},
			base: config.Global{Locales: []language.Tag{language.English}},
			want: Config{Locales: []language.Tag{language.English}},
		},
		{
			name:    "locale not in global config",
			args:    []string{"--locale", "en"},
			base:    config.Global{Locales: []language.Tag{language.Und}},
			wantErr: true,
		},
		{
			name: "target without locale",
			args: []string{"filename.md"},
			base: config.Global{Locales: []language.Tag{language.English}},
			want: Config{
				Locales: []language.Tag{language.English},
				Target:  "filename.md",
			},
		},
		{
			name: "target with locale",
			args: []string{"--locale", "ko", "filename.md"},
			base: config.Global{
				Locales: []language.Tag{language.English, language.Korean},
			},
			want: Config{
				Locales: []language.Tag{language.Korean},
				Target:  "filename.md",
			},
		},
		{
			name:    "target is not markdown",
			args:    []string{"filename.img"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := Config{Global: tt.base}
			err := cfg.loadAppConfig(tt.args)

			if tt.wantErr && err == nil {
				t.Fatalf("expected error, got none")
			} else if !tt.wantErr && err != nil {
				t.Fatalf("expected no error, got: %v", err)
			} else if tt.wantErr {
				return
			}

			if !slices.Equal(cfg.Locales, tt.want.Locales) {
				t.Errorf("Locales = %v, want %v", cfg.Locales, tt.want.Locales)
			}

			if cfg.Target != tt.want.Target {
				t.Errorf("Target = %v, want %v", cfg.Target, tt.want.Target)
			}
		})
	}
}

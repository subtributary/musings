package main

import (
	"testing"
)

func TestLoadAppConfig(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		args []string
		env  map[string]string
		want Config
	}{
		{
			name: "no configs",
			args: []string{},
			env:  map[string]string{},
			want: Config{
				BindAddress:   ":8080",
				LiveTemplates: false,
			},
		},
		{
			name: "environment variables",
			args: []string{},
			env: map[string]string{
				"MUSINGS_BIND": ":1",
			},
			want: Config{
				BindAddress:   ":1",
				LiveTemplates: false,
			},
		},
		{
			name: "command arguments",
			args: []string{"--bind", ":2", "--live-templates"},
			env:  map[string]string{},
			want: Config{
				BindAddress:   ":2",
				LiveTemplates: true,
			},
		},
		{
			name: "command arguments override environment variables",
			args: []string{"--bind", ":2", "--live-templates"},
			env: map[string]string{
				"MUSINGS_BIND": ":1",
			},
			want: Config{
				BindAddress:   ":2",
				LiveTemplates: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg, err := LoadAppConfig(tt.args, func(key string) string {
				return tt.env[key]
			})

			if err != nil {
				t.Fatalf("Expected no error, got: %v", err)
			}

			if cfg.BindAddress != tt.want.BindAddress {
				t.Errorf("BindAddress: want %v, got %v", tt.want.BindAddress, cfg.BindAddress)
			}

			if cfg.LiveTemplates != tt.want.LiveTemplates {
				t.Fatalf("LiveTemplates: want %v, got %v", tt.want.LiveTemplates, cfg.LiveTemplates)
			}
		})
	}
}

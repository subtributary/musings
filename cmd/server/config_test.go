package main

import (
	"bytes"
	"errors"
	"flag"
	"io"
	"testing"
)

func TestArgsParser(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		args []string
		env  map[string]string
		want ConfigArgs
	}{
		{
			name: "no configs",
			args: []string{},
			env:  map[string]string{},
			want: ConfigArgs{
				BindAddress: ":8080",
			},
		},
		{
			name: "environment variables",
			args: []string{},
			env: map[string]string{
				"MUSINGS_BIND": ":1",
			},
			want: ConfigArgs{
				BindAddress: ":1",
			},
		},
		{
			name: "command arguments",
			args: []string{"--bind", ":2"},
			env:  map[string]string{},
			want: ConfigArgs{
				BindAddress: ":2",
			},
		},
		{
			name: "command arguments override environment variables",
			args: []string{"--bind", ":2"},
			env: map[string]string{
				"MUSINGS_BIND": ":1",
			},
			want: ConfigArgs{
				BindAddress: ":2",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg, err := ArgsParser{
				Stdout: io.Discard,
				Stderr: io.Discard,
				Getenv: func(key string) string { return tt.env[key] },
			}.Parse(tt.args)
			if err != nil {
				t.Fatalf("expected no error, got: %v", err)
			}

			if cfg.BindAddress != tt.want.BindAddress {
				t.Errorf("BindAddress = %v, want %v", cfg.BindAddress, tt.want.BindAddress)
			}
		})
	}
}

func TestArgsParser_Output(t *testing.T) {
	t.Run("help outputs to stdout", func(t *testing.T) {
		t.Parallel()

		var stdout, stderr bytes.Buffer

		_, err := ArgsParser{
			Getenv: func(string) string { return "" },
			Stdout: &stdout,
			Stderr: &stderr,
		}.Parse([]string{"--help"})

		switch {
		case err == nil:
			t.Fatalf("expected flag.ErrHelp error, got no error")
		case !errors.Is(err, flag.ErrHelp):
			t.Fatalf("expected flag.ErrHelp error, got: %v", err)
		}

		if length := stdout.Len(); length == 0 {
			t.Errorf("stdout.Len() = %d, want > 0", length)
		}
		if length := stderr.Len(); length != 0 {
			t.Errorf("stderr.Len() = %d, want 0", length)
		}
	})

	t.Run("error outputs to stderr", func(t *testing.T) {
		t.Parallel()

		var stdout, stderr bytes.Buffer

		_, err := ArgsParser{
			Getenv: func(string) string { return "" },
			Stdout: &stdout,
			Stderr: &stderr,
		}.Parse([]string{"invalid"})
		if err == nil {
			t.Errorf("expected error, got none")
		}

		if length := stdout.Len(); length != 0 {
			t.Errorf("stdout.Len() = %d, want 0", length)
		}
		if length := stderr.Len(); length == 0 {
			t.Errorf("stderr.Len() = %d, want > 0", length)
		}
	})

	t.Run("success does not output", func(t *testing.T) {
		t.Parallel()

		var stdout, stderr bytes.Buffer

		_, err := ArgsParser{
			Getenv: func(string) string { return "" },
			Stdout: &stdout,
			Stderr: &stderr,
		}.Parse([]string{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if length := stdout.Len(); length != 0 {
			t.Errorf("stdout.Len() = %d, want 0", length)
		}
		if length := stderr.Len(); length != 0 {
			t.Errorf("stderr.Len() = %d, want 0", length)
		}
	})
}

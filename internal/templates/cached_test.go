package templates_test

import (
	"testing"
	"testing/fstest"
	"time"

	"github.com/subtributary/musings/internal/templates"
)

func TestCachedStore(t *testing.T) {
	older := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
	newer := time.Date(2026, 1, 2, 12, 0, 0, 0, time.UTC)
	files := fstest.MapFS{
		"index.gohtml": {
			Data:    []byte(`Hello, {{.Name}}! {{template "partials/footer" .}}`),
			ModTime: older,
		},
		"post.gohtml": {
			Data:    []byte(`Hello, {{.Name}}!`),
			ModTime: older,
		},
		"partials/footer.gohtml": {
			Data:    []byte(`Footer`),
			ModTime: newer,
		},
		"ignored.txt": {
			Data: []byte(`ignored`),
		},
	}

	tests := []struct {
		name    string
		target  string
		modTime time.Time
		wantErr bool
	}{
		{
			name:    "template with dependencies",
			target:  "index",
			modTime: newer,
		},
		{
			name:    "template with no dependencies",
			target:  "post",
			modTime: newer,
		},
		{
			name:    "non-template file",
			target:  "ignored",
			wantErr: true,
		},
		{
			name:    "non-template file with extension",
			target:  "ignored.txt",
			wantErr: true,
		},
		{
			name:    "missing file",
			target:  "missing",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := templates.NewCachedStore(templates.Funcs{})
			if err := store.LoadFS(files); err != nil {
				t.Fatalf("LoadFS() error = %v", err)
			}

			tmpl, err := store.Lookup(tt.target)
			if tt.wantErr && err == nil {
				t.Fatalf("Lookup(): want error, got nil")
			}
			if err != nil {
				return
			}

			modTime := tmpl.LastModified()
			if !modTime.Equal(tt.modTime) {
				t.Errorf("LastModified() modTime = %v, want %v", modTime, tt.modTime)
			}
		})
	}
}

func TestCachedStore_ParseError(t *testing.T) {
	files := fstest.MapFS{
		"bad.gohtml": {Data: []byte(`{{if}}`)},
	}

	store := templates.NewCachedStore(templates.Funcs{})
	if err := store.LoadFS(files); err == nil {
		t.Fatal("LoadFS() error = nil, want parse error")
	}
}

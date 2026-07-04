package templates_test

import (
	"testing"
	"testing/fstest"

	"github.com/subtributary/musings/internal/templates"
)

func TestCachedStore(t *testing.T) {
	files := fstest.MapFS{
		"index.gohtml": {
			Data: []byte(`Hello, {{.Name}}! {{template "partials/footer" .}}`),
		},
		"post.gohtml": {
			Data: []byte(`Hello, {{.Name}}!`),
		},
		"partials/footer.gohtml": {
			Data: []byte(`Footer`),
		},
		"ignored.txt": {
			Data: []byte(`ignored`),
		},
	}

	tests := []struct {
		name    string
		target  string
		wantErr bool
	}{
		{
			name:   "template with dependencies",
			target: "index",
		},
		{
			name:   "template with no dependencies",
			target: "post",
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

			_, err := store.Lookup(tt.target)
			if tt.wantErr && err == nil {
				t.Fatalf("Lookup(): want error, got nil")
			}
			if err != nil {
				return
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

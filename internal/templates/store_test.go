package templates_test

import (
	"testing"
	"testing/fstest"
	"time"

	"github.com/subtributary/musings/internal/templates"
)

func TestStore(t *testing.T) {
	contentFS := fstest.MapFS{}
	staticFS := fstest.MapFS{}
	templateFS := fstest.MapFS{
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
		modTime time.Time
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
			store, err := templates.NewStore(templateFS, contentFS, staticFS)
			if err != nil {
				t.Fatalf("NewStore() error = %v", err)
			}

			_, err = store.Lookup(tt.target)
			if tt.wantErr && err == nil {
				t.Fatalf("Lookup(): want error, got nil")
			}
		})
	}
}

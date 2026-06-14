package templates_test

import (
	"strings"
	"testing"
	"testing/fstest"
	"time"

	"github.com/subtributary/musings/internal/templates"
)

func TestFuncs_Versioned(t *testing.T) {
	t.Parallel()

	modTime := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
	funcs := templates.Funcs{
		StaticDir: fstest.MapFS{
			"js/script.js": {Data: []byte(""), ModTime: modTime},
		},
	}

	t.Run("without versioned the URL is unchanged", func(t *testing.T) {
		got, err := render(t, funcs, `<script src="/_static/js/script.js"></script>`)
		if err != nil {
			t.Fatalf("Execute() error = %v", err)
		}

		want := `<script src="/_static/js/script.js"></script>`
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("with versioned the version is appended", func(t *testing.T) {
		got, err := render(t, funcs, `<script src='{{"/_static/js/script.js" | versioned}}'></script>`)
		if err != nil {
			t.Fatalf("Execute() error = %v", err)
		}

		const timestamp = "20060102150405"
		want := `<script src='/_static/js/script.js?v=` + modTime.Format(timestamp) + `'></script>`
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("missing asset produces an error", func(t *testing.T) {
		_, err := render(t, funcs, `{{"/_static/js/missing.js" | versioned}}`)
		if err == nil {
			t.Fatal("Execute(): want error, got nil")
		}
	})
}

// render loads body as a template, executes it, and returns the output.
func render(t *testing.T, funcs templates.Funcs, body string) (string, error) {
	t.Helper()

	store := templates.NewCachedStore(funcs)
	err := store.LoadFS(fstest.MapFS{"page.gohtml": {Data: []byte(body)}})
	if err != nil {
		t.Fatalf("LoadFS() error = %v", err)
	}

	tmpl, err := store.Lookup("page")
	if err != nil {
		t.Fatalf("Lookup() error = %v", err)
	}

	var sb strings.Builder
	err = tmpl.Execute(&sb, nil)
	return sb.String(), err
}

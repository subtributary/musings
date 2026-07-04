package templates_test

import (
	"html/template"
	"strconv"
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
		ContentFS: fstest.MapFS{
			"post.md": {Data: []byte(""), ModTime: modTime},
		},
		StaticFS: fstest.MapFS{
			"js/script.js": {Data: []byte(""), ModTime: modTime},
		},
	}

	t.Run("without versioned the URL is unchanged", func(t *testing.T) {
		got, err := render(t, funcs, `{{"/_static/js/script.js"}}`)
		if err != nil {
			t.Fatalf("Execute() error = %v", err)
		}

		want := `/_static/js/script.js`
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("with versioned the version is appended", func(t *testing.T) {
		got, err := render(t, funcs, `{{"/_static/js/script.js" | versioned}}`)
		if err != nil {
			t.Fatalf("Execute() error = %v", err)
		}

		version := strconv.FormatInt(modTime.Unix(), 16)
		want := `/_static/js/script.js?v=` + version
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

	t.Run("look in content when no static prefix", func(t *testing.T) {
		got, err := render(t, funcs, `{{"/post.md" | versioned}}`)
		if err != nil {
			t.Fatalf("Execute() error = %v", err)
		}

		version := strconv.FormatInt(modTime.Unix(), 16)
		want := "/post.md?v=" + version
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})
}

// render loads body as a template, executes it, and returns the output.
func render(t *testing.T, funcs templates.Funcs, body string) (string, error) {
	t.Helper()

	tmpl := template.New("")
	funcs.ApplyTo(tmpl)

	if _, err := tmpl.Parse(body); err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	var sb strings.Builder
	err := tmpl.Execute(&sb, nil)
	return sb.String(), err
}

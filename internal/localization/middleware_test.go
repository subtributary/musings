package localization_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/subtributary/musings/internal/localization"
)

func TestLocalizeRoute(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		locales      []localization.Locale
		reqPath      string
		wantRedirect string
	}{
		/* no configured locale */
		{
			name:    "no configured locale, no path locale",
			locales: []localization.Locale{},
			reqPath: "/index.html",
		},
		{
			name:    "no configured locale, special path",
			locales: []localization.Locale{},
			reqPath: "/_shared/index.html",
		},

		/* und locale */
		{
			name:    "configured locale is und, no path locale",
			locales: []localization.Locale{},
			reqPath: "/index.html",
		},
		{
			name:    "configured locale is und, special path",
			locales: []localization.Locale{},
			reqPath: "/_shared/index.html",
		},

		/* one configured locale */
		{
			name:         "one configured locale, no path locale",
			locales:      []localization.Locale{en},
			reqPath:      "/index.html",
			wantRedirect: "/en/index.html",
		},
		{
			name:         "one configured locale, unsupported path locale",
			locales:      []localization.Locale{en},
			reqPath:      "/ko/index.html",
			wantRedirect: "/en/ko/index.html",
		},
		{
			name:    "one configured locale, path has locale",
			locales: []localization.Locale{en},
			reqPath: "/en/index.html",
		},
		{
			name:    "one configured locale, special path",
			locales: []localization.Locale{en},
			reqPath: "/_shared/index.html",
		},

		/* configured locale is regional */
		{
			name:         "regional configured locale, no path locale",
			locales:      []localization.Locale{zhHans},
			reqPath:      "/index.html",
			wantRedirect: "/zh-Hans/index.html",
		},
		{
			name:    "regional configured locale, path has locale",
			locales: []localization.Locale{zhHans},
			reqPath: "/zh-Hans/index.html",
		},
		{
			name:         "regional configured locale, path is lowercase",
			locales:      []localization.Locale{zhHans},
			reqPath:      "/zh-hans/index.html",
			wantRedirect: "/zh-Hans/index.html",
		},

		/* special cases */
		{
			name:         "redirect keeps query",
			reqPath:      "/index?q=search",
			locales:      []localization.Locale{en},
			wantRedirect: "/en/index?q=search",
		},
	}

	nopHandler := func(w http.ResponseWriter, r *http.Request) {}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r := chi.NewRouter()
			r.Use(localization.LocalizedRoute(tt.locales))
			r.Get("/_shared/index.html", nopHandler)
			r.Get("/index.html", nopHandler)

			req := httptest.NewRequest("GET", tt.reqPath, nil)
			rec := httptest.NewRecorder()
			r.ServeHTTP(rec, req)

			if tt.wantRedirect == "" && rec.Code != http.StatusOK {
				t.Fatalf("Code: got %v, want %v", rec.Code, http.StatusOK)
			}

			if loc := rec.Header().Get("Location"); loc != tt.wantRedirect {
				t.Fatalf("Location: got %v, want %v", loc, tt.wantRedirect)
			}
		})
	}
}

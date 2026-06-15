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
		name    string
		reqPath string
		locales []localization.Locale
		newPath string
	}{
		{
			name:    "no path locale, no configured locale",
			reqPath: "/index.html",
			locales: []localization.Locale{},
		},
		{
			name:    "no path locale, one configured locale",
			reqPath: "/index.html",
			locales: []localization.Locale{en},
			newPath: "/en/index.html",
		},
		{
			name:    "no path locale, configured locale has region",
			reqPath: "/index.html",
			locales: []localization.Locale{zhHans},
			newPath: "/zh-Hans/index.html",
		},
		{
			name:    "unsupported path locale, on configured locale",
			reqPath: "/ko/index.html",
			locales: []localization.Locale{en},
			newPath: "/en/ko/index.html",
		},
		{
			name:    "path locale",
			reqPath: "/en/index.html",
			locales: []localization.Locale{en},
		},
		{
			name:    "path locale with region",
			reqPath: "/zh-Hans/index.html",
			locales: []localization.Locale{zhHans},
		},
		{
			name:    "path locale with lowercase region",
			reqPath: "/zh-hans/index.html",
			locales: []localization.Locale{zhHans},
			newPath: "/zh-Hans/index.html",
		},
		{
			name:    "redirect keeps query",
			reqPath: "/index?q=search",
			locales: []localization.Locale{en},
			newPath: "/en/index?q=search",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r := chi.NewRouter()
			r.Use(localization.LocalizedRoute(tt.locales))
			r.Get("/index.html", func(w http.ResponseWriter, r *http.Request) {})

			req := httptest.NewRequest("GET", tt.reqPath, nil)
			rec := httptest.NewRecorder()
			r.ServeHTTP(rec, req)

			if tt.newPath == "" && rec.Code != http.StatusOK {
				t.Fatalf("Code: got %v, want %v", rec.Code, http.StatusOK)
			}

			if loc := rec.Header().Get("Location"); loc != tt.newPath {
				t.Fatalf("Location: got %v, want %v", loc, tt.newPath)
			}
		})
	}
}

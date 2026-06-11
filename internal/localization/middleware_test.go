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
			name:    "redirect if no path locale",
			reqPath: "/index.html",
			locales: []localization.Locale{en},
			newPath: "/en/index.html",
		},
		{
			name:    "redirect if unsupported path locale",
			reqPath: "/ko/index.html",
			locales: []localization.Locale{en},
			newPath: "/en/ko/index.html",
		},
		{
			name:    "no redirect if supported path locale",
			reqPath: "/en/index.html",
			locales: []localization.Locale{en},
		},
		{
			name:    "no redirect if no locales",
			reqPath: "/index.html",
			locales: []localization.Locale{},
		},
		{
			name:    "redirect keeps query",
			reqPath: "/index?q=search",
			locales: []localization.Locale{en},
			newPath: "/en/index?q=search",
		},
		{
			name:    "lowercase locale when redirected",
			reqPath: "/index.html",
			locales: []localization.Locale{zhHans},
			newPath: "/zh-hans/index.html",
		},
		{
			// This test fails. I will leave it for now.
			// I might want to lowercase in other middleware.
			name:    "redirect to lowercase locale if uppercase",
			reqPath: "/zh-Hans/index.html",
			locales: []localization.Locale{zhHans},
			newPath: "/zh-hans/index.html",
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

func TestLocalizedRouteAcceptLanguage(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		acceptLanguage string
		supported      []localization.Locale
		expected       localization.Locale
	}{
		{
			name:           "invalid locale defaults to first",
			acceptLanguage: "xx",
			supported:      []localization.Locale{en, ar},
			expected:       en,
		},
		{
			name:           "match parent locale when user accepts only child locale",
			acceptLanguage: "zh-Hans",
			supported:      []localization.Locale{en, zh},
			expected:       zh,
		},
		{
			name:           "match child locale when user accepts only parent locale",
			acceptLanguage: "zh",
			supported:      []localization.Locale{en, zhHans},
			expected:       zhHans,
		},
		{
			name:           "prefer zh-Hans over zh-Hant when both supported",
			acceptLanguage: "zh",
			supported:      []localization.Locale{en, zhHans, zhHant},
			expected:       zhHans,
		},
		{
			name:           "prefer exact zh over zh-Hans when user accepts only zh",
			acceptLanguage: "zh",
			supported:      []localization.Locale{zh, zhHans, zhHant},
			expected:       zh,
		},
	}

	createRouter := func(locales []localization.Locale) *chi.Mux {
		r := chi.NewRouter()
		r.Use(localization.LocalizedRoute(locales))
		r.Get("/index.html", func(w http.ResponseWriter, r *http.Request) {})
		return r
	}

	makeRequest := func(acceptLanguage string) *http.Request {
		req := httptest.NewRequest("GET", "/index.html", nil)
		req.Header.Set("Accept-Language", acceptLanguage)
		return req
	}

	allLocs := []localization.Locale{ar, en, zh, zhHans, zhHant}
	middleware := localization.NewLocalizedRouteMiddleware(allLocs)
	extractLocale := middleware.ExtractLocale

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r := createRouter(tt.supported)
			req := makeRequest(tt.acceptLanguage)
			rec := httptest.NewRecorder()
			r.ServeHTTP(rec, req)

			if rec.Code < 300 || rec.Code >= 400 {
				t.Fatalf("got status %d, expected redirect", rec.Code)
			}

			loc := rec.Header().Get("Location")
			if loc == "" {
				t.Fatalf("empty Location header")
			}

			locale, _ := extractLocale(loc)
			if locale != tt.expected {
				t.Errorf("locale: got %v, want %v", locale, tt.expected)
			}
		})
	}
}

func TestExtractLocale(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		urlPath   string
		localeTag string
		trailing  string
	}{
		{
			name:      "empty path",
			urlPath:   "",
			localeTag: "und",
			trailing:  "/",
		},
		{
			name:      "root path",
			urlPath:   "/",
			localeTag: "und",
			trailing:  "/",
		},
		{
			name:      "localized root path",
			urlPath:   "/en/",
			localeTag: "en",
			trailing:  "/",
		},
		{
			name:      "localized root path without trailing slash",
			urlPath:   "/en",
			localeTag: "en",
			trailing:  "/",
		},
		{
			name:      "localized root path without leading slash",
			urlPath:   "en/",
			localeTag: "en",
			trailing:  "/",
		},
		{
			name:      "localized root path without slashes",
			urlPath:   "en",
			localeTag: "en",
			trailing:  "/",
		},
		{
			name:      "localized file path",
			urlPath:   "/en/index.html",
			localeTag: "en",
			trailing:  "/index.html",
		},
		{
			name:      "localized nested file path",
			urlPath:   "/en/sub/index.html",
			localeTag: "en",
			trailing:  "/sub/index.html",
		},
		{
			name:      "invalid locale",
			urlPath:   "/xx/index.html",
			localeTag: "und",
			trailing:  "/xx/index.html",
		},
		{
			name:      "locale with region",
			urlPath:   "/zh-hans/index.html",
			localeTag: "zh-Hans",
			trailing:  "/index.html",
		},
	}

	locales := []localization.Locale{en, zhHans}
	m := localization.NewLocalizedRouteMiddleware(locales)

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			locale, trailing := m.ExtractLocale(tt.urlPath)

			if tt.localeTag != locale.Tag {
				t.Errorf("ExtractLocale(%q): locale = %v, want %v",
					tt.urlPath, locale.Tag, tt.localeTag)
			}

			if tt.trailing != trailing {
				t.Errorf("ExtractLocale(%q): trailing = %v, want %v",
					tt.urlPath, trailing, tt.trailing)
			}
		})
	}
}

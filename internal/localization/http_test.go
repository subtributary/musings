package localization

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"golang.org/x/text/language"
)

func TestLocalizeRoute(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		reqPath       string
		supportedTags []language.Tag
		newPath       string
	}{
		{
			name:          "redirect if no path locale",
			reqPath:       "/index.html",
			supportedTags: []language.Tag{language.English},
			newPath:       "/en/index.html",
		},
		{
			name:          "redirect if unsupported path locale",
			reqPath:       "/ko/index.html",
			supportedTags: []language.Tag{language.English},
			newPath:       "/en/ko/index.html",
		},
		{
			name:          "no redirect if supported path locale",
			reqPath:       "/en/index.html",
			supportedTags: []language.Tag{language.English},
		},
		{
			name:          "no redirect if undefined locale is supported",
			reqPath:       "/index.html",
			supportedTags: []language.Tag{language.Und},
		},
		{
			name:          "redirect keeps query",
			reqPath:       "/index?q=search",
			supportedTags: []language.Tag{language.English},
			newPath:       "/en/index?q=search",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r := chi.NewRouter()
			r.Use(LocalizedRoute(tt.supportedTags))
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
		name             string
		acceptLanguage   string
		supportedLocales []language.Tag
		expectedLocale   language.Tag
	}{
		{
			name:             "invalid locale defaults to first",
			acceptLanguage:   "xx",
			supportedLocales: []language.Tag{language.English, language.Korean},
			expectedLocale:   language.English,
		},
		{
			name:             "match parent locale when user accepts child locale",
			acceptLanguage:   "zh-Hans",
			supportedLocales: []language.Tag{language.Chinese},
			expectedLocale:   language.Chinese,
		},
		{
			name:             "match child locale when user accepts parent locale",
			acceptLanguage:   "zh",
			supportedLocales: []language.Tag{language.TraditionalChinese},
			expectedLocale:   language.TraditionalChinese,
		},
		{
			name:             "prefer zh-Hans over zh-Hant when both supported",
			acceptLanguage:   "zh",
			supportedLocales: []language.Tag{language.SimplifiedChinese, language.TraditionalChinese},
			expectedLocale:   language.SimplifiedChinese,
		},
		{
			name:             "prefer exact zh over zh-Hans when zh is supported",
			acceptLanguage:   "zh",
			supportedLocales: []language.Tag{language.Chinese, language.SimplifiedChinese},
			expectedLocale:   language.Chinese,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r := chi.NewRouter()
			r.Use(LocalizedRoute(tt.supportedLocales))
			r.Get("/index.html", func(w http.ResponseWriter, r *http.Request) {})

			req := httptest.NewRequest("GET", "/index.html", nil)
			req.Header.Set("Accept-Language", tt.acceptLanguage)
			rec := httptest.NewRecorder()
			r.ServeHTTP(rec, req)

			if rec.Code < 300 || rec.Code >= 400 {
				t.Fatalf("got status %d, expected redirect", rec.Code)
			}

			loc := rec.Header().Get("Location")
			if loc == "" {
				t.Fatalf("empty Location header")
			}

			tag, _ := ParsePath(loc)
			if tag != tt.expectedLocale {
				t.Errorf("locale: got %v, want %v", tag, tt.expectedLocale)
			}
		})
	}
}

func TestParsePath(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		urlPath  string
		tag      language.Tag
		trailing string
	}{
		{
			name:     "empty path",
			urlPath:  "",
			tag:      language.Und,
			trailing: "/",
		},
		{
			name:     "root path",
			urlPath:  "/",
			tag:      language.Und,
			trailing: "/",
		},
		{
			name:     "localized root path",
			urlPath:  "/en/",
			tag:      language.English,
			trailing: "/",
		},
		{
			name:     "localized root path without trailing slash",
			urlPath:  "/en",
			tag:      language.English,
			trailing: "/",
		},
		{
			name:     "localized root path without leading slash",
			urlPath:  "en/",
			tag:      language.English,
			trailing: "/",
		},
		{
			name:     "localized root path without slashes",
			urlPath:  "en",
			tag:      language.English,
			trailing: "/",
		},
		{
			name:     "localized file path",
			urlPath:  "/en/index.html",
			tag:      language.English,
			trailing: "/index.html",
		},
		{
			name:     "localized nested file path",
			urlPath:  "/en/sub/index.html",
			tag:      language.English,
			trailing: "/sub/index.html",
		},
		{
			name:     "invalid locale",
			urlPath:  "/xx/index.html",
			tag:      language.Und,
			trailing: "/xx/index.html",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tag, trailing := ParsePath(tt.urlPath)

			if tt.tag != tag {
				t.Errorf("ParsePath(%q): tag = %v, want %v", tt.urlPath, tag, tt.tag)
			}

			if tt.trailing != trailing {
				t.Errorf("ParsePath(%q): trailing = %v, want %v", tt.urlPath, trailing, tt.trailing)
			}
		})
	}
}

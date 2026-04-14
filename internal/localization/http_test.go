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

			tag, _, err := ParsePath(loc)
			if err != nil {
				t.Fatalf("parse location %q: %v", loc, err)
			}
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
		hasError bool
	}{
		{
			name:     "localized file path is processed",
			urlPath:  "/en/index.html",
			tag:      language.English,
			trailing: "/index.html",
		},
		{
			name:     "localized root path is processed",
			urlPath:  "/en/",
			tag:      language.English,
			trailing: "/",
		},
		{
			name:     "localized root path without trailing slash is processed",
			urlPath:  "/en",
			tag:      language.English,
			trailing: "/",
		},
		{
			name:     "invalid locale is an error",
			urlPath:  "/xx/index.html",
			hasError: true,
		},
		{
			name:     "root path is an error",
			urlPath:  "/",
			hasError: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tag, trailing, err := ParsePath(tt.urlPath)

			if tt.hasError && err == nil {
				t.Fatalf("ParsePath(%s): expected error but got none", tt.urlPath)
			} else if !tt.hasError && err != nil {
				t.Fatalf("ParsePath(%s): expected no error but got one: %v", tt.urlPath, err)
			}

			if tt.tag != tag {
				t.Errorf("ParsePath(%s): expected tag %v but got %v", tt.urlPath, tt.tag, tag)
			}

			if tt.trailing != trailing {
				t.Errorf("ParsePath(%s): expected trailing %v but got %v", tt.urlPath, tt.tag, tag)
			}
		})
	}
}

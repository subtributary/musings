package main

import (
	"fmt"
	"io/fs"
	"net/http/httptest"
	"testing"
	"testing/fstest"
	"time"

	"github.com/subtributary/musings/internal/localization"
	"github.com/subtributary/musings/internal/templates"
)

func TestViewOptions(t *testing.T) {
	t.Parallel()

	create := func(name string, opts ...ViewOption) *ViewOptions {
		options := &ViewOptions{}
		for _, opt := range opts {
			opt(options)
		}
		options.viewName = name
		return options
	}

	t.Run("fully built", func(t *testing.T) {
		if err := create("view",
			WithData("data"),
			WithDataModified(time.Now()),
			WithLocale(localization.UndLocale),
			WithPath("path"),
		).Validate(); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	t.Run("missing optional params", func(t *testing.T) {
		if err := create("view",
			WithLocale(localization.UndLocale),
			WithPath("path"),
		).Validate(); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	t.Run("missing view name", func(t *testing.T) {
		if err := create("",
			WithLocale(localization.UndLocale),
			WithPath("path"),
		).Validate(); err == nil {
			t.Fatalf("expected error, got none")
		}
	})

	t.Run("missing locale", func(t *testing.T) {
		if err := create("view",
			WithPath("path"),
		).Validate(); err == nil {
			t.Fatalf("expected error, got none")
		}
	})

	t.Run("missing path", func(t *testing.T) {
		if err := create("view",
			WithLocale(localization.UndLocale),
		).Validate(); err == nil {
			t.Fatalf("expected error, got none")
		}
	})
}

func TestViewFactory(t *testing.T) {
	t.Parallel()

	older := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
	templateTime := time.Date(2026, 1, 2, 12, 0, 0, 0, time.UTC)
	newer := time.Date(2026, 1, 3, 12, 0, 0, 0, time.UTC)

	en := localization.Locale{Tag: "en"}
	fr := localization.Locale{Tag: "fr"}

	tests := []struct {
		name         string
		dataModified time.Time
		locales      []localization.Locale
		locale       localization.Locale
		template     string
		wantHtml     string
		wantModified time.Time
	}{
		{
			name:     "no locales",
			locale:   localization.UndLocale,
			template: "{{.Locale.Tag}} {{range .LocaleOptions}}{{.Tag}} {{end}}",
			wantHtml: "und ",
		},
		{
			name:     "locales",
			locales:  []localization.Locale{en, fr},
			locale:   en,
			template: "{{.Locale.Tag}} {{range .LocaleOptions}}{{.Tag}} {{end}}",
			wantHtml: "en en fr ",
		},
		{
			name:         "old data",
			dataModified: older,
			wantModified: templateTime,
			locale:       localization.UndLocale, // required, not tested here
		},
		{
			name:         "new data",
			dataModified: newer,
			wantModified: newer,
			locale:       localization.UndLocale, // required, not tested here
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tfs := fstest.MapFS{"test.gohtml": {
				Data:    []byte(tt.template),
				ModTime: templateTime,
			}}
			factory, err := newViewFactory(t, tfs, tt.locales)
			if err != nil {
				t.Fatalf("Error creating test view factory: %v", err)
			}

			rec := httptest.NewRecorder()
			err = factory.CreateAndServe(rec, "test",
				WithLocale(tt.locale),
				WithPath("/test"),
			)
			if err != nil {
				t.Fatalf("Error serving view: %v", err)
			}

			got := rec.Body.String()
			if got != tt.wantHtml {
				t.Fatalf("HTML = %v, want %v", got, tt.wantHtml)
			}
		})
	}
}

func newViewFactory(t *testing.T, templatesFS fs.FS, locales []localization.Locale) (*ViewFactory, error) {
	t.Helper()

	templateStore := templates.NewCachedStore()
	if err := templateStore.LoadFS(templatesFS); err != nil {
		return nil, fmt.Errorf("load templates: %w", err)
	}

	translations := localization.Store{}
	for _, locale := range locales {
		translations[locale.Tag] = localization.Translations{"title": "hello"}
	}

	return &ViewFactory{
		locales:       locales,
		templateStore: templateStore,
		translations:  translations,
	}, nil
}

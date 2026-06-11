package main

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/subtributary/musings/internal/config"
	"github.com/subtributary/musings/internal/localization"
	"github.com/subtributary/musings/internal/templates"
)

type ViewOption func(*View) error

func WithData(data any) ViewOption {
	return func(view *View) error {
		view.model.Data = data
		return nil
	}
}

// WithDataModified uses its argument for the `Last-Modified` header if it is
// more recent than the template modification date; otherwise it uses the
// template modification date.
func WithDataModified(when time.Time) ViewOption {
	return func(view *View) error {
		if when.After(view.modified) {
			view.modified = when
		}
		return nil
	}
}

type ViewFactory struct {
	locales       []localization.Locale
	templateStore templates.Store
	translations  localization.Store
}

func LoadViewFactory(cfg Config) (*ViewFactory, error) {
	f := &ViewFactory{
		locales: cfg.Locales,
	}

	templateStore, err := templates.NewStore(config.TemplatesPath, cfg.LiveTemplates)
	if err != nil {
		return nil, fmt.Errorf("load templates: %v", err)
	}
	f.templateStore = templateStore

	translations, err := localization.LoadStore(config.DataPath)
	if err != nil {
		return nil, fmt.Errorf("load translations: %v", err)
	}
	f.translations = translations

	return f, nil
}

func (f *ViewFactory) Create(r *http.Request, name string, opts ...ViewOption) (View, error) {
	v := View{}

	locale := localization.LocaleFromContext(r.Context())
	path := currentRoutePath(r)

	f.setLocale(&v, locale, path)
	f.setTranslations(&v, locale)
	if err := f.setTemplate(&v, name); err != nil {
		return v, err
	}
	f.setRoot(&v, locale)

	v.modified = v.tmpl.LastModified()

	for _, opt := range opts {
		if err := opt(&v); err != nil {
			return v, err
		}
	}

	return v, nil
}

func (f *ViewFactory) Serve(w http.ResponseWriter, r *http.Request, name string, opts ...ViewOption) error {
	v, err := f.Create(r, name, opts...)
	if err != nil {
		return fmt.Errorf("create view: %v", err)
	}

	return v.Serve(w)
}

func (f *ViewFactory) setLocale(v *View, current localization.Locale, path string) {
	// These defaults work when localization is disabled.
	v.model.LocaleOptions = make([]LocaleOption, len(f.locales))
	v.model.Locale = LocaleOption{Locale: current, IsCurrent: true, URL: "/"}

	// If localization is enabled, set up locales and find current one.
	for i, locale := range f.locales {
		tag := strings.ToLower(locale.Tag)
		localizedPath, _ := url.JoinPath("/", tag, path)
		option := LocaleOption{
			Locale:    locale,
			IsCurrent: locale == current,
			URL:       localizedPath,
		}

		v.model.LocaleOptions[i] = option
		if option.IsCurrent {
			v.model.Locale = option
		}
	}
}

func (f *ViewFactory) setRoot(v *View, locale localization.Locale) {
	if locale == localization.UndLocale {
		v.model.RootURL = "/"
	}
	v.model.RootURL = "/" + locale.Tag + "/"
}

func (f *ViewFactory) setTemplate(v *View, name string) (err error) {
	v.tmpl, err = f.templateStore.Lookup(name)
	return
}

func (f *ViewFactory) setTranslations(v *View, locale localization.Locale) {
	v.model.Translations = f.translations.For(locale)
}

type LocaleOption struct {
	localization.Locale
	IsCurrent bool
	URL       string
}

type ViewModel struct {
	LocaleOptions []LocaleOption
	Locale        LocaleOption
	Translations  localization.Translations
	RootURL       string
	Data          any
}

type View struct {
	modified time.Time
	model    ViewModel
	tmpl     templates.Template
}

func (v View) Serve(w http.ResponseWriter) error {
	// Write to a buffer so that errors do not leave it partially written.
	var buf bytes.Buffer
	if err := v.tmpl.Execute(&buf, v.model); err != nil {
		return fmt.Errorf("execute template: %w", err)
	}

	// The `Last-Modified` header only has seconds precision.
	modified := v.modified.UTC().Truncate(time.Second)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Last-Modified", modified.Format(http.TimeFormat))
	if _, err := buf.WriteTo(w); err != nil {
		return fmt.Errorf("write response: %w", err)
	}
	return nil
}

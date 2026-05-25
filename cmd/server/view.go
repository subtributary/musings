package main

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/subtributary/musings/internal/config"
	"github.com/subtributary/musings/internal/localization"
	"github.com/subtributary/musings/internal/templates"
	"golang.org/x/text/language"
	"golang.org/x/text/language/display"
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
	locales []language.Tag
	store   templates.Store
}

func LoadViewFactory(cfg Config) (*ViewFactory, error) {
	f := &ViewFactory{
		locales: cfg.Locales,
	}

	store, err := templates.NewStore(config.TemplatesPath, cfg.LiveTemplates)
	if err != nil {
		return nil, fmt.Errorf("load templates: %v", err)
	}
	f.store = store

	return f, nil
}

func (f *ViewFactory) Create(r *http.Request, name string, opts ...ViewOption) (View, error) {
	v := View{}

	locale := localization.LocaleFromContext(r.Context())
	path := currentRoutePath(r)

	f.setLanguage(&v, locale, path)
	f.setTranslations(&v, locale)
	if err := f.setTemplate(&v, name); err != nil {
		return v, err
	}

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

func (f *ViewFactory) setLanguage(v *View, current language.Tag, path string) {
	v.model.LanguageOptions = make([]LanguageOption, 0, len(f.locales))
	v.model.Language = LanguageOption{IsCurrent: true, URL: "/"}

	for _, tag := range f.locales {
		code := tag.String()
		localizedPath, _ := url.JoinPath("/", code, path)
		option := LanguageOption{
			Code:      code,
			Label:     display.Self.Name(tag),
			IsCurrent: tag == current,
			URL:       localizedPath,
		}

		v.model.LanguageOptions = append(v.model.LanguageOptions, option)
		if option.IsCurrent {
			v.model.Language = option
		}
	}
}

func (f *ViewFactory) setTemplate(v *View, name string) (err error) {
	v.tmpl, err = f.store.Lookup(name)
	return
}

func (f *ViewFactory) setTranslations(v *View, language language.Tag) {
	v.model.Translations = localization.LoadFor(language)
}

type LanguageOption struct {
	Code      string
	Label     string
	IsCurrent bool
	URL       string
}

type ViewModel struct {
	LanguageOptions []LanguageOption
	Language        LanguageOption
	Translations    localization.Strings
	Data            any
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

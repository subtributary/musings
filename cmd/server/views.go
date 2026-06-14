package main

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/subtributary/musings/internal/localization"
	"github.com/subtributary/musings/internal/templates"
)

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

func (v *View) Serve(w http.ResponseWriter) error {
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

type ViewOption func(*ViewOptions)

func WithData(data any) ViewOption {
	return func(opts *ViewOptions) {
		opts.data = data
	}
}

func WithDataModified(when time.Time) ViewOption {
	return func(opts *ViewOptions) { opts.dataModified = when }
}

func WithLocale(locale localization.Locale) ViewOption {
	return func(opts *ViewOptions) { opts.locale = locale }
}

func WithPath(path string) ViewOption {
	return func(opts *ViewOptions) { opts.path = path }
}

type ViewOptions struct {
	data         any
	dataModified time.Time
	locale       localization.Locale
	path         string
	viewName     string
}

func (opts *ViewOptions) Validate() error {
	if opts.locale.Tag == "" {
		return errors.New("locale is not set")
	}
	if opts.path == "" {
		return errors.New("path is not set")
	}
	if opts.viewName == "" {
		return errors.New("view name is not set")
	}
	return nil
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

	templateStore, err := templates.NewStore(TemplatesPath, cfg.LiveTemplates)
	if err != nil {
		return nil, fmt.Errorf("load templates: %w", err)
	}
	f.templateStore = templateStore

	translations, err := localization.NewStore(DataPath, cfg.LiveTemplates)
	if err != nil {
		return nil, fmt.Errorf("load translations: %w", err)
	}
	f.translations = translations

	return f, nil
}

func (f *ViewFactory) Create(name string, opts ...ViewOption) (*View, error) {
	options := &ViewOptions{}
	for _, opt := range opts {
		opt(options)
	}
	options.viewName = name

	if err := options.Validate(); err != nil {
		return nil, fmt.Errorf("invalid options: %w", err)
	}

	vm, err := f.createVM(options)
	if err != nil {
		return nil, err
	}
	v, err := f.createView(options, vm)
	if err != nil {
		return nil, err
	}

	return v, nil
}

func (f *ViewFactory) CreateAndServe(w http.ResponseWriter, name string, opts ...ViewOption) error {
	v, err := f.Create(name, opts...)
	if err != nil {
		return fmt.Errorf("create view: %w", err)
	}
	return v.Serve(w)
}

func (f *ViewFactory) createView(opts *ViewOptions, vm ViewModel) (*View, error) {
	tmpl, err := f.templateStore.Lookup(opts.viewName)
	if err != nil {
		return nil, fmt.Errorf("load template: %w", err)
	}

	modified := tmpl.LastModified()
	if !opts.dataModified.IsZero() && opts.dataModified.After(modified) {
		modified = opts.dataModified
	}

	return &View{
		model:    vm,
		modified: modified,
		tmpl:     tmpl,
	}, nil
}

func (f *ViewFactory) createVM(opts *ViewOptions) (ViewModel, error) {
	vm := ViewModel{Data: opts.data}

	translations, err := f.translations.For(opts.locale)
	if err != nil {
		return vm, fmt.Errorf("load translations: %w", err)
	}
	vm.Translations = translations

	// Defaults for no localization.
	vm.LocaleOptions = make([]LocaleOption, len(f.locales))
	vm.Locale = LocaleOption{Locale: opts.locale, IsCurrent: true, URL: "/"}
	vm.RootURL = "/"

	// Localized websites use the locale prefix as the root URL.
	if opts.locale != localization.UndLocale {
		vm.RootURL = "/" + opts.locale.Tag + "/"
	}

	// If locales are configured, set up locales and find current one.
	for i, locale := range f.locales {
		localizedPath, _ := url.JoinPath("/", locale.Tag, opts.path)
		option := LocaleOption{
			Locale:    locale,
			IsCurrent: locale == opts.locale,
			URL:       localizedPath,
		}

		vm.LocaleOptions[i] = option
		if option.IsCurrent {
			vm.Locale = option
		}
	}

	return vm, nil
}

package main

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/subtributary/musings/internal/localization"
	"github.com/subtributary/musings/internal/templates"
)

type LocaleOption struct {
	localization.Locale
	IsCurrent bool
	URL       string
}

type ViewModel struct {
	LocaleOptions []*LocaleOption
	Locale        *LocaleOption
	Translations  localization.Translations
	RootURL       string
	Data          any
}

type View struct {
	modified time.Time
	model    ViewModel
	tmpl     templates.Template
}

func (v *View) SetData(data any, dataModified time.Time) error {
	if v.model.Data != nil {
		// Setting data twice may leave dateModified in a bad state,
		// so we do not allow it.
		return errors.New("data has already been set")
	}

	v.model.Data = data

	if dataModified.After(v.modified) {
		v.modified = dataModified
	}

	return nil
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

type ViewFactory struct {
	locales       []localization.Locale
	templateStore templates.Store
	translations  localization.Store
}

func NewViewFactory(liveTemplates bool, locales []localization.Locale, staticRoot *os.Root) (*ViewFactory, error) {
	f := &ViewFactory{locales: locales}

	funcs := templates.Funcs{
		StaticDir: staticRoot.FS(),
	}

	templateStore, err := templates.NewStore(TemplatesPath, funcs, liveTemplates)
	if err != nil {
		return nil, fmt.Errorf("load templates: %w", err)
	}
	f.templateStore = templateStore

	translations, err := localization.NewStore(DataPath, liveTemplates)
	if err != nil {
		return nil, fmt.Errorf("load translations: %w", err)
	}
	f.translations = translations

	return f, nil
}

func (f *ViewFactory) Create(r *http.Request, name string) (*View, error) {
	locale := localization.LocaleFromContext(r.Context())
	reqPath := chi.RouteContext(r.Context()).RoutePath

	view, err := f.createView(name)
	if err != nil {
		return nil, fmt.Errorf("create view: %w", err)
	}

	vm, err := f.createVM(locale, reqPath)
	if err != nil {
		return nil, fmt.Errorf("create view model: %w", err)
	}

	view.model = vm
	return view, nil
}

func (f *ViewFactory) createView(name string) (*View, error) {
	tmpl, err := f.templateStore.Lookup(name)
	if err != nil {
		return nil, fmt.Errorf("load template: %w", err)
	}

	return &View{
		modified: tmpl.LastModified(),
		tmpl:     tmpl,
	}, nil
}

func (f *ViewFactory) createVM(reqLocale localization.Locale, relPath string) (vm ViewModel, err error) {
	vm.Translations, err = f.translations.For(reqLocale)
	if err != nil {
		return vm, fmt.Errorf("load translations: %w", err)
	}

	// Defaults for no localization.
	vm.LocaleOptions = make([]*LocaleOption, len(f.locales))
	vm.Locale = &LocaleOption{Locale: reqLocale, IsCurrent: true, URL: "/"}
	vm.RootURL = "/"

	// Localized websites use the locale prefix as the root URL.
	if reqLocale.Tag != localization.Und {
		vm.RootURL = "/" + reqLocale.Tag + "/"
	}

	// If locales are configured, set up locales and find current one.
	for i, cfgLocale := range f.locales {
		localizedPath, _ := url.JoinPath("/", cfgLocale.Tag, relPath)
		option := &LocaleOption{
			Locale:    cfgLocale,
			IsCurrent: cfgLocale.Tag == reqLocale.Tag,
			URL:       localizedPath,
		}

		vm.LocaleOptions[i] = option
		if option.IsCurrent {
			vm.Locale = option
		}
	}

	return vm, nil
}

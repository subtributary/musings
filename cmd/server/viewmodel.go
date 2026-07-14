package main

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/subtributary/musings/internal/app"
	"github.com/subtributary/musings/internal/localization"
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

func (vm ViewModel) FormatDate(when time.Time) string {
	return vm.Locale.FormatDate(when)
}

type ViewModelFactory struct {
	locales      []localization.Locale
	translations localization.Store
}

func NewModelViewFactory(locales []localization.Locale) (*ViewModelFactory, error) {
	translations, err := localization.NewStore(app.DataPath)
	if err != nil {
		return nil, fmt.Errorf("load translations: %w", err)
	}

	return &ViewModelFactory{
		locales:      locales,
		translations: translations,
	}, nil
}

func (f *ViewModelFactory) Create(r *http.Request, data any) ViewModel {
	reqLocale := localization.LocaleFromContext(r.Context())
	reqPath := chi.RouteContext(r.Context()).RoutePath

	vm := ViewModel{
		Data:         data,
		Translations: f.translations.For(reqLocale),
	}

	// Defaults for no localization.
	vm.LocaleOptions = make([]LocaleOption, len(f.locales))
	vm.Locale = LocaleOption{Locale: reqLocale, IsCurrent: true, URL: "/"}
	vm.RootURL = "/"

	// Localized websites use the locale prefix as the root URL.
	if reqLocale.Tag != localization.Und {
		vm.RootURL = "/" + reqLocale.Tag + "/"
	}

	// If locales are configured, set up locales and find current one.
	for i, cfgLocale := range f.locales {
		localizedPath, _ := url.JoinPath("/", cfgLocale.Tag, reqPath)
		option := LocaleOption{
			Locale:    cfgLocale,
			IsCurrent: cfgLocale.Tag == reqLocale.Tag,
			URL:       localizedPath,
		}

		vm.LocaleOptions[i] = option
		if option.IsCurrent {
			vm.Locale = option
		}
	}

	return vm
}

package main

import (
	"net/url"

	"github.com/subtributary/musings/internal/localization"
	"golang.org/x/text/language"
	"golang.org/x/text/language/display"
)

type LanguageOption struct {
	Code      string
	Label     string
	IsCurrent bool
	URL       string
}

func newLanguageOption(tag language.Tag, current language.Tag, path string) LanguageOption {
	code := tag.String()
	path, _ = url.JoinPath("/", code, path)
	return LanguageOption{
		Code:      code,
		Label:     display.Self.Name(tag),
		IsCurrent: tag == current,
		URL:       path,
	}
}

type ViewModel struct {
	LanguageOptions []LanguageOption
	Language        LanguageOption
	Translations    localization.Strings
	Data            any
}

type ViewModelParams struct {
	CurrentLocale    language.Tag
	SupportedLocales []language.Tag
	CurrentPath      string
	Data             any
}

func NewViewModel(params ViewModelParams) (vm ViewModel) {
	vm.LanguageOptions = make([]LanguageOption, 0, len(params.SupportedLocales))
	vm.Language = newLanguageOption(language.Und, language.Und, params.CurrentPath)
	vm.Translations = localization.LoadFor(params.CurrentLocale)
	vm.Data = params.Data

	for _, tag := range params.SupportedLocales {
		option := newLanguageOption(tag, params.CurrentLocale, params.CurrentPath)
		vm.LanguageOptions = append(vm.LanguageOptions, option)
		if option.IsCurrent {
			vm.Language = option
		}
	}

	return
}

package app

import (
	"net/url"

	"golang.org/x/text/language"
	"golang.org/x/text/language/display"
)

type LanguageOption struct {
	Code      string
	Label     string
	IsCurrent bool
	URL       string
}

func newLanguageOption(tag language.Tag, isCurrent bool, path string) LanguageOption {
	code := tag.String()
	path, _ = url.JoinPath("/", code, path)
	return LanguageOption{
		Code:      code,
		Label:     display.Self.Name(tag),
		IsCurrent: isCurrent,
		URL:       path,
	}
}

type ViewModel struct {
	LanguageOptions []LanguageOption
	Language        LanguageOption
	Data            any
}

type ViewModelParams struct {
	CurrentLocale    language.Tag
	SupportedLocales []language.Tag
	CurrentPath      string
	Data             any
}

func NewViewModel(params ViewModelParams) (vm ViewModel) {
	vm.Data = params.Data

	vm.LanguageOptions = make([]LanguageOption, 0, len(params.SupportedLocales))
	for _, tag := range params.SupportedLocales {
		if tag == params.CurrentLocale {
			vm.Language = newLanguageOption(tag, true, params.CurrentPath)
			vm.LanguageOptions = append(vm.LanguageOptions, vm.Language)
		} else {
			vm.Language = newLanguageOption(tag, false, params.CurrentPath)
		}
	}

	return
}

package localization

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"golang.org/x/text/language"
	"golang.org/x/text/language/display"
)

const Und = "und"

var UndLocale = Locale{Tag: Und}

type Locale struct {
	Tag         string `json:"tag"`
	dateFormat  string
	digits      []rune
	Direction   string `json:"direction"`
	NativeName  string `json:"native_name"`
	WritingMode string `json:"writing_mode"`
}

func (loc *Locale) FormatDate(when time.Time) string {
	formatted := when.Format(loc.dateFormat)

	// Replace digits in date string with locale-specific digits.
	var replaced strings.Builder
	for _, r := range formatted {
		if r >= '0' && r <= '9' {
			replaced.WriteRune(loc.digits[r-'0'])
		} else {
			replaced.WriteRune(r)
		}
	}

	return replaced.String()
}

func (loc *Locale) UnmarshalJSON(data []byte) error {
	type Alias Locale
	cfg := &struct {
		Alias
		DateFormat string `json:"date_format"`
		Digits     string `json:"digits"`
	}{}

	if err := json.Unmarshal(data, &cfg); err != nil {
		return err
	}

	if cfg.DateFormat == "" {
		cfg.dateFormat = "2006-01-02"
	} else {
		cfg.dateFormat = cfg.DateFormat
	}

	if cfg.Digits == "" {
		cfg.digits = []rune("0123456789")
	} else {
		cfg.digits = []rune(cfg.Digits)
		if len(cfg.digits) != 10 {
			return fmt.Errorf("digit count is not 10: %s", cfg.Digits)
		}
	}

	tag, err := language.Parse(cfg.Tag)
	if err != nil {
		return fmt.Errorf("parse tag %s: %w", cfg.Tag, err)
	}
	if normTag := tag.String(); cfg.Tag != normTag {
		log.Printf("Note: Locale %s will be treated as %s.", cfg.Tag, normTag)
		cfg.Tag = normTag
	}

	switch cfg.Direction {
	case "ltr", "rtl", "auto":
	case "":
		cfg.Direction = "auto"
	default:
		return fmt.Errorf("invalid direction: %s", cfg.Direction)
	}

	if cfg.NativeName == "" {
		cfg.NativeName = display.Self.Name(tag)
		if cfg.NativeName == "" {
			return fmt.Errorf("native_name required for locale %s", cfg.Tag)
		}
	}

	switch cfg.WritingMode {
	case "horizontal-tb", "vertical-rl", "vertical-lr":
	case "":
		cfg.WritingMode = "horizontal-tb"
	default:
		return fmt.Errorf("invalid writing mode: %s", cfg.WritingMode)
	}

	*loc = Locale(cfg.Alias)
	return nil
}

func normalizeTag(input string) (string, error) {
	tag, err := language.Parse(input)
	if err != nil {
		return "", fmt.Errorf("parse tag %s: %w", input, err)
	}
	return tag.String(), nil
}

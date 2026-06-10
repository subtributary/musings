package localization

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"

	"golang.org/x/text/language"
	"golang.org/x/text/language/display"
)

var UndLocale = Locale{Tag: "und"}

type Locale struct {
	Tag         string `json:"tag"`
	NativeName  string `json:"native_name"`
	Direction   string `json:"direction"`
	WritingMode string `json:"writing_mode"`
}

func (loc *Locale) UnmarshalJSON(data []byte) error {
	type Alias Locale

	var state Alias
	if err := json.Unmarshal(data, &state); err != nil {
		return err
	}

	tag, err := language.Parse(state.Tag)
	if err != nil {
		return fmt.Errorf("parse tag %q: %w", state.Tag, err)
	}
	if tag == language.Und {
		return errors.New("undefined tag is invalid")
	}

	normalizedTag := tag.String()
	if !strings.EqualFold(state.Tag, normalizedTag) {
		log.Printf("Warning: Configured locale %q will be treated as %q.", state.Tag, tag.String())
	}
	state.Tag = normalizedTag

	if state.NativeName == "" {
		state.NativeName = display.Self.Name(tag)
		if state.NativeName == "" {
			return fmt.Errorf("native_name required for locale %q", state.Tag)
		}
	}

	switch state.Direction {
	case "ltr", "rtl", "auto":
	case "":
		state.Direction = "auto"
	default:
		return fmt.Errorf("invalid direction %q", state.Direction)
	}

	switch state.WritingMode {
	case "horizontal-tb", "vertical-rl", "vertical-lr":
	case "":
		state.WritingMode = "horizontal-tb"
	default:
		return fmt.Errorf("invalid writing mode %q", state.WritingMode)
	}

	*loc = Locale(state)
	return nil
}

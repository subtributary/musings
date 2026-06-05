package tokens

import (
	"unicode"
	"unicode/utf8"
)

// Scripts breaks text into script tokens.
// Unlike other tokenizers, it returns the script name along with the text.
type Scripts struct{}

type ScriptToken struct {
	Script string
	Text   string
}

func (t Scripts) Tokens(text string) []ScriptToken {
	var results []ScriptToken

	var script string
	var length int
	for ; len(text) > 0; text = text[length:] {
		script, length = currentScript(text)

		// Skip text that is entirely in the Common script.
		// This can happen at the start of some texts.
		if script == "" {
			continue
		}

		results = append(results, ScriptToken{
			Script: script,
			Text:   text[:length],
		})
	}

	return results
}

// currentScript returns the current script's name and length.
// Common script characters are consumed without breaking.
func currentScript(text string) (string, int) {
	c, n := utf8.DecodeRuneInString(text)
	initScript := detectScript(c)

	i := n
	for ; i < len(text); i += n {
		c, n = utf8.DecodeRuneInString(text[i:])
		script := detectScript(c)

		if script == "" {
			// Common script doesn't break.
			continue
		}

		if script != initScript {
			break
		}
	}

	return initScript, i
}

// detectScript returns the Unicode script of the character.
// If the character should inherit the previous script, "" is returned.
func detectScript(c rune) string {
	if !unicode.In(c, unicode.Common, unicode.Inherited) {
		for name, rt := range unicode.Scripts {
			if unicode.Is(rt, c) {
				return name
			}
		}
	}
	return ""
}

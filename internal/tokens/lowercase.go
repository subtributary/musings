package tokens

import "strings"

type Lowercase struct{}

func (n Lowercase) Tokens(text string) []string {
	// This lowercasing algorithm is not accurate for all characters; for
	// example, the Turkish "I" is incorrectly lowercased to "i". Such errors
	// do not matter for search applications since they are consistent and do
	// not remove information.
	return []string{strings.ToLower(text)}
}

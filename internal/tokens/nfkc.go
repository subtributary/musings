package tokens

import "golang.org/x/text/unicode/norm"

type NFKC struct{}

func (n NFKC) Tokens(text string) []string {
	return []string{norm.NFKC.String(text)}
}

package tokens

type NGram struct {
	MinN int `json:"min_n"`
	MaxN int `json:"max_n"`
}

func (t NGram) Tokens(text string) []string {
	runes := []rune(text)

	// Text shorter than the minimum is just returned.
	if len(runes) == 0 {
		return nil
	}
	if len(runes) < t.MinN {
		return []string{text}
	}

	// maxN can't be longer than the text.
	maxN := min(t.MaxN, len(runes))

	// Return overlapping chunks of sizes from minN to maxN.
	var result []string
	for n := t.MinN; n <= maxN; n++ {
		for i := 0; i+n <= len(runes); i++ {
			result = append(result, string(runes[i:i+n]))
		}
	}
	return result
}

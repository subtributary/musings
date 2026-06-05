package tokens

import (
	"encoding/json"
	"fmt"
	"maps"
	"slices"
)

type Analyzer interface {
	Tokens(token string) []string
}

// analyzerSaveState holds Analyzer data in a way that can be (un)marshaled.
type analyzerSaveState struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

func newAnalyzerSaveState(a Analyzer) (s analyzerSaveState, err error) {
	switch a.(type) {
	case Lowercase:
		s.Type = "lowercase"
	case NFKC:
		s.Type = "nfkc"
	case NGram:
		if s.Data, err = json.Marshal(a.(NGram)); err != nil {
			err = fmt.Errorf("marshal ngram: %v", err)
		}
		s.Type = "ngram"
	case UAX29:
		s.Type = "uax29"
	default:
		err = fmt.Errorf("unknown analyzer type")
	}
	return
}

func (s analyzerSaveState) ToAnalyzer() (Analyzer, error) {
	switch s.Type {
	case "lowercase":
		return Lowercase{}, nil
	case "nfkc":
		return NFKC{}, nil
	case "ngram":
		var ngram NGram
		if err := json.Unmarshal(s.Data, &ngram); err != nil {
			return nil, err
		}
		return ngram, nil
	case "uax29":
		return UAX29{}, nil
	default:
		return nil, fmt.Errorf("unknown analyzer type: %s", s.Type)
	}
}

// Smart tokenizes and normalizes text by delegating work to other analyzers.
// Which are used depends on the scripts detected while parsing the text.
type Smart struct {
	subs map[string][]Analyzer
}

func (t *Smart) Analyzers(script string) []Analyzer {
	return t.subs[script]
}

func (t *Smart) SetAnalyzers(script string, analyzers []Analyzer) {
	if t.subs == nil {
		t.subs = make(map[string][]Analyzer)
	}
	t.subs[script] = analyzers
}

func (t *Smart) Scripts() []string {
	return slices.Collect(maps.Keys(t.subs))
}

func (t *Smart) Tokens(text string) []string {
	var results []string

	for _, script := range (Scripts{}).Tokens(text) {
		analyzers, ok := t.subs[script.Script]
		if !ok {
			// If no analyzers are configured for the script, ignore the text.
			continue
		}

		// Init to the whole text of the script segment.
		scriptTokens := []string{script.Text}

		// Iteratively break down the text into subtokens.
		for _, analyzer := range analyzers {
			var layerTokens []string
			for _, token := range scriptTokens {
				layerTokens = append(layerTokens, analyzer.Tokens(token)...)
			}
			scriptTokens = layerTokens
		}

		results = append(results, scriptTokens...)
	}

	return results
}

func (t *Smart) MarshalJSON() ([]byte, error) {
	state := make(map[string][]analyzerSaveState)
	for script, subs := range t.subs {
		var subsState []analyzerSaveState
		for _, sub := range subs {
			subState, err := newAnalyzerSaveState(sub)
			if err != nil {
				return nil, err
			}
			subsState = append(subsState, subState)
		}
		state[script] = subsState
	}
	return json.Marshal(state)
}

func (t *Smart) UnmarshalJSON(data []byte) error {
	var state map[string][]analyzerSaveState
	if err := json.Unmarshal(data, &state); err != nil {
		return fmt.Errorf("unmarshal smart tokenizer: %v", err)
	}

	t.subs = make(map[string][]Analyzer)
	for script, subsState := range state {
		var analyzers []Analyzer
		for _, subState := range subsState {
			analyzer, err := subState.ToAnalyzer()
			if err != nil {
				return err
			}
			analyzers = append(analyzers, analyzer)
		}
		t.subs[script] = analyzers
	}

	return nil
}

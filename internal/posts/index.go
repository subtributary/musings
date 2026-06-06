package posts

import (
	"cmp"
	"encoding/json"
	"fmt"
	"iter"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/computerghost/bm25f"
	"github.com/subtributary/musings/internal/tokens"
	"golang.org/x/text/language"
)

const (
	fieldTitle   = "title"
	fieldContent = "content"

	metadataTitle = "title"
	metadataPath  = "path"
)

type Index struct {
	ranker    *bm25f.BM25F
	corpus    *bm25f.Corpus
	locale    language.Tag
	tokenizer *tokens.Smart
}

func NewIndex(locale language.Tag) *Index {
	idx := &Index{}

	idx.ranker = bm25f.New()
	_ = idx.ranker.SetWeight(fieldTitle, 5.0)
	_ = idx.ranker.SetWeight(fieldContent, 1.0)

	idx.corpus = bm25f.NewCorpus()

	idx.locale = locale

	idx.tokenizer = &tokens.Smart{}
	idx.tokenizer.SetAnalyzers("Adlam", []tokens.Analyzer{tokens.UAX29{}, tokens.NFKC{}, tokens.Lowercase{}})
	idx.tokenizer.SetAnalyzers("Ahom", []tokens.Analyzer{tokens.UAX29{}, tokens.NFKC{}, tokens.NGram{MinN: 2, MaxN: 3}})
	idx.tokenizer.SetAnalyzers("Arabic", []tokens.Analyzer{tokens.UAX29{}, tokens.NFKC{}})
	idx.tokenizer.SetAnalyzers("Armenian", []tokens.Analyzer{tokens.UAX29{}, tokens.NFKC{}, tokens.Lowercase{}})
	idx.tokenizer.SetAnalyzers("Balinese", []tokens.Analyzer{tokens.UAX29{}, tokens.NFKC{}, tokens.NGram{MinN: 2, MaxN: 3}})
	idx.tokenizer.SetAnalyzers("Batak", []tokens.Analyzer{tokens.UAX29{}, tokens.NFKC{}, tokens.NGram{MinN: 2, MaxN: 3}})
	idx.tokenizer.SetAnalyzers("Bengali", []tokens.Analyzer{tokens.UAX29{}, tokens.NFKC{}})
	idx.tokenizer.SetAnalyzers("Bopomofo", []tokens.Analyzer{tokens.UAX29{}, tokens.NFKC{}, tokens.NGram{MinN: 2, MaxN: 3}})
	idx.tokenizer.SetAnalyzers("Braille", []tokens.Analyzer{tokens.UAX29{}, tokens.NFKC{}, tokens.NGram{MinN: 2, MaxN: 3}})
	idx.tokenizer.SetAnalyzers("Buginese", []tokens.Analyzer{tokens.UAX29{}, tokens.NFKC{}})
	idx.tokenizer.SetAnalyzers("Buhid", []tokens.Analyzer{tokens.UAX29{}, tokens.NFKC{}})
	idx.tokenizer.SetAnalyzers("Canadian_Aboriginal", []tokens.Analyzer{tokens.UAX29{}, tokens.NFKC{}})
	idx.tokenizer.SetAnalyzers("Chakma", []tokens.Analyzer{tokens.UAX29{}, tokens.NFKC{}})
	idx.tokenizer.SetAnalyzers("Cham", []tokens.Analyzer{tokens.UAX29{}, tokens.NFKC{}, tokens.NGram{MinN: 2, MaxN: 3}})
	idx.tokenizer.SetAnalyzers("Cherokee", []tokens.Analyzer{tokens.UAX29{}, tokens.NFKC{}, tokens.Lowercase{}})
	idx.tokenizer.SetAnalyzers("Cyrillic", []tokens.Analyzer{tokens.UAX29{}, tokens.NFKC{}, tokens.Lowercase{}})
	idx.tokenizer.SetAnalyzers("Devanagari", []tokens.Analyzer{tokens.UAX29{}, tokens.NFKC{}})
	idx.tokenizer.SetAnalyzers("Duployan", []tokens.Analyzer{tokens.UAX29{}, tokens.NFKC{}})
	idx.tokenizer.SetAnalyzers("Ethiopic", []tokens.Analyzer{tokens.UAX29{}, tokens.NFKC{}})
	idx.tokenizer.SetAnalyzers("Georgian", []tokens.Analyzer{tokens.UAX29{}, tokens.NFKC{}, tokens.Lowercase{}})
	idx.tokenizer.SetAnalyzers("Grantha", []tokens.Analyzer{tokens.UAX29{}, tokens.NFKC{}})
	idx.tokenizer.SetAnalyzers("Greek", []tokens.Analyzer{tokens.UAX29{}, tokens.NFKC{}, tokens.Lowercase{}})
	idx.tokenizer.SetAnalyzers("Gujarati", []tokens.Analyzer{tokens.UAX29{}, tokens.NFKC{}})
	idx.tokenizer.SetAnalyzers("Gunjala_Gondi", []tokens.Analyzer{tokens.UAX29{}, tokens.NFKC{}})
	idx.tokenizer.SetAnalyzers("Gurmukhi", []tokens.Analyzer{tokens.UAX29{}, tokens.NFKC{}})
	idx.tokenizer.SetAnalyzers("Han", []tokens.Analyzer{tokens.UAX29{}, tokens.NFKC{}, tokens.NGram{MinN: 1, MaxN: 2}})
	idx.tokenizer.SetAnalyzers("Hangul", []tokens.Analyzer{tokens.UAX29{}, tokens.NFKC{}})
	idx.tokenizer.SetAnalyzers("Hanifi_Rohingya", []tokens.Analyzer{tokens.UAX29{}, tokens.NFKC{}})
	idx.tokenizer.SetAnalyzers("Hanunoo", []tokens.Analyzer{tokens.UAX29{}, tokens.NFKC{}})
	idx.tokenizer.SetAnalyzers("Hebrew", []tokens.Analyzer{tokens.UAX29{}, tokens.NFKC{}})
	idx.tokenizer.SetAnalyzers("Hiragana", []tokens.Analyzer{tokens.UAX29{}, tokens.NFKC{}, tokens.NGram{MinN: 2, MaxN: 3}})
	idx.tokenizer.SetAnalyzers("Javanese", []tokens.Analyzer{tokens.UAX29{}, tokens.NFKC{}, tokens.NGram{MinN: 2, MaxN: 3}})
	idx.tokenizer.SetAnalyzers("Kannada", []tokens.Analyzer{tokens.UAX29{}, tokens.NFKC{}})
	idx.tokenizer.SetAnalyzers("Katakana", []tokens.Analyzer{tokens.UAX29{}, tokens.NFKC{}, tokens.NGram{MinN: 2, MaxN: 3}})
	idx.tokenizer.SetAnalyzers("Kayah_Li", []tokens.Analyzer{tokens.UAX29{}, tokens.NFKC{}, tokens.NGram{MinN: 2, MaxN: 3}})
	idx.tokenizer.SetAnalyzers("Khmer", []tokens.Analyzer{tokens.UAX29{}, tokens.NFKC{}, tokens.NGram{MinN: 2, MaxN: 3}})
	idx.tokenizer.SetAnalyzers("Khojki", []tokens.Analyzer{tokens.UAX29{}, tokens.NFKC{}})
	idx.tokenizer.SetAnalyzers("Lao", []tokens.Analyzer{tokens.UAX29{}, tokens.NFKC{}, tokens.NGram{MinN: 2, MaxN: 3}})
	idx.tokenizer.SetAnalyzers("Latin", []tokens.Analyzer{tokens.UAX29{}, tokens.NFKC{}, tokens.Lowercase{}})
	idx.tokenizer.SetAnalyzers("Lepcha", []tokens.Analyzer{tokens.UAX29{}, tokens.NFKC{}})
	idx.tokenizer.SetAnalyzers("Limbu", []tokens.Analyzer{tokens.UAX29{}, tokens.NFKC{}})
	idx.tokenizer.SetAnalyzers("Lisu", []tokens.Analyzer{tokens.UAX29{}, tokens.NFKC{}})
	idx.tokenizer.SetAnalyzers("Malayalam", []tokens.Analyzer{tokens.UAX29{}, tokens.NFKC{}})
	idx.tokenizer.SetAnalyzers("Mandaic", []tokens.Analyzer{tokens.UAX29{}, tokens.NFKC{}})
	idx.tokenizer.SetAnalyzers("Masaram_Gondi", []tokens.Analyzer{tokens.UAX29{}, tokens.NFKC{}})
	idx.tokenizer.SetAnalyzers("Medefaidrin", []tokens.Analyzer{tokens.UAX29{}, tokens.NFKC{}, tokens.Lowercase{}})
	idx.tokenizer.SetAnalyzers("Meetei_Mayek", []tokens.Analyzer{tokens.UAX29{}, tokens.NFKC{}})
	idx.tokenizer.SetAnalyzers("Mende_Kikakui", []tokens.Analyzer{tokens.UAX29{}, tokens.NFKC{}})
	idx.tokenizer.SetAnalyzers("Miao", []tokens.Analyzer{tokens.UAX29{}, tokens.NFKC{}, tokens.NGram{MinN: 2, MaxN: 3}})
	idx.tokenizer.SetAnalyzers("Modi", []tokens.Analyzer{tokens.UAX29{}, tokens.NFKC{}, tokens.Lowercase{}})
	idx.tokenizer.SetAnalyzers("Mongolian", []tokens.Analyzer{tokens.UAX29{}, tokens.NFKC{}})
	idx.tokenizer.SetAnalyzers("Mro", []tokens.Analyzer{tokens.UAX29{}, tokens.NFKC{}})
	idx.tokenizer.SetAnalyzers("Myanmar", []tokens.Analyzer{tokens.UAX29{}, tokens.NFKC{}, tokens.NGram{MinN: 2, MaxN: 3}})
	idx.tokenizer.SetAnalyzers("Nag_Mundari", []tokens.Analyzer{tokens.UAX29{}, tokens.NFKC{}})
	idx.tokenizer.SetAnalyzers("New_Tai_Lue", []tokens.Analyzer{tokens.UAX29{}, tokens.NFKC{}, tokens.NGram{MinN: 2, MaxN: 3}})
	idx.tokenizer.SetAnalyzers("Newa", []tokens.Analyzer{tokens.UAX29{}, tokens.NFKC{}})
	idx.tokenizer.SetAnalyzers("Nko", []tokens.Analyzer{tokens.UAX29{}, tokens.NFKC{}})
	idx.tokenizer.SetAnalyzers("Nushu", []tokens.Analyzer{tokens.UAX29{}, tokens.NFKC{}, tokens.NGram{MinN: 2, MaxN: 3}})
	idx.tokenizer.SetAnalyzers("Nyiakeng_Puachue_Hmong", []tokens.Analyzer{tokens.UAX29{}, tokens.NFKC{}})
	idx.tokenizer.SetAnalyzers("Ol_Chiki", []tokens.Analyzer{tokens.UAX29{}, tokens.NFKC{}})
	idx.tokenizer.SetAnalyzers("Old_Hungarian", []tokens.Analyzer{tokens.UAX29{}, tokens.NFKC{}, tokens.Lowercase{}})
	idx.tokenizer.SetAnalyzers("Oriya", []tokens.Analyzer{tokens.UAX29{}, tokens.NFKC{}})
	idx.tokenizer.SetAnalyzers("Osage", []tokens.Analyzer{tokens.UAX29{}, tokens.NFKC{}, tokens.Lowercase{}})
	idx.tokenizer.SetAnalyzers("Pahawh_Hmong", []tokens.Analyzer{tokens.UAX29{}, tokens.NFKC{}})
	idx.tokenizer.SetAnalyzers("Pau_Cin_Hau", []tokens.Analyzer{tokens.UAX29{}, tokens.NFKC{}})
	idx.tokenizer.SetAnalyzers("Rejang", []tokens.Analyzer{tokens.UAX29{}, tokens.NFKC{}, tokens.NGram{MinN: 2, MaxN: 3}})
	idx.tokenizer.SetAnalyzers("Samaritan", []tokens.Analyzer{tokens.UAX29{}, tokens.NFKC{}})
	idx.tokenizer.SetAnalyzers("Saurashtra", []tokens.Analyzer{tokens.UAX29{}, tokens.NFKC{}})
	idx.tokenizer.SetAnalyzers("Shavian", []tokens.Analyzer{tokens.UAX29{}, tokens.NFKC{}})
	idx.tokenizer.SetAnalyzers("Siddham", []tokens.Analyzer{tokens.UAX29{}, tokens.NFKC{}})
	idx.tokenizer.SetAnalyzers("SignWriting", []tokens.Analyzer{tokens.UAX29{}, tokens.NFKC{}})
	idx.tokenizer.SetAnalyzers("Sinhala", []tokens.Analyzer{tokens.UAX29{}, tokens.NFKC{}})
	idx.tokenizer.SetAnalyzers("Sogdian", []tokens.Analyzer{tokens.UAX29{}, tokens.NFKC{}})
	idx.tokenizer.SetAnalyzers("Sora_Sompeng", []tokens.Analyzer{tokens.UAX29{}, tokens.NFKC{}})
	idx.tokenizer.SetAnalyzers("Sundanese", []tokens.Analyzer{tokens.UAX29{}, tokens.NFKC{}})
	idx.tokenizer.SetAnalyzers("Syloti_Nagri", []tokens.Analyzer{tokens.UAX29{}, tokens.NFKC{}})
	idx.tokenizer.SetAnalyzers("Syriac", []tokens.Analyzer{tokens.UAX29{}, tokens.NFKC{}})
	idx.tokenizer.SetAnalyzers("Tagalog", []tokens.Analyzer{tokens.UAX29{}, tokens.NFKC{}})
	idx.tokenizer.SetAnalyzers("Tagbanwa", []tokens.Analyzer{tokens.UAX29{}, tokens.NFKC{}})
	idx.tokenizer.SetAnalyzers("Tai_Le", []tokens.Analyzer{tokens.UAX29{}, tokens.NFKC{}, tokens.NGram{MinN: 2, MaxN: 3}})
	idx.tokenizer.SetAnalyzers("Tai_Tham", []tokens.Analyzer{tokens.UAX29{}, tokens.NFKC{}, tokens.NGram{MinN: 2, MaxN: 3}})
	idx.tokenizer.SetAnalyzers("Tai_Viet", []tokens.Analyzer{tokens.UAX29{}, tokens.NFKC{}, tokens.NGram{MinN: 2, MaxN: 3}})
	idx.tokenizer.SetAnalyzers("Tamil", []tokens.Analyzer{tokens.UAX29{}, tokens.NFKC{}})
	idx.tokenizer.SetAnalyzers("Tangsa", []tokens.Analyzer{tokens.UAX29{}, tokens.NFKC{}})
	idx.tokenizer.SetAnalyzers("Telugu", []tokens.Analyzer{tokens.UAX29{}, tokens.NFKC{}})
	idx.tokenizer.SetAnalyzers("Thaana", []tokens.Analyzer{tokens.UAX29{}, tokens.NFKC{}, tokens.Lowercase{}})
	idx.tokenizer.SetAnalyzers("Thai", []tokens.Analyzer{tokens.UAX29{}, tokens.NFKC{}, tokens.NGram{MinN: 2, MaxN: 3}})
	idx.tokenizer.SetAnalyzers("Tibetan", []tokens.Analyzer{tokens.UAX29{}, tokens.NFKC{}})
	idx.tokenizer.SetAnalyzers("Tifinagh", []tokens.Analyzer{tokens.UAX29{}, tokens.NFKC{}})
	idx.tokenizer.SetAnalyzers("Tirhuta", []tokens.Analyzer{tokens.UAX29{}, tokens.NFKC{}})
	idx.tokenizer.SetAnalyzers("Toto", []tokens.Analyzer{tokens.UAX29{}, tokens.NFKC{}})
	idx.tokenizer.SetAnalyzers("Vai", []tokens.Analyzer{tokens.UAX29{}, tokens.NFKC{}})
	idx.tokenizer.SetAnalyzers("Wancho", []tokens.Analyzer{tokens.UAX29{}, tokens.NFKC{}})
	idx.tokenizer.SetAnalyzers("Warang_Citi", []tokens.Analyzer{tokens.UAX29{}, tokens.NFKC{}, tokens.Lowercase{}})
	idx.tokenizer.SetAnalyzers("Yezidi", []tokens.Analyzer{tokens.UAX29{}, tokens.NFKC{}})
	idx.tokenizer.SetAnalyzers("Yi", []tokens.Analyzer{tokens.UAX29{}, tokens.NFKC{}, tokens.NGram{MinN: 1, MaxN: 2}})

	return idx
}

func LoadIndex(dataPath string, locale language.Tag) (*Index, error) {
	path := indexDataFile(dataPath, locale)
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read index file: %v", err)
	}

	var index Index
	if err := json.Unmarshal(data, &index); err != nil {
		return nil, fmt.Errorf("unmarshal index: %v", err)
	}

	index.locale = locale

	return &index, nil
}

func (idx *Index) Save(dataPath string) error {
	data, err := json.Marshal(idx)
	if err != nil {
		return fmt.Errorf("marshal index: %v", err)
	}

	path := indexDataFile(dataPath, idx.locale)
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("save index file: %v", err)
	}
	return nil
}

func indexDataFile(dataPath string, locale language.Tag) string {
	if locale == language.Und {
		return filepath.Join(dataPath, "index.json")
	}
	localeStr := strings.ToLower(locale.String())
	return filepath.Join(dataPath, "index."+localeStr+".json")
}

type indexSaveState struct {
	BM        *bm25f.BM25F  `json:"bm"`
	Corpus    *bm25f.Corpus `json:"corpus"`
	Tokenizer *tokens.Smart `json:"tokenizer"`
}

func (idx *Index) MarshalJSON() ([]byte, error) {
	return json.Marshal(indexSaveState{
		BM:        idx.ranker,
		Corpus:    idx.corpus,
		Tokenizer: idx.tokenizer,
	})
}

func (idx *Index) UnmarshalJSON(data []byte) error {
	var state indexSaveState
	if err := json.Unmarshal(data, &state); err != nil {
		return err
	}

	if state.BM == nil {
		return fmt.Errorf("index is missing BM25F state")
	}
	if state.Corpus == nil {
		return fmt.Errorf("index is missing corpus")
	}
	if state.Tokenizer == nil {
		return fmt.Errorf("index is missing tokenizer")
	}

	idx.ranker = state.BM
	idx.corpus = state.Corpus
	idx.tokenizer = state.Tokenizer
	return nil
}

// Upsert adds the parsed post to the index.
//
// This method is not safe for concurrent calls to Search, List, or Save.
func (idx *Index) Upsert(path string, post ParsedPost) {
	path = strings.TrimSuffix(path, ".md")

	title := idx.tokenizer.Tokens(post.Title)
	content := idx.tokenizer.Tokens(string(post.Content))

	idx.corpus.Upsert(path, bm25f.NewDocument(
		bm25f.WithField(fieldTitle, title),
		bm25f.WithField(fieldContent, content),
		bm25f.WithMetadata(metadataTitle, post.Title),
		bm25f.WithMetadata(metadataPath, path),
	))
}

type IndexedPost struct {
	Path  string
	Title string
}

func indexedPostFromMetadata(meta interface{ Metadata(string) (string, bool) }) (IndexedPost, bool) {
	path, ok := meta.Metadata("path")
	if !ok || path == "" {
		return IndexedPost{}, false
	}

	title, ok := meta.Metadata("title")
	if !ok {
		return IndexedPost{}, false
	}

	return IndexedPost{Path: path, Title: title}, true
}

func (idx *Index) List() iter.Seq[IndexedPost] {
	var results []IndexedPost
	for _, id := range idx.corpus.DocumentIDs() {
		doc, _ := idx.corpus.Document(id)
		title, _ := doc.Metadata(metadataTitle)
		path, _ := doc.Metadata(metadataPath)
		results = append(results, IndexedPost{Path: path, Title: title})
	}

	slices.SortFunc(results, func(a, b IndexedPost) int {
		if c := cmp.Compare(a.Path, b.Path); c != 0 {
			return c
		}
		return cmp.Compare(a.Title, b.Title)
	})

	return func(yield func(p IndexedPost) bool) {
		for _, post := range results {
			if !yield(post) {
				return
			}
		}
	}
}

func (idx *Index) Search(query string) iter.Seq[IndexedPost] {
	if strings.TrimSpace(query) == "" {
		return idx.List()
	}

	queryToks := idx.tokenizer.Tokens(query)
	scores := idx.ranker.Score(idx.corpus, queryToks)

	scores = slices.DeleteFunc(scores, func(r bm25f.Result) bool {
		return r.Score == 0
	})

	slices.SortFunc(scores, func(a, b bm25f.Result) int {
		if c := cmp.Compare(b.Score, a.Score); c != 0 {
			return c
		}
		return cmp.Compare(a.ID, b.ID)
	})

	return func(yield func(p IndexedPost) bool) {
		for _, result := range scores {
			post, _ := indexedPostFromMetadata(&result)
			if !yield(post) {
				return
			}
		}
	}
}

package search_test

import (
	"slices"
	"testing"

	"github.com/subtributary/musings/internal/search"
)

func TestDocumentTokenizer(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		text string
		want []string
	}{
		/* basic parsing */
		{
			name: "empty",
			text: "",
			want: []string{},
		},
		{
			name: "one word",
			text: "word",
			want: []string{"word"},
		},
		{
			name: "uppercase",
			text: "CAPITALS",
			want: []string{"capitals"},
		},
		{
			name: "numbers and letters",
			text: "88 wb9 10wb o_o",
			want: []string{"88", "wb9", "10wb", "o_o"},
		},
		{
			name: "emojis",
			text: "melt \U0001fae0 zwj\U0001f468\u200d\U0001f9b0nospace",
			want: []string{"melt", "zwj", "nospace"},
		},
		{
			name: "format and extend characters",
			text: "눈〯 ac\u0327ai hy\u00adphen z\u200bw",
			want: []string{"눈〯", "açai", "hyphen", "zw"},
		},
		{
			// Not extensive, just check if norm.NFKC seems to be used.
			name: "NFKC",
			text: "\u212b n\u0303 ﬀ",
			want: []string{"\u00e5", "\u00f1", "ff"},
		},
		{
			// These are wrong, but any tokens with them probably don't matter.
			name: "unimportant normalizations",
			text: "½",
			want: []string{"1\u20442"},
		},

		/* scripts */
		{
			name: "script change",
			text: "English를 배웠다",
			want: []string{"english", "를", "배웠다"},
		},
		{
			name: "scripts with different tokenization algorithms",
			text: "中文를 배웠다",
			want: []string{"中", "文", "를", "배웠다"},
		},
		{
			name: "Han",
			text: "北京大学",
			want: []string{"北", "京", "大", "学"},
		},
		{
			name: "Kana",
			text: "カタカナ",
			want: []string{"カタ", "タカ", "カナ", "カタカ", "タカナ"},
		},
		{
			name: "Hebrew",
			text: "אָלֶף־בֵּית עִבְרִי",
			want: []string{"אָלֶף", "בֵּית", "עִבְרִי"},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			subject := search.NewDocumentTokenizer(tt.text)
			got := slices.Collect(subject.Tokens())
			if !slices.Equal(got, tt.want) {
				t.Errorf("Tokens: got %#v, want: %#v", got, tt.want)
			}
		})
	}
}

func TestNGramTokenizer(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		text string
		minN int
		maxN int
		want []string
	}{
		{
			name: "unigram",
			text: "1234",
			minN: 1,
			maxN: 1,
			want: []string{"1", "2", "3", "4"},
		},
		{
			name: "bigram",
			text: "1234",
			minN: 2,
			maxN: 2,
			want: []string{"12", "23", "34"},
		},
		{
			name: "trigram",
			text: "1234",
			minN: 3,
			maxN: 3,
			want: []string{"123", "234"},
		},
		{
			name: "bigram and trigram",
			text: "1234",
			minN: 2,
			maxN: 3,
			want: []string{"12", "23", "34", "123", "234"},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			subject := search.NewNGramTokenizer(tt.text, tt.minN, tt.maxN)
			got := slices.Collect(subject.Tokens())
			if !slices.Equal(got, tt.want) {
				t.Errorf("Tokens: got %#v, want: %#v", got, tt.want)
			}
		})
	}
}

func TestQueryTokenizer(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		text string
		want []string
	}{
		{
			name: "empty",
			text: "",
			want: []string{},
		},
		{
			name: "one word",
			text: "word",
			want: []string{"word"},
		},
		{
			name: "tag and non-tag",
			text: "#word word#",
			want: []string{"#word", "word"},
		},
		{
			name: "empty tags",
			text: "# #",
			want: []string{},
		},
		{
			name: "uppercase",
			text: "CAPITALS #CAPITALS",
			want: []string{"capitals", "#capitals"},
		},
		{
			name: "numbers and letters",
			text: "88 wb9 10wb o_o #88 #wb9 #10wb #_o",
			want: []string{
				"88", "wb9", "10wb", "o_o",
				"#88", "#wb9", "#10wb", "#_o",
			},
		},
		{
			// Not extensive, just check if norm.NFKC seems to be used.
			name: "NFKC",
			text: "\u212b n\u0303 ﬀ #\u212b #ﬀ",
			want: []string{"\u00e5", "\u00f1", "ff", "#\u00e5", "#ff"},
		},
		{
			// These are wrong, but any tokens with them probably don't matter.
			name: "unimportant normalizations",
			text: "½ #½",
			want: []string{"1\u20442", "#1\u20442"},
		},
		{
			name: "script change",
			text: "English를 배웠다 #English배웠다",
			want: []string{"english", "를", "배웠다", "#english배웠다"},
		},
		{
			name: "scripts with different tokenization algorithms",
			text: "中文를 배웠다 #中文배웠다",
			want: []string{"中", "文", "를", "배웠다", "#中文배웠다"},
		},
		{
			name: "Han",
			text: "北京大学 #北京大学",
			want: []string{"北", "京", "大", "学", "#北京大学"},
		},
		{
			name: "Kana",
			text: "カタカナ #カタカナ",
			want: []string{"カタ", "タカ", "カナ", "カタカ", "タカナ", "#カタカナ"},
		},
		{
			name: "Hebrew",
			text: "אָלֶף־בֵּית עִבְרִי #אָלֶף_בֵּית_עִבְרִי",
			// This script is RTL, so this array may appear visually backwards.
			want: []string{"אָלֶף", "בֵּית", "עִבְרִי", "#אָלֶף_בֵּית_עִבְרִי"},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			subject := search.NewQueryTokenizer(tt.text)
			got := slices.Collect(subject.Tokens())
			if !slices.Equal(got, tt.want) {
				t.Errorf("Tokens: got %#v, want: %#v", got, tt.want)
			}
		})
	}
}

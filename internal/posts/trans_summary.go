package posts

import (
	"bytes"
	"strings"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

var summaryKey = parser.NewContextKey()

func getSummary(pc parser.Context) string {
	if summary := pc.Get(summaryKey); summary != nil {
		return summary.(string)
	}
	return ""
}

func setSummary(pc parser.Context, summary string) {
	pc.Set(summaryKey, summary)
}

// summaryTransformer stores the first paragraph of the document so that it can
// be referenced via getSummary.
type summaryTransformer struct{}

func (t *summaryTransformer) Transform(doc *ast.Document, reader text.Reader, pc parser.Context) {
	for n := doc.FirstChild(); n != nil; n = n.NextSibling() {
		paragraph, ok := n.(*ast.Paragraph)
		if !ok {
			continue
		}

		summary := paragraphPlainText(paragraph, reader.Source())
		if summary == "" {
			continue
		}

		setSummary(pc, summary)
		break
	}
}

func paragraphPlainText(paragraph *ast.Paragraph, source []byte) string {
	var buf bytes.Buffer

	err := ast.Walk(paragraph, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}

		switch n := n.(type) {
		case *ast.Text:
			buf.Write(n.Value(source))
			if n.SoftLineBreak() || n.HardLineBreak() {
				buf.WriteByte(' ')
			}
		case *ast.String:
			buf.Write(n.Value)
		}

		return ast.WalkContinue, nil
	})

	// We should never get an error since we never return one in the walker.
	if err != nil {
		panic("Unexpected: error walking AST:" + err.Error())
	}

	return strings.TrimSpace(buf.String())
}

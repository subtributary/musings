package posts

import (
	"bytes"
	"log"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

var titleKey = parser.NewContextKey()

func getTitle(pc parser.Context) string {
	if title := pc.Get(titleKey); title != nil {
		return title.(string)
	}
	return ""
}

func setTitle(pc parser.Context, title string) {
	pc.Set(titleKey, title)
}

// removeH1Transformer removes the first primary heading from the document and
// uses its text to set the document title in the context.
type removeH1Transformer struct{}

func (t *removeH1Transformer) Transform(doc *ast.Document, reader text.Reader, pc parser.Context) {
	err := ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}

		heading, ok := n.(*ast.Heading)
		if !ok || heading.Level != 1 {
			return ast.WalkContinue, nil
		}

		// Only use top-level H1 as title
		if heading.Parent() != doc {
			return ast.WalkContinue, nil
		}

		setTitle(pc, headingPlainText(heading, reader.Source()))

		heading.Parent().RemoveChild(heading.Parent(), heading)
		return ast.WalkStop, nil
	})

	// We should never get an error since we never return one in the walker.
	if err != nil {
		log.Fatalf("Unexpected error walking AST: %v", err)
	}
}

// headingPlainText extracts the text from a heading, ignoring any formatting.
func headingPlainText(heading *ast.Heading, source []byte) string {
	var buf bytes.Buffer

	err := ast.Walk(heading, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
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
		case *ast.CodeSpan:
			for i := 0; i < n.Lines().Len(); i++ {
				segment := n.Lines().At(i)
				buf.Write(segment.Value(source))
			}
		}

		return ast.WalkContinue, nil
	})

	// We should never get an error since we never return one in the walker.
	if err != nil {
		log.Fatalf("Unexpected error walking AST: %v", err)
	}

	return buf.String()
}

package posts

import (
	"bytes"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"time"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
	"go.abhg.dev/goldmark/frontmatter"
)

type ParsedPost struct {
	Title     string
	Content   template.HTML
	Bylines   []string
	Published *time.Time
}

type Parser struct {
	docParser goldmark.Markdown
	markdown  goldmark.Markdown
}

func NewParser() Parser {
	docParser := goldmark.New(
		goldmark.WithExtensions(
			&frontmatter.Extender{},
		),
		goldmark.WithParserOptions(
			parser.WithASTTransformers(
				util.Prioritized(&removeH1Transformer{}, 0),
			),
		),
	)
	return Parser{
		docParser: docParser,
		markdown:  goldmark.New(),
	}
}

func (s Parser) ParseFile(dir fs.FS, path string) (ParsedPost, error) {
	if content, err := fs.ReadFile(dir, path); err != nil {
		return ParsedPost{}, fmt.Errorf("read file: %v", err)
	} else {
		return s.ParseContent(content)
	}
}

func (s Parser) ParseContent(content []byte) (ParsedPost, error) {
	context := parser.NewContext()

	var buf bytes.Buffer
	if err := s.docParser.Convert(content, &buf, parser.WithContext(context)); err != nil {
		return ParsedPost{}, fmt.Errorf("parse content: %w", err)
	}
	parsedContent := string(buf.Bytes())

	fm, err := parseFrontmatter(context)
	if err != nil {
		return ParsedPost{}, fmt.Errorf("parse frontmatter: %w", err)
	}

	return ParsedPost{
		Title:     getTitle(context),
		Content:   template.HTML(parsedContent),
		Bylines:   fm.Bylines,
		Published: parsePostTime(fm.Published),
	}, nil
}

// parsePostTime parses a time string in RFC3339 or a supported format.
// If the string cannot be parsed, nil is returned.
func parsePostTime(value string) *time.Time {
	if t, err := time.Parse(time.RFC3339, value); err == nil {
		return &t
	}
	if t, err := time.Parse(time.DateTime, value); err == nil {
		return &t
	}
	if t, err := time.Parse(time.DateOnly, value); err == nil {
		return &t
	}
	return nil
}

type postFrontmatter struct {
	Bylines   []string
	Published string
}

func parseFrontmatter(context parser.Context) (result postFrontmatter, err error) {
	if fm := frontmatter.Get(context); fm != nil {
		err = fm.Decode(&result)
	}
	return
}

var titleKey = parser.NewContextKey()

func getTitle(pc parser.Context) string {
	if title := pc.Get(titleKey); title != nil {
		return title.(string)
	}
	return ""
}

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

		pc.Set(titleKey, headingPlainText(heading, reader.Source()))

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

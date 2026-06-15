package posts

import (
	"bytes"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/url"
	"path"
	"strconv"
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
	Published time.Time

	// Modified is is only set by Parser.ParseFile which uses the return value
	// of WithModTime's argument. By default it is the time of parsing.
	//
	// Store overrides the behavior to use the actual modified time.
	Modified time.Time
}

type ParserOption func(p *Parser)

func WithModTime(modTime func(name string) time.Time) ParserOption {
	return func(p *Parser) {
		p.modTime = modTime

		transformer := versionAssetsTransformer{modTime: modTime}
		p.docParser.Parser().AddOptions(
			parser.WithASTTransformers(
				util.Prioritized(&transformer, 100),
			),
		)
	}
}

type Parser struct {
	docParser goldmark.Markdown
	modTime   func(name string) time.Time
}

func NewParser(opts ...ParserOption) Parser {
	p := Parser{}

	p.docParser = goldmark.New(
		goldmark.WithExtensions(
			&frontmatter.Extender{},
		),
		goldmark.WithParserOptions(
			parser.WithASTTransformers(
				util.Prioritized(&removeH1Transformer{}, 0),
			),
		),
	)

	p.modTime = func(_ string) time.Time {
		return time.Now()
	}

	for _, opt := range opts {
		opt(&p)
	}

	return p
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
// If the string cannot be parsed, zero time is returned.
func parsePostTime(value string) time.Time {
	if t, err := time.Parse(time.RFC3339, value); err == nil {
		return t
	}
	if t, err := time.Parse(time.DateTime, value); err == nil {
		return t
	}
	if t, err := time.Parse(time.DateOnly, value); err == nil {
		return t
	}
	return time.Time{}
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

// versionAssetsTransformer appends a version timestamp to the query string of
// any link or image that points at a local asset. A local asset is a relative
// URL whose path has a file extension; relative URLs without an extension are
// links to posts, which aren't cached, so they're left alone.
type versionAssetsTransformer struct {
	modTime func(name string) time.Time
}

func (t *versionAssetsTransformer) Transform(doc *ast.Document, reader text.Reader, pc parser.Context) {
	err := ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}

		switch node := n.(type) {
		case *ast.Link:
			node.Destination = t.versionAssetURL(node.Destination)
		case *ast.Image:
			node.Destination = t.versionAssetURL(node.Destination)
		}

		return ast.WalkContinue, nil
	})

	// We should never get an error since we never return one in the walker.
	if err != nil {
		log.Fatalf("Unexpected error walking AST: %v", err)
	}
}

// versionAssetURL appends a version timestamp if it's a local asset URL.
func (t *versionAssetsTransformer) versionAssetURL(dest []byte) []byte {
	u, err := url.Parse(string(dest))
	if err != nil {
		return dest // Not a URL we can work with; leave it alone.
	}

	// Only version local, relative URLs.
	if u.IsAbs() || u.Host != "" {
		return dest
	}

	// No file extension means it's a link to a post, not a cached asset.
	if path.Ext(u.Path) == "" {
		return dest
	}

	modTime := t.modTime(u.Path)
	version := strconv.FormatInt(modTime.Unix(), 16)

	query := u.Query()
	query.Set("v", version)
	u.RawQuery = query.Encode()

	return []byte(u.String())
}

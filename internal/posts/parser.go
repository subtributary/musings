package posts

import (
	"bytes"
	"fmt"
	"html/template"
	"io/fs"
	"strings"
	"time"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/util"
	"go.abhg.dev/goldmark/frontmatter"
)

type ParsedPost struct {
	Title     string
	Content   template.HTML
	Bylines   []string
	Published time.Time
}

type Parser struct {
	docParser goldmark.Markdown
}

// ModTimeFunc returns the modification time of a file. The `name` argument
// will be a path relative to the content directory and prefixed with a '/'.
type ModTimeFunc func(name string) (time.Time, bool)

func NewParser(modTime ModTimeFunc) Parser {
	p := Parser{}

	p.docParser = goldmark.New(
		goldmark.WithExtensions(
			&frontmatter.Extender{},
		),
		goldmark.WithParserOptions(
			parser.WithASTTransformers(
				util.Prioritized(&removeH1Transformer{}, 0),
				util.Prioritized(&versionAssetsTransformer{modTime: modTime}, 100),
			),
		),
	)

	return p
}

func (s Parser) ParseFile(dir fs.FS, name string) (ParsedPost, error) {
	name, _ = strings.CutPrefix(name, "/")

	content, err := fs.ReadFile(dir, name)
	if err != nil {
		return ParsedPost{}, fmt.Errorf("read file: %w", err)
	}

	return s.ParseContent(name, content)
}

func (s Parser) ParseContent(name string, content []byte) (ParsedPost, error) {
	ctx := parser.NewContext()
	setName(ctx, name)

	var buf bytes.Buffer
	err := s.docParser.Convert(content, &buf, parser.WithContext(ctx))
	if err != nil {
		return ParsedPost{}, fmt.Errorf("parse content: %w", err)
	}
	parsedContent := string(buf.Bytes())

	fm, err := parseFrontmatter(ctx)
	if err != nil {
		return ParsedPost{}, fmt.Errorf("parse frontmatter: %w", err)
	}

	return ParsedPost{
		Title:     getTitle(ctx),
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

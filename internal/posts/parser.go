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
	Bylines   []string
	Content   template.HTML
	Published time.Time
	Summary   string
	Thumbnail string
	Title     string
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
				util.Prioritized(&summaryTransformer{}, 50),
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

	summary := getSummary(ctx)
	if fm.Summary != "" {
		summary = fm.Summary
	}

	return ParsedPost{
		Bylines:   fm.Bylines,
		Content:   template.HTML(parsedContent),
		Published: fm.PublishedTime(),
		Summary:   summary,
		Thumbnail: fm.Thumbnail,
		Title:     getTitle(ctx),
	}, nil
}

type postFrontmatter struct {
	Bylines   []string
	Published string
	Summary   string
	Thumbnail string
}

// PublishedTime parses a time string in RFC3339 or a supported format.
// If Published cannot be parsed, zero time is returned.
func (fm *postFrontmatter) PublishedTime() time.Time {
	if t, err := time.Parse(time.RFC3339, fm.Published); err == nil {
		return t
	}
	if t, err := time.Parse(time.DateTime, fm.Published); err == nil {
		return t
	}
	if t, err := time.Parse(time.DateOnly, fm.Published); err == nil {
		return t
	}
	return time.Time{}
}

func parseFrontmatter(context parser.Context) (result postFrontmatter, err error) {
	if fm := frontmatter.Get(context); fm != nil {
		err = fm.Decode(&result)
	}
	return
}

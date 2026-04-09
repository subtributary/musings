package posts

import (
	"bytes"
	"fmt"
	"html/template"
	"io/fs"

	"github.com/yuin/goldmark"
)

type PostData struct {
	HtmlContent template.HTML
}

type Parser struct {
	markdown goldmark.Markdown
}

func NewParser() Parser {
	return Parser{
		markdown: goldmark.New(),
	}
}

func (s Parser) Parse(fileSystem fs.FS, name string) (PostData, error) {
	contents, err := fs.ReadFile(fileSystem, name)
	if err != nil {
		return PostData{}, fmt.Errorf("could not read file %q: %w", name, err)
	}

	buffer := bytes.Buffer{}
	err = s.markdown.Convert(contents, &buffer)
	if err != nil {
		return PostData{}, fmt.Errorf("could not parse file %q: %w", name, err)
	}
	html := buffer.String()

	return PostData{
		HtmlContent: template.HTML(html),
	}, nil
}

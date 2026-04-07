package posts

import (
	"io/fs"

	"github.com/yuin/goldmark"
)

type Parser struct {
	markdown goldmark.Markdown
}

func NewParser() Parser {
	return Parser{
		markdown: goldmark.New(),
	}
}

func (s Parser) Parse(file fs.File) (PostData, error) {
	return PostData{}, nil
}

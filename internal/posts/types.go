package posts

import (
	"html/template"
)

type PostData struct {
	Content string

	HtmlContent template.HTML // Deprecated. Posts shouldn't know about template.*.
}

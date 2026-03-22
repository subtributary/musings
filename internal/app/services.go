package app

import (
	"github.com/subtributary/musings/internal/files"
	"github.com/subtributary/musings/internal/markdown"
	"github.com/subtributary/musings/internal/templates"
)

// Services contains application layer dependencies.
// Dependencies for lower layers should not be included here.
type Services struct {
	ContentStore     *files.Store
	MarkdownStore    *markdown.Store
	StaticStore      *files.Store
	TemplateProvider templates.TemplateProvider
}

package main

import (
	"github.com/subtributary/musings/internal/posts"
	"github.com/subtributary/musings/internal/templates"
)

// Services contains application layer dependencies.
// Dependencies for lower layers should not be included here.
type Services struct {
	PostParser    posts.Parser
	TemplateStore templates.Store
}

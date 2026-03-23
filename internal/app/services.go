package app

import (
	"github.com/subtributary/musings/internal/store"
)

// Services contains application layer dependencies.
// Dependencies for lower layers should not be included here.
type Services struct {
	ContentStore     *store.StaticStore
	PostsStore       *store.PostsStore
	StaticStore      *store.StaticStore
	TemplateProvider store.TemplateProvider
}
